package goravel

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"os"
)

// CreateDirIfNotExists creates a new directory if it does not exist
func (g *Goravel) CreateDirIfNotExists(path string) error {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		err = os.Mkdir(path, 0755) // default permissions
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateFileIfNotExists creates a new file at path if it does not exist
func (g *Goravel) CreateFileIfNotExists(path string) error {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return err
		}

		defer func(file *os.File) {
			_ = file.Close()
		}(file)
	}
	return nil
}

const (
	randomString = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321_+"
)

// RandomString generates a random string of length n from values in the const randomString
func (g *Goravel) RandomString(n int) string {
	// Create a slice s of runes with length n
	s, r := make([]rune, n), []rune(randomString)

	// Loop through each index of slice s
	for i := range s {
		// Generate a random prime number p with bit length equal to the length of r
		p, _ := rand.Prime(rand.Reader, len(r))

		// Get the uint64 representation of the prime number
		x, y := p.Uint64(), uint64(len(r))

		// Use the modulo operation to ensure the index is within bounds of r
		s[i] = r[x%y]
	}

	// Convert the slice of runes s to a string and return it
	return string(s)
}

type Encryption struct {
	Key []byte
}

// Encrypt encrypts a string using the AES encryption algorithm
func (e *Encryption) Encrypt(text string) (string, error) {
	plaintext := []byte(text)

	block, err := aes.NewCipher(e.Key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize] // Initialization vector
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// Encode the encrypted text to base64
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a string using the AES encryption algorithm
func (e *Encryption) Decrypt(cryptoText string) (string, error) {
	encetptedText, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(e.Key)
	if err != nil {
		return "", err
	}

	if len(encetptedText) < aes.BlockSize {
		return "", err
	}

	iv := encetptedText[:aes.BlockSize]
	ciphertext := encetptedText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}
