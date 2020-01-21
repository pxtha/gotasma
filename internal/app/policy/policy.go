package policy

import "github.com/casbin/casbin"

type (
	CasbinConfig struct {
		ConfigPath string `envconfig:"CONFIG_PATH" default:"configs/casbin.conf"`
	}
	Service struct {
		enforcer *casbin.Enforcer
	}
)

func New(enforcer *casbin.Enforcer) (*Service, err) {
	enforcer.EnableAutoSave(true)
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}
	return &Service{
		enforcer: enforcer,
	}, nil
}
