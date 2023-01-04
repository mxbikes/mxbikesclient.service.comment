package handler

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/mxbikes/mxbikesclient.service.comment/repository"
	protobuffer "github.com/mxbikes/protobuf/comment"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const log_failedConn = "an error '%s' was not expected when opening a stub database connection"

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf(log_failedConn, err)
	}

	return db, mock
}

// will test get comment by modId empty uuid
func TestGetCommentByModIDEmptyUUID(t *testing.T) {
	// Arrange
	db, _ := NewMock()
	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		log.Fatalf(log_failedConn, err)
	}

	repo := repository.NewRepository(gdb)
	handler := New(repo, *logrus.New())

	// Act
	_, err = handler.GetCommentByModID(context.Background(), &protobuffer.GetCommentByModIDRequest{ModID: ""})

	// Assert
	assert.Error(t, err, "rpc error: code = Internal desc = Error request value ID, is not a valid UUID!")
}

// will test get comment by modId empty uuid
func TestGetCommentByModIDWrongUUID(t *testing.T) {
	// Arrange
	db, _ := NewMock()
	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		log.Fatalf(log_failedConn, err)
	}

	repo := repository.NewRepository(gdb)
	handler := New(repo, *logrus.New())

	// Act
	_, err = handler.GetCommentByModID(context.Background(), &protobuffer.GetCommentByModIDRequest{ModID: "I am not a uuid"})

	// Assert
	assert.Error(t, err, "rpc error: code = Internal desc = Error request value ID, is not a valid UUID!")
}

