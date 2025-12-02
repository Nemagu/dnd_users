package email

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Nemagu/dnd/internal/application"
	"golang.org/x/crypto/argon2"
)

type emailToken struct {
	EncryptedEmail string    `json:"e"`
	CreatedAt      time.Time `json:"c"`
	Salt           string    `json:"s"`
	Nonce          string    `json:"n"`
}

func (t emailToken) String() string {
	return fmt.Sprintf("%s,%s,%s,%s", t.EncryptedEmail, t.CreatedAt, t.Salt, t.Nonce)
}

type BcryptEmailCrypter struct {
	secretKey []byte
	lifetime  time.Duration
}

func MustNewBcryptEmailCrypter(secretKey string, lifetime time.Duration) *BcryptEmailCrypter {
	if len(secretKey) < 32 {
		panic("секретный ключ для сервиса шифрования email слишком короткий")
	}
	if lifetime < time.Minute {
		panic(
			"время жизни токена активации слишком мало",
		)
	}
	return &BcryptEmailCrypter{
		secretKey: []byte(secretKey),
		lifetime:  lifetime,
	}
}

func (s *BcryptEmailCrypter) Encrypt(email string) (string, error) {
	salt, err := makeSalt()
	if err != nil {
		return "", err
	}

	key, err := s.deriveKeyFrom(salt)
	if err != nil {
		return "", err
	}

	crypter, err := makeCrypterFrom(key)
	if err != nil {
		return "", err
	}

	nonce, err := makeNonce(crypter)
	if err != nil {
		return "", err
	}

	encryptedEmail := crypter.Seal(nil, nonce, []byte(email), nil)

	token := emailToken{
		EncryptedEmail: base64.StdEncoding.EncodeToString(encryptedEmail),
		CreatedAt:      time.Now(),
		Salt:           base64.StdEncoding.EncodeToString(salt),
		Nonce:          base64.StdEncoding.EncodeToString(nonce),
	}

	marshaledToken, err := json.Marshal(token)
	if err != nil {
		return "", fmt.Errorf("%w: %s", application.ErrInternal, err)
	}

	return base64.URLEncoding.EncodeToString(marshaledToken), nil
}

func (s *BcryptEmailCrypter) Decrypt(token string) (string, error) {
	marshaledToken, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return "", fmt.Errorf("%w: %s", application.ErrInternal, err)
	}

	var t emailToken
	if err := json.Unmarshal(marshaledToken, &t); err != nil {
		return "", fmt.Errorf("%w: %s", application.ErrInternal, err)
	}

	createdAt := t.CreatedAt
	encryptedEmail, err := base64.StdEncoding.DecodeString(t.EncryptedEmail)
	if err != nil {
		return "", fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	encryptedEmail = []byte(encryptedEmail)
	salt, err := base64.StdEncoding.DecodeString(t.Salt)
	if err != nil {
		return "", fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	salt = []byte(salt)
	nonce, err := base64.StdEncoding.DecodeString(t.Nonce)
	if err != nil {
		return "", fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	nonce = []byte(nonce)

	if time.Since(createdAt) > s.lifetime {
		return "", fmt.Errorf(
			"%w: токен активации устарел",
			application.ErrValidation,
		)
	}

	key, err := s.deriveKeyFrom(salt)
	if err != nil {
		return "", fmt.Errorf("%w: %s", application.ErrInternal, err)
	}

	crypter, err := makeCrypterFrom(key)
	if err != nil {
		return "", fmt.Errorf("%w: %s", application.ErrInternal, err)
	}

	decryptedEmail, err := crypter.Open(nil, nonce, encryptedEmail, nil)
	if err != nil {
		return "", fmt.Errorf("%w: %s", application.ErrInternal, err)
	}

	return string(decryptedEmail), nil
}

func (s *BcryptEmailCrypter) deriveKeyFrom(salt []byte) ([]byte, error) {
	key := argon2.IDKey([]byte(s.secretKey), salt, 1, 64*1024, 4, 32)
	return key, nil
}

func makeSalt() ([]byte, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	return salt, nil
}

func makeNonce(crypter cipher.AEAD) ([]byte, error) {
	nonce := make([]byte, crypter.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	return nonce, nil
}

func makeCrypterFrom(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	return cipher.NewGCM(block)
}
