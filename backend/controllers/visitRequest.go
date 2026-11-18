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

type visitRequestHandler struct {
	service services.VisitRequestService
}

func NewVisitRequestHandler(service services.VisitRequestService) *visitRequestHandler {
	return &visitRequestHandler{
		service: service,
	}
}

func (handler *visitRequestHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decode := json.NewDecoder(r.Body)
	decode.DisallowUnknownFields()

	visitRequest := models.VisitRequest{}
	if err := decode.Decode(&visitRequest); err != nil {
		utils.HttpResponse(w, http.StatusBadRequest, internalerrors.InvalidObject)
		return
	}

	visitRequest.RemoteAddress = utils.GetIPAddress(r)
	if err := visitRequest.VerifyVisitRequest(); err != nil {
		utils.HttpResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	visitRequest.Format()

	// é barrado quando for solicitado mesmo email ou phone
	visitRequestFound, err := services.NewVisitRequestService().FindByEmailOrPhone(visitRequest.Email, visitRequest.Phone)
	if err != nil {
		utils.HttpResponse(w, http.StatusBadRequest, internalerrors.ErrorSearchingForObject)
		return
	} else if !reflect.ValueOf(visitRequestFound).IsZero() {
		utils.HttpResponse(w, http.StatusBadRequest, internalerrors.VisitRequestEmailOrPhoneExists)
		return
	}

	visitRequest.RemoteAddress = utils.GetIPAddress(r)
	if id, err := services.NewVisitRequestService().Create(&visitRequest); err != nil {
		utils.HttpResponse(w, http.StatusBadRequest, internalerrors.ErrOnRegisteringObject)
		return
	} else {
		if objectID, valid := id.(primitive.ObjectID); valid {
			visitRequest.ID = &objectID
		}
	}

	if err := utils.SendEmail(visitRequest.Name, visitRequest.Email, visitRequest.Phone, visitRequest.Address, fmt.Sprintf("Pedido de visita: %s", visitRequest.Name)); err != nil {
		utils.HttpResponse(w, http.StatusInternalServerError, internalerrors.ErrOnSendEmail)
		return
	}

	utils.HttpResponse(w, http.StatusCreated, nil)
}

func (handler *visitRequestHandler) Find(w http.ResponseWriter, r *http.Request) {
	filter := utils.ReturnFilter(r)

	query := r.URL.Query()
	visitRequestFilter := models.VisitRequest{}

	visitRequestFilter.Address = query.Get("address")
	visitRequestFilter.Email = query.Get("email")
	visitRequestFilter.Name = query.Get("name")
	visitRequestFilter.Phone = query.Get("phone")

	result, err := services.NewVisitRequestService().Find(filter, visitRequestFilter)
	if err != nil {
		utils.HttpResponse(w, http.StatusBadRequest, internalerrors.ErrorSearchingForObject)
		return
	}

	utils.HttpResponse(w, http.StatusCreated, result)
}
