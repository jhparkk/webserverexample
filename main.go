package main

import (
	"net/http"

	"jhpark.sinsiway.com/webserver/myapp"
)

func main() {
	http.ListenAndServe(":3000", myapp.NewHandler())
}
