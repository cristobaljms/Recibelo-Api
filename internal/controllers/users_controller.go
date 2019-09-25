package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"recibe_me/configs/constants"
	"recibe_me/internal/helpers"
	"recibe_me/internal/models"
)

func EditUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	defer r.Body.Close()

	user := models.User{}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		helpers.Response(w, http.StatusBadRequest, constants.ERR_INSERT_DATA, err, nil)
		return
	}

	if !bson.IsObjectIdHex(params["userId"]){
		helpers.Response(w, http.StatusNotFound, constants.ERR_USER_NOT_FOUND, "El ID es inválido.", nil)
		return
	}

    userID := bson.ObjectIdHex(params["userId"])

    updatedUserData := bson.M{
      "$set": map[string]interface{}{
        "first_name": user.FirstName,
        "last_name":  user.LastName,
        "phone":      user.Phone,
      },
    }

	if err := helpers.UsersCollection.Update(bson.M{"_id" : userID}, updatedUserData); err != nil {
		helpers.Response(w, http.StatusInternalServerError, constants.ERR_USER_NOT_FOUND, err, nil)
		return
	}

    if err := helpers.UsersCollection.Find(bson.M{"_id": userID}).One(&user); err != nil {
      helpers.Response(w, http.StatusInternalServerError, constants.ERR_USER_NOT_FOUND, err, nil)
    }

    user.Password = ""

    helpers.Response(w, http.StatusOK, constants.SUCCESS, nil, user)
}
