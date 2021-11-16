package model

import (
	"time"

	"github.com/woolen-sheep/Flicker-BE/constant"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const colNameCard = "card"

type CardInterface interface {
	AddCard(card Card) (string, error)
	GetCard(id string) (Card, bool, error)
	DeleteCard(card Card) (bool, error)
	UpdateCard(card Card) error
}

// Card struct in model layer
type Card struct {
	ID primitive.ObjectID `bson:"_id"`
	// CardsetID is ID of set that the card belongs to
	CardsetID primitive.ObjectID `bson:"cardset_id"`
	Question  string             `bson:"question,omitempty"`
	Answer    string             `bson:"answer,omitempty"`
	Image     string             `bson:"image,omitempty"`
	Audio     string             `bson:"audio,omitempty"`
	Status    int                `bson:"status,omitempty"`

	// CreateTime is the first time of adding the card
	CreateTime int64 `bson:"create_time"`
	// LastUpdateTime is the last time of updating the card
	LastUpdateTime int64 `bson:"update_time"`
}

func (m *model) cardC() *mongo.Collection {
	return m.db.Collection(colNameCard)
}

// AddCard inserts a new card into database
func (m *model) AddCard(card Card) (string, error) {
	card.ID = primitive.NewObjectID()

	card.CreateTime = time.Now().Unix()
	card.LastUpdateTime = card.CreateTime

	if _, err := m.cardC().InsertOne(m.ctx, &card); err != nil {
		return "", err
	}
	return card.ID.Hex(), nil
}

// UpdateCard updates card info in database, empty fields will be ignore
func (m *model) UpdateCard(card Card) error {
	card.LastUpdateTime = time.Now().Unix()
	filter := bson.M{
		"_id":    card.ID,
		"status": constant.StatusNormal,
	}
	_, err := m.cardC().UpdateOne(m.ctx, filter, bson.M{"$set": card})
	return err
}

// DeleteCard by id, returns the card exist and error
func (m *model) DeleteCard(card Card) (bool, error) {
	err := m.UpdateCard(card)
	if err == ErrNotFound {
		return false, nil
	}
	if err != nil {
		return true, err
	}
	return true, err
}

// GetCard by id, returns the card struct, whether the card exist and error
func (m *model) GetCard(id string) (Card, bool, error) {
	card := Card{}
	cardID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return card, false, err
	}
	filter := bson.M{
		"_id":    cardID,
		"status": constant.StatusNormal,
	}
	err = m.cardC().FindOne(m.ctx, filter).Decode(&card)
	if err == mongo.ErrNoDocuments {
		return card, false, nil
	}
	if err != nil {
		return card, false, err
	}
	return card, true, err
}
