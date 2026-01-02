package router

import (
	"user/controller"
	authentication "user/helper/Authentication"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/signup", controller.Signup).Methods("POST")
	r.HandleFunc("/login", controller.Login).Methods("POST")
	user := r.PathPrefix("/user").Subrouter()
	user.Use(authentication.AuthMiddleware)
	user.HandleFunc("/upload", controller.Upload).Methods("POST")
	user.HandleFunc("/download", controller.Download).Methods("POST")
	return r
}
