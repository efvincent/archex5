package API

import (
	"fmt"
	"log"
	"net/http"

	"github.com/efvincent/archex5/models"
	"github.com/gorilla/mux"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	pm := models.MakeProduct()
	w.Write([]byte(fmt.Sprintf("<h1>Product: %s</h1>", pm.SKU)))
}
func Run(host string, port string) {
	router := mux.NewRouter()
	router.HandleFunc("/", homeHandler)
	addr := fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("Server running. Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