// will test get comment by modId empty uuid
func TestGetCommentByModID(t *testing.T) {
	// Arrange
	var modID = uuid.New()

	db, mock := NewMock()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "comments" WHERE mod_id = $1 AND "comments"."deleted_at" IS NULL`)).
		WithArgs(modID).
		WillReturnRows(sqlmock.
			NewRows([]string{"ID", "ModID", "UserID", "Text"}).
			AddRow(uuid.New().String(), modID.String(), uuid.New().String(), "Good Job!"))

	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		log.Fatalf(log_failedConn, err)
	}

	repo := repository.NewRepository(gdb)
	handler := New(repo, *logrus.New())

	// Act
	result, err := handler.GetCommentByModID(context.Background(), &protobuffer.GetCommentByModIDRequest{ModID: modID.String()})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, result.Comments[0].ModID, modID.String())
}

// will test update comment with wrong uuid
func TestUpdateCommentValidationUuidFailed(t *testing.T) {
	// Arrange
	request := &protobuffer.UpdateCommentRequest{
		ID:     "I am not a uuid",
		ModID:  uuid.NewString(),
		UserID: uuid.NewString(),
		Text:   "dLooks Nice",
	}

	db, _ := NewMock()

	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		log.Fatalf(log_failedConn, err)
	}

	repo := repository.NewRepository(gdb)
	handler := New(repo, *logrus.New())

	// Act
	_, err = handler.UpdateComment(context.Background(), request)

	// Assert
	assert.Equal(t, err.Error(), errors.New("Key: 'Comment.ID' Error:Field validation for 'ID' failed on the 'uuid4' tag").Error())
}

// will test update comment with max text
func TestUpdateCommentValidationTextExceedMaxFailed(t *testing.T) {
	// Arrange
	request := &protobuffer.UpdateCommentRequest{
		ID:     uuid.NewString(),
		ModID:  uuid.NewString(),
		UserID: uuid.NewString(),
		Text:   "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Donec quam felis, ultricies nec, pellentesque eu, pretium quis, sem. Nulla consequat massa quis enim. Donec pede justo, fringilla vel, aliquet nec, vulputate eget, arcu. In enim justo, rhoncus ut, imperdiet a, venenatis vitae, justo. Nullam dictum felis eu pede mollis pretium. Integer tincidunt. Cras dapibus. Vivamus elementum semper nisi. Aenean vulputate eleifend tellus. Aenean leo ligula, porttitor eu, consequat vitae, eleifend ac, enim. Aliquam lorem ante, dapibus in, viverra quis, feugiat a, tellus. Phasellus viverra nulla ut metus varius laoreet. Quisque rutrum. Aenean imperdiet. Etiam ultricies nisi vel augue. Curabitur ullamcorper ultricies nisi. Nam eget dui. Etiam rhoncus. Maecenas tempus, tellus eget condimentum rhoncus, sem quam semper libero, sit amet adipiscing sem neque sed ipsum. Nam quam nunc, blandit vel, luctus pulvinar, hendrerit id, lorem. Maecenas nec odio et ante tincidunt tempus. Donec vitae sapien ut libero venenatis faucibus. Nullam quis ante. Etiam sit amet orci eget eros faucibus tincidunt. Duis leo. Sed fringilla mauris sit amet nibh. Donec sodales sagittis magna. Sed consequat, leo eget bibendum sodales, augue velit cursus nunc, quis gravida magna mi a libero. Fusce vulputate eleifend sapien. Vestibulum purus quam, scelerisque ut, mollis sed, nonummy id, metus. Nullam accumsan lorem in dui. Cras ultricies mi eu turpis hendrerit fringilla. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; In ac dui quis mi consectetuer",
	}

	db, _ := NewMock()

	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		log.Fatalf(log_failedConn, err)
	}

	repo := repository.NewRepository(gdb)
	handler := New(repo, *logrus.New())

	// Act
	_, err = handler.UpdateComment(context.Background(), request)

	// Assert
	assert.Equal(t, err.Error(), errors.New("Key: 'Comment.Text' Error:Field validation for 'Text' failed on the 'max' tag").Error())
}

// will test update comment with min text
func TestUpdateCommentValidationTextExceedMinFailed(t *testing.T) {
	// Arrange
	request := &protobuffer.UpdateCommentRequest{
		ID:     uuid.NewString(),
		ModID:  uuid.NewString(),
		UserID: uuid.NewString(),
		Text:   "",
	}

	db, _ := NewMock()

	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		log.Fatalf(log_failedConn, err)
	}

	repo := repository.NewRepository(gdb)
	handler := New(repo, *logrus.New())

	// Act
	_, err = handler.UpdateComment(context.Background(), request)

	// Assert
	assert.Equal(t, err.Error(), errors.New("Key: 'Comment.Text' Error:Field validation for 'Text' failed on the 'min' tag").Error())
}

// will test update comment
func TestUpdateComment(t *testing.T) {
	// Arrange
	request := &protobuffer.UpdateCommentRequest{
		ID:     uuid.NewString(),
		ModID:  uuid.NewString(),
		UserID: uuid.NewString(),
		Text:   "dLooks Nice",
	}

	db, mock := NewMock()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "comments" SET "created_at"=$1,"updated_at"=$2,"deleted_at"=$3,"mod_id"=$4,"user_id"=$5,"text"=$6 WHERE "comments"."deleted_at" IS NULL AND "id" = $7`)).
		WithArgs(AnyTime{}, AnyTime{}, nil, request.ModID, request.UserID, request.Text, request.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		log.Fatalf(log_failedConn, err)
	}

	repo := repository.NewRepository(gdb)
	handler := New(repo, *logrus.New())

	// Act
	result, err := handler.UpdateComment(context.Background(), request)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

