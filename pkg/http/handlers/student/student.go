package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/kushvardhan/Students-Api/pkg/utils/response"
	"github.com/kushvardhan/Students-Api/types"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if err != nil {
			if errors.Is(err, io.EOF) {
				response.WriteJson(w, http.StatusBadRequest, "request body is empty")
				return
			}
			response.WriteJson(w, http.StatusBadRequest, err.Error())
			return
		}

		slog.Info("Creating a student", "student", student)

		response.WriteJson(w, http.StatusCreated, map[string]string{
			"success": "OK",
		})
	}
}
