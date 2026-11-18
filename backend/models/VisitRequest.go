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

type VisitRequest struct {
	ID                    *primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name                  string              `bson:"name" json:"name,omitempty"`
	Phone                 string              `bson:"phone" json:"phone,omitempty"`
	Email                 string              `bson:"email" json:"email,omitempty"`
	Address               string              `bson:"address" json:"address,omitempty"`
	RemoteAddress         string              `bson:"remote_address" json:"remote_address,omitempty"`
	RemoteAddressFrontEnd string              `bson:"remote_address_front_end" json:"remote_address_front_end,omitempty"`
	CreatedAt             *time.Time          `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt             *time.Time          `bson:"updated_at" json:"updated_at,omitempty"`
	UpdatedByID           *primitive.ObjectID `bson:"updated_by_id" json:"updated_by_id,omitempty"`
	DeletedAt             *time.Time          `bson:"deleted_at" json:"deleted_at,omitempty"`
	DeletedByID           *primitive.ObjectID `bson:"deleted_by_id" json:"deleted_by_id,omitempty"`
}

func (r *VisitRequest) VerifyVisitRequest() error {
	if strings.TrimSpace(r.Name) == "" {
		return errors.New(internalerrors.InvalidName)
	}
	if strings.TrimSpace(r.Email) == "" {
		return errors.New(internalerrors.InvalidEmail)
	}
	if strings.TrimSpace(r.Phone) == "" {
		return errors.New(internalerrors.InvalidPhone)
	}
	if strings.TrimSpace(r.Address) == "" {
		return errors.New(internalerrors.InvalidAddress)
	}
	if err := checkmail.ValidateFormat(r.Email); err != nil {
		return errors.New(internalerrors.InvalidEmail)
	}
	if _, err := libphonenumber.Parse(r.Phone, "BR"); err != nil {
		return errors.New(internalerrors.InvalidPhone)
	}
	if isValid, err := utils.VerifyPhoneNumber(r.Phone); !isValid || err != nil {
		return errors.New(internalerrors.InvalidPhoneFormat)
	}

	return nil
}

func (r *VisitRequest) Format() {
	r.Name = strings.TrimSpace(r.Name)
	r.Email = strings.TrimSpace(r.Email)
	r.Phone = strings.TrimSpace(r.Phone)
	r.Address = strings.TrimSpace(r.Address)
}

func (u *VisitRequest) TableName() string {
	return "visit_requests"
}
