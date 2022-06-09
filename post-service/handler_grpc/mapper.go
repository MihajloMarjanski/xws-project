package handler_grpc

import (
	"post-service/model"
	"strconv"

	pb "github.com/MihajloMarjanski/xws-project/common/proto/post_service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func mapProtoToLikeDTO(like *pb.Like) model.LikeDTO {
	userId, _ := strconv.ParseUint(like.User, 0, 32)
	postId, _ := primitive.ObjectIDFromHex(like.Post)
	likePb := model.LikeDTO{
		PostID: postId,
		UserID: uint(userId),
	}
	return likePb
}

func mapPostsToProto(posts []model.Post) []*pb.PostResp {
	var postResponse []*pb.PostResp

	for _, post := range posts {
		postResp := pb.PostResp{
			Id:        post.ID.Hex(),
			User:      int64(post.UserID),
			Title:     post.Title,
			Text:      post.Text,
			Img:       post.Img,
			Link:      post.Link,
			CreatedAt: timestamppb.New(post.CreatedAt),
		}
		for _, comment := range post.Commnets {
			commentResp := pb.CommentResp{
				User:      int64(comment.UserID),
				Text:      comment.Text,
				CreatedAt: timestamppb.New(comment.CreatedAt),
			}
			postResp.Comments = append(postResp.Comments, &commentResp)
		}
		for _, like := range post.Likes {
			likePb := pb.LikeResp{
				User: int64(like),
			}
			postResp.Like = append(postResp.Like, &likePb)
		}
		for _, dislikes := range post.Dislikes {
			dislikePb := pb.LikeResp{
				User: int64(dislikes),
			}
			postResp.Dislike = append(postResp.Dislike, &dislikePb)
		}

		postResponse = append(postResponse, &postResp)
	}

	return postResponse
}
