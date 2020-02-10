package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gotasma/internal/app/status"
	"github.com/gotasma/internal/app/types"
	"github.com/gotasma/internal/app/user"
)

func before(t *testing.T) (*MockRepository, *MockPolicyService, *user.Service) {
	mockRepo := NewMockRepository(gomock.NewController(t))
	mockPolicy := NewMockPolicyService(gomock.NewController(t))
	services := user.New(mockRepo, mockPolicy)
	status.Init("../../../configs/status.yml")
	return mockRepo, mockPolicy, services
}

func TestRegister(t *testing.T) {

	mockRepo, _, services := before(t)

	req := &types.RegisterRequest{
		Email:     "pxthang@gmail.com",
		FirstName: "Thang",
		LastName:  "Pham",
		Role:      0,
		Password:  "1234",
	}

	user := &types.User{
		Email:     "pxthang@gmail.com",
		FirstName: "Thang",
		LastName:  "Pham",
		Role:      0,
		Password:  "password",
		CreaterID: "1234ab",
		UserID:    "1234ab",
	}

	testCases := []struct {
		desc   string
		expect func()
		input  *types.RegisterRequest
		err    error
	}{
		{
			desc:   "Validate fail",
			expect: func() {},
			input: &types.RegisterRequest{
				Email:     "pxthanggmail.com",
				FirstName: "Thang",
				LastName:  "Pham",
				Role:      0,
				Password:  "1234",
			},
			err: status.Gen().BadRequest,
		},
		{
			desc: "Database err",
			expect: func() {
				mockRepo.EXPECT().FindByEmail(gomock.Any(), req.Email).Return(nil, status.User().NotFoundUser)
			},
			input: req,
			err:   status.User().NotFoundUser,
		},
		{
			desc: "Duplicate user",
			expect: func() {
				mockRepo.EXPECT().FindByEmail(gomock.Any(), req.Email).Return(user, nil)
			},
			input: req,
			err:   status.User().DuplicatedEmail,
		},
		{
			desc: "Fail to insert PM",
			expect: func() {
				mockRepo.EXPECT().FindByEmail(gomock.Any(), req.Email).Return(nil, nil)
				mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return("", status.Gen().Internal)
			},
			input: req,
			err:   status.Gen().Internal,
		},

		{
			desc: "Create PM",
			expect: func() {
				mockRepo.EXPECT().FindByEmail(gomock.Any(), req.Email).Return(nil, nil)
				mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(user.UserID, nil)
			},
			input: req,
			err:   nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			tC.expect()
			_, err := services.Register(context.TODO(), tC.input)
			if !errors.Is(err, tC.err) {
				t.Errorf("\n got err = %v, \n wants err = %v", err, tC.err)
			}
		})
	}
}

func TestCreateDev(t *testing.T) {
	mockRepo, mockPolicy, services := before(t)

	req := &types.RegisterRequest{
		Email:     "pxthang_devs@gmail.com",
		FirstName: "Thang",
		LastName:  "Pham",
		Role:      1,
		Password:  "1234",
	}

	testCases := []struct {
		desc   string
		expect func()
		input  *types.RegisterRequest
		err    error
	}{
		{
			desc: "Not authorize",
			expect: func() {
				mockPolicy.EXPECT().Validate(gomock.Any(), gomock.Any(), gomock.Any()).Return(status.Policy().Unauthorized)
			},
			input: req,
			err:   status.Policy().Unauthorized,
		},
		{
			desc: "Validate fail",
			expect: func() {
				mockPolicy.EXPECT().Validate(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			input: &types.RegisterRequest{
				Email:     "pxthanggmail.com",
				FirstName: "",
				LastName:  "",
				Role:      1,
				Password:  "",
			},
			err: status.Gen().BadRequest,
		},
		{
			desc: "Devs duplicate",
			expect: func() {
				mockPolicy.EXPECT().Validate(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockRepo.EXPECT().FindByEmail(gomock.Any(), req.Email).Return(nil, status.Gen().Internal)
			},
			input: req,
			err:   status.Gen().Internal,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			tC.expect()
			_, err := services.CreateDev(context.TODO(), tC.input)
			if !errors.Is(err, tC.err) {
				t.Errorf("\n got err = %v, \n wants err = %v", err, tC.err)
			}
		})
	}
}
