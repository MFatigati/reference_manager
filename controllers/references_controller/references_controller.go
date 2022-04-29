package referencescontroller

import (
	"fmt"
	"net/http"
	"regexp"
)

var apiRoute = regexp.MustCompile("api")
var formEditRoute = regexp.MustCompile("edit")
var formNewRoute = regexp.MustCompile("new")
var formDeleteRoute = regexp.MustCompile("delete")

// func apiPostNewRef(w http.ResponseWriter, r *http.Request) {
// 	headerContentTtype := r.Header.Get("Content-Type")
// 	if headerContentTtype != "application/json" {
// 		error_controller.ErrorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
// 		return
// 	}
// 	var newRef db_controller.NewReference
// 	var unmarshalErr *json.UnmarshalTypeError

// 	decoder := json.NewDecoder(r.Body)
// 	decoder.DisallowUnknownFields()
// 	err := decoder.Decode(&newRef)
// 	if err != nil {
// 		if errors.As(err, &unmarshalErr) {
// 			error_controller.ErrorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
// 		} else {
// 			error_controller.ErrorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
// 		}
// 		return
// 	}
// 	id, err := db_controller.AddReference(newRef)
// 	if err != nil {
// 		error_controller.ErrorResponse(w, "Unable to add reference "+err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	message := fmt.Sprintf("Successfully added new reference (ID: %v)", id)
// 	w.Write([]byte(message))
// }

// func formPostNewRef(w http.ResponseWriter, r *http.Request) {
// 	//fmt.Println("Form submitted.")
// 	var newRef db_controller.NewReference
// 	newRef.Title = r.FormValue("title")
// 	newRef.Author_last = r.FormValue("last")
// 	newRef.Author_first = r.FormValue("first")
// 	date, err := strconv.Atoi(r.FormValue("date"))
// 	if err != nil {
// 		error_controller.ErrorResponse(w, "Unable to add reference "+err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	newRef.Publication_date = date
// 	listID, err := strconv.Atoi(r.FormValue("id"))
// 	if err != nil {
// 		error_controller.ErrorResponse(w, "Unable to add reference "+err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	newRef.List = listID
// 	id, err := db_controller.AddReference(newRef)
// 	if err != nil {
// 		error_controller.ErrorResponse(w, "Unable to add reference "+err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	message := fmt.Sprintf("Successfully added new reference (ID: %v)", id)
// 	w.Write([]byte(message))
// }

// func showNewRefForm(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")
// 	html, err := os.ReadFile("./lib/templates/addRef.html")
// 	if err != nil {
// 		error_controller.ErrorResponse(w, "Unable to render this page", http.StatusInternalServerError)
// 	}
// 	fmt.Fprint(w, string(html))
// }

// func sendAllRefsJSON(w http.ResponseWriter, r *http.Request) {
// 	references, err := db_controller.SelectAllReferences()
// 	if err != nil {
// 		error_controller.ErrorResponse(w, "Unable to get references", http.StatusInternalServerError)
// 	}
// 	// referencesJSON, err := json.MarshalIndent(references, "", "  ")
// 	// fmt.Println(string(referencesJSON))
// 	// if err != nil {
// 	// 	fmt.Printf("Error: %s", err)
// 	// 	return
// 	// }
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(references)
// }

// func showEditRefForm(w http.ResponseWriter, r *http.Request) {
// 	query := r.URL.Query()
// 	fmt.Println(query)
// 	refId := string(query["refId"][0])
// 	refTitle := string(query["refTitle"][0])
// 	refAuthorLast := string(query["refAuthorLast"][0])
// 	refAuthorFirst := string(query["refAuthorFirst"][0])
// 	refDate := string(query["refDate"][0])
// 	listID := string(query["listID"][0])

// 	t, _ := template.ParseFiles("./lib/templates/editRef.html")
// 	t.Execute(w, map[string]string{
// 		"refID":          refId,
// 		"refTitle":       refTitle,
// 		"refAuthorLast":  refAuthorLast,
// 		"refAuthorFirst": refAuthorFirst,
// 		"refDate":        refDate,
// 		"listID":         listID,
// 	})
// }

// func modifyExistingRefAPI(w http.ResponseWriter, r *http.Request) {
// 	headerContentTtype := r.Header.Get("Content-Type")
// 	if headerContentTtype != "application/json" {
// 		error_controller.ErrorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
// 		return
// 	}

// 	data := make(map[string]interface{})
// 	body, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		panic(err)
// 	}
// 	err = json.Unmarshal([]byte(body), &data)
// 	if err != nil {
// 		panic(err)
// 	}

// 	var id float64

// 	for key, val := range data {
// 		if key == "id" {
// 			id = val.(float64)
// 		}
// 	}

// 	affectedID := int(math.Trunc(id))

// 	for key, val := range data {
// 		if key != "id" {
// 			column := key
// 			switch v := val.(type) {
// 			case float64:
// 				value := strconv.FormatInt(int64(val.(float64)), 10)
// 				db_controller.EditEntry("texts", column, value, affectedID)
// 			case string:
// 				value := v
// 				db_controller.EditEntry("texts", column, value, affectedID)
// 			default:
// 				error_controller.ErrorResponse(w, "Error updating entry: "+err.Error(), http.StatusBadRequest)
// 			}
// 		}
// 	}

// 	message := fmt.Sprintf("Successfully updated reference (ID: %v)", affectedID)
// 	w.Write([]byte(message))
// }

