package password

import (
	"fmt"
	"math/rand"
	"testing"
)

func initHasher() *BcryptPasswordHasher {
	return MustNewBcryptPasswordHasher(14)
}

func TestHasher_Hash_Success(t *testing.T) {
	cases := []string{
		"it_is_strong_password1&",
		"this_is_very_strong_passw0rd#",
		"paskwoghn*8",
	}
	hasher := initHasher()
	for _, c := range cases {
		t.Run(fmt.Sprintf("test_password_%s", c), func(t *testing.T) {
			if _, err := hasher.Hash(c); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestHasher_Hash_Bad(t *testing.T) {
	cases := []string{
		"qewrtyuiop[]sgdfhghjkl;cvbnm,dfghjkcvbnmcfgvjnkhfbksldsjvdnvjkldfnjjkvdsflvjjdfvjksfvnjklsfvlkfnv",
	}
	hasher := initHasher()
	for _, c := range cases {
		t.Run(fmt.Sprintf("test_password_%s", c), func(t *testing.T) {
			if _, err := hasher.Hash(c); err == nil {
				t.Errorf("password %s did not raise error", c)
			}
		})
	}
}

func TestHasher_Compare_Success(t *testing.T) {
	cases := []string{
		"it_is_strong_password1&",
		"this_is_very_strong_passw0rd#",
		"paskwoghn*8",
	}
	hasher := initHasher()
	for _, c := range cases {
		t.Run(fmt.Sprintf("test_password_%s", c), func(t *testing.T) {
			hash, err := hasher.Hash(c)
			if err != nil {
				t.Error(err)
			}
			compared, err := hasher.ComparePassword(c, hash)
			if err != nil {
				t.Error(err)
			}
			if !compared {
				t.Errorf(
					"password %s was not compared successfully with hash %s",
					c,
					hash,
				)
			}
		})
	}
}

func TestHasher_Compare_Bad(t *testing.T) {
	cases := []string{
		"it_is_strong_password1&",
		"this_is_very_strong_passw0rd#",
		"paskwoghn*8",
	}
	hasher := initHasher()
	for _, c := range cases {
		t.Run(fmt.Sprintf("test_password_%s", c), func(t *testing.T) {
			hash, err := hasher.Hash(fmt.Sprintf("%s%d", c, rand.Int()))
			if err != nil {
				t.Error(err)
			}
			compared, err := hasher.ComparePassword(c, hash)
			if err != nil {
				t.Error(err)
			}
			if compared {
				t.Errorf(
					"password %s was compared successfully with hash %s",
					c,
					hash,
				)
			}
		})
	}
}
