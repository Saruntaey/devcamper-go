package main

import (
	"devcamper/config"
	"devcamper/controllers"
	"devcamper/models"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	// connect to DB
	conn := config.ConnDB()

	// mount models to DB
	conn.Register(&models.Bootcamp{}, "bootcamps")

	r := httprouter.New()

	// bootcamp router
	bc := controllers.NewBootcamp(conn)
	r.GET("/api/v1/bootcamps", bc.GetBootcamps)
	r.GET("/api/v1/bootcamps/:id", bc.GetBootcamp)
	r.POST("/api/v1/bootcamps", bc.CreateBootcamp)
	r.PUT("/api/v1/bootcamps/:id", bc.UpdateBootcamp)
	r.DELETE("/api/v1/bootcamps/:id", bc.DeleteBootcamp)

	http.ListenAndServe(":8080", r)
}
