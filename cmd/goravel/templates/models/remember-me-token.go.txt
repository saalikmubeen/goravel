package models

import (
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/saalikmubeen/goravel"

	up "github.com/upper/db/v4"
)

type RememberMeToken struct {
	ID              int       `db:"id,omitempty"`
	UserID          int       `db:"user_id"`
	RememberMeToken string    `db:"remember_me_token"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

func (t *RememberMeToken) Table() string {
	return "remember_me_tokens"
}

func (t *RememberMeToken) InsertToken(userID int, hash ...string) (string, error) {

	token, err := t.GenerateToken()

	if err != nil {
		return "", err
	}

	collection := upper.Collection(t.Table())
	rememberToken := RememberMeToken{
		UserID:          userID,
		RememberMeToken: token,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	_, err = collection.Insert(rememberToken)
	if err != nil {
		return "", err
	}
	return token, nil
}

// GetByToken returns a token by it's remember me token
func (t *RememberMeToken) GetByToken(token string) (*RememberMeToken, error) {
	var rememberMeToken RememberMeToken
	collection := upper.Collection(t.Table())
	res := collection.Find(up.Cond{"remember_me_token": token})
	err := res.One(&rememberMeToken)
	if err != nil {
		return nil, err
	}

	return &rememberMeToken, nil
}

func (t *RememberMeToken) Delete(rememberMeToken string) error {
	collection := upper.Collection(t.Table())
	res := collection.Find(up.Cond{"remember_me_token": rememberMeToken})
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

func (t *RememberMeToken) GenerateToken() (string, error) {
	g := &goravel.Goravel{}

	randomString := g.RandomString(12)

	hasher := sha256.New()
	_, err := hasher.Write([]byte(randomString))
	if err != nil {
		return "", err
	}

	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	return sha, nil

}
