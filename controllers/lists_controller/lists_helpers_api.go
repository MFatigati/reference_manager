package listscontroller

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

func addListAPI(w http.ResponseWriter, r *http.Request) {
	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		error_controller.ErrorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}
	var newList db_controller.NewList
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&newList)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			error_controller.ErrorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			error_controller.ErrorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}
	id, err := db_controller.AddList(newList)
	if err != nil {
		error_controller.ErrorResponse(w, "Unable to add list "+err.Error(), http.StatusBadRequest)
		return
	}
	message := fmt.Sprintf("Successfully added new list (ID: %v)", id)
	w.Write([]byte(message))
}

func deleteListAPI(w http.ResponseWriter, r *http.Request) {
	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		error_controller.ErrorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}
	var list db_controller.List
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&list)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			error_controller.ErrorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			error_controller.ErrorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}
	title, err := db_controller.DeleteList(list)
	if err != nil {
		error_controller.ErrorResponse(w, "Unable to delete list "+err.Error(), http.StatusBadRequest)
		return
	}
	message := fmt.Sprintf("Successfully deleted list (Title: %v)", title)
	w.Write([]byte(message))
}

func modifyExistingListAPI(w http.ResponseWriter, r *http.Request) {
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
				db_controller.EditEntry("lists", column, value, affectedID)
			case string:
				value := v
				db_controller.EditEntry("lists", column, value, affectedID)
			default:
				error_controller.ErrorResponse(w, "Error updating entry: "+err.Error(), http.StatusBadRequest)
			}
		}
	}

	message := fmt.Sprintf("Successfully updated list (ID: %v)", affectedID)
	w.Write([]byte(message))
}

func sendAllListsJSON(w http.ResponseWriter, r *http.Request) {
	lists, err := db_controller.SelectAllLists()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	listsJSON, err := json.MarshalIndent(lists, "", "  ")
	fmt.Println(string(listsJSON))
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(lists)
}
