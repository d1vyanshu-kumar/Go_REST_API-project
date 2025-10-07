package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

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

// for getting the student by ID we need to create a new handler function
func GetByID(storage storage.Storage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// we are going to get the ID from the URL path in go there is inbuilt servemux we are going to use here:

		id := r.PathValue("id") // this is going to return the string value of the ID from the URL path.
		slog.Info("getting a student by ID", slog.String("ID", id))

	 // now we need to a method over storage to get the student.

	 intId, err := strconv.ParseInt(id, 10, 64)

	 if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
		return
	 }
	  student, err := storage.GetStudentByID(intId)

	  if err != nil {
		slog.Error("error getting user", slog.String("id", id))
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	  }

	  response.WriteJSON(w, http.StatusOK, student)

	}

}


func GetList(storage storage.Storage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// so we need to write some logic here to get the list of the students from the database.

		slog.Info("getting list of students")

		storages, err := storage.GetStudents()

		if err != nil {
			slog.Error("error getting list of students", slog.String("error", err.Error()))
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		// and this write json method is going to write/ encode the json response to the client.
		response.WriteJSON(w, http.StatusOK, storages)
	}

}