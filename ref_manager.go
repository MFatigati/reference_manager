package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	db_controller "example.com/ref_manager/controllers/db_controller"
	lists_controller "example.com/ref_manager/controllers/lists_controller"
	references_controller "example.com/ref_manager/controllers/references_controller"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/lists", http.StatusSeeOther)
}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	port := os.Getenv("PORT")

	fs := http.FileServer(http.Dir("./public"))

	db_controller.ResetDB()
	http.Handle("/public/", http.StripPrefix("/public/", fs))
	http.HandleFunc("/lists/", lists_controller.ListsHandler)
	http.HandleFunc("/lists/delete/", lists_controller.ListHandler)
	http.HandleFunc("/list/new", lists_controller.ListHandler)
	http.HandleFunc("/list/edit", lists_controller.ListHandler)
	http.HandleFunc("/list/", lists_controller.ListHandler)
	http.HandleFunc("/reference/delete/", references_controller.ReferenceHandler)
	http.HandleFunc("/reference/new", references_controller.ReferenceHandler)
	http.HandleFunc("/reference/edit", references_controller.ReferenceHandler)
	http.HandleFunc("/references", references_controller.ReferencesHandler)
	http.HandleFunc("/api/references", references_controller.ReferencesHandler)
	http.HandleFunc("/api/reference", references_controller.ReferenceHandler)
	http.HandleFunc("/api/lists", lists_controller.ListsHandler)
	http.HandleFunc("/api/list", lists_controller.ListHandler)
	http.HandleFunc("/api/resetdb", db_controller.DBHandler)
	http.HandleFunc("/resetdb", db_controller.DBHandler)
	http.HandleFunc("/", homeHandler)
	log.Printf("Listening on port %v", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
