package routes

import (
	"net/http"

	"github.com/Sai7xp/gomuxmongo/controllers"
	"github.com/gorilla/mux"
)

func Router() (router *mux.Router) {
	router = mux.NewRouter()
	/// Default home route
	router.HandleFunc("/", controllers.HomeHandler).Methods(http.MethodGet)

	/// for fetching all cats
	router.HandleFunc("/api/getAllCats", controllers.GetAllCatsHandler).Methods(http.MethodGet)

	/// for adding a new Cat Details
	router.HandleFunc("/api/addCat", controllers.AddCatHandler).Methods(http.MethodPost)

	/// for removing cat details from db
	router.HandleFunc("/api/deleteCat/{catId}", controllers.DeleteCatHandler).Methods(http.MethodDelete)

	/// for updating cat details
	router.HandleFunc("/api/updateCatName/{catId}", controllers.UpdateCatNameHandler).Methods(http.MethodPut)

	return
}
