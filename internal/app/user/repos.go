package user

import (
	"context"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	"praslar.com/gotasma/internal/app/types"
	"praslar.com/gotasma/internal/pkg/db"
	"praslar.com/gotasma/internal/pkg/util/timeutil"
)

type (
	MongoDBRepository struct {
		session *mgo.Session
	}
)

func NewMongoDBRespository(session *mgo.Session) *MongoDBRepository {
	return &MongoDBRepository{
		session: session,
	}
}

func (r *MongoDBRepository) Create(ctx context.Context, user *types.User) (string, error) {
	s := r.session.Clone()
	defer s.Close()
	user.ID = db.NewID()
	user.CreatedAt = timeutil.Now()
	user.UpdateAt = user.CreatedAt

	if err := r.collection(s).Insert(user); err != nil {
		return "", err
	}
	return user.ID, nil
}

func (r *MongoDBRepository) FindByEmail(ctx context.Context, email string) (*types.User, error) {
	selector := bson.M{"email": email}
	s := r.session.Clone()
	defer s.Close()
	var user *types.User
	if err := r.collection(s).Find(selector).One(&user); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *MongoDBRepository) collection(s *mgo.Session) *mgo.Collection {
	return s.DB("").C("user")
}
