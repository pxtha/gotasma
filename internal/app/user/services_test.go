package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/globalsign/mgo"
	"github.com/golang/mock/gomock"
	"github.com/gotasma/internal/app/auth"
	"github.com/gotasma/internal/app/status"
	"github.com/gotasma/internal/app/types"
	"github.com/gotasma/internal/app/user"
)

type (
	req struct {
		Email    string
		Password string
	}
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
				mockRepo.EXPECT().FindByEmail(gomock.Any(), req.Email).Return(nil, status.Gen().Internal)
			},
			input: req,
			err:   status.Gen().Internal,
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
			desc: "PM Created",
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
		Email:     "pxthangdevs@gmail.com",
		FirstName: "Thang",
		LastName:  "Pham",
		Role:      1,
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
	ctx := auth.NewContext(context.Background(), &types.User{
		Role:   0,
		UserID: "1234ab",
	})
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
			desc: "Database error finding devs",
			expect: func() {
				mockPolicy.EXPECT().Validate(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockRepo.EXPECT().FindByEmail(gomock.Any(), req.Email).Return(nil, status.Gen().Internal)
			},
			input: req,
			err:   status.Gen().Internal,
		},
		{
			desc: "Dev duplicate",
			expect: func() {
				mockPolicy.EXPECT().Validate(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				mockRepo.EXPECT().FindByEmail(gomock.Any(), req.Email).Return(user, nil)
			},
			input: req,
			err:   status.User().DuplicatedEmail,
		},
		{
			desc: "Fails to insert new devs account",
			expect: func() {
				mockPolicy.EXPECT().Validate(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				mockRepo.EXPECT().FindByEmail(gomock.Any(), req.Email).Return(nil, nil)

				mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return("", status.Gen().Internal)
			},
			input: req,
			err:   status.Gen().Internal,
		},
		{
			desc: "Dev created",
			expect: func() {

				mockPolicy.EXPECT().Validate(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

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
			_, err := services.CreateDev(ctx, tC.input)
			if !errors.Is(err, tC.err) {
				t.Errorf("\n got err = %v, \n wants err = %v", err, tC.err)
			}
		})
	}
}

func TestAuth(t *testing.T) {

	mockRepo, _, services := before(t)
	userWrongPass := &types.User{
		Email:     "thang@gmail.com",
		FirstName: "Thang",
		LastName:  "Pham",
		CreaterID: "1234ab",
		UserID:    "1234ab",
		Password:  "123",
	}
	reqUserWrongpass := req{
		Email:    "thang@gmail.com",
		Password: "123",
	}
	user := req{
		Email:    "pxthang@gmail.com",
		Password: "1234",
	}
	pass, _ := services.GeneratePassword(user.Password)
	userInfo := &types.User{
		Email:     "thang@gmail.com",
		FirstName: "Thang",
		LastName:  "Pham",
		CreaterID: "1234ab",
		UserID:    "1234ab",
		Password:  pass,
	}
	testCases := []struct {
		desc   string
		expect func()
		input  req
		err    error
	}{
		{
			desc: "Database err",
			expect: func() {
				mockRepo.EXPECT().FindByEmail(gomock.Any(), user.Email).Return(nil, status.Gen().Internal)
			},
			input: user,
			err:   status.Gen().Internal,
		},
		{
			desc: "Not found",
			expect: func() {
				mockRepo.EXPECT().FindByEmail(gomock.Any(), user.Email).Return(nil, mgo.ErrNotFound)
			},
			input: user,
			err:   status.User().NotFoundUser,
		},
		{
			desc: "Invalid password",
			expect: func() {
				mockRepo.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(userWrongPass, nil)
			},
			input: reqUserWrongpass,
			err:   status.Auth().InvalidUserPassword,
		},
		{
			desc: "Auth ok",
			expect: func() {
				mockRepo.EXPECT().FindByEmail(gomock.Any(), user.Email).Return(userInfo, nil)
			},
			input: user,
			err:   nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			tC.expect()
			_, err := services.Auth(context.TODO(), user.Email, user.Password)
			if !errors.Is(err, tC.err) {
				t.Errorf("\n got err = %v, \n wants err = %v", err, tC.err)
			}
		})
	}
}

