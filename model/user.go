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
	AddUser(user User) (string, error)
	GetUser(id string) (User, bool, error)
	GetUserByMail(mail string) (user User, err error)
	UpdateUser(user User) error
}

// User struct in model layer
type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	Mail     string             `bson:"mail,omitempty"`
	Username string             `bson:"username,omitempty"`
	Password string             `bson:"password,omitempty"`
	// CreateTime is time of signing up
	CreateTime int64 `bson:"create_time"`
	// Favorite card sets of user
	Favorite []primitive.ObjectID `bson:"favorite,omitempty"`
	// Avatar image url
	Avatar string `bson:"avatar,omitempty"`
}

func (m *model) userC() *mongo.Collection {
	return m.db.Collection(colNameUser)
}

// AddUser inserts a new user into database and will return `ErrExist`
// when user with same mail exist in database
func (m *model) AddUser(user User) (string, error) {
	var res = ""
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
	return user.ID.Hex(), err
}

// GetUser by id, returns the user struct, whether the user exist and error
func (m *model) GetUser(id string) (User, bool, error) {
	user := User{}
	userID, err := primitive.ObjectIDFromHex(id)
	err = m.userC().FindOne(m.ctx, bson.M{"_id": userID}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return user, false, nil
	}
	if err != nil {
		return user, false, err
	}
	return user, true, err
}

// GetUserByMail will get user by mail
func (m *model) GetUserByMail(mail string) (user User, err error) {
	filter := bson.D{{"mail", mail}}
	err = m.userC().FindOne(m.ctx, filter).Decode(&user)
	return user, err
}

// UpdateUser updates user info in database, empty fields will be ignore
func (m *model) UpdateUser(user User) error {
	_, err := m.userC().UpdateOne(m.ctx, bson.M{"_id": user.ID}, bson.M{"$set": user})
	return err
}
