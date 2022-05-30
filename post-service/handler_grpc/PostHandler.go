package handler_grpc

import (
	"context"
	"fmt"
	"post-service/model"
	"post-service/service"
	"strconv"
	"strings"

	pb "github.com/MihajloMarjanski/xws-project/common/proto/post_service"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/metadata"
)

type PostHandler struct {
	pb.UnimplementedPostServiceServer
	postService *service.PostService
}

func Verify(accessToken string) (*model.Claims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&model.Claims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return []byte("tajni_kljuc_za_jwt_hash"), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*model.Claims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func GetUserID(ctx context.Context) uint {
	md, _ := metadata.FromIncomingContext(ctx)
	values := md["authorization"]
	accessToken := values[0]
	words := strings.Fields(accessToken)

	claims, _ := Verify(words[1])
	id, _ := strconv.ParseUint(claims.Id, 10, 64)
	return uint(id)
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
	post.UserID = GetUserID(ctx)
	createdPost := handler.postService.CreatePost(post.Title, post.Text, post.Img, post.Link, post.UserID)
	response := &pb.CreatePostResponse{
		Id: createdPost.ID.Hex(),
	}
	return response, nil
}

func (handler *PostHandler) AddComment(ctx context.Context, request *pb.AddCommentRequest) (*pb.AddCommnetResponse, error) {
	comment := mapProtoToCommentDTO(request.Comment)
	comment.UserID = GetUserID(ctx)
	handler.postService.AddComment(&comment)
	response := &pb.AddCommnetResponse{
		Id: comment.Text,
	}
	return response, nil
}

func (handler *PostHandler) AddLike(ctx context.Context, request *pb.AddLikeRequest) (*pb.AddLikeResponse, error) {
	like := mapProtoToLikeDTO(request.Like)
	like.UserID = GetUserID(ctx)
	handler.postService.AddLike(&like)
	response := &pb.AddLikeResponse{
		Id: strconv.FormatUint(uint64(like.UserID), 10),
	}
	return response, nil
}

func (handler *PostHandler) AddDislike(ctx context.Context, request *pb.AddDislikeRequest) (*pb.AddLikeResponse, error) {
	like := mapProtoToLikeDTO(request.Dislike)
	like.UserID = GetUserID(ctx)
	handler.postService.AddDislike(&like)
	response := &pb.AddLikeResponse{
		Id: strconv.FormatUint(uint64(like.UserID), 10),
	}
	return response, nil
}

func (handler *PostHandler) GetPostsForUser(ctx context.Context, request *pb.User) (*pb.PostsResponse, error) {
	id, _ := strconv.ParseUint(request.User, 0, 32)
	posts := handler.postService.GetPostsForUser(uint(id))
	pbPosts := mapPostsToProto(posts)
	response := &pb.PostsResponse{Post: pbPosts}
	return response, nil
}
