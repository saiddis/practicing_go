package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (cfg *ApiConfig) Chat(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Name    string `json:"name"`
		Message string `json:"message"`
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("error reading request body: %v", err))
		return
	}
	err = json.Unmarshal(body, &params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error unmarshalling request body: %v", err))
		return
	}

	err = cfg.NC.Publish(fmt.Sprintf("chat.%s", params.Name), []byte(params.Message))
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error publishing subject: %v", err))
		return
	}

	respondWithJSON(w, 201, fmt.Sprintf("%s: %s", params.Name, params.Message))
}
