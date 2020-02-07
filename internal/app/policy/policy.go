package policy

import (
	"context"
	"strconv"

	"github.com/casbin/casbin"
	"github.com/gotasma/internal/app/auth"
	"github.com/gotasma/internal/app/status"
	"github.com/gotasma/internal/app/types"

	"github.com/sirupsen/logrus"
)

type (
	CasbinConfig struct {
		ConfigPath string `envconfig:"CONFIG_PATH" default:"configs/casbin.conf"`
		PolicyPath string `envconfig:"CONFIG_PATH" default:"configs/casbin_policy.csv"`
	}
	Service struct {
		enforcer *casbin.Enforcer
	}
)

//NewFileCasbinEnforcer Static conf and policy file
func NewFileCasbinEnforcer(conf CasbinConfig) *casbin.Enforcer {
	enforcer := casbin.NewEnforcer(conf.ConfigPath, conf.PolicyPath)
	return enforcer
}

func New(enforcer *casbin.Enforcer) (*Service, error) {
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}
	return &Service{
		enforcer: enforcer,
	}, nil
}

func (s *Service) isAllowed(ctx context.Context, sub string, obj string, act string) bool {
	ok, err := s.enforcer.EnforceSafe(sub, obj, act)

	return err == nil && ok
}

func (s *Service) Validate(ctx context.Context, obj string, act string) error {
	//Set sub by role PM or DEV
	sub := types.PolicySubjectAny

	user := auth.FromContext(ctx)
	if user != nil {
		if user.Role == types.PM {
			return nil
		}
		sub = strconv.Itoa(int(user.Role))
	}

	//Check casbin file for policy
	//TODO set role for dev to VIEW only project of devs
	//TODO dev can ADD new task - this task automate assign to dev, modify task assigned to dev, dev cannot assign another dev to task

	if !s.isAllowed(ctx, sub, obj, act) {
		logrus.Errorf("the user is not authorized to do the action")
		return status.Policy().Unauthorized
	}

	return nil
}
