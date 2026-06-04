package broker

import (
	"fmt"
	"log"
	"net/http"
)

// sets up the server
func Run() {
	fmt.Println("heyo")

	http.HandleFunc("/append", AppendHandler)
	log.Fatal(http.ListenAndServe(":5001", nil))

}
