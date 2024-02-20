package endpoint

import "net/http"

func NotFound() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		returnError(w, http.StatusNotFound, "Endpoint Not Found")
	}
}
