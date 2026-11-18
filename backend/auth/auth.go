package auth

import (
	"errors"
	"fmt"
	"net/http"
	"olhar-biblico/config"
	"olhar-biblico/models"
	internalerrors "olhar-biblico/models/InternalErrors"
	"olhar-biblico/utils"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateToken(usuarioID, remoteAddress string) (string, error) {
	permissions := jwt.MapClaims{}
	permissions["authorized"] = true
	permissions["exp"] = time.Now().Add(time.Hour * 24).Unix()
	permissions["usuario_id"] = usuarioID
	permissions["remote_address"] = remoteAddress

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissions)
	return token.SignedString([]byte(config.SecretKey))
}

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := validateToken(r); err != nil {
			utils.HttpResponse(w, http.StatusUnauthorized, internalerrors.ErrValidatingToken)
			return
		}

		if authenticationData, err := ExtractAuthenticationData(r); err != nil {
			utils.HttpResponse(w, http.StatusUnauthorized, internalerrors.ErrExtractAuthDataToken)
			return
		} else {
			if authenticationData.RemoteAddress != utils.GetIPAddress(r) {
				utils.HttpResponse(w, http.StatusUnauthorized, internalerrors.DifferentIPAddress)
				return
			}
		}

		next(w, r)
	}
}

func validateToken(r *http.Request) error {
	tokenString := extractToken(r)
	token, err := jwt.Parse(tokenString, returnVerificationKey)
	if err != nil {
		return err
	}

	//só pra ver se os claims existem
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}

	return errors.New(internalerrors.InvalidToken)
}

func extractToken(r *http.Request) string {
	token := r.Header.Get("Authorization")

	if len(strings.Split(token, " ")) == 2 {
		return strings.Split(token, " ")[1]
	}

	return ""
}

func returnVerificationKey(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("método de assinatura inesperado! %v", token.Header["alg"])
	}
	return []byte(config.SecretKey), nil
}

func ExtractAuthenticationData(r *http.Request) (models.AuthenticationData, error) {
	tokenString := extractToken(r)
	token, err := jwt.Parse(tokenString, returnVerificationKey)
	if err != nil {
		return models.AuthenticationData{}, err
	}

	if permissions, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		authenticationData := models.AuthenticationData{}
		if authenticationData.Authorized, err = strconv.ParseBool(fmt.Sprint(permissions["authorized"])); err != nil {
			return models.AuthenticationData{}, err
		}

		if authenticationData.Exp, err = strconv.ParseInt(fmt.Sprintf("%.0f", permissions["exp"]), 10, 64); err != nil {
			return models.AuthenticationData{}, err
		}

		if authenticationData.UsuarioID = fmt.Sprint(permissions["usuario_id"]); authenticationData.UsuarioID == "" {
			return models.AuthenticationData{}, errors.New("ID de usuário inválido")
		}

		if authenticationData.RemoteAddress = fmt.Sprint(permissions["remote_address"]); authenticationData.RemoteAddress == "" {
			return models.AuthenticationData{}, errors.New("endereço de IP de usuário é inválido")
		}

		return authenticationData, nil
	}

	return models.AuthenticationData{}, err
}
