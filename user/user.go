package user

import (
	"log"

	"github.com/knq/envcfg"
	"github.com/knq/firebase"
)

// User represents firebase database model for user data ref.
type User struct {
	FirstName string                   `json:"first_name"`
	LastName  string                   `json:"last_name"`
	Created   firebase.ServerTimestamp `json:"created"`
}

// Service is implementation of user service server.
type Service struct {
	// configPath is envcfg config file path.
	configPath string

	// config is server environment config.
	config *envcfg.Envcfg

	// authRef is synoday-auth firebase database ref.
	authRef *firebase.DatabaseRef

	// dataRef is synoday-data firebase database ref.
	dataRef *firebase.DatabaseRef
}

// NewService creates new instance of user service.
func NewService(opts ...Option) *Service {
	var err error
	s := new(Service)

	// set options
	for _, o := range opts {
		err = o(s)
		if err != nil {
			log.Fatal(err)
		}
	}

	if s.config == nil {
		s.config, err = envcfg.New()
		if err != nil {
			log.Fatal(err)
		}
	}

	// setup database
	s.authRef, err = firebase.NewDatabaseRef(
		firebase.GoogleServiceAccountCredentialsJSON([]byte(s.config.GetKey("firebase.authcreds"))),
	)
	if err != nil {
		log.Fatal(err)
	}
	s.dataRef, err = firebase.NewDatabaseRef(
		firebase.GoogleServiceAccountCredentialsJSON([]byte(s.config.GetKey("firebase.datacreds"))),
	)
	if err != nil {
		log.Fatal(err)
	}

	return s
}
