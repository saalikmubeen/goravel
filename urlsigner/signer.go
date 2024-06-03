package urlsigner

import (
	"fmt"
	"strings"
	"time"
)

type Signer struct {
	Secret []byte // Secret is the key used to sign the token
}

func (s *Signer) GenerateTokenFromString(data string) string {
	var urlToSign = ""

	crypt := New(s.Secret, Timestamp)

	if strings.Contains(data, "?") {
		urlToSign = fmt.Sprintf("%s&hash=", data)
	} else {
		urlToSign = fmt.Sprintf("%s?hash=", data)
	}

	tokenBytes := crypt.Sign([]byte(urlToSign))
	token := string(tokenBytes)

	return token
}

func (s *Signer) VerifyToken(token string) bool {
	crypt := New(s.Secret, Timestamp)

	_, err := crypt.Unsign([]byte(token))

	if err != nil {
		return false
	}

	return true
}

func (s *Signer) Expired(data string, minutesUntilExpire int) bool {
	crypt := New(s.Secret, Timestamp)

	parts := strings.Split(data, "&hash=")
	if len(parts) != 2 {
		return true
	}

	token := parts[1]

	ts := crypt.Parse([]byte(token))

	return time.Since(ts.Timestamp) > time.Duration(minutesUntilExpire)*time.Minute
}
