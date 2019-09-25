package controllers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"recibe_me/configs"
	"recibe_me/configs/constants"
	"recibe_me/internal/helpers"
	"recibe_me/internal/models"
	"recibe_me/internal/services"
	"recibe_me/pkg/crypto"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/mailjet/mailjet-apiv3-go"
	"gopkg.in/mgo.v2/bson"
)

// SignUp es la funcion de registro de usuarios
func SignUp(w http.ResponseWriter, r *http.Request) {
	userData := &models.User{}

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(userData); err != nil {
		helpers.Response(w, http.StatusInternalServerError, constants.ERR_DECODE, err, nil)
		return
	}

	// Se valida la informacion del usuario
	if validErrs := userData.Validate(); len(validErrs) > 0 {
		err := map[string]interface{}{"validationError": validErrs}
		helpers.Response(w, http.StatusBadRequest, constants.ERR_INVALID_DATA, err, nil)
		return
	}

	// Codificamos la contraseña
	encPass, err := crypto.EncodePassword([]byte(userData.Password), configs.SecurityCfg.PasswordEncKey)
	if err != nil {
		helpers.Response(w, http.StatusInternalServerError, constants.ERR_DECODE, err, nil)
		return
	}

	userData.Password = base64.StdEncoding.EncodeToString([]byte(encPass))
	userData.Email = strings.ToLower(userData.Email)
	userData.VerificationCode = helpers.EncodeToString(5)
	userData.CreateAt = int32(time.Now().Unix())
	userData.Verified = false

	// Verificamos que el correo no exista
	result, err := helpers.UsersCollection.Find(bson.M{"email": userData.Email}).Count()
	if result > 0 {
		helpers.Response(w, http.StatusBadRequest, constants.ERR_USER_EXIST, err, nil)
		return
	}

	// Insertamos el usuario en la base de datos
	err = helpers.UsersCollection.Insert(userData)
	if err != nil {
		helpers.Response(w, http.StatusInternalServerError, constants.ERR_INSERT_DATA, err, nil)
		return
	}

	// Enviamos por correo el codigo de verificacion
	resp, err := SendVerificationCode(*userData)
	if err != nil {
		helpers.Response(w, http.StatusInternalServerError, constants.ERR_SEND_EMAIL, err, resp)
		return
	}

	helpers.Response(w, http.StatusOK, constants.SUCCESS, err, nil)
}

// Login es la funcion para autenticacion
func Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var userAuthData models.UserAuth

	err := decoder.Decode(&userAuthData)
	if err != nil {
		helpers.Response(w, http.StatusInternalServerError, constants.ERR_DECODE, err, nil)
		return
	}

	userAuthData.Email = strings.ToLower(userAuthData.Email)

	user := models.User{}

	// Buscamos el usuario en la base de datos
	err = helpers.UsersCollection.Find(bson.M{"email": userAuthData.Email}).One(&user)
	if err != nil {
		helpers.Response(w, http.StatusUnauthorized, constants.ERR_USER_NOT_FOUND, err, nil)
		return
	}

	// Decodificamos la contraseña
	decPass, err := base64.StdEncoding.DecodeString(user.Password)
	if err != nil {
		helpers.Response(w, http.StatusBadRequest, constants.ERR_DECODE, err, nil)
		return
	}

	// Chequeamos la contraseña
	isValid, err := crypto.CheckPassword([]byte(userAuthData.Password), []byte(decPass), configs.SecurityCfg.PasswordEncKey)
	if !isValid {
		helpers.Response(w, http.StatusUnauthorized, constants.ERR_PASS_NOT_VALID, err, nil)
		return
	}

	user.Password = ""

	// Generamos el token de autenticacion
	token, err := crypto.CreateTokenString(&crypto.Claims{
		Type: constants.User,
		StandardClaims: jwt.StandardClaims{
			Id: user.ID.Hex(),
		},
	}, configs.SecurityCfg.TokenSecret, configs.SecurityCfg.TokenDuration)

	helpers.Response(w, http.StatusOK, constants.SUCCESS, err, models.UserLogResponse{Token: token, User: user})
}

// ResendVerificationCode esta funcion se encarga de reenviar un codigo de verificacion de cuentas de usuario
func ResendVerificationCode(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if !bson.IsObjectIdHex(params["userId"]) {
		helpers.Response(w, http.StatusNotFound, constants.ERR_USER_NOT_FOUND, nil, nil)
		return
	}

	userID := bson.ObjectIdHex(params["userId"])

	userData := models.User{}

	// Se busca la informacion del usuario
	err := helpers.UsersCollection.FindId(userID).One(&userData)
	if err != nil {
		helpers.Response(w, http.StatusBadRequest, constants.ERR_USER_NOT_FOUND, err, nil)
		return
	}

	// Se actualiza el codigo de verificación en la base de datos
	var verificationCode = helpers.EncodeToString(5)
	err = helpers.UsersCollection.UpdateId(userID, bson.M{"$set": bson.M{"verification_code": verificationCode}})
	if err != nil {
		helpers.Response(w, http.StatusInternalServerError, constants.ERR_UPDATE_DATA, err, nil)
		return
	}

	// Se reenvia el codigo al correo del usuario
	userData.VerificationCode = verificationCode
	resp, err := SendVerificationCode(userData)
	if err != nil {
		helpers.Response(w, http.StatusInternalServerError, constants.ERR_SEND_EMAIL, err, resp)
		return
	}

	helpers.Response(w, http.StatusOK, constants.SUCCESS, nil, nil)
}

