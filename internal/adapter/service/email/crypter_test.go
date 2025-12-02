package email

import (
	"testing"
	"time"
)

func createCrypter() *BcryptEmailCrypter {
	return MustNewBcryptEmailCrypter(
		"EBVBniiRU1aa9oyAxSShaJoNCKHac8tRrLf9wpOqXx3DFhH2fqyzT39WSlO",
		300*time.Second,
	)
}

func generateToken(crypter *BcryptEmailCrypter, email string) string {
	token, err := crypter.Encrypt(email)
	if err != nil {
		panic(err)
	}
	return token
}

func TestCrypter_Encrypt_Success(t *testing.T) {
	cases := []string{
		"normal@mail.com",
		"kek@mail.com",
		"lol@mail.com",
		"ohohoh@mail.com",
	}
	crypter := createCrypter()
	for _, c := range cases {
		t.Run("test_bcrypt_encode_"+c, func(t *testing.T) {
			if _, err := crypter.Encrypt(c); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestCrypter_Decrypt_Success(t *testing.T) {
	crypter := createCrypter()
	cases := []string{
		generateToken(crypter, "normal@mail.com"),
		generateToken(crypter, "kek@mail.com"),
		generateToken(crypter, "lol@mail.com"),
		generateToken(crypter, "ohohoh@mail.com"),
	}
	for _, c := range cases {
		t.Run("test_bcrypt_encode_"+c, func(t *testing.T) {
			if _, err := crypter.Decrypt(c); err != nil {
				t.Error(err)
			}
		})
	}
}
