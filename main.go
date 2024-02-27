/*
* Created on 27 Feb 2024
* @author Sai Sumanth
 */

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Sai7xp/gomuxmongo/controllers"
	"github.com/Sai7xp/gomuxmongo/routes"
)

func main() {
	fmt.Println("MongoDB CRUD using Go!")

	/// initialize Database
	controllers.Init()

	/// get router from routes
	router := routes.Router()

	fmt.Println("Server is up and running at port 6000")
	log.Fatal(http.ListenAndServe(":6000", router))
}
