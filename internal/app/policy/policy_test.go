package policy_test

import (
	"testing"

	"golang.org/x/net/context"
	"praslar.com/gotasma/internal/app/auth"
	"praslar.com/gotasma/internal/app/policy"
	"praslar.com/gotasma/internal/app/types"
)

func TestPolicy(t *testing.T) {
	conf := policy.CasbinConfig{
		ConfigPath: "/home/pxthang/Code/Golang/src/praslar.com/gotasma/configs/casbin.conf",
		PolicyPath: "/home/pxthang/Code/Golang/src/praslar.com/gotasma/configs/casbin_policy.csv",
	}
	enforcer := policy.NewFileCasbinEnforcer(conf)
	policySrv, err := policy.New(enforcer)
	if err != nil {
		t.Error(err)
	}
	ctx := context.Background()
	ctx = auth.NewContext(ctx, &types.User{
		Role: types.Role(1),
	})
	if err := policySrv.Validate(ctx, "/api/v1/users/", "GET"); err != nil {
		t.Error(err)
	}
}