// will test delete comment
func TestDeleteComment(t *testing.T) {
	// Arrange
	commentID := uuid.New()

	db, mock := NewMock()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "comments" SET "deleted_at"=$1 WHERE "comments"."id" = $2 AND "comments"."deleted_at" IS NULL`)).
		WithArgs(AnyTime{}, commentID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		log.Fatalf(log_failedConn, err)
	}

	repo := repository.NewRepository(gdb)
	handler := New(repo, *logrus.New())

	// Act
	_, err = handler.DeleteComment(context.Background(), &protobuffer.DeleteCommentRequest{ID: commentID.String()})

	// Assert
	assert.NoError(t, err)
}

// will test create comment with wrong uuid
func TestCreateCommentValidationUuidFailed(t *testing.T) {
	// Arrange
	request := &protobuffer.CreateCommentRequest{
		ModID:  "I am not a uuid",
		UserID: uuid.NewString(),
		Text:   "dLooks Nice",
	}

	db, _ := NewMock()

	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		log.Fatalf(log_failedConn, err)
	}

	repo := repository.NewRepository(gdb)
	handler := New(repo, *logrus.New())

	// Act
	_, err = handler.CreateComment(context.Background(), request)

	// Assert
	assert.Equal(t, err.Error(), errors.New("Key: 'Comment.ModID' Error:Field validation for 'ModID' failed on the 'uuid4' tag").Error())
}

// will test update comment with max text
func TestCreateCommentValidationTextExceedMaxFailed(t *testing.T) {
	// Arrange
	request := &protobuffer.CreateCommentRequest{
		ModID:  uuid.NewString(),
		UserID: uuid.NewString(),
		Text:   "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Donec quam felis, ultricies nec, pellentesque eu, pretium quis, sem. Nulla consequat massa quis enim. Donec pede justo, fringilla vel, aliquet nec, vulputate eget, arcu. In enim justo, rhoncus ut, imperdiet a, venenatis vitae, justo. Nullam dictum felis eu pede mollis pretium. Integer tincidunt. Cras dapibus. Vivamus elementum semper nisi. Aenean vulputate eleifend tellus. Aenean leo ligula, porttitor eu, consequat vitae, eleifend ac, enim. Aliquam lorem ante, dapibus in, viverra quis, feugiat a, tellus. Phasellus viverra nulla ut metus varius laoreet. Quisque rutrum. Aenean imperdiet. Etiam ultricies nisi vel augue. Curabitur ullamcorper ultricies nisi. Nam eget dui. Etiam rhoncus. Maecenas tempus, tellus eget condimentum rhoncus, sem quam semper libero, sit amet adipiscing sem neque sed ipsum. Nam quam nunc, blandit vel, luctus pulvinar, hendrerit id, lorem. Maecenas nec odio et ante tincidunt tempus. Donec vitae sapien ut libero venenatis faucibus. Nullam quis ante. Etiam sit amet orci eget eros faucibus tincidunt. Duis leo. Sed fringilla mauris sit amet nibh. Donec sodales sagittis magna. Sed consequat, leo eget bibendum sodales, augue velit cursus nunc, quis gravida magna mi a libero. Fusce vulputate eleifend sapien. Vestibulum purus quam, scelerisque ut, mollis sed, nonummy id, metus. Nullam accumsan lorem in dui. Cras ultricies mi eu turpis hendrerit fringilla. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; In ac dui quis mi consectetuer",
	}

	db, _ := NewMock()

	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		log.Fatalf(log_failedConn, err)
	}

	repo := repository.NewRepository(gdb)
	handler := New(repo, *logrus.New())

	// Act
	_, err = handler.CreateComment(context.Background(), request)

	// Assert
	assert.Equal(t, err.Error(), errors.New("Key: 'Comment.Text' Error:Field validation for 'Text' failed on the 'max' tag").Error())
}

// will test update comment with min text
func TestCreateCommentValidationTextExceedMinFailed(t *testing.T) {
	// Arrange
	request := &protobuffer.CreateCommentRequest{
		ModID:  uuid.NewString(),
		UserID: uuid.NewString(),
		Text:   "",
	}

	db, _ := NewMock()

	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		log.Fatalf(log_failedConn, err)
	}

	repo := repository.NewRepository(gdb)
	handler := New(repo, *logrus.New())

	// Act
	_, err = handler.CreateComment(context.Background(), request)

	// Assert
	assert.Equal(t, err.Error(), errors.New("Key: 'Comment.Text' Error:Field validation for 'Text' failed on the 'min' tag").Error())
}

// will test update comment
func TestCreateComment(t *testing.T) {
	// Arrange
	newId := uuid.New()
	request := &protobuffer.CreateCommentRequest{
		ModID:  uuid.NewString(),
		UserID: uuid.NewString(),
		Text:   "dLooks Nice",
	}

	db, mock := NewMock()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "comments" ("created_at","updated_at","deleted_at","mod_id","user_id","text") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
		WithArgs(AnyTime{}, AnyTime{}, nil, request.ModID, request.UserID, request.Text).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newId))
	mock.ExpectCommit()

	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		log.Fatalf(log_failedConn, err)
	}

	repo := repository.NewRepository(gdb)
	handler := New(repo, *logrus.New())

	// Act
	result, err := handler.CreateComment(context.Background(), request)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, result.ID, newId.String())
}

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}
