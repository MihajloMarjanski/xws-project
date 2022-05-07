package handler_grpc

import (
	"context"
	"post-service/service"

	pb "github.com/MihajloMarjanski/xws-project/common/proto/post_service"
)

type PostHandler struct {
	pb.UnimplementedPostServiceServer
	postService *service.PostService
}

func New() (*PostHandler, error) {

	postService, err := service.New()
	if err != nil {
		return nil, err
	}

	return &PostHandler{
		postService: postService,
	}, nil
}

func (handler *PostHandler) CreatePost(ctx context.Context, request *pb.CreatePostRequest) (*pb.CreatePostResponse, error) {
	post := mapProtoToPost(request.Post)
	createdPost := handler.postService.CreatePost(post.Title, post.Text, post.Img, post.Link, post.UserID)
	response := &pb.CreatePostResponse{
		Id: int64(createdPost.ID),
	}
	return response, nil
}
