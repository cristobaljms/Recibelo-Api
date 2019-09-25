package server

import (
	"net/http"
	ctr "recibe_me/internal/controllers"
	mid "recibe_me/internal/middlewares"

	"github.com/gorilla/mux"
)

// Route is a Route type
type Route struct {
	Name       string
	Method     string
	Pattern    string
	HandleFunc http.HandlerFunc
}

// Routes is a array of Route
type Routes []Route

// NewRouter ...
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Name(route.Name).
			Path(route.Pattern).
			Handler(route.HandleFunc)
	}
	return router
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		ctr.Index,
	},
	Route{
		"Issues",
		"POST",
		"/issues",
		mid.Authenticate(ctr.IssueAdd),
	},
	Route{
		"listDeliveries",
		"GET",
		"/deliveries",
		mid.Authenticate(ctr.ListDeliveries),
	},
	Route{
		"getDelivery",
		"GET",
		"/deliveries/{id}",
		mid.Authenticate(ctr.GetDelivery),
	},
	Route{
		"deliveries.ratings.store",
		"POST",
		"/deliveries/{id}/ratings",
		mid.Authenticate(ctr.Rate),
	},
	Route{
		"AddDelivery",
		"POST",
		"/adddelivery",
		ctr.AddDelivery,
	},
	Route{
		"register.delivery",
		"POST",
		"/register-delivery",
		ctr.RegisterDelivery,
	},
	Route{
		"Signup",
		"POST",
		"/signup",
		ctr.SignUp,
	},
	Route{
		"Login",
		"POST",
		"/login",
		ctr.Login,
	},
	Route{
		"ResendVerificationCode",
		"GET",
		"/resend-verification-code/{userId}",
		ctr.ResendVerificationCode,
	},
	Route{
		"VerificationAccount",
		"GET",
		"/verification-account/{verificationCode}/{userId}",
		ctr.VerificationAccount,
	},
	Route{
		"SendPasswordRecoveryCode",
		"GET",
		"/send-password-recovery-code/{email}",
		ctr.SendPasswordRecoveryCode,
	},
	Route{
		"PasswordRecovery",
		"GET",
		"/password-recovery/{passwordRecoveryCode}/{userId}",
		ctr.PasswordRecovery,
	},
	Route{
		"EditUser",
		"PUT",
		"/users/{userId}",
		ctr.EditUser,
	},
	Route{
		"PasswordUpdate",
		"PUT",
		"/password-update",
		mid.Authenticate(ctr.PasswordUpdate),
	},
}
