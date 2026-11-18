package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"olhar-biblico/models"
	internalerrors "olhar-biblico/models/InternalErrors"
	"olhar-biblico/services"
	"olhar-biblico/utils"
	"reflect"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type bookRequestHandler struct {
	service services.BookRequestService
}

func NewBookRequestHandler(service services.BookRequestService) *bookRequestHandler {
	return &bookRequestHandler{
		service: service,
	}
}

func (handler *bookRequestHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decode := json.NewDecoder(r.Body)
	decode.DisallowUnknownFields()

	bookRequest := models.BookRequest{}
	if err := decode.Decode(&bookRequest); err != nil {
		utils.HttpResponse(w, http.StatusBadRequest, internalerrors.InvalidObject)
		return
	}

	bookRequest.RemoteAddress = utils.GetIPAddress(r)
	if err := bookRequest.VerifyBookRequest(); err != nil {
		utils.HttpResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	bookRequest.Format()

	// é barrado quando for solicitado mesmo email ou phone
	bookRequestFound, err := services.NewBookRequestService().FindByEmailOrPhone(bookRequest.Email, bookRequest.Phone)
	if err != nil {
		utils.HttpResponse(w, http.StatusBadRequest, internalerrors.ErrorSearchingForObject)
		return
	} else if !reflect.ValueOf(bookRequestFound).IsZero() {
		utils.HttpResponse(w, http.StatusBadRequest, internalerrors.BookRequestEmailOrPhoneExists)
		return
	}

	bookRequest.RemoteAddress = utils.GetIPAddress(r)
	if id, err := services.NewBookRequestService().Create(&bookRequest); err != nil {
		utils.HttpResponse(w, http.StatusBadRequest, internalerrors.ErrOnRegisteringObject)
		return
	} else {
		if objectID, valid := id.(primitive.ObjectID); valid {
			bookRequest.ID = &objectID
		}
	}

	if err := utils.SendEmail(bookRequest.Name, bookRequest.Email, bookRequest.Phone, bookRequest.Address, fmt.Sprintf("Pedido de livro, O Grande Conflito: %s", bookRequest.Name)); err != nil {
		utils.HttpResponse(w, http.StatusInternalServerError, internalerrors.ErrOnSendEmail)
		return
	}

	utils.HttpResponse(w, http.StatusCreated, nil)
}

func (handler *bookRequestHandler) Find(w http.ResponseWriter, r *http.Request) {
	filter := utils.ReturnFilter(r)

	query := r.URL.Query()
	bookRequestFilter := models.BookRequest{}

	bookRequestFilter.Address = query.Get("address")
	bookRequestFilter.Email = query.Get("email")
	bookRequestFilter.Name = query.Get("name")
	bookRequestFilter.Phone = query.Get("phone")

	result, err := services.NewBookRequestService().Find(filter, bookRequestFilter)
	if err != nil {
		utils.HttpResponse(w, http.StatusBadRequest, internalerrors.ErrorSearchingForObject)
		return
	}

	utils.HttpResponse(w, http.StatusCreated, result)
}
