package model

import (
	"time"

	"github.com/woolen-sheep/Flicker-BE/constant"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const colNameComment = "comment"

type CommentInterface interface {
	AddComment(comment Comment) (string, error)
	GetComments(cardID string) ([]Comment, bool, error)
	UpdateComment(comment Comment) error
	DeleteComment(comment Comment) (bool, error)
	GetCommentWithOwner(commentID, owner string) (Comment, error)
}

// Comment struct in model layer
type Comment struct {
	ID      primitive.ObjectID `bson:"_id"`
	OwnerID primitive.ObjectID `bson:"owner_id,omitempty"`
	CardID  primitive.ObjectID `bson:"card_id,omitempty"`
	Content string             `bson:"comment,omitempty"`
	Status  int                `bson:"status,omitempty"`

	// Add more properties for comment

	// CreateTime is the first time of adding the comment
	CreateTime int64 `bson:"create_time"`
	// LastUpdateTime is the last time of updating the comment
	LastUpdateTime int64 `bson:"update_time"`
}

func (m *model) commentC() *mongo.Collection {
	return m.db.Collection(colNameComment)
}

// AddComment inserts a new comment into database
func (m *model) AddComment(comment Comment) (string, error) {
	comment.ID = primitive.NewObjectID()

	comment.CreateTime = time.Now().Unix()
	comment.LastUpdateTime = comment.CreateTime

	if _, err := m.commentC().InsertOne(m.ctx, &comment); err != nil {
		return "", err
	}
	return comment.ID.Hex(), nil
}

// GetComments by id, returns the comments slice, whether the comment exist and error
func (m *model) GetComments(cardID string) ([]Comment, bool, error) {
	var comments []Comment
	cardObjectID, err := primitive.ObjectIDFromHex(cardID)
	if err != nil {
		return comments, false, err // error in ObjectIDFromHex
	}

	filter := bson.M{
		"card_id": cardObjectID,
		"status":  constant.StatusNormal,
	}

	cursor, err := m.commentC().Find(m.ctx, filter)
	if err == mongo.ErrNoDocuments {
		return comments, false, nil
	}
	if err != nil {
		return comments, true, err // comment exists, but error in Find
	}

	defer cursor.Close(m.ctx)
	for cursor.Next(m.ctx) {
		var cmt Comment
		if err = cursor.Decode(&cmt); err != nil {
			return comments, true, err // comment exists, but error in Decode
		}
		comments = append(comments, cmt)
	}

	return comments, true, err
}

// UpdateComment updates comment info in database, empty fields will be ignore
func (m *model) UpdateComment(comment Comment) error {
	comment.LastUpdateTime = time.Now().Unix()
	filter := bson.M{
		"_id":    comment.ID,
		"status": constant.StatusNormal,
	}
	_, err := m.commentC().UpdateOne(m.ctx, filter, bson.M{"$set": comment})
	return err
}

// DeleteComment by id, returns the comment exist and error
func (m *model) DeleteComment(comment Comment) (bool, error) {
	err := m.UpdateComment(comment)
	if err == ErrNotFound {
		return false, nil
	}
	if err != nil {
		return true, err
	}
	return true, err
}

// GetCommentWithOwner by comment id & owner id, returns the comment struct, whether the comment exists and error
func (m *model) GetCommentWithOwner(commentID, owner string) (Comment, error) {
	comment := Comment{}
	commentObjID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return comment, err
	}
	ownerID, err := primitive.ObjectIDFromHex(owner)
	if err != nil {
		return comment, err
	}
	filter := bson.M{
		"_id":      commentObjID,
		"owner_id": ownerID,
		"status":   constant.StatusNormal,
	}
	err = m.commentC().FindOne(m.ctx, filter).Decode(&comment)
	if err == mongo.ErrNoDocuments {
		return comment, nil
	}
	return comment, err
}
