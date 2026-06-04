package broker

import (
	"io"
	"net/http"
)

// individual handler func fo each route
func AppendHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Giving \n")
}
