package user

import (
	"github.com/knq/envcfg"
)

// Option is user service server option.
type Option func(*Service) error

// ConfigFile is an option that sets the file path to read data from.
func ConfigFile(path string) Option {
	return func(s *Service) error {
		var err error
		s.config, err = envcfg.New(
			envcfg.ConfigFile(path),
		)
		return err
	}
}
