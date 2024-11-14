package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (cfg *ApiConfig) Greet(w http.ResponseWriter, r *http.Request) {
	var parameters struct {
		Name  string `json:"name"`
		Greet string `json:"greet"`
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("error reading request body: %v", err))
		return
	}
	err = json.Unmarshal(body, &parameters)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error unmarshalling request body: %v", err))
		return
	}

	err = cfg.NC.Publish(fmt.Sprintf("greet.%s", parameters.Name), []byte(parameters.Greet))
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error publishing subject: %v", err))
		return
	}

	respondWithJSON(w, 201, fmt.Sprintf("%s: %s", parameters.Name, parameters.Greet))
}
