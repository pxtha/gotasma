package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/globalsign/mgo"
	gomock "github.com/golang/mock/gomock"
	"github.com/gotasma/internal/app/status"
	"github.com/gotasma/internal/app/types"
	"github.com/gotasma/internal/app/user"
	"github.com/gotasma/internal/pkg/uuid"
)

func TestRegister(t *testing.T) {
	//status file
	status.Init("../../../configs/status.yml")
	mockRepo := NewMockRepository(gomock.NewController(t))
	mockPolicy := NewMockPolicyService(gomock.NewController(t))
	services := user.New(mockRepo, mockPolicy)

	req := types.RegisterRequest{
		Email:     "newuser@gmail.com",
		FirstName: "Thang",
		LastName:  "Pham",
		Password:  "1234",
		Role:      0,
	}

	password, err := services.GeneratePassword(req.Password)

	if err != nil {
		t.Errorf("failed to generate password: %w", err)
	}
	userID := uuid.New()

	user := &types.User{
		Password:  password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Role:      types.PM,
		UserID:    userID,
		CreaterID: userID,
	}

	testCases := []struct {
		desc   string
		expect func()
		input  *types.RegisterRequest
		output *types.User
		err    error
	}{
		{
			desc: "Validator not ok",
			input: &types.RegisterRequest{
				Email:     "newusergmail.com",
				FirstName: "Thang",
				LastName:  "Pham",
				Password:  "1234",
				Role:      0,
			},
			expect: func() {},
			output: nil,
			err:    status.Gen().BadRequest,
		},
		{
			desc:  "Database error",
			input: &req,
			expect: func() {
				mockRepo.EXPECT().FindByEmail(gomock.Any(), req.Email).Return(nil, status.Gen().Internal)
			},
			output: nil,
			err:    status.Gen().Internal,
		},
		{
			desc:  "Duplicate email",
			input: &req,
			expect: func() {
				mockRepo.EXPECT().FindByEmail(gomock.Any(), req.Email).Return(user, nil)
			},
			output: nil,
			err:    status.User().DuplicatedEmail,
		},
		{
			desc:  "Fail to insert PM",
			input: &req,
			expect: func() {
				mockRepo.EXPECT().FindByEmail(gomock.Any(), req.Email).Return(nil, mgo.ErrNotFound)
				mockRepo.EXPECT().Create(gomock.Any(), req).Return("", status.Gen().Internal)
			},
			output: nil,
			err:    status.Gen().Internal,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			tC.expect()
			output, err := services.Register(context.TODO(), tC.input)
			if output != tC.output {
				t.Errorf("got output: %v want output: %v", output, tC.output)
			}
			if !errors.Is(err, tC.err) {
				t.Errorf("got err: %v want err: %v", err, tC.err)
			}
		})
	}
}
