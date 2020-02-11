package policy_test

import (
	"context"
	"testing"

	"github.com/gotasma/internal/app/auth"
	"github.com/gotasma/internal/app/policy"
	"github.com/gotasma/internal/app/status"
	"github.com/gotasma/internal/app/types"
)

func TestCreateEnforcerFailed(t *testing.T) {
	status.Init("../../../configs/status.yml")
	enforcer := policy.NewFileCasbinEnforcer(policy.CasbinConfig{
		ConfigPath: "../../../configs/casbin.conf",
		PolicyPath: "testdata/casbin_policy_not_found.csv",
	})
	_, err := policy.New(enforcer)
	if err == nil {
		t.Fatalf("got service created, want service creation failed: %v", err)
	}
}

func TestIsAllowed(t *testing.T) {
	status.Init("../../../configs/status.yml")
	enforcer := policy.NewFileCasbinEnforcer(policy.CasbinConfig{
		ConfigPath: "../../../configs/casbin.conf",
		PolicyPath: "testdata/casbin_policy.csv",
	})
	srv, err := policy.New(enforcer)
	if err != nil {
		t.Fatalf("failed to create service, err: %v", err)
	}
	testCases := []struct {
		name string
		ctx  context.Context
		obj  string
		act  string
		err  error
	}{
		{
			name: "PM should be allowed to do anything",
			ctx: auth.NewContext(context.Background(), &types.User{
				UserID: "123",
				Role:   types.PM,
			}),
			obj: "/api/v1/users/dev",
			act: types.PolicyActionAny,
			err: nil,
		},
		{
			name: "dev cannot delete dev",
			ctx: auth.NewContext(context.Background(), &types.User{
				UserID: "124",
				Role:   types.DEV,
			}),
			obj: "/api/v1/users/dev",
			act: types.PolicyActionAny,
			err: status.Policy().Unauthorized,
		},
		{
			name: "dev is allowed to create dev",
			ctx: auth.NewContext(context.Background(), &types.User{
				UserID: "124",
				Role:   types.DEV,
			}),
			obj: "/api/v1/users/dev",
			act: types.PolicyActionAny,
			err: nil,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			if err := srv.Validate(c.ctx, c.obj, c.act); err != c.err {
				t.Errorf("got err=%v, want err=%v", err, c.err)
			}
		})
	}
}