// VerificationAccount esta funcion recibe el codigo de verificacion y el id del usuario y hace la validacion
func VerificationAccount(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if !bson.IsObjectIdHex(params["userId"]) {
		helpers.Response(w, http.StatusNotFound, constants.ERR_USER_NOT_FOUND, nil, nil)
		return
	}
	userID := bson.ObjectIdHex(params["userId"])

	code := params["verificationCode"]

	userData := models.User{}

	// Se busca la informacion del usuario
	err := helpers.UsersCollection.FindId(userID).One(&userData)
	if err != nil {
		helpers.Response(w, http.StatusBadRequest, constants.ERR_USER_NOT_FOUND, err, nil)
		return
	}

	// Se valida el codigo ingresado
	if userData.VerificationCode != code {
		helpers.Response(w, http.StatusBadRequest, constants.ERR_INVALID_CODE, nil, nil)
		return
	}

	// Se actualiza el codigo de verificación en la base de datos
	err = helpers.UsersCollection.UpdateId(userID, bson.M{"$set": bson.M{"verified": true}})
	if err != nil {
		helpers.Response(w, http.StatusInternalServerError, constants.ERR_UPDATE_DATA, err, nil)
		return
	}

	helpers.Response(w, http.StatusOK, constants.SUCCESS, nil, nil)
}

// PasswordRecovery Recibe el UserID y el codigo de recuperacion de contraseña y hace la validacion
func PasswordRecovery(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if !bson.IsObjectIdHex(params["userId"]) {
		helpers.Response(w, http.StatusNotFound, constants.ERR_USER_NOT_FOUND, nil, nil)
		return
	}

	userID := bson.ObjectIdHex(params["userId"])

	code := params["passwordRecoveryCode"]

	userData := models.User{}

	// Se busca la informacion del usuario
	err := helpers.UsersCollection.FindId(userID).One(&userData)
	if err != nil {
		helpers.Response(w, http.StatusBadRequest, constants.ERR_USER_NOT_FOUND, err, nil)
		return
	}

	// Se valida el codigo ingresado
	if userData.VerificationCode != code {
		helpers.Response(w, http.StatusBadRequest, constants.ERR_INVALID_CODE, nil, nil)
		return
	}

	newPassword := helpers.EncodeToString(8)

	// Codificamos la contraseña
	encPass, err := crypto.EncodePassword([]byte(newPassword), configs.SecurityCfg.PasswordEncKey)
	if err != nil {
		helpers.Response(w, http.StatusInternalServerError, constants.ERR_DECODE, err, nil)
		return
	}

	err = helpers.UsersCollection.UpdateId(userID, bson.M{"$set": bson.M{"password": base64.StdEncoding.EncodeToString([]byte(encPass))}})
	if err != nil {
		helpers.Response(w, http.StatusInternalServerError, constants.ERR_UPDATE_DATA, err, nil)
		return
	}

	// Se prepara la informacion para hacer el envio
	var to = services.Recipient{Name: userData.FirstName, Email: userData.Email}
	var from = services.Recipient{Name: "Cristobal Muñoz", Email: "cmunoz21x@gmail.com"}

	var info = services.Info{
		ApiKeyPrivate: configs.SecurityCfg.MjApiKeyPrivate,
		ApiKeyPublic:  configs.SecurityCfg.MjApiKeyPublic,
		FromRecipient: from,
		ToRecipient:   to,
		Code:          "",
		Password: newPassword,
	}

	// se envia el codigo
	_, err = info.SendNewPassword()
	if err != nil {
		helpers.Response(w, http.StatusInternalServerError, constants.ERR_SEND_EMAIL, err, nil)
		return
	}

	helpers.Response(w, http.StatusOK, constants.SUCCESS, nil, nil)
}

// SendVerificationCode Envio el codigo de verificacion al usuario
func SendVerificationCode(userData models.User) (*mailjet.ResultsV31, error) {
	var to = services.Recipient{Name: userData.FirstName, Email: userData.Email}
	var from = services.Recipient{Name: "Cristobal Muñoz", Email: "cmunoz21x@gmail.com"}

	var info = services.Info{
		ApiKeyPrivate: configs.SecurityCfg.MjApiKeyPrivate,
		ApiKeyPublic:  configs.SecurityCfg.MjApiKeyPublic,
		FromRecipient: from,
		ToRecipient:   to,
		Code:          userData.VerificationCode,
	}

	resp, err := info.SendVerificationCode()
	if err != nil {
		return nil, err
	}

	return resp, err
}

