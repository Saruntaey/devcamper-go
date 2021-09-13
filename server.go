package main

import (
	"devcamper/config"
	"devcamper/controllers"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	// connect to DB
	s := config.ConnDB()

	r := httprouter.New()

	// bootcamp router
	bc := controllers.NewBootcamp(s)
	r.GET("/api/v1/bootcamps", bc.GetBootcamps)
	r.GET("/api/v1/bootcamps/:id", bc.GetBootcamp)
	r.POST("/api/v1/bootcamps", bc.CreateBootcamp)
	r.PUT("/api/v1/bootcamps/:id", bc.UpdateBootcamp)
	r.DELETE("/api/v1/bootcamps/:id", bc.DeleteBootcamp)

	http.ListenAndServe(":8080", r)
}
