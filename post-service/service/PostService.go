package service

import (
	"fmt"
	"post-service/model"
	"post-service/repo"
	"time"
)

type PostService struct {
	postRepo *repo.PostRepository
}

func New() (*PostService, error) {

	postRepo, err := repo.New()
	if err != nil {
		return nil, err
	}

	return &PostService{
		postRepo: postRepo,
	}, nil
}

func (service *PostService) CreatePost(title string, text string, img string, link string, userId uint) model.Post {
	post := model.Post{
		ID:        1,
		UserID:    userId,
		Title:     title,
		Text:      text,
		Img:       img,
		Link:      link,
		CreatedAt: time.Now(),
	}
	fmt.Println(post)

	return post

}
