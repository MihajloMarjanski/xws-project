package service

import (
	"fmt"
	"os"
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

	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	fmt.Println("hostname:", name)

	post := model.Post{
		UserID:    userId,
		Title:     title,
		Text:      text,
		Img:       img,
		Link:      link,
		CreatedAt: time.Now().UTC(),
		Commnets:  []model.Comment{},
		Likes:     []uint{},
		Dislikes:  []uint{},
	}
	service.postRepo.CreatePost(&post)
	return post

}

func (service *PostService) AddComment(comment *model.CommentDTO) error {
	return service.postRepo.AddComment(comment)
}

func (service *PostService) AddLike(like *model.LikeDTO) error {
	return service.postRepo.AddLike(like)
}

func (service *PostService) AddDislike(like *model.LikeDTO) error {
	return service.postRepo.AddDislike(like)
}

func (service *PostService) GetPostsForUser(userID uint) []model.Post {
	return service.postRepo.GetPostsForUser(&userID)
}
