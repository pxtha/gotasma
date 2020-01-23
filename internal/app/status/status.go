package status

import (
	"os"
	"sync"

	"gopkg.in/yaml.v2"

	"praslar.com/gotasma/internal/pkg/status"
)

type (
	Status    = status.Status
	GenStatus struct {
		Success    Status
		NotFound   Status
		Timeout    status.Timeout
		BadRequest Status `yaml:"bad_request"`
		Internal   Status
	}
	HolidayStatus struct {
		DuplicatedHoliday Status `yaml:"duplicated_holiday"`
	}
	PolicyStatus struct {
		Unauthorized Status
	}
	UserStatus struct {
		DuplicatedEmail Status `yaml:"duplicated_email"`
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
		if err != nil {
			panic(err)
		}
		all = &statuses{}
		if err := yaml.NewDecoder(f).Decode(all); err != nil {
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
	if all == nil {
		Init(conf)
	}
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

func Hoiday() HolidayStatus {
	return load().Holiday
}
