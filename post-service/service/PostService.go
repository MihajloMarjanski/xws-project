package service

import (
	"fmt"
	pbReq "github.com/MihajloMarjanski/xws-project/common/proto/requests_service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"os"
	"path/filepath"
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

	//for _, user := range service.FindAllConnections(userId) {
	service.SendNotification(userId, "A new post has been created by user '")
	//}
	return post
}

func (service *PostService) SendNotification(id uint, message string) []model.User {
	var res []model.User
	crtTlsPath, _ := filepath.Abs("./service.pem")

	creds, err6 := credentials.NewClientTLSFromFile(crtTlsPath, "")
	if err6 != nil {
		//log.Fatalf("could not process the credentials: %v", err6)
	}

	conn, err := grpc.Dial("request-service:8200", grpc.WithTransportCredentials(creds))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pbReq.NewRequestsServiceClient(conn)
	users, err := client.FindConnections(context.Background(), &pbReq.FindConnectionsRequest{Id: int64(id)})
	if err != nil {
		panic(err)
	}
	for _, user := range users.Users {
		_, err = client.SendNotification(context.Background(), &pbReq.SendNotificationRequest{SenderId: int64(id), ReceiverId: user.Id, Message: message})
		if err != nil {
			panic(err)
		}
	}

	return res
}

//func (s *PostService) SendNotification(senderID, receiverId uint, message string) {
//	conn, err := grpc.Dial("localhost:8200", grpc.WithInsecure())
//	if err != nil {
//		panic(err)
//	}
//	defer conn.Close()
//	client := pbReq.NewRequestsServiceClient(conn)
//	_, err = client.SendNotification(context.Background(), &pbReq.SendNotificationRequest{SenderId: int64(senderID), ReceiverId: int64(receiverId), Message: message})
//	if err != nil {
//		panic(err)
//	}
//}

func mapUser(user *pbReq.User) model.User {
	res := model.User{
		ID:        uint(user.Id),
		UserName:  user.Username,
		Biography: user.Biography,
		Name:      user.Name,
	}
	return res
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
