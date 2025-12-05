package email

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func createCrypter() *BcryptEmailCrypter {
	return MustNewBcryptEmailCrypter(
		"EBVBniiRU1aa9oyAxSShaJoNCKHac8tRrLf9wpOqXx3DFhH2fqyzT39WSlO",
		300*time.Second,
	)
}

func generateEmailToken(crypter *BcryptEmailCrypter, email string) string {
	token, err := crypter.EncryptEmail(email)
	if err != nil {
		panic(err)
	}
	return token
}

func generateEmailUserIDToken(crypter *BcryptEmailCrypter, email string, userID uuid.UUID) string {
	token, err := crypter.EncryptEmailUserID(email, userID)
	if err != nil {
		panic(err)
	}
	return token
}

func TestCrypter_EncryptEmail_Success(t *testing.T) {
	cases := []string{
		"normal@mail.com",
		"kek@mail.com",
		"lol@mail.com",
		"ohohoh@mail.com",
	}
	crypter := createCrypter()
	for _, c := range cases {
		t.Run("test_bcrypt_encode_"+c, func(t *testing.T) {
			if _, err := crypter.EncryptEmail(c); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestCrypter_DecryptEmail_Success(t *testing.T) {
	crypter := createCrypter()
	cases := []string{
		generateEmailToken(crypter, "normal@mail.com"),
		generateEmailToken(crypter, "kek@mail.com"),
		generateEmailToken(crypter, "lol@mail.com"),
		generateEmailToken(crypter, "ohohoh@mail.com"),
	}
	for _, c := range cases {
		t.Run("test_bcrypt_encode_"+c, func(t *testing.T) {
			if _, err := crypter.DecryptEmail(c); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestCrypter_EncryptEmilUserID_Success(t *testing.T) {
	cases := []struct {
		Email  string
		UserID uuid.UUID
	}{
		{
			Email:  "normal@mail.com",
			UserID: uuid.New(),
		},
		{
			Email:  "kek@mail.com",
			UserID: uuid.New(),
		},
		{
			Email:  "lol@mail.com",
			UserID: uuid.New(),
		},
		{
			Email:  "ohohoh@mail.com",
			UserID: uuid.New(),
		},
	}
	crypter := createCrypter()
	for _, c := range cases {
		t.Run("test_bcrypt_encode_email_user_id_"+c.Email, func(t *testing.T) {
			if _, err := crypter.EncryptEmailUserID(c.Email, c.UserID); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestCrypter_DecryptEmailUserID_Success(t *testing.T) {
	crypter := createCrypter()
	cases := []string{
		generateEmailUserIDToken(crypter, "normal@mail.com", uuid.New()),
		generateEmailUserIDToken(crypter, "kek@mail.com", uuid.New()),
		generateEmailUserIDToken(crypter, "lol@mail.com", uuid.New()),
		generateEmailUserIDToken(crypter, "ohohoh@mail.com", uuid.New()),
	}
	for _, c := range cases {
		t.Run("test_bcrypt_encode_email_user_id_"+c, func(t *testing.T) {
			if _, err := crypter.DecryptEmail(c); err != nil {
				t.Error(err)
			}
		})
	}
}
