package routes

import (
	"log"
	"net/http"

	"myapp/controller"

	"github.com/gorilla/mux"
)

func InitializeRoutes() {
	r := mux.NewRouter()

	r.HandleFunc("/student/add", controller.AddStudent).Methods("POST")
	r.HandleFunc("/student/{sid}", controller.GetStud).Methods("GET")
	r.HandleFunc("/students", controller.GetAllStuds).Methods("GET")
	r.HandleFunc("/student/{sid}", controller.UpdateStud).Methods("PUT")
	r.HandleFunc("/student/{sid}", controller.DeleteStud).Methods("DELETE")

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
