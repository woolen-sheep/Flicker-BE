package model

import (
	"time"

	"github.com/woolen-sheep/Flicker-BE/constant"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const colNameCardset = "cardset"

type CardsetInterface interface {
	AddCardset(cardset Cardset) (string, error)
	GetCardset(id string) (Cardset, bool, error)
	DeleteCardset(cardset Cardset) (bool, error)
	UpdateCardset(cardset Cardset) error
	GetCardsetWithOwner(id, owner string) (Cardset, error)
}

// Cardset struct in model layer
type Cardset struct {
	ID          primitive.ObjectID `bson:"_id"`
	OwnerID     primitive.ObjectID `bson:"owner_id,omitempty"`
	Name        string             `bson:"name,omitempty"`
	Description string             `bson:"description,omitempty"`
	Access      int                `bson:"access,omitempty"`
	Status      int                `bson:"status,omitempty"`

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
	filter := bson.M{
		"_id":      cardset.ID,
		"owner_id": cardset.OwnerID,
		"status":   constant.StatusNormal,
	}
	res, err := m.cardsetC().UpdateOne(m.ctx, filter, bson.M{"$set": cardset})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}

// DeleteCardset by id, returns the cardset exist and error
func (m *model) DeleteCardset(cardset Cardset) (bool, error) {
	err := m.UpdateCardset(cardset)
	if err == ErrNotFound {
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
	filter := bson.M{
		"_id":    cardsetID,
		"status": constant.StatusNormal,
	}
	err = m.cardsetC().FindOne(m.ctx, filter).Decode(&cardset)
	if err == mongo.ErrNoDocuments {
		return cardset, false, nil
	}
	if err != nil {
		return cardset, false, err
	}
	return cardset, true, err
}

// GetCardsetWithOwner by id and owner id, returns the cardset struct, whether the cardset exist and error
func (m *model) GetCardsetWithOwner(id, owner string) (Cardset, error) {
	cardset := Cardset{}
	cardsetID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return cardset, err
	}
	ownerID, err := primitive.ObjectIDFromHex(owner)
	if err != nil {
		return cardset, err
	}
	filter := bson.M{
		"_id":      cardsetID,
		"owner_id": ownerID,
		"status":   constant.StatusNormal,
	}
	err = m.cardsetC().FindOne(m.ctx, filter).Decode(&cardset)
	if err == mongo.ErrNoDocuments {
		return cardset, nil
	}
	return cardset, err
}
