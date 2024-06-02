package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	apiMail "github.com/ainsleyclark/go-mail"
	"github.com/vanng822/go-premailer/premailer"
	smtpMail "github.com/xhit/go-simple-mail/v2"
)

// Mail holds the information necessary to connect to an SMTP server
type Mail struct {
	Domain      string
	Templates   string // Templates is the path to the email templates
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string       // "tls", "ssl", or "none"
	FromAddress string       // default from address
	FromName    string       // default from name
	Jobs        chan Message // Jobs is the channel that holds the messages/mails to be sent
	Results     chan Result  // Results is the channel that holds the results of the sent messages
	API         string       // "smtp", "mailgun", "sparkpost", "sendgrid"
	APIKey      string
	APIUrl      string // e.g https://api.mailgun.net
}

// Message is the type for an email message
type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Template    string      // Template is the name of the email template to be used
	Attachments []string    // Attachments is a slice of file paths to attach to the email
	Data        interface{} // Data is the data to be passed to the email template
}

// Result contains information regarding the status of the sent email message
type Result struct {
	Success bool
	Error   error
}

// ListenForMail listens to Jobs channel channel and sends mail
// when it receives a payload. It runs continually in the background,
// in a separate goroutine and sends error/success messages back on the
// Results channel.
// Note that if api and api key are set, it will prefer using
// an api to send mail
func (m *Mail) ListenForMail() {
	for {
		msg := <-m.Jobs
		err := m.Send(msg)
		if err != nil {
			m.Results <- Result{false, err}
		} else {
			m.Results <- Result{true, nil}
		}
	}
}

// Send sends an email message using correct method. If API values are set,
// it will send using the appropriate api; otherwise, it sends via smtp
func (m *Mail) Send(msg Message) error {
	if len(m.API) > 0 && len(m.APIKey) > 0 && len(m.APIUrl) > 0 && m.API != "smtp" {
		return m.ChooseAPI(msg)
	}
	return m.SendUsingSMTP(msg)
}

// ChooseAPI chooses api (specified in .env) to use to send the mail
// Options: "mailgun", "sparkpost", "sendgrid"
func (m *Mail) ChooseAPI(msg Message) error {
	switch m.API {
	case "mailgun", "sparkpost", "sendgrid":
		return m.SendUsingAPI(msg, m.API)
	default:
		return fmt.Errorf("unknown api %s; only mailgun, sparkpost or sendgrid accepted", m.API)
	}
}

// SendUsingAPI sends a message using the appropriate API. It can be called directly, if necessary.
// transport can be one of sparkpost, sendgrid, or mailgun
func (m *Mail) SendUsingAPI(msg Message, transport string) error {
	if msg.From == "" {
		msg.From = m.FromAddress
	}

	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	cfg := apiMail.Config{
		URL:         m.APIUrl,
		APIKey:      m.APIKey,
		Domain:      m.Domain,
		FromAddress: msg.From,
		FromName:    msg.FromName,
	}

	driver, err := apiMail.NewClient(transport, cfg)
	if err != nil {
		return err
	}

	formattedMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		return err
	}

	plainMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		return err
	}

	tx := &apiMail.Transmission{
		Recipients: []string{msg.To},
		Subject:    msg.Subject,
		HTML:       formattedMessage,
		PlainText:  plainMessage,
	}

	// add attachments
	err = m.addAPIAttachments(msg, tx)
	if err != nil {
		return err
	}

	_, err = driver.Send(tx)
	if err != nil {
		return err
	}

	return nil
}

// addAPIAttachments adds attachments, if any, to mail being sent via api
func (m *Mail) addAPIAttachments(msg Message, tx *apiMail.Transmission) error {
	if len(msg.Attachments) > 0 {
		var attachments []apiMail.Attachment

		for _, x := range msg.Attachments {
			var attach apiMail.Attachment
			content, err := os.ReadFile(x)
			if err != nil {
				return err
			}

			fileName := filepath.Base(x)
			attach.Bytes = content
			attach.Filename = fileName
			attachments = append(attachments, attach)
		}

		tx.Attachments = attachments
	}

	return nil
}

// SendUsingSMTP builds and sends an email message using SMTP. This is called by ListenForMail,
// and can also be called directly when necessary
func (m *Mail) SendUsingSMTP(msg Message) error {
	formattedMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		return err
	}

	plainMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		return err
	}

	server := smtpMail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	email := smtpMail.NewMSG()
	email.SetFrom(msg.From).
		AddTo(msg.To).
		SetSubject(msg.Subject)

	email.SetBody(smtpMail.TextHTML, formattedMessage)
	email.AddAlternative(smtpMail.TextPlain, plainMessage)

	if len(msg.Attachments) > 0 {
		for _, x := range msg.Attachments {
			email.AddAttachment(x)
		}
	}

	err = email.Send(smtpClient)
	if err != nil {
		return err
	}

	return nil
}

// getEncryption returns the appropriate encryption type based on a string value
func (m *Mail) getEncryption(e string) smtpMail.Encryption {
	switch e {
	case "tls":
		return smtpMail.EncryptionSTARTTLS
	case "ssl":
		return smtpMail.EncryptionSSL
	case "none":
		return smtpMail.EncryptionNone
	default:
		return smtpMail.EncryptionSTARTTLS
	}
}

// buildHTMLMessage creates the html version of the message
func (m *Mail) buildHTMLMessage(msg Message) (string, error) {
	templateToRender := fmt.Sprintf("%s/%s.html.tmpl", m.Templates, msg.Template)

	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.Data); err != nil {
		return "", err
	}

	formattedMessage := tpl.String()
	formattedMessage, err = m.inlineCSS(formattedMessage)
	if err != nil {
		return "", err
	}

	return formattedMessage, nil
}

// buildPlainTextMessage creates the plaintext version of the message
func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {
	templateToRender := fmt.Sprintf("%s/%s.plain.tmpl", m.Templates, msg.Template)

	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.Data); err != nil {
		return "", err
	}

	plainMessage := tpl.String()

	return plainMessage, nil
}

// inlineCSS takes html input as a string, and inlines css where possible
func (m *Mail) inlineCSS(s string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(s, &options)
	if err != nil {
		return "", err
	}

	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}
