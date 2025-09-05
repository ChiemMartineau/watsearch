package main

import (
	"net/http"

	"github.com/Samuel-Martineau/watsearch/internal/handler"
)

func main() {
	mux := handler.Mux()

	http.ListenAndServe("localhost:1234", mux)
}
