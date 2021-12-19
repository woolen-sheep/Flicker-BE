package model

import (
	"time"

	"github.com/woolen-sheep/Flicker-BE/constant"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const colNameCardset = "cardset"

type CardsetInterface interface {
	AddCardset(cardset Cardset) (string, error)
	GetCardset(id string) (Cardset, bool, error)
	DeleteCardset(cardset Cardset) (bool, error)
	UpdateCardset(cardset Cardset) error
	GetCardsetWithOwner(id, owner string) (Cardset, error)
	GetCardsetByOwner(owner string) ([]Cardset, error)
	GetCardsetByIDList(ids []primitive.ObjectID) ([]Cardset, error)
	GetCardsetByKeyword(keyword string, skip, limit int) ([]Cardset, error)
	GetRandomCardset(count int) ([]Cardset, error)
	UpdateFavoriteCount(id string, liked bool) error
	UpdateVisitCount(id string) error
}

// Cardset struct in model layer
type Cardset struct {
	ID            primitive.ObjectID `bson:"_id"`
	OwnerID       primitive.ObjectID `bson:"owner_id,omitempty"`
	Name          string             `bson:"name,omitempty"`
	Description   string             `bson:"description,omitempty"`
	Access        int                `bson:"access,omitempty"`
	Status        int                `bson:"status,omitempty"`
	FavoriteCount int                `bson:"favorite_count,omitempty"`
	VisitCount    int                `bson:"visit_count,omitempty"`

	// CreateTime is the first time of adding the cardset
	CreateTime int64 `bson:"create_time,omitempty"`
	// LastUpdateTime is the last time of updating the cardset
	LastUpdateTime int64 `bson:"update_time,omitempty"`
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

// GetCardsetWithOwner by id and owner id, returns the cardset struct, whether the cardset exists and error
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

// GetCardsetByOwner owner id, returns the cardset list created by the user.
func (m *model) GetCardsetByOwner(owner string) ([]Cardset, error) {
	cardset := []Cardset{}
	ownerID, err := primitive.ObjectIDFromHex(owner)
	if err != nil {
		return cardset, err
	}
	filter := bson.M{
		"owner_id": ownerID,
		"status":   constant.StatusNormal,
	}
	res, err := m.cardsetC().Find(m.ctx, filter)
	if err != nil {
		return cardset, err
	}
	err = res.All(m.ctx, &cardset)
	if err == mongo.ErrNoDocuments {
		return cardset, nil
	}
	return cardset, err
}

// GetCardsetByIDList returns the card struct, whether the card exist and error.
func (m *model) GetCardsetByIDList(ids []primitive.ObjectID) ([]Cardset, error) {
	card := []Cardset{}
	filter := bson.M{
		"_id": bson.M{
			"$in": ids,
		},
		"status": constant.StatusNormal,
	}
	res, err := m.cardsetC().Find(m.ctx, filter)
	if err == mongo.ErrNoDocuments {
		return card, nil
	}
	if err != nil {
		return card, err
	}
	err = res.All(m.ctx, &card)
	return card, err
}

// GetCardsetByKeyword returns cardsets which contains keyword in name or description.
func (m *model) GetCardsetByKeyword(keyword string, skip, limit int) ([]Cardset, error) {
	card := []Cardset{}
	filter := bson.M{
		"$or": bson.A{
			bson.M{
				"name": bson.M{"$regex": keyword},
			},
			bson.M{
				"description": bson.M{"$regex": keyword},
			},
		},
		"access": constant.CardsetAccessPublic,
		"status": constant.StatusNormal,
	}
	opt := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
	res, err := m.cardsetC().Find(m.ctx, filter, opt)
	if err == mongo.ErrNoDocuments {
		return card, nil
	}
	if err != nil {
		return card, err
	}
	err = res.All(m.ctx, &card)
	return card, err
}

// GetRandomCardset returns cardsets which contains keyword in name or description.
func (m *model) GetRandomCardset(count int) ([]Cardset, error) {
	card := []Cardset{}
	cur, err := m.cardsetC().Aggregate(m.ctx, mongo.Pipeline{
		{{
			Key: "$match", Value: bson.M{
				"access": constant.CardsetAccessPublic,
				"status": constant.StatusNormal,
			},
		}},
		{{
			Key: "$sample", Value: bson.M{
				"size": count,
			},
		}},
	})
	if err == mongo.ErrNoDocuments {
		return card, nil
	}
	if err != nil {
		return card, err
	}
	err = cur.All(m.ctx, &card)
	return card, err
}

// UpdateFavoriteCount will increase favorite count when liked is false
// and decrease otherwise.
func (m *model) UpdateFavoriteCount(id string, liked bool) error {
	cardsetID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": cardsetID,
	}

	count := 1
	if liked {
		count = -1
	}

	update := bson.M{
		"$inc": bson.M{
			"favorite_count": count,
		},
	}

	_, err = m.cardsetC().UpdateOne(m.ctx, filter, update)
	return err
}

// UpdateVisitCount will increase visit count.
func (m *model) UpdateVisitCount(id string) error {
	cardsetID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": cardsetID,
	}

	update := bson.M{
		"$inc": bson.M{
			"visit_count": 1,
		},
	}

	_, err = m.cardsetC().UpdateOne(m.ctx, filter, update)
	return err
}
