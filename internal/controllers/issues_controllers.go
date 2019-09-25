package controllers

import (
	"encoding/json"
	"net/http"
	"recibe_me/configs/constants"
	"recibe_me/internal/helpers"
	"recibe_me/internal/models"
)

// IssueAdd adds a Issue
func IssueAdd(w http.ResponseWriter, r *http.Request) {

	issueData := &models.Issue{}

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(issueData); err != nil {
		helpers.Response(w, http.StatusInternalServerError, constants.ERR_DECODE, err, nil)
		return
	}

	// Se valida la informacion del usuario
	if validErrs := issueData.Validate(); len(validErrs) > 0 {
		err := map[string]interface{}{"validationError": validErrs}
		helpers.Response(w, http.StatusBadRequest, constants.ERR_INVALID_DATA, err, nil)
		return
	}

	// Se guarda el Issue en la base de datos
	if err := helpers.IssuesCollection.Insert(issueData); err != nil {
		helpers.Response(w, http.StatusBadRequest, constants.ERR_INSERT_DATA, err, nil)
		return
	}

	helpers.Response(w, http.StatusOK, constants.SUCCESS, nil, nil)
}
