package handler_grpc

import (
	"post-service/model"
	"strconv"

	pb "github.com/MihajloMarjanski/xws-project/common/proto/post_service"
)

func mapProtoToPost(post *pb.Post) model.Post {
	userId, _ := strconv.ParseUint(post.User, 0, 32)
	postPb := model.Post{
		Title:  post.Title,
		Text:   post.Text,
		Img:    post.Img,
		Link:   post.Link,
		UserID: uint(userId),
	}
	return postPb
}
