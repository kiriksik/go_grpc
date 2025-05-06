package sqlite

import (
	"context"
	"database/sql"
	"errors"

	// "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/kiriksik/go_grpc/internal/domain/models"
	"github.com/kiriksik/go_grpc/internal/services/storage"
	sqlite3 "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	stmt, err := s.db.Prepare("INSERT INTO users(email, pass_hash) VALUES(?, ?)")
	if err != nil {
		return 0, err
	}

	res, err := stmt.ExecContext(ctx, email, passHash)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, storage.ErrUserExists
		}
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Storage) GetUser(ctx context.Context, email string) (models.User, error) {
	stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email = ?")
	if err != nil {
		return models.User{}, err
	}

	row := stmt.QueryRowContext(ctx, email)

	var user models.User
	err = row.Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, storage.ErrUserNotFound
		}
		return models.User{}, err
	}
	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	stmt, err := s.db.Prepare("SELECT is_admin FROM users WHERE id = ?")
	if err != nil {
		return false, err
	}
	row := stmt.QueryRowContext(ctx, userID)

	var isAdmin bool
	err = row.Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, storage.ErrUserNotFound
		}
		return false, err
	}
	return isAdmin, nil
}

func (s *Storage) App(ctx context.Context, appID int) (models.App, error) {
	stmt, err := s.db.Prepare("SELECT id, name, secret FROM apps WHERE id = ?")
	if err != nil {
		return models.App{}, err
	}
	row := stmt.QueryRowContext(ctx, appID)

	var app models.App
	err = row.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, storage.ErrAppNotFound
		}
		return models.App{}, err
	}
	return app, nil
}
