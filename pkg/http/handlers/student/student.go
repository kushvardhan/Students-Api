package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/kushvardhan/Students-Api/pkg/utils/response"
	"github.com/kushvardhan/Students-Api/types"
	"golang.org/x/mod/sumdb/storage"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if err != nil {
			if errors.Is(err, io.EOF) {
				response.WriteJson(w, http.StatusBadRequest, "request body is empty")
				return
			}
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err := validator.New().Struct(student); err != nil{
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return 
		}

		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)

		if err != nil{
			response.WriteJson(w, http.StatusInternalServerError, err)
		}


		slog.Info("Creating a student", "student", student)


		response.WriteJson(w, http.StatusCreated, map[string]int64{
			"id": lastId,
		})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		id:= r.PathValue("id")
		slog.Info("getting a student", slog.String("id",id))

		intId, err := strconv.ParseInt(id, 10, 64);
		if err != nil{
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		student, err := storage.GetStudentById(intId)

		if err != nil{
			slog.Error("error getting user", slog.String("id", id))
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, student)
	}
}