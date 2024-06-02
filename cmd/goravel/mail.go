package main

import (
	"errors"
	"strings"

	"github.com/fatih/color"
)

func handleMail(arg3 string) error {

	if arg3 == "" {
		exitGracefully(errors.New("you must give your email template a name"), "Usage: goravel make mail <template-name>")
	}

	name := strings.ToLower(arg3)

	// To files:
	htmlMail := gor.RootPath + "/mail/" + name + ".html.tmpl"
	plainMail := gor.RootPath + "/mail/" + name + ".plain.tmpl"

	err := copyFilefromTemplate("templates/mailer/mail.html.tmpl", htmlMail)
	if err != nil {
		return err
	}

	err = copyFilefromTemplate("templates/mailer/mail.plain.tmpl", plainMail)
	if err != nil {
		return err
	}

	color.Green("Email template files created successfully!")
	color.Yellow(`The following templates have been added:
	  1. %s
	  2. %s
		`, "/mail/"+name+".html.tmpl", "/mail/"+name+".plain.tmpl")

	return nil
}
