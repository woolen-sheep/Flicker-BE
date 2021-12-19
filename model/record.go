package model

import (
	"time"

	"github.com/woolen-sheep/Flicker-BE/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const colNameRecord = "record"

type RecordInterface interface {
	UpdateRecord(record Record) error
	GetRecords(cardsetID, ownerID primitive.ObjectID) ([]Record, error)
	ClearRecords(cardsetID, ownerID primitive.ObjectID) error
	GetLastStudyTime(cardsetID, ownerID primitive.ObjectID) int64
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
// Only record in different natural day (UTC+8) will increase `study_times`.
func (m *model) UpdateRecord(record Record) error {
	record.ID = primitive.NewObjectID()
	record.LastStudy = time.Now().Unix()

	oldRecord, err := m.GetRecord(record.CardsetID, record.CardID, record.OwnerID)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	// only first study in a day will increase study_times
	incValue := 1
	if err == nil {
		currentDay := util.GetDayUnix(time.Now())
		if oldRecord.LastStudy > currentDay {
			incValue = 0
		}
	}

	filter := bson.M{
		"owner_id":   record.OwnerID,
		"cardset_id": record.CardsetID,
		"card_id":    record.CardID,
	}
	update := bson.D{
		{Key: "$setOnInsert", Value: record},
		{Key: "$inc", Value: bson.M{"study_times": incValue}},
	}
	boolTrue := true
	opt := options.UpdateOptions{
		Upsert: &boolTrue,
	}
	_, err = m.recordC().UpdateOne(m.ctx, filter, update, &opt)
	return err
}

func (m *model) GetRecord(cardsetID, cardID, ownerID primitive.ObjectID) (Record, error) {
	records := Record{}
	filter := bson.M{
		"cardset_id": cardsetID,
		"card_id":    cardID,
		"owner_id":   ownerID,
	}
	err := m.recordC().FindOne(m.ctx, filter).Decode(&records)
	return records, err
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

// ClearRecords of certain cardset and user
func (m *model) ClearRecords(cardsetID, ownerID primitive.ObjectID) error {
	filter := bson.M{
		"cardset_id": cardsetID,
		"owner_id":   ownerID,
	}
	_, err := m.recordC().DeleteMany(m.ctx, filter)
	return err
}

//GetLastStudyTime of certain cardset and user
func (m *model) GetLastStudyTime(cardsetID, ownerID primitive.ObjectID) int64 {
	record := []Record{}
	filter := bson.M{
		"cardset_id": cardsetID,
		"owner_id":   ownerID,
	}

	opt := options.Find().SetLimit(1).SetSort(bson.M{"time": -1})
	res, err := m.recordC().Find(m.ctx, filter, opt)
	if err != nil {
		return 0
	}

	err = res.All(m.ctx, &record)
	if err != nil || len(record) == 0 {
		return 0
	}

	return record[0].LastStudy
}