func TestFindAllDev(t *testing.T) {

	mockRepo, mockPolicy, services := before(t)
	ctx := auth.NewContext(context.Background(), &types.User{
		Role:   0,
		UserID: "1234ab",
	})
	devs := []*types.User{
		{
			Email:     "pxthang",
			FirstName: "thang",
			LastName:  "pham",
			Role:      1,
			CreaterID: "123",
			UserID:    "123",
		},
		{
			Email:     "pxthang2",
			FirstName: "thang2",
			LastName:  "pham2",
			Role:      1,
			CreaterID: "123",
			UserID:    "123",
		},
	}
	pm := auth.FromContext(ctx)

	testCases := []struct {
		desc   string
		expect func()
		err    error
	}{
		{
			desc: "Not authorized",
			expect: func() {
				mockPolicy.EXPECT().Validate(gomock.Any(), gomock.Any(), gomock.Any()).Return(status.Policy().Unauthorized)
			},
			err: status.Policy().Unauthorized,
		},
		{
			desc: "Not found",
			expect: func() {
				mockPolicy.EXPECT().Validate(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockRepo.EXPECT().FindAllDev(gomock.Any(), pm.UserID).Return(nil, mgo.ErrNotFound)
			},
			err: status.User().NotFoundUser,
		},
		{
			desc: "Found devs",
			expect: func() {
				mockPolicy.EXPECT().Validate(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockRepo.EXPECT().FindAllDev(gomock.Any(), pm.UserID).Return(devs, nil)
			},
			err: nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			tC.expect()
			_, err := services.FindAllDev(ctx)
			if !errors.Is(err, tC.err) {
				t.Errorf("\n got err = %v, \n wants err = %v", err, tC.err)
			}
		})
	}
}

func TestDelete(t *testing.T) {

	mockRepo, mockPolicy, services := before(t)
	userInfo := &types.User{
		Email:  "thang@gmail.com",
		UserID: "1234",
		Role:   0,
	}
	devInfo := &types.User{
		Email:  "thang@gmail.com",
		UserID: "12345",
		Role:   1,
	}
	testCases := []struct {
		desc   string
		expect func()
		err    error
		input  string
	}{
		{
			desc: "Not authorized",
			expect: func() {
				mockPolicy.EXPECT().Validate(gomock.Any(), gomock.Any(), gomock.Any()).Return(status.Policy().Unauthorized)
			},
			err:   status.Policy().Unauthorized,
			input: "1234",
		},
		{
			desc: "Database error ",
			expect: func() {
				mockPolicy.EXPECT().Validate(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockRepo.EXPECT().FindByID(gomock.Any(), "1234").Return(nil, status.Gen().Internal)
			},
			err:   status.Gen().Internal,
			input: "1234",
		},
		{
			desc: "Not found user",
			expect: func() {
				mockPolicy.EXPECT().Validate(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockRepo.EXPECT().FindByID(gomock.Any(), "1234").Return(nil, mgo.ErrNotFound)
			},
			err:   status.User().NotFoundUser,
			input: "1234",
		},
		{
			desc: "Try to delete PM",
			expect: func() {
				mockPolicy.EXPECT().Validate(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockRepo.EXPECT().FindByID(gomock.Any(), "1234").Return(userInfo, nil)
			},
			err:   status.Sercurity().InvalidAction,
			input: "1234",
		},
		{
			desc: "Delete fail",
			expect: func() {
				mockPolicy.EXPECT().Validate(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockRepo.EXPECT().FindByID(gomock.Any(), "12345").Return(devInfo, nil)
				mockRepo.EXPECT().Delete(gomock.Any(), "12345").Return(status.Gen().Internal)
			},
			err:   status.Gen().Internal,
			input: "12345",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			tC.expect()
			err := services.Delete(context.TODO(), tC.input)
			if !errors.Is(err, tC.err) {
				t.Errorf("\n got err = %v, \n wants err = %v", err, tC.err)
			}
		})
	}
}
func TestCheckUserExist(t *testing.T) {

	mockRepo, _, services := before(t)
	userInfo := &types.User{
		Email:  "thang@gmail.com",
		UserID: "1234",
		Role:   0,
	}
	testCases := []struct {
		desc   string
		expect func()
		err    error
		input  string
	}{
		{
			desc: "Database error ",
			expect: func() {
				mockRepo.EXPECT().FindByID(gomock.Any(), "1234").Return(nil, status.Gen().Internal)
			},
			err:   status.Gen().Internal,
			input: "1234",
		},
		{
			desc: "Not found user",
			expect: func() {
				mockRepo.EXPECT().FindByID(gomock.Any(), "1234").Return(nil, mgo.ErrNotFound)
			},
			err:   status.User().NotFoundUser,
			input: "1234",
		},
		{
			desc: "Found user",
			expect: func() {
				mockRepo.EXPECT().FindByID(gomock.Any(), "1234").Return(userInfo, nil)
			},
			err:   nil,
			input: "1234",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			tC.expect()
			_, err := services.CheckUsersExist(context.TODO(), tC.input)
			if !errors.Is(err, tC.err) {
				t.Errorf("\n got err = %v, \n wants err = %v", err, tC.err)
			}
		})
	}
}
func TestGetDevInfo(t *testing.T) {

	mockRepo, _, services := before(t)
	userInfo := []*types.User{
		{
			Email:  "thang12@gmail.com",
			UserID: "12",
			Role:   1,
		},
		{
			Email:  "thang34@gmail.com",
			UserID: "34",
			Role:   2,
		},
		{
			Email:  "thang56@gmail.com",
			UserID: "56",
			Role:   3,
		},
	}
	ids := []string{"12", "34", "56"}
	testCases := []struct {
		desc   string
		expect func()
		err    error
		input  []string
	}{
		{
			desc: "Database error ",
			expect: func() {
				mockRepo.EXPECT().FindDevsByID(gomock.Any(), ids).Return(nil, status.Gen().Internal)
			},
			err:   status.Gen().Internal,
			input: ids,
		},
		{
			desc: "Not found user",
			expect: func() {
				mockRepo.EXPECT().FindDevsByID(gomock.Any(), ids).Return(nil, mgo.ErrNotFound)
			},
			err:   status.User().NotFoundUser,
			input: ids,
		},
		{
			desc: "Found user",
			expect: func() {
				mockRepo.EXPECT().FindDevsByID(gomock.Any(), ids).Return(userInfo, nil)
			},
			err:   nil,
			input: ids,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			tC.expect()
			_, err := services.GetDevsInfo(context.TODO(), tC.input)
			if !errors.Is(err, tC.err) {
				t.Errorf("\n got err = %v, \n wants err = %v", err, tC.err)
			}
		})
	}
}
