package models

import (
	protobuffer "github.com/mxbikes/protobuf/comment"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	ID     string `gorm:"type:uuid;default:uuid_generate_v4()" validate:"omitempty,uuid4"`
	ModID  string `gorm:"type:uuid;" validate:"uuid4,required"`
	UserID string `gorm:"type:varchar(50);not null;default:null" validate:"uuid4,required"`
	Text   string `gorm:"type:varchar(250);not null;default:null" validate:"min=1,max=250"`
}

func CommentToProto(comment *Comment) *protobuffer.Comment {
	return &protobuffer.Comment{
		ID:        comment.ID,
		ModID:     comment.ModID,
		UserID:    comment.UserID,
		Text:      comment.Text,
		Create_At: timestamppb.New(comment.CreatedAt),
	}
}

func CommentsToProto(comments []*Comment) []*protobuffer.Comment {
	result := make([]*protobuffer.Comment, 0, len(comments))
	for _, projection := range comments {
		result = append(result, CommentToProto(projection))
	}
	return result
}
