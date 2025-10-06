package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/d1vyanshu-kumar/students-api/internal/storage"
	"github.com/d1vyanshu-kumar/students-api/internal/types"
	"github.com/d1vyanshu-kumar/students-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

// and in near future if we want to add a new dependecy  we can inject here inside a new function.
// this is time where we are going to inject the storage dependency here so that we can use it inside the handler function to store the student data in the database.
// and here are are gving the interface type not the concrete type because we want to follow the dependency inversion principle.
func New(storage storage.Storage) http.HandlerFunc { // now go and pass into the main.go file
	return  func(w http.ResponseWriter, r *http.Request) {

		slog.Info("creating a student")
		
		var student types.Student        // decode for the incoming request body this is the way in Go decode then read
		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			// we are going to give a json response
			 response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(errors.New("body cannot be empty")))
			 return
		}

		// now if there is no empty body but there is some other error here is how can we catch in genrall error:

		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// request validation....

		if err := validator.New().Struct(student); err != nil {
			validErr := err.(validator.ValidationErrors) // type assertion
			response.WriteJSON(w, http.StatusBadRequest, response.ValidError(validErr))
			return
		}

		lastID, err := storage.CreateStudent(student.Name, student.Email, student.Age) // now we can use the storage dependency here to store the student data in the database.

		slog.Info("student created", slog.Int64("id", lastID))
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err) // 500 status code database error
			return
		}

		// write some dummy data to the response
		response.WriteJSON(w, http.StatusCreated, map[string]int64{
			"ID": lastID,
		 })
	}

}