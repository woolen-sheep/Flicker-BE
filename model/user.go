package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const colNameUser = "user"

type UserInterface interface {
	userC() *mongo.Collection
	AddUser(user User) (primitive.ObjectID, error)
	//UpdateUser(id primitive.ObjectID, toUpdate bson.M) error
}

// User struct in model layer
type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	Mail     string             `bson:"mail"`
	Username string             `bson:"username"`
	Password string             `bson:"password"`
	// CreateTime is time of signing up
	CreateTime int64 `bson:"create_time"`
	// Favorite card sets of user
	Favorite []primitive.ObjectID `bson:"favorite"`
	// Avatar image url
	Avatar string `bson:"avatar"`
}

func (m *model) userC() *mongo.Collection {
	return m.db.Collection(colNameUser)
}

// AddUser inserts a new user into database and will return `ErrExist`
// when user with same mail exist in database
func (m *model) AddUser(user User) (primitive.ObjectID, error) {
	res := primitive.ObjectID{}
	user.ID = primitive.NewObjectID()
	// set init value, otherwise $addToSet will have problems
	user.Favorite = []primitive.ObjectID{}
	user.CreateTime = time.Now().Unix()

	filter := bson.M{
		"mail": user.Mail,
	}
	update := bson.M{
		"$setOnInsert": user,
	}
	boolTrue := true
	opt := options.UpdateOptions{
		Upsert: &boolTrue,
	}
	result, err := m.userC().UpdateOne(m.ctx, filter, update, &opt)
	if err != nil {
		return res, err
	}
	if result.UpsertedCount != 1 {
		return res, ErrExist
	}
	return user.ID, err
}
