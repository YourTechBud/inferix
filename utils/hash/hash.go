package hash

import "fmt"

type (
	Hasher interface {
		VerifyConfig() error
		Hash(text string) (string, error)
		Compare(hash, text string) error
	}

	Options struct {
		Algo   HashAlgo `mapstructure:"algo"`
		Bcrypt Bcrypt   `mapstructure:"bcrypt"`
	}

	HashAlgo string
)

const (
	HashAlgo_Bcrypt HashAlgo = "bcrypt"
)

var DefaultHasher Hasher

func InitialiseHasher(opts Options) error {
	switch opts.Algo {
	case HashAlgo_Bcrypt:
		DefaultHasher = &opts.Bcrypt
	default:
		return fmt.Errorf("hasher type %s is not supported", opts.Algo)
	}

	// Verify the configuration
	if err := DefaultHasher.VerifyConfig(); err != nil {
		return err
	}

	// Return nil if everything is fine
	return nil
}
