package holiday

import (
	"context"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	"praslar.com/gotasma/internal/app/types"
)

type (
	MongoDBRepository struct {
		session *mgo.Session
	}
)

func NewMongoRepository(session *mgo.Session) *MongoDBRepository {
	return &MongoDBRepository{
		session: session,
	}
}

func (r *MongoDBRepository) Create(ctx context.Context, holiday *types.Holiday) error {
	s := r.session.Clone()
	defer s.Close()
	holiday.CreatedAt = time.Now()
	holiday.UpdateAt = holiday.CreatedAt
	if err := r.collection(s).Insert(holiday); err != nil {
		return err
	}
	return nil
}

func (r *MongoDBRepository) FindByTitle(ctx context.Context, title string, createrID string) (*types.Holiday, error) {
	selector := bson.M{"title": title, "creater_id": createrID}
	s := r.session.Clone()
	defer s.Close()
	var holiday *types.Holiday
	if err := r.collection(s).Find(selector).One(&holiday); err != nil {
		return nil, err
	}
	return holiday, nil
}

func (r *MongoDBRepository) collection(s *mgo.Session) *mgo.Collection {
	return s.DB("").C("holiday")
}

func (r *MongoDBRepository) Delete(ctx context.Context, id string) error {
	s := r.session.Clone()
	defer s.Close()
	if err := r.collection(s).Remove(bson.M{"holiday_id": id}); err != nil {
		return err
	}
	return nil
}

func (r *MongoDBRepository) FindAll(ctx context.Context, createrID string) ([]*types.Holiday, error) {
	selector := bson.M{"creater_id": createrID}
	s := r.session.Clone()
	defer s.Close()
	var holidays []*types.Holiday
	if err := r.collection(s).Find(selector).All(&holidays); err != nil {
		return nil, err
	}
	return holidays, nil
}
