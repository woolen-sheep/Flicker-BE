package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const colNameCardset = "cardset"

type CardsetInterface interface {
	AddCardset(cardset Cardset) (string, error)
	GetCardset(id string) (Cardset, bool, error)
	DeleteCardset(id string) (bool, error)
	UpdateCardset(cardset Cardset) error
}

// Cardset struct in model layer
type Cardset struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string             `bson:"name,omitempty"`
	Description string             `bson:"description,omitempty"`
	Access      int                `bson:"access,omitempty"`
	Cards       []string           `bson:"cards,omitempty"`

	// CreateTime is the first time of adding the cardset
	CreateTime int64 `bson:"create_time"`
	// LastUpdateTime is the last time of updating the cardset
	LastUpdateTime int64 `bson:"update_time"`
}

func (m *model) cardsetC() *mongo.Collection {
	return m.db.Collection(colNameCardset)
}

// AddCardset inserts a new cardset into database and will return `ErrExist`
// when cardset with same mail exist in database
func (m *model) AddCardset(cardset Cardset) (string, error) {
	var res = ""
	cardset.ID = primitive.NewObjectID()
	cardset.CreateTime = time.Now().Unix()
	cardset.LastUpdateTime = cardset.CreateTime

	filter := bson.M{
		"name": cardset.Name,
	}
	update := bson.M{
		"$setOnInsert": cardset,
	}
	boolTrue := true
	opt := options.UpdateOptions{
		Upsert: &boolTrue,
	}
	result, err := m.cardsetC().UpdateOne(m.ctx, filter, update, &opt)
	if err != nil {
		return res, err
	}
	if result.UpsertedCount != 1 {
		return res, ErrExist
	}
	return cardset.ID.Hex(), err
}

// UpdateCardset updates cardset info in database, empty fields will be ignore
func (m *model) UpdateCardset(cardset Cardset) error {
	cardset.LastUpdateTime = time.Now().Unix()
	_, err := m.cardsetC().UpdateOne(m.ctx, bson.M{"_id": cardset.ID}, bson.M{"$set": cardset})
	return err
}

// DeleteCardset by id, returns the cardset exist and error
func (m *model) DeleteCardset(id string) (bool, error) {
	cardsetID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}
	res, err := m.cardsetC().DeleteOne(m.ctx, bson.M{"_id": cardsetID})
	if res.DeletedCount == 0 {
		return false, nil
	}
	if err != nil {
		return true, err
	}
	return true, err
}

// GetCardset by id, returns the cardset struct, whether the cardset exist and error
func (m *model) GetCardset(id string) (Cardset, bool, error) {
	cardset := Cardset{}
	cardsetID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return cardset, false, err
	}
	err = m.cardsetC().FindOne(m.ctx, bson.M{"_id": cardsetID}).Decode(&cardset)
	if err == mongo.ErrNoDocuments {
		return cardset, false, nil
	}
	if err != nil {
		return cardset, false, err
	}
	return cardset, true, err
}
