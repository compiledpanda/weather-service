package endpoint

import "net/http"

type GetConditionsResponse struct {
	Temperature int    `json:"temperature"`
	Units       string `json:"units"`
	Condition   string `json:"condition"`
	FeelsLike   string `json:"feelsLike"`
}

func GetConditions() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		// TODO create function
		res := GetConditionsResponse{}

		returnJSON(w, http.StatusOK, res)
	}
}
