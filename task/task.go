package task

import (
	"log"

	"github.com/knq/envcfg"
	"github.com/knq/firebase"
)

// Task represents firebase database model for task data ref.
type Task struct {
	TaskName  string `json:"task_name"`
	URL       string `json:"url,omitempty"`
	Tags      string `json:"tags,omitempty"`
	Notes     string `json:"notes,omitempty"`
	NotesMD   string `json:"notes_md,omitempty"`
	Date      string `json:"date"`
	Completed bool   `json:"completed"`
}

// Service is implementation of task service server.
type Service struct {
	// configPath is envcfg config file path.
	configPath string

	// config is server environment config.
	config *envcfg.Envcfg

	// dataRef is synoday-data firebase database ref.
	dataRef *firebase.DatabaseRef
}

// NewService creates new instance of task service server.
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
	s.dataRef, err = firebase.NewDatabaseRef(
		firebase.GoogleServiceAccountCredentialsJSON([]byte(s.config.GetKey("firebase.datacreds"))),
	)
	if err != nil {
		log.Fatal(err)
	}

	return s
}
