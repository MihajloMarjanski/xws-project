package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	UserID    uint               `bson:"userId" json:"userId"`
	Title     string             `bson:"title" json:"title"`
	Text      string             `bson:"text" json:"text"`
	Img       string             `bson:"img" json:"img"`
	Link      string             `bson:"link" json:"link"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	Commnets  []Comment          `bson:"comments" json:"comments"`
	Likes     []uint             `bson:"likes" json:"likes"`
	Dislikes  []uint             `bson:"dilikes" json:"dilikes"`
}

type Comment struct {
	UserID    uint      `bson:"userId" json:"userId"`
	Text      string    `bson:"text" json:"text"`
	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
}

type CommentDTO struct {
	PostID primitive.ObjectID `bson:"postId" json:"postId"`
	UserID uint               `bson:"userId" json:"userId"`
	Text   string             `bson:"text" json:"text"`
}
