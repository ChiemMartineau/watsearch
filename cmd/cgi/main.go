package main

import (
	"net/http/cgi"

	"github.com/Samuel-Martineau/watsearch/internal/handler"
)

func main() {
	mux := handler.Mux()

	cgi.Serve(mux)
}
