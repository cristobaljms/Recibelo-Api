package models

import (
	"net/url"

	"gopkg.in/mgo.v2/bson"
)

type Issue struct {
	ID          bson.ObjectId `json:"id"          bson:"_id,omitempty"`
	DeliveryID  string        `json:"delivery_id" bson:"delivery_id"`
	Description string        `json:"description" bson:"description"`
	IssueType   string        `json:"issue_type"  bson:"issue_type" `
}

func (issue *Issue) Validate() url.Values {
	errs := url.Values{}

	if issue.Description == "" {
		errs.Add("description", "El campo de descripción es requerido")
	}

	if issue.DeliveryID == "" {
		errs.Add("delivery_id", "El campo de delivery_id es requerido")
	}

	if len(issue.Description) < 6 {
		errs.Add("description", "Descripción muy corta, al menos 6 ")
	}

	if issue.IssueType == "" {
		errs.Add("issue_type", "El campo issue_type es requerido")
	}

	return errs
}
