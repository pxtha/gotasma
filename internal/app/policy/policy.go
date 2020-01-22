package policy

import (
	"context"
	"strconv"

	"github.com/casbin/casbin"
	"github.com/sirupsen/logrus"
	"praslar.com/gotasma/internal/app/auth"
	"praslar.com/gotasma/internal/app/status"
	"praslar.com/gotasma/internal/app/types"
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

	sub := types.PolicySubjectAny

	user := auth.FromContext(ctx)
	if user != nil {
		if user.Role == types.PM {
			return nil
		}

		sub = strconv.Itoa(int(user.Role))
	}

	if !s.isAllowed(ctx, sub, obj, act) {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"sub": sub, "action": act, "obj": obj}).Errorf("the user is not authorized to do the action")
		return status.Policy().Unauthorized
	}

	return nil
}
