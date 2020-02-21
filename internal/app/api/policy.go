package api

import (
	"github.com/gotasma/internal/app/policy"
	"github.com/sirupsen/logrus"

	envconfig "github.com/gotasma/internal/pkg/env"
)

func newPolicyService() (*policy.Service, error) {
	var conf policy.CasbinConfig
	envconfig.LoadWithPrefix("CASBIN", &conf)
	logrus.Info(conf)
	enforcer := policy.NewFileCasbinEnforcer(conf)
	return policy.New(enforcer)
}
