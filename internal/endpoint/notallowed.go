package endpoint

import "net/http"

func NotAllowed() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		returnError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
	}
}
