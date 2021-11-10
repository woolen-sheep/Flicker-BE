package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	ID          primitive.ObjectID   `bson:"_id"`
	Name        string               `bson:"name,omitempty"`
	Description string               `bson:"description,omitempty"`
	Access      int                  `bson:"access,omitempty"`
	Cards       []primitive.ObjectID `bson:"cards,omitempty"`

	// CreateTime is the first time of adding the cardset
	CreateTime int64 `bson:"create_time"`
	// LastUpdateTime is the last time of updating the cardset
	LastUpdateTime int64 `bson:"update_time"`
}

func (m *model) cardsetC() *mongo.Collection {
	return m.db.Collection(colNameCardset)
}

// AddCardset inserts a new cardset into database
func (m *model) AddCardset(cardset Cardset) (string, error) {
	cardset.ID = primitive.NewObjectID()
	cardset.Cards = make([]primitive.ObjectID, 0)

	cardset.CreateTime = time.Now().Unix()
	cardset.LastUpdateTime = cardset.CreateTime

	if _, err := m.cardsetC().InsertOne(m.ctx, &cardset); err != nil {
		return "", err
	}
	return cardset.ID.Hex(), nil
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
