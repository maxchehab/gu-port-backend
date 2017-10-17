package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

var routes = Routes{
     Route{
          "Register",
          "POST",
          "/register",
          Register,
     },
	Route{
		"LoginHandle",
		"POST",
		"/login",
		LoginHandle,
	},
	Route{
		"PagePagination",
		"GET",
		"/pages/{count}/{offset}",
		PagePagination,
	},
	Route{
		"PageShow",
		"GET",
		"/users/{userID}/pages/{pageID}",
		PageShow,
	},
}
