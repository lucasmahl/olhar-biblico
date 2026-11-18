package models

import (
	"errors"
	internalerrors "olhar-biblico/models/InternalErrors"
	"olhar-biblico/utils"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/ttacon/libphonenumber"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          *primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name        string              `bson:"name" json:"name,omitempty"`
	Phone       string              `bson:"phone" json:"phone,omitempty"`
	Email       string              `bson:"email" json:"email,omitempty"`
	Password    string              `bson:"password" json:"password,omitempty"`
	CreatedAt   *time.Time          `bson:"created_at" json:"created_at,omitempty"`
	CreatedByID *primitive.ObjectID `bson:"created_by_id" json:"created_by_id,omitempty"`
	UpdatedAt   *time.Time          `bson:"updated_at" json:"updated_at,omitempty"`
	UpdatedByID *primitive.ObjectID `bson:"updated_by_id" json:"updated_by_id,omitempty"`
	DeletedAt   *time.Time          `bson:"deleted_at" json:"deleted_at,omitempty"`
	DeletedByID *primitive.ObjectID `bson:"deleted_by_id" json:"deleted_by_id,omitempty"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) VerifyUser() error {
	if strings.TrimSpace(u.Name) == "" {
		return errors.New(internalerrors.InvalidName)
	}
	if strings.TrimSpace(u.Email) == "" {
		return errors.New(internalerrors.InvalidEmail)
	}
	if len(strings.TrimSpace(u.Password)) < 12 {
		return errors.New(internalerrors.MinPassword)
	}
	if err := checkmail.ValidateFormat(u.Email); err != nil {
		return errors.New(internalerrors.InvalidEmail)
	}
	if _, err := libphonenumber.Parse(u.Phone, "BR"); err != nil {
		return errors.New(internalerrors.InvalidPhone)
	}
	if isValid, err := utils.VerifyPhoneNumber(u.Phone); !isValid || err != nil {
		return errors.New(internalerrors.InvalidPhoneFormat)
	}

	return nil
}

func (u *User) VerifyUserEmailAndPassword() error {
	if strings.TrimSpace(u.Email) == "" {
		return errors.New(internalerrors.InvalidEmail)
	}
	if strings.TrimSpace(u.Password) == "" {
		return errors.New(internalerrors.InvalidPassword)
	}
	if err := checkmail.ValidateFormat(u.Email); err != nil {
		return errors.New(internalerrors.InvalidEmail)
	}

	return nil
}

func (u *User) Format() {
	u.Name = strings.TrimSpace(u.Name)
	u.Email = strings.TrimSpace(u.Email)
	u.Phone = strings.TrimSpace(u.Phone)
	u.Password = strings.TrimSpace(u.Password)
}

func (u *User) HashPassword() error {
	hash, erro := utils.Hash(u.Password)
	if erro != nil {
		return erro
	}
	u.Password = string(hash)

	return nil
}
