package stores

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/zett-8/go-clean-echo/db"
	"github.com/zett-8/go-clean-echo/models"
	"testing"
)

func TestAuthorStore_GetSuccessCase(t *testing.T) {
	mockDB, mock := db.Mock()
	defer mockDB.Close()

	authors := []models.Author{
		{ID: 1, Name: "test1", Country: "US"},
		{ID: 2, Name: "test2", Country: "UK"},
	}

	rows := mock.NewRows([]string{"id", "name", "country"})
	for _, a := range authors {
		rows.AddRow(a.ID, a.Name, a.Country)
	}

	mock.
		ExpectQuery("SELECT id, name, country from authors").
		WillReturnRows(rows)

	s := New(mockDB)

	r, err := s.AuthorStore.Get()

	assert.NoError(t, err)
	assert.Equal(t, authors, r)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthorStore_CreateSuccessCase(t *testing.T) {
	mockDB, mock := db.Mock()
	defer mockDB.Close()

	a := &models.Author{
		ID:      1,
		Name:    "test",
		Country: "US",
	}

	mock.NewRows([]string{"id", "name", "country"})
	pr := mock.ExpectPrepare("INSERT INTO authors (name, country) VALUES ($1, $2) RETURNING id")
	pr.ExpectQuery().WithArgs(a.Name, a.Country).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	s := New(mockDB)

	r, err := s.Create(a)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), r)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthorStore_UpdateByIdSuccessCase(t *testing.T) {
	mockDB, mock := db.Mock()
	defer mockDB.Close()

	a := &models.Author{
		ID:      1,
		Name:    "test",
		Country: "US",
	}

	mock.NewRows([]string{"id", "name", "country"}).AddRow(a.ID, a.Name, a.Country)

	a.Name = "new name"
	a.Country = "new country"

	pr := mock.ExpectPrepare("UPDATE authors SET name = $1, country = $2 WHERE authors.id = $3 RETURNING id")
	pr.ExpectQuery().WithArgs(a.Name, a.Country, a.ID).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(a.ID))

	s := New(mockDB)

	r, err := s.UpdateById(a)

	assert.NoError(t, err)
	assert.Equal(t, int64(a.ID), r)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthorStore_DeleteByIdSuccessCase(t *testing.T) {
	mockDB, mock := db.Mock()
	defer mockDB.Close()

	authors := []models.Author{
		{ID: 1, Name: "test1", Country: "US"},
		{ID: 2, Name: "test2", Country: "UK"},
	}

	deletingID := authors[0].ID

	rows := mock.NewRows([]string{"id", "name", "country"})
	for _, a := range authors {
		rows.AddRow(a.ID, a.Name, a.Country)
	}
	mock.
		ExpectExec("DELETE FROM authors WHERE authors.id = $1 RETURNING authors.id").
		WithArgs(deletingID).
		WillReturnResult(sqlmock.NewResult(int64(deletingID), 1))
	mock.
		ExpectExec("DELETE FROM authors WHERE authors.id = $1 RETURNING authors.id").
		WithArgs(deletingID).
		WillReturnResult(sqlmock.NewResult(int64(deletingID), 0))

	s := New(mockDB)

	assert.NoError(t, s.AuthorStore.DeleteById(deletingID))
	assert.Equal(t, s.AuthorStore.DeleteById(deletingID), sql.ErrNoRows)
	assert.NoError(t, mock.ExpectationsWereMet())
}
