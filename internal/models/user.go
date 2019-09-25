package models

import (
	"net/url"
	"regexp"

	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID               bson.ObjectId `json:"id"                 bson:"_id,omitempty" `
	FirstName        string        `json:"first_name"         bson:"first_name"`
	LastName         string        `json:"last_name"          bson:"last_name"`
	Phone            string        `json:"phone"              bson:"phone"`
	Email            string        `json:"email"              bson:"email"`
	Password         string        `json:"password,omitempty" bson:"password"`
	CreateAt         int32         `json:"created_at"         bson:"created_at"`
	Verified         bool          `json:"verified"           bson:"verified"`
	VerificationCode string        `json:"verification_code"  bson:"verification_code"`
}

type UserAuth struct {
	Email    string
	Password string
}

type UserLogResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

func (user *User) Validate() url.Values {
	errs := url.Values{}

	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if user.FirstName == "" {
		errs.Add("first_name", "El campo nombre es requerido")
	}

	if user.LastName == "" {
		errs.Add("last_name", "El campo apellido es requerido")
	}

	if user.Phone == "" {
		errs.Add("phone", "El campo telefono es requerido")
	}

	if user.Email == "" {
		errs.Add("email", "El campo correo es requerido")
	}

	if !re.MatchString(user.Email) {
		errs.Add("email", "Correo invalido")
	}

	if user.Password == "" {
		errs.Add("password", "El campo correo es requerido")
	}

	if len(user.Password) < 6 || len(user.Password) > 16 {
		errs.Add("password", "El tamaño debe estar entre 6 y 16 caracteres")
	}

	return errs
}
