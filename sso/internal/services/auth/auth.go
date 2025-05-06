package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/kiriksik/go_grpc/internal/domain/models"
	"github.com/kiriksik/go_grpc/internal/lib/jwt"
	"github.com/kiriksik/go_grpc/internal/services/storage"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID       = errors.New("invalid app id")
	ErrUserExists         = errors.New("user exists")
)

type Auth struct {
	log      *slog.Logger
	storage  Storage
	tokenTTL time.Duration
}

type Storage interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error)
	GetUser(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
	App(ctx context.Context, appID int) (models.App, error)
}

func New(log *slog.Logger, storage Storage, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:      log,
		storage:  storage,
		tokenTTL: tokenTTL,
	}
}

// func (a *Auth) SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error) {
// 	return 0, nil
// }

// func (a *Auth) GetUser(ctx context.Context, email string) (models.User, error) {
// 	return models.User{}, nil
// }

// func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
// 	return false, nil
// }

// func (a *Auth) App(ctx context.Context, appID int) (models.App, error) {
// 	return models.App{}, nil
// }

func (a *Auth) Login(ctx context.Context, email string, password string, appID int) (string, error) {
	const op = "authService.Login"
	log := a.log.With(slog.String("op", op), slog.String("email", email))

	log.Info("authorization user")

	user, err := a.storage.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found" + err.Error())
			return "", ErrInvalidCredentials
		}
		a.log.Error("failed to get user" + err.Error())
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials" + err.Error())
		return "", ErrInvalidCredentials
	}

	app, err := a.storage.App(ctx, appID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			a.log.Warn("app not found" + err.Error())
			return "", ErrInvalidAppID
		}
		return "", err
	}

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("invalid token" + err.Error())
		return "", err
	}
	return token, nil
}

func (a *Auth) RegisterNewUser(ctx context.Context, email string, password string) (int64, error) {
	const op = "authService.RegisterNewUser"
	log := a.log.With(slog.String("op", op), slog.String("email", email))
	log.Info("register new user")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("error while hashing password" + err.Error())
		return 0, err
	}

	uid, err := a.storage.SaveUser(ctx, email, hashedPassword)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			a.log.Warn("user already exists" + err.Error())
			return 0, ErrUserExists
		}
		log.Error("error while saving user" + err.Error())
		return 0, err
	}
	return uid, nil
}
func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "authService.IsAdmin"
	log := a.log.With(slog.String("op", op), slog.Int64("userID", userID))
	log.Info("checking user for admin")

	isAdmin, err := a.storage.IsAdmin(ctx, int64(userID))
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found" + err.Error())
			return false, ErrInvalidCredentials
		}
		log.Error("error while checking user for admin" + err.Error())
		return false, err
	}
	return isAdmin, nil
}
