package service

import (
	"encoding/base64"
	"errors"
	"log"

	"github.com/regmarmcem/echo-session-demo/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) Signup(email, password string) (*model.User, error) {
	user := model.User{}
	s.db.Where("email = ?", email).First(&user)
	if user.ID != 0 {
		err := errors.New("user already exists")
		return nil, err
	}

	p, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		err = errors.New("failed to encrypt password")
		return nil, err
	}

	user = model.User{Email: email, Password: base64.StdEncoding.EncodeToString(p)}
	s.db.Create(&user)
	return &user, nil
}

func (s *UserService) Signin(email string, password string) (*model.User, error) {
	user := model.User{}
	s.db.Where("email = ?", email).First(&user)
	if user.ID == 0 {
		err := errors.New("user not found")
		return nil, err
	}

	p, err := base64.StdEncoding.DecodeString(user.Password)
	if err != nil {
		log.Println(err)
		err := errors.New("cannot decode password in database")
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(p, []byte(password))
	if err != nil {
		log.Println(err)
		err := errors.New("password mismatch")
		return nil, err
	}

	return &user, nil
}
