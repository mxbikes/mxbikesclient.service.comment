package repository

import (
	"database/sql"
	"database/sql/driver"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/mxbikes/mxbikesclient.service.comment/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func NewMockRepository(db gorm.ConnPool) *postgresRepository {
	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return NewRepository(gdb)
}

// will test get by mod id
func TestRepository_GetByModID(t *testing.T) {
	// Arrange
	var modID = uuid.New()

	db, mock := NewMock()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "comments" WHERE mod_id = $1 AND "comments"."deleted_at" IS NULL`)).
		WithArgs(modID).
		WillReturnRows(sqlmock.
			NewRows([]string{"ID", "ModID", "UserID", "Text"}).
			AddRow(uuid.New().String(), modID.String(), uuid.New().String(), "Good Job!"))

	repo := NewMockRepository(db)

	// Act
	l, err := repo.SearchByModID(modID.String())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, len(l), 1)
	assert.Equal(t, l[0].ModID, modID.String())
}

// will test insert comment
func TestRepository_Insert(t *testing.T) {
	// Arrange
	newId := uuid.New()
	comment := &models.Comment{
		ModID:  uuid.NewString(),
		UserID: uuid.NewString(),
		Text:   "Looks Nice",
	}

	db, mock := NewMock()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "comments" ("created_at","updated_at","deleted_at","mod_id","user_id","text") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
		WithArgs(AnyTime{}, AnyTime{}, nil, comment.ModID, comment.UserID, comment.Text).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newId))
	mock.ExpectCommit()

	repo := NewMockRepository(db)

	// Act
	err := repo.Save(comment)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, comment.ID, newId.String())
}

// will test update comment
func TestRepository_Update(t *testing.T) {
	// Arrange
	comment := &models.Comment{
		ID:     uuid.NewString(),
		ModID:  uuid.NewString(),
		UserID: uuid.NewString(),
		Text:   "dLooks Nice",
	}

	db, mock := NewMock()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "comments" SET "created_at"=$1,"updated_at"=$2,"deleted_at"=$3,"mod_id"=$4,"user_id"=$5,"text"=$6 WHERE "comments"."deleted_at" IS NULL AND "id" = $7`)).
		WithArgs(AnyTime{}, AnyTime{}, nil, comment.ModID, comment.UserID, comment.Text, comment.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	repo := NewMockRepository(db)

	// Act
	err := repo.Save(comment)

	// Assert
	assert.NoError(t, err)
}

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}
