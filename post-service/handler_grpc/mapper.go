package handler_grpc

import (
	"post-service/model"
	"strconv"

	pb "github.com/MihajloMarjanski/xws-project/common/proto/post_service"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func mapProtoToCommentDTO(comment *pb.Comment) model.CommentDTO {
	userId, _ := strconv.ParseUint(comment.User, 0, 32)
	postId, _ := primitive.ObjectIDFromHex(comment.Post)
	commentPb := model.CommentDTO{
		PostID: postId,
		Text:   comment.Text,
		UserID: uint(userId),
	}
	return commentPb
}
