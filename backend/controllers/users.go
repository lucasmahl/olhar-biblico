package controllers

import (
	"encoding/json"
	"net/http"
	"olhar-biblico/auth"
	"olhar-biblico/models"
	internalerrors "olhar-biblico/models/InternalErrors"
	"olhar-biblico/services"
	"olhar-biblico/utils"
	"reflect"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *userHandler {
	return &userHandler{
		service: service,
	}
}

func (handler *userHandler) Login(w http.ResponseWriter, r *http.Request) {
	decode := json.NewDecoder(r.Body)
	decode.DisallowUnknownFields()

	user := models.User{}
	if err := decode.Decode(&user); err != nil {
		utils.HttpResponse(w, http.StatusBadRequest, internalerrors.InvalidObject)
		return
	}

	if err := user.VerifyUserEmailAndPassword(); err != nil {
		utils.HttpResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	user.Format()

	userFinded, err := services.NewUserService().FindByEmail(user.Email)
	if err != nil {
		utils.HttpResponse(w, http.StatusBadRequest, internalerrors.ErrorSearchingForObject)
		return
	}

	if reflect.ValueOf(userFinded).IsZero() {
		utils.HttpResponse(w, http.StatusBadRequest, internalerrors.UserNotFound)
		return
	}

	if userFinded.DeletedAt != nil {
		utils.HttpResponse(w, http.StatusBadRequest, internalerrors.UserInactive)
		return
	}

	if err = utils.VerifyPassword(userFinded.Password, user.Password); err != nil {
		utils.HttpResponse(w, http.StatusBadRequest, internalerrors.EmailOrPasswordIncorrect)
		return
	}

	token, err := auth.CreateToken(userFinded.ID.Hex(), utils.GetIPAddress(r))
	if err != nil {
		utils.HttpResponse(w, http.StatusInternalServerError, internalerrors.ErrCreateToken)
		return
	}

	res := struct {
		ID    primitive.ObjectID
		Token string
	}{
		ID:    *userFinded.ID,
		Token: token,
	}

	utils.HttpResponse(w, http.StatusOK, res)
}

func (handler *userHandler) Create(w http.ResponseWriter, r *http.Request) {
	decode := json.NewDecoder(r.Body)
	decode.DisallowUnknownFields()

	user := models.User{}
	if err := decode.Decode(&user); err != nil {
		utils.HttpResponse(w, http.StatusBadRequest, internalerrors.InvalidObject)
		return
	}

	if err := user.VerifyUser(); err != nil {
		utils.HttpResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	user.Format()

	// não pode ter usuário com mesmo e-mail
	if userFinded, err := services.NewUserService().FindByEmail(user.Email); err != nil {
		utils.HttpResponse(w, http.StatusInternalServerError, internalerrors.ErrorSearchingForObject)
		return
	} else if !reflect.ValueOf(userFinded).IsZero() {
		utils.HttpResponse(w, http.StatusBadRequest, internalerrors.EmailExists)
		return
	}

	if err := user.HashPassword(); err != nil {
		utils.HttpResponse(w, http.StatusBadRequest, internalerrors.ErrFormatPassword)
		return
	}

	authenticationData, err := auth.ExtractAuthenticationData(r)
	if err != nil {
		utils.HttpResponse(w, http.StatusUnauthorized, internalerrors.ErrExtractAuthDataToken)
		return
	}

	if userID, err := primitive.ObjectIDFromHex(authenticationData.UsuarioID); err != nil {
		utils.HttpResponse(w, http.StatusInternalServerError, internalerrors.ErrConvertObjectIDToString)
		return
	} else {
		user.CreatedByID = &userID
	}

	user.ID = nil //para evitar erro caso usuário informe id
	id, err := services.NewUserService().Create(&user)
	if err != nil {
		utils.HttpResponse(w, http.StatusBadRequest, internalerrors.ErrOnRegisteringObject)
		return
	}

	utils.HttpResponse(w, http.StatusCreated, struct{ ID interface{} }{ID: id})
}

func (handler *userHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	decode := json.NewDecoder(r.Body)
	decode.DisallowUnknownFields()

	var user struct {
		Password           string `bson:"password" json:"password,omitempty"`
		NewPassword        string `bson:"new_password" json:"new_password,omitempty"`
		ConfirmNewPassword string `bson:"confirm_new_password" json:"confirm_new_password,omitempty"`
	}

	if err := decode.Decode(&user); err != nil {
		utils.HttpResponse(w, http.StatusBadRequest, internalerrors.InvalidObject)
		return
	}

	if user.NewPassword != user.ConfirmNewPassword {
		utils.HttpResponse(w, http.StatusBadRequest, internalerrors.ConfirmNewPassword)
		return
	} else if user.Password == user.NewPassword {
		utils.HttpResponse(w, http.StatusBadRequest, internalerrors.OldPasswordEqual)
		return
	}

	authenticationData, err := auth.ExtractAuthenticationData(r)
	if err != nil {
		utils.HttpResponse(w, http.StatusUnauthorized, internalerrors.ErrExtractAuthDataToken)
		return
	}

	userWithoutPassword, err := services.NewUserService().FindUserWithoutPasswordByID(authenticationData.UsuarioID)
	if err != nil {
		utils.HttpResponse(w, http.StatusBadRequest, internalerrors.UserNotFound)
		return
	}

	userFinded, err := services.NewUserService().FindByEmail(userWithoutPassword.Email)
	if err != nil {
		utils.HttpResponse(w, http.StatusBadRequest, internalerrors.UserNotFound)
		return
	}

	if reflect.ValueOf(userFinded).IsZero() {
		utils.HttpResponse(w, http.StatusBadRequest, internalerrors.UserNotFound)
		return
	}

	if authenticationData.UsuarioID != userFinded.ID.Hex() {
		utils.HttpResponse(w, http.StatusBadRequest, internalerrors.ChangePasswordNotAllowed)
		return
	}

	if userFinded.DeletedAt != nil {
		utils.HttpResponse(w, http.StatusBadRequest, internalerrors.UserInactive)
		return
	}

	if err = utils.VerifyPassword(userFinded.Password, user.Password); err != nil {
		utils.HttpResponse(w, http.StatusBadRequest, internalerrors.EmailOrPasswordIncorrect)
		return
	}

	userFinded.Password = user.NewPassword
	userFinded.Format()

	if err := userFinded.VerifyUser(); err != nil {
		utils.HttpResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := userFinded.HashPassword(); err != nil {
		utils.HttpResponse(w, http.StatusBadRequest, internalerrors.ErrFormatPassword)
		return
	}

	if _, err := services.NewUserService().UpdatePassword(&userFinded); err != nil {
		utils.HttpResponse(w, http.StatusBadRequest, internalerrors.ErrOnUpdatingObject)
		return
	}

	utils.HttpResponse(w, http.StatusOK, map[string]string{"message": "Senha alterada com sucesso!"})
}
