package status

import (
	"os"
	"sync"

	"gopkg.in/yaml.v2"

	"github.com/gotasma/internal/pkg/status"
	"github.com/sirupsen/logrus"
)

type (
	//Status format from status pkg
	Status = status.Status

	GenStatus struct {
		Success    Status
		NotFound   Status `yaml:"not_found"`
		BadRequest Status `yaml:"bad_request"`
		Internal   Status
	}

	ProjectStatus struct {
		NotFoundProject  Status `yaml:"not_found_project"`
		DuplicateProject Status `yaml:"duplicated_project"`
		AlreadyInProject Status `yaml:"already_in_project"`
		NotInProject     Status `yaml:"not_in_project"`
	}
	HolidayStatus struct {
		InvalidHoliday    Status `yaml:"invalid_holiday"`
		DuplicatedHoliday Status `yaml:"duplicated_holiday"`
		NotFoundHoliday   Status `yaml:"not_found_holiday"`
	}
	PolicyStatus struct {
		Unauthorized Status
	}
	UserStatus struct {
		DuplicatedEmail Status `yaml:"duplicated_email"`
		NotFoundUser    Status `yaml:"not_found_user"`
	}
	AuthStatus struct {
		InvalidUserPassword Status `yaml:"invalid_user_password"`
	}
	SercurityStatus struct {
		InvalidAction Status `yaml:"invalid_action"`
	}
	statuses struct {
		Gen       GenStatus
		User      UserStatus
		Auth      AuthStatus
		Policy    PolicyStatus
		Sercurity SercurityStatus
		Holiday   HolidayStatus
		Project   ProjectStatus
	}
)

var (
	all  *statuses
	once sync.Once
)

// Init load statuses from the given config file.
// Init panics if cannot access or error while parsing the config file.
func Init(conf string) {
	once.Do(func() {
		f, err := os.Open(conf)
		logrus.Infof("Succesful open status file")
		if err != nil {
			logrus.Errorf("Fail to open status file, %v", err)
			panic(err)
		}
		all = &statuses{}
		if err := yaml.NewDecoder(f).Decode(all); err != nil {
			logrus.Errorf("Fail to parse status file data to statuses struct, %v", err)
			panic(err)
		}
	})
}

// all return all registered statuses.
// all will load statuses from configs/Status.yml if the statuses has not initalized yet.
func load() *statuses {
	conf := os.Getenv("STATUS_PATH")
	if conf == "" {
		conf = "configs/status.yml"
	}
	Init(conf)
	return all
}

func Gen() GenStatus {
	return load().Gen
}

func User() UserStatus {
	return load().User
}

func Success() Status {
	return Gen().Success
}

func Auth() AuthStatus {
	return load().Auth
}

func Policy() PolicyStatus {
	return load().Policy
}

func Sercurity() SercurityStatus {
	return load().Sercurity
}

func Holiday() HolidayStatus {
	return load().Holiday
}

func Project() ProjectStatus {
	return load().Project
}
