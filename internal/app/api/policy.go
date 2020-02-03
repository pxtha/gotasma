package api

import "github.com/gotasma/internal/app/policy"

import envconfig "github.com/gotasma/internal/pkg/env"

func newPolicyService() (*policy.Service, error) {
	var conf policy.CasbinConfig
	envconfig.Load(&conf)
	enforcer := policy.NewFileCasbinEnforcer(conf)
	return policy.New(enforcer)
}
