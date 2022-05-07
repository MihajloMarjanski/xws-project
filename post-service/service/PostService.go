package service

import (
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
