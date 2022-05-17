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

	uri := fmt.Sprintf("mongodb://postdb:27017")
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

func (repo *PostRepository) GetPostsForUser(userID *uint) []model.Post {

	var posts []model.Post
	filter := bson.M{"userId": userID}
	cur, err := repo.posts.Find(context.TODO(), filter)

	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var elem model.Post
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		posts = append(posts, elem)

	}

	return posts

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

func (repo *PostRepository) AddLike(like *model.LikeDTO) error {
	post := repo.GetById(like.PostID)
	index1 := findIndex(post.Likes, func(n uint) bool {
		return n == like.UserID
	})
	index2 := findIndex(post.Dislikes, func(n uint) bool {
		return n == like.UserID
	})
	if index1 == -1 {
		post.Likes = append(post.Likes, like.UserID)
	} else {
		newLikes := RemoveIndex(post.Likes, index1)
		post.Likes = newLikes
	}
	if index2 != -1 {
		newDislikes := RemoveIndex(post.Dislikes, index2)
		post.Dislikes = newDislikes
	}

	filter := bson.M{"_id": post.ID}
	_, err := repo.posts.ReplaceOne(context.TODO(), filter, post)

	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (repo *PostRepository) AddDislike(dislike *model.LikeDTO) error {
	post := repo.GetById(dislike.PostID)
	index := findIndex(post.Dislikes, func(n uint) bool {
		return n == dislike.UserID
	})
	index2 := findIndex(post.Likes, func(n uint) bool {
		return n == dislike.UserID
	})
	if index == -1 {
		post.Dislikes = append(post.Dislikes, dislike.UserID)
	} else {
		newDislikes := RemoveIndex(post.Dislikes, index)
		post.Dislikes = newDislikes
	}
	if index2 != -1 {
		newLikes := RemoveIndex(post.Likes, index2)
		post.Likes = newLikes
	}

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
func RemoveIndex(s []uint, index int) []uint {
	return append(s[:index], s[index+1:]...)
}
func findIndex[T any](slice []T, matchFunc func(T) bool) int {
	for index, element := range slice {
		if matchFunc(element) {
			return index
		}
	}

	return -1 // not found
}
