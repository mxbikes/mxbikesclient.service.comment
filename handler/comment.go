package handler

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gogo/status"
	"github.com/google/uuid"
	"github.com/mxbikes/mxbikesclient.service.comment/models"
	"github.com/mxbikes/mxbikesclient.service.comment/repository"
	protobuffer "github.com/mxbikes/protobuf/comment"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
)

type Mod struct {
	protobuffer.UnimplementedCommentServiceServer
	repository repository.ModRepository
	logger     logrus.Logger
	validate   *validator.Validate
}

// Return a new handler
func New(postgres repository.ModRepository, logger logrus.Logger) *Mod {
	return &Mod{repository: postgres, validate: validator.New(), logger: logger}
}

func (e *Mod) GetCommentByModID(ctx context.Context, req *protobuffer.GetCommentByModIDRequest) (*protobuffer.GetCommentByModIDResponse, error) {
	// Check if valid uuid
	_, err := uuid.Parse(req.ModID)
	if err != nil {
		e.logger.WithFields(logrus.Fields{"prefix": "SERVICE.Comment_GetCommentByModID"}).Errorf("request ModID is not a valid UUID: {%s}", req.ModID)
		return nil, status.Error(codes.Internal, "Error request value ModID, is not a valid UUID!")
	}

	// Get Requested Comment
	comments, err := e.repository.SearchByModID(req.ModID)
	if err != nil {
		return nil, err
	}

	e.logger.WithFields(logrus.Fields{"prefix": "SERVICE.Comment_GetCommentByModID"}).Infof("mod with id: {%s} ", req.ModID)

	return &protobuffer.GetCommentByModIDResponse{Comments: models.CommentsToProto(comments)}, nil
}

func (e *Mod) UpdateComment(ctx context.Context, req *protobuffer.UpdateCommentRequest) (*protobuffer.UpdateCommentResponse, error) {
	comment := &models.Comment{
		ID:     req.ID,
		ModID:  req.ModID,
		UserID: req.UserID,
		Text:   req.Text,
	}

	// Validate
	err := e.validate.Struct(comment)
	if err != nil {
		e.logger.WithFields(logrus.Fields{"prefix": "SERVICE.Comment_UpdateComment"}).Errorf("request validation is not a valid: {%s}", err)
		return nil, err.(validator.ValidationErrors)
	}

	// Get Requested Comment
	err = e.repository.Save(comment)
	if err != nil {
		return nil, err
	}

	e.logger.WithFields(logrus.Fields{"prefix": "SERVICE.Comment_UpdateComment"}).Infof("mod with id: {%s} ", comment.ModID)

	return &protobuffer.UpdateCommentResponse{}, nil
}

func (e *Mod) DeleteComment(ctx context.Context, req *protobuffer.DeleteCommentRequest) (*protobuffer.DeleteCommentResponse, error) {
	// Validate
	_, err := uuid.Parse(req.ID)
	if err != nil {
		e.logger.WithFields(logrus.Fields{"prefix": "SERVICE.Comment_DeleteComment"}).Errorf("request ID is not a valid UUID: {%s}", req.ID)
		return nil, status.Error(codes.Internal, "Error request value ID, is not a valid UUID!")
	}

	err = e.repository.Delete(req.ID)
	if err != nil {
		return nil, err
	}

	e.logger.WithFields(logrus.Fields{"prefix": "SERVICE.Comment_DeleteComment"}).Infof("mod with id: {%s} ", req.ID)

	return &protobuffer.DeleteCommentResponse{}, nil
}

func (e *Mod) CreateComment(ctx context.Context, req *protobuffer.CreateCommentRequest) (*protobuffer.CreateCommentResponse, error) {
	comment := &models.Comment{
		ModID:  req.ModID,
		UserID: req.UserID,
		Text:   req.Text,
	}

	// Validate
	err := e.validate.Struct(comment)
	if err != nil {
		e.logger.WithFields(logrus.Fields{"prefix": "SERVICE.Comment_CreateComment"}).Errorf("request validation is not a valid: {%s}", err)
		return nil, err.(validator.ValidationErrors)
	}

	// Get Requested Comment
	err = e.repository.Save(comment)
	if err != nil {
		return nil, err
	}

	e.logger.WithFields(logrus.Fields{"prefix": "SERVICE.Comment_CreateComment"}).Infof("mod with id: {%s} ", comment.ID)

	return &protobuffer.CreateCommentResponse{ID: comment.ID}, nil
}
