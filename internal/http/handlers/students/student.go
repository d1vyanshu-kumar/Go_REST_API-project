package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/d1vyanshu-kumar/students-api/internal/types"
	"github.com/d1vyanshu-kumar/students-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

// and in near future if we want to add a new dependecy  we can inject here inside a new function.
func New() http.HandlerFunc {
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

		// write some dummy data to the response
		response.WriteJSON(w, http.StatusCreated, map[string]string{
			"message": "student created successfully",
		 })
	}

}