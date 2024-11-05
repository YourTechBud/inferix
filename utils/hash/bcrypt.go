package hash

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Bcrypt struct {
	Cost int `mapstructure:"cost"`
}

func (b *Bcrypt) VerifyConfig() error {
	if b.Cost < bcrypt.MinCost || b.Cost > bcrypt.MaxCost {
		return fmt.Errorf("bcrypt cost should be between %d and %d", bcrypt.MinCost, bcrypt.MaxCost)
	}

	return nil
}

func (b *Bcrypt) Hash(text string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), b.Cost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (b *Bcrypt) Compare(hash, text string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(text))
}
