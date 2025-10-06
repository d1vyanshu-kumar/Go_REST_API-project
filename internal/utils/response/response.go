package response

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status  string `json:"status"`
	Error string      
}

const (
	StatusOk = "ok"
	StatusError = "Error"
)

func WriteJSON(w http.ResponseWriter, status int, data any) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// here wr are writing thr response so what we can do here is to convert the struct data to json formate.use Encode here 
	return json.NewEncoder(w).Encode(data)

}

func GeneralError(err error) Response {

	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}

func ValidError(errs validator.ValidationErrors) Response {

	//. we need to formate that error in a way that is user friendly. add you can see that those error are slice so we are going to use a for loop here.

	var errorMessages []string


	for _, err := range errs {

		switch err.Tag() {
		case "required":
			errorMessages = append(errorMessages, err.Field()+" is required")
		case "email":
			errorMessages = append(errorMessages, err.Field()+" must be a valid email")
		case "gte":
			errorMessages = append(errorMessages, err.Field()+" must be greater than or equal to "+err.Param())
		case "lte":
			errorMessages = append(errorMessages, err.Field()+" must be less than or equal to "+err.Param())
		default:
			errorMessages = append(errorMessages, err.Field()+" is not valid")
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errorMessages, ", "),
	}
}