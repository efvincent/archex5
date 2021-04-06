package API

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Gorilla!</h1>"))
}
func Run(host string, port string) {
	router := mux.NewRouter()
	router.HandleFunc("/", homeHandler)
	addr := fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("Server running. Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