// SendPasswordRecoveryCode envia el codigo de recuperacion de contraseña al usuario
func SendPasswordRecoveryCode(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userData := models.User{}

	// Se busca la informacion del usuario
	if err := helpers.UsersCollection.Find(bson.M{"email": params["email"]}).One(&userData); err != nil {
		helpers.Response(w, http.StatusBadRequest, constants.ERR_USER_NOT_FOUND, err, nil)
		return
	}

	// Se actualiza el codigo de verificación en la base de datos
	var verificationCode = helpers.EncodeToString(5)

	if err := helpers.UsersCollection.UpdateId(userData.ID, bson.M{"$set": bson.M{"verification_code": verificationCode}}); err != nil {
		helpers.Response(w, http.StatusInternalServerError, constants.ERR_UPDATE_DATA, err, nil)
		return
	}

	// Se prepara la informacion para hacer el envio
	var to = services.Recipient{Name: userData.FirstName, Email: userData.Email}
	var from = services.Recipient{Name: "Cristobal Muñoz", Email: "cmunoz21x@gmail.com"}

	var info = services.Info{
		ApiKeyPrivate: configs.SecurityCfg.MjApiKeyPrivate,
		ApiKeyPublic:  configs.SecurityCfg.MjApiKeyPublic,
		FromRecipient: from,
		ToRecipient:   to,
		Code:          verificationCode,
	}

	// se envia el codigo
	_, err := info.SendPasswordRecoveyCode()
	if err != nil {
		helpers.Response(w, http.StatusInternalServerError, constants.ERR_INVALID_CODE, err, nil)
		return
	}

	helpers.Response(w, http.StatusOK, constants.SUCCESS, nil, userData)
}

// PasswordUpdate update user password's
func PasswordUpdate(responseWriter http.ResponseWriter, request *http.Request) {
	// Se obtiene el usuario autenticado
	user, err := helpers.GetUserFromRequest(request)
	var result map[string]interface{}

	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&result)
	if err != nil {
		helpers.Response(responseWriter, http.StatusInternalServerError, constants.ERR_DECODE, err, nil)
		return
	}

	// Se verifica que se han enviado los campos requeridos
	if result["old_password"] == nil || strings.TrimSpace(result["old_password"].(string)) == "" {
		helpers.Response(responseWriter, http.StatusUnprocessableEntity, constants.ERR_INVALID_DATA, "El campo 'old_password' es obligatorio.", nil)
		return
	}

	// Se verifica que se han enviado los campos requeridos
	if result["new_password"] == nil || strings.TrimSpace(result["new_password"].(string)) == "" {
		helpers.Response(responseWriter, http.StatusUnprocessableEntity, constants.ERR_INVALID_DATA, "El campo 'new_password' es obligatorio.", nil)
		return
	}

	// Se verifica que se han enviado los campos requeridos
	if result["new_password_confirmation"] == nil || strings.TrimSpace(result["new_password_confirmation"].(string)) == "" {
		helpers.Response(responseWriter, http.StatusUnprocessableEntity, constants.ERR_INVALID_DATA, "El campo 'new_password_confirmation' es obligatorio.", nil)
		return
	}

	oldPassword := result["old_password"].(string)
	newPassword := result["new_password"].(string)
	newPasswordConfirmation := result["new_password_confirmation"].(string)

	decodedPassword, _ := base64.StdEncoding.DecodeString(user.Password)
	// Se valida la contraseña
	isValid, err := crypto.CheckPassword([]byte(oldPassword), []byte(decodedPassword), configs.SecurityCfg.PasswordEncKey)
	if !isValid || err != nil {
		helpers.Response(responseWriter, http.StatusUnauthorized, constants.ERR_PASS_NOT_VALID, err, nil)
		return
	}

	// Se verifica que la nueva contraseña coincida
	if newPassword != newPasswordConfirmation {
		helpers.Response(responseWriter, http.StatusUnprocessableEntity, constants.ERR_INVALID_DATA, "Error al confirmar la nueva contraseña: Las contraseñas no coinciden.", nil)
		return
	}

	// Se valida el tamaño de la nueva contraseña
	if len(newPassword) < 6 || len(newPassword) > 16 {
		helpers.Response(responseWriter, http.StatusUnprocessableEntity, constants.ERR_INVALID_DATA, "El tamaño debe estar entre 6 y 16 caracteres", nil)
		return
	}

	// Se encripta la nueva contraseña
	encPass, err := crypto.EncodePassword([]byte(newPassword), configs.SecurityCfg.PasswordEncKey)
	if err != nil {
		helpers.Response(responseWriter, http.StatusInternalServerError, constants.ERR_DECODE, err, nil)
		return
	}

	// Se actualiza la contraseña del usuario
	err = helpers.UsersCollection.UpdateId(user.ID, bson.M{"$set": bson.M{"password": base64.StdEncoding.EncodeToString([]byte(encPass))}})

	if err != nil {
		helpers.Response(responseWriter, http.StatusInternalServerError, constants.ERR_INTERNAL_ERROR, err, nil)
		return
	}

	// Se limpia el campo password
	user.Password = ""

	// Se envía la respuesta al usuario
	helpers.Response(responseWriter, http.StatusOK, constants.SUCCESS, nil, user)
}
