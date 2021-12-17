package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const colNameRecord = "record"

type RecordInterface interface {
	UpdateRecord(record Record) error
	GetRecords(cardsetID, ownerID primitive.ObjectID) ([]Record, error)
}

// Card struct in model layer
type Record struct {
	ID        primitive.ObjectID `bson:"_id"`
	OwnerID   primitive.ObjectID `bson:"owner_id,omitempty"`
	CardsetID primitive.ObjectID `bson:"cardset_id,omitempty"`
	// CardID is ID of card that the record belongs to
	CardID     primitive.ObjectID `bson:"card_id,omitempty"`
	LastStudy  int64              `bson:"last_study,omitempty"`
	StudyTimes int                `bson:"study_times,omitempty"`
	// Status 0 means not finish study and 1 otherwise
	Status int `bson:"status,omitempty"`
}

func (m *model) recordC() *mongo.Collection {
	return m.db.Collection(colNameRecord)
}

// UpdateRecord inserts a new study record into database when it's not exist
// and update `last_study` and `study_times` otherwise.
func (m *model) UpdateRecord(record Record) error {
	record.ID = primitive.NewObjectID()

	record.LastStudy = time.Now().Unix()
	filter := bson.M{
		"owner_id":   record.OwnerID,
		"cardset_id": record.CardsetID,
		"card_id":    record.CardID,
	}
	update := bson.D{
		{Key: "$setOnInsert", Value: record},
		{Key: "$inc", Value: bson.M{"study_times": 1}},
	}
	boolTrue := true
	opt := options.UpdateOptions{
		Upsert: &boolTrue,
	}
	_, err := m.recordC().UpdateOne(m.ctx, filter, update, &opt)
	return err
}

// GetRecords of certain cardset and user
func (m *model) GetRecords(cardsetID, ownerID primitive.ObjectID) ([]Record, error) {
	records := []Record{}
	filter := bson.M{
		"cardset_id": cardsetID,
		"owner_id":   ownerID,
	}
	res, err := m.recordC().Find(m.ctx, filter)
	if err != nil {
		return records, err
	}
	err = res.All(m.ctx, &records)
	return records, err
}
