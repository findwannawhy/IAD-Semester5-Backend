package repository

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/findwannawhy/IAD-Semester5/internal/app/ds"
	"github.com/findwannawhy/IAD-Semester5/internal/app/dto"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (r *Repository) GetUserByID(id uuid.UUID) (ds.User, error) {
	user := ds.User{}
	sub := r.db.Where("id = ?", id).Find(&user)
	if sub.Error != nil {
		return ds.User{}, sub.Error
	}
	if sub.RowsAffected == 0 {
		return ds.User{}, ErrNotFound
	}
	err := sub.First(&user).Error
	if err != nil {
		return ds.User{}, err
	}
	return user, nil
}

func (r *Repository) GetUserByLogin(login string) (ds.User, error) {
	user := ds.User{}
	sub := r.db.Where("login = ?", login).Find(&user)
	if sub.Error != nil {
		return ds.User{}, sub.Error
	}
	if sub.RowsAffected == 0 {
		return ds.User{}, ErrNotFound
	}
	err := sub.First(&user).Error
	if err != nil {
		return ds.User{}, err
	}
	return user, nil
}

func (r *Repository) CreateUser(user dto.UserRequest) (dto.UserResponse, error) {
	if user.Login == "" {
		return dto.UserResponse{}, errors.New("login is empty")
	}
	if user.Password == "" {
		return dto.UserResponse{}, errors.New("password is empty")
	}
	if _, err := r.GetUserByLogin(user.Login); err == nil {
		return dto.UserResponse{}, errors.New("user already exists")
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return dto.UserResponse{}, err
	}
	user.Password = hashedPassword

	userDS := ds.User{
		ID: uuid.New(),
		Login: user.Login,
		Password: hashedPassword,
		IsModerator: false,
	}

	sub := r.db.Create(&userDS)
	if sub.Error != nil {
		return dto.UserResponse{}, sub.Error
	}
	return dto.UserResponse{
		ID: userDS.ID,
		Login: userDS.Login,
		IsModerator: userDS.IsModerator,
	}, nil
}

func (r *Repository) UpdateProfile(login string, userJSON ds.User) (dto.UserResponse, error) {
	currUser, err := r.GetUserByLogin(login)
	if err != nil {
		return dto.UserResponse{}, err
	}

	if userJSON.Login != "" {
		currUser.Login = userJSON.Login
	}

	if userJSON.Password != "" {
		hashedPassword, err := HashPassword(userJSON.Password)
		if err != nil {
			return dto.UserResponse{}, err
		}
		currUser.Password = hashedPassword
	}

	if userJSON.IsModerator && !currUser.IsModerator {
		userJSON.IsModerator = false
	}
	currUser.IsModerator = userJSON.IsModerator

	err = r.db.Save(&currUser).Error
	if err != nil {
		return dto.UserResponse{}, err
	}
	return dto.UserResponse{
		ID: currUser.ID,
		Login: currUser.Login,
		IsModerator: currUser.IsModerator,
	}, nil
}

func (r *Repository) SignIn(userJSON dto.UserRequest) (string, error) {
	user, err := r.GetUserByLogin(userJSON.Login)
	if err != nil {
		return "", err
	}

	if !CheckPasswordHash(userJSON.Password, user.Password) {
		return "", errors.New("invalid password")
	}

	token, err := GenerateToken(user.ID, user.IsModerator)
	if err != nil {
		return "", err
	}

	return token, nil
}

func GenerateToken(id uuid.UUID, isModerator bool) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user_id"] = id.String()
	claims["is_moderator"] = isModerator
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	fmt.Println(password)
	
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}