// func deleteRefAPI(w http.ResponseWriter, r *http.Request) {
// 	headerContentTtype := r.Header.Get("Content-Type")
// 	if headerContentTtype != "application/json" {
// 		error_controller.ErrorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
// 		return
// 	}
// 	var ref db_controller.Reference
// 	var unmarshalErr *json.UnmarshalTypeError

// 	decoder := json.NewDecoder(r.Body)
// 	decoder.DisallowUnknownFields()
// 	err := decoder.Decode(&ref)
// 	if err != nil {
// 		if errors.As(err, &unmarshalErr) {
// 			error_controller.ErrorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
// 		} else {
// 			error_controller.ErrorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
// 		}
// 		return
// 	}
// 	title, err := db_controller.DeleteRef(ref)
// 	if err != nil {
// 		error_controller.ErrorResponse(w, "Unable to delete reference "+err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	message := fmt.Sprintf("Successfully deleted reference (Title: %v)", title)
// 	w.Write([]byte(message))
// }

// func viewAllRefs(w http.ResponseWriter, r *http.Request) {
// 	refs, err := db_controller.SelectAllReferences()
// 	if err != nil {
// 		fmt.Printf("Error: %s", err)
// 		return
// 	}
// 	t, _ := template.ParseFiles("./lib/templates/allRefs.html")
// 	t.Execute(w, refs)
// }

// func modifyExistingRefForm(w http.ResponseWriter, r *http.Request) {
// 	refTitle := r.FormValue("title")
// 	refAuthorLast := r.FormValue("last")
// 	refAuthorFirst := r.FormValue("first")
// 	refDate := r.FormValue("date")
// 	// refID := r.FormValue("id")
// 	// refDate, err := strconv.Atoi(r.FormValue("date"))
// 	// if err != nil {
// 	// 	fmt.Printf("Error parsing refDate: %s", err)
// 	// 	return
// 	// }
// 	refID, err := strconv.Atoi(r.FormValue("refID"))
// 	if err != nil {
// 		fmt.Printf("Error parsing refID: %s", err)
// 		return
// 	}
// 	listID := r.FormValue("listID")

// 	data := make(map[string]interface{})
// 	data["Title"] = refTitle
// 	data["Author_last"] = refAuthorLast
// 	data["Author_first"] = refAuthorFirst
// 	data["Publication_date"] = refDate
// 	data["List"] = listID

// 	affectedID := refID

// 	fmt.Println(data, listID, affectedID)

// 	for key, val := range data {
// 		if key != "id" {
// 			column := key
// 			switch v := val.(type) {
// 			// case float64:
// 			// 	//value := strconv.FormatInt(int64(val.(float64)), 10)
// 			// 	value := v
// 			// 	db_controller.EditEntry("texts", column, value, affectedID)
// 			case string:
// 				value := v
// 				db_controller.EditEntry("texts", column, value, affectedID)
// 			default:
// 				error_controller.ErrorResponse(w, "Error updating entry: "+err.Error(), http.StatusBadRequest)
// 			}
// 		}
// 	}

// 	http.Redirect(w, r, "/list/"+listID, http.StatusSeeOther)
// }

// func deleteRefForm(w http.ResponseWriter, r *http.Request) {
// 	refID, err := strconv.Atoi(r.URL.Path[len("/reference/delete/"):])
// 	if err != nil {
// 		fmt.Printf("Error: %s", err)
// 		return
// 	}
// 	references, err := db_controller.SelectAllReferences()
// 	if err != nil {
// 		fmt.Printf("Error: %s", err)
// 		return
// 	}
// 	var currentRef db_controller.Reference
// 	for _, ref := range references {
// 		if ref.ID == refID {
// 			currentRef.ID = ref.ID
// 			currentRef.Title = ref.Title
// 			currentRef.List = ref.List
// 		}
// 	}
// 	title, err := db_controller.DeleteRef(currentRef)
// 	if err != nil {
// 		fmt.Printf("Error deleting reference %v: %s", title, err)
// 		return
// 	}
// 	fmt.Println("Deleted refeference.")
// 	redirectToList := strconv.Itoa(currentRef.List)
// 	http.Redirect(w, r, "/list/"+redirectToList, http.StatusSeeOther)
// }

func ReferenceHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if formNewRoute.FindStringSubmatch(r.URL.Path) != nil {
			showNewRefForm(w, r)
		} else if formEditRoute.FindStringSubmatch(r.URL.Path) != nil {
			showEditRefForm(w, r)
		} else if formDeleteRoute.FindStringSubmatch(r.URL.Path) != nil {
			deleteRefForm(w, r)
		}
	case "POST":
		if formNewRoute.FindStringSubmatch(r.URL.Path) != nil {
			formPostNewRef(w, r)
			return
		} else if apiRoute.FindStringSubmatch(r.URL.Path) != nil {
			apiPostNewRef(w, r)
			return
		} else if formEditRoute.FindStringSubmatch(r.URL.Path) != nil {
			modifyExistingRefForm(w, r)
		}
	case "PUT":
		modifyExistingRefAPI(w, r)
	case "DELETE":
		deleteRefAPI(w, r)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported for this route.")
	}
}

func ReferencesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if apiRoute.FindStringSubmatch(r.URL.Path) != nil {
			sendAllRefsJSON(w, r)
		} else {
			viewAllRefs(w, r)
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET method is supported for this route.")
	}
}
