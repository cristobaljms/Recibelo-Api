package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"recibe_me/configs/constants"
	"recibe_me/internal/helpers"
	"recibe_me/internal/models"
	"strconv"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

// GetDelivery returns a Delivery
func GetDelivery(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deliveryID := params["id"]

	if !bson.IsObjectIdHex(deliveryID) {
		helpers.Response(w, http.StatusNotFound, constants.ERR_USER_NOT_FOUND, nil, nil)
		return
	}

	var result models.Delivery

	// Se consulta el envio
	err := helpers.DeliveriesCollection.Find(deliveryID).One(&result)
	if err != nil {
		helpers.Response(w, http.StatusNotFound, constants.ERR_USER_NOT_FOUND, err, nil)
		return
	}
	helpers.Response(w, http.StatusOK, constants.SUCCESS, err, result)
}

// ListDeliveries returns a List of Deliveries by user id
func ListDeliveries(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetUserFromRequest(r)

	var results []models.Delivery

	// Se consultas los envios por usuario
	err = helpers.DeliveriesCollection.Find(bson.M{"user_id": user.ID}).Sort("-_id").All(&results)
	if err != nil {
		helpers.Response(w, http.StatusInternalServerError, constants.ERR_USER_NOT_FOUND, err, nil)
		return
	}

	helpers.Response(w, http.StatusOK, constants.SUCCESS, err, results)
}

// Rate a Delivery Service
func Rate(responseWriter http.ResponseWriter, request *http.Request) {

	var result map[string]interface{}

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&result)
	if err != nil {
		helpers.Response(responseWriter, http.StatusInternalServerError, constants.ERR_DECODE, err, nil)
		return
	}

	if result["rating"] == nil {
		helpers.Response(responseWriter, http.StatusUnprocessableEntity, constants.ERR_INVALID_DATA, "El campo 'rating' es obligatorio.", nil)
		return
	}

	rating, err := strconv.Atoi(fmt.Sprintf("%v", result["rating"]))

	if err != nil {
		helpers.Response(responseWriter, http.StatusUnprocessableEntity, constants.ERR_INVALID_DATA, "Calificación Inválida: Debe ser un número entero.", nil)
		return
	}

	if rating < 1 || rating > 5 {
		helpers.Response(responseWriter, http.StatusUnprocessableEntity, constants.ERR_INVALID_DATA, "Calificación Inválida: Debe ser un número entero entre 1 y 5.", nil)
		return
	}

	params := mux.Vars(request)
	deliveryID := params["id"]

	if !bson.IsObjectIdHex(deliveryID) {
		helpers.Response(responseWriter, http.StatusNotFound, constants.ERR_INVALID_DATA, "El ID es inválido.", nil)
		return
	}

	oid := bson.ObjectIdHex(deliveryID)

	delivery := models.Delivery{}

	err = helpers.DeliveriesCollection.FindId(oid).One(&delivery)

	if err != nil {
		helpers.Response(responseWriter, http.StatusNotFound, constants.ERR_NOT_FOUND, err, nil)
		return
	}

	if delivery.Rated {
		helpers.Response(responseWriter, http.StatusConflict, constants.ERR_UPDATE_DATA, "Este envío ya ha sido calificado.", nil)
		return
	}

	err = helpers.DeliveriesCollection.UpdateId(oid, bson.M{"$set": bson.M{"rating": rating, "rated": true}})

	if err != nil {
		helpers.Response(responseWriter, http.StatusInternalServerError, constants.ERR_INTERNAL_ERROR, err, nil)
		return
	}

	delivery.Rating = int64(rating)
	delivery.Rated = true

	helpers.Response(responseWriter, http.StatusOK, constants.SUCCESS, nil, delivery)
}

// AddDelivery ..
func AddDelivery(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var deliveriesData = models.Delivery{}

	err := decoder.Decode(&deliveriesData)

	if err != nil {
		helpers.Response(w, http.StatusInternalServerError, constants.ERR_DECODE, err, nil)
		return
	}
	defer r.Body.Close()

	resultuser := 1

	result, err := helpers.DeliveriesCollection.Find(bson.M{"tracker_id": deliveriesData.TrackerID}).Count()
	if resultuser == 0 {
		//helpers.Response(writer, http.StatusBadRequest, constants.ERR_USER_NOT_FOUND,err1,nil)
	} else if result > 0 {

		helpers.Response(w, http.StatusBadRequest, constants.ERR_DELIVERY_EXIST, err, nil)
	} else {
		err = helpers.DeliveriesCollection.Insert(deliveriesData)

		if err != nil {
			helpers.Response(w, http.StatusInternalServerError, constants.ERR_INSERT_DATA, err, nil)
		}
		helpers.Response(w, http.StatusOK, constants.SUCCESS, err, nil)
	}

}

// RegisterDelivery to User
func RegisterDelivery(writer http.ResponseWriter, read *http.Request) {

	decoder := json.NewDecoder(read.Body)
	var teamdelivery = models.TeamDelivery{}
	var deliveriesData = models.Delivery{}

	err := decoder.Decode(&teamdelivery)
	defer read.Body.Close()
	if err != nil {
		helpers.Response(writer, http.StatusInternalServerError, constants.ERR_DECODE, err, nil)
		return
	} else if teamdelivery.UserID == "" {
		helpers.Response(writer, http.StatusUnprocessableEntity, constants.ERR_INVALID_DATA, "Falta Id del Usuario", nil)
		return
	} else if teamdelivery.TrackerID == "" {
		helpers.Response(writer, http.StatusUnprocessableEntity, constants.ERR_INVALID_DATA, "Falta ID del Tracker ", nil)
	} else {

		resultuser, err1 := helpers.UsersCollection.Find(bson.M{"_id": teamdelivery.UserID}).Count()
		result := helpers.DeliveriesCollection.Find(bson.M{"tracker_id": teamdelivery.TrackerID}).One(&deliveriesData)

		if resultuser == 0 {
			helpers.Response(writer, http.StatusBadRequest, "Usuario no Existe", err1, nil)
		} else if result == nil {
			result := helpers.DeliveriesCollection.Update(bson.M{"_id": deliveriesData.ID}, bson.M{"$set": bson.M{"user_id": teamdelivery.UserID}})
			if result != nil {
				helpers.Response(writer, http.StatusBadRequest, "Error al agregar ID del Usuario", err, nil)
			}
			helpers.Response(writer, http.StatusOK, constants.SUCCESS, err, nil)
		} else {
			helpers.Response(writer, http.StatusBadRequest, "Error Envio no Existe", err, nil)
		}
	}
}
