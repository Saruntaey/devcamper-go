package main

import (
	"devcamper/config"
	"devcamper/controllers"
	"devcamper/models"
	"fmt"
	"log"
	"net/http"
	"os"

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

	port := os.Getenv("PORT")
	port = fmt.Sprint(":", port)
	fmt.Printf("Listening on port %s\n", port)
	log.Fatalln(http.ListenAndServe(port, r))
}
