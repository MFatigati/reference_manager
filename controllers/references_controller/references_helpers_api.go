package referencescontroller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"

	db_controller "example.com/ref_manager/controllers/db_controller"
	error_controller "example.com/ref_manager/controllers/error_controller"
)

func apiPostNewRef(w http.ResponseWriter, r *http.Request) {
	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		error_controller.ErrorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}
	var newRef db_controller.NewReference
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&newRef)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			error_controller.ErrorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			error_controller.ErrorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}
	id, err := db_controller.AddReference(newRef)
	if err != nil {
		error_controller.ErrorResponse(w, "Unable to add reference "+err.Error(), http.StatusBadRequest)
		return
	}
	message := fmt.Sprintf("Successfully added new reference (ID: %v)", id)
	w.Write([]byte(message))
}

func sendAllRefsJSON(w http.ResponseWriter, r *http.Request) {
	references, err := db_controller.SelectAllReferences()
	if err != nil {
		error_controller.ErrorResponse(w, "Unable to get references", http.StatusInternalServerError)
	}
	// referencesJSON, err := json.MarshalIndent(references, "", "  ")
	// fmt.Println(string(referencesJSON))
	// if err != nil {
	// 	fmt.Printf("Error: %s", err)
	// 	return
	// }
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(references)
}

func modifyExistingRefAPI(w http.ResponseWriter, r *http.Request) {
	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		error_controller.ErrorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}

	data := make(map[string]interface{})
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		panic(err)
	}

	var id float64

	for key, val := range data {
		if key == "id" {
			id = val.(float64)
		}
	}

	affectedID := int(math.Trunc(id))

	for key, val := range data {
		if key != "id" {
			column := key
			switch v := val.(type) {
			case float64:
				value := strconv.FormatInt(int64(val.(float64)), 10)
				db_controller.EditEntry("texts", column, value, affectedID)
			case string:
				value := v
				db_controller.EditEntry("texts", column, value, affectedID)
			default:
				error_controller.ErrorResponse(w, "Error updating entry: "+err.Error(), http.StatusBadRequest)
			}
		}
	}

	message := fmt.Sprintf("Successfully updated reference (ID: %v)", affectedID)
	w.Write([]byte(message))
}

func deleteRefAPI(w http.ResponseWriter, r *http.Request) {
	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		error_controller.ErrorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}
	var ref db_controller.Reference
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&ref)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			error_controller.ErrorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			error_controller.ErrorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}
	title, err := db_controller.DeleteRef(ref)
	if err != nil {
		error_controller.ErrorResponse(w, "Unable to delete reference "+err.Error(), http.StatusBadRequest)
		return
	}
	message := fmt.Sprintf("Successfully deleted reference (Title: %v)", title)
	w.Write([]byte(message))
}
