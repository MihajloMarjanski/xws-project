package repo

import (
	"context"
	"fmt"
	"log"
	"post-service/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostRepository struct {
	posts *mongo.Collection
}

func New() (*PostRepository, error) {

	uri := fmt.Sprintf("mongodb://localhost:27017")
	options := options.Client().ApplyURI(uri)
	client, _ := mongo.Connect(context.TODO(), options)

	posts := client.Database("post").Collection("post")

	repo := &PostRepository{
		posts: posts,
	}

	return repo, nil
}

func (repo *PostRepository) CreatePost(post *model.Post) error {
	post.ID = primitive.NewObjectID()
	result, err := repo.posts.InsertOne(context.TODO(), post)
	if err != nil {
		return err
	}
	post.ID = result.InsertedID.(primitive.ObjectID)

	return nil

}

func (repo *PostRepository) AddComment(comment *model.CommentDTO) error {
	post := repo.GetById(comment.PostID)
	createdComment := model.Comment{
		UserID:    comment.UserID,
		Text:      comment.Text,
		CreatedAt: time.Now().UTC(),
	}
	post.Commnets = append(post.Commnets, createdComment)

	filter := bson.M{"_id": post.ID}
	_, err := repo.posts.ReplaceOne(context.TODO(), filter, post)

	if err != nil {
		log.Fatal(err)
	}
	return nil

}

func (repo *PostRepository) GetById(id primitive.ObjectID) *model.Post {
	filter := bson.M{"_id": id}
	post, _ := repo.filterOne(filter)
	return post
}
func (repo *PostRepository) filterOne(filter interface{}) (post *model.Post, err error) {
	result := repo.posts.FindOne(context.TODO(), filter)
	err = result.Decode(&post)
	return
}
