package referencescontroller

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	db_controller "example.com/ref_manager/controllers/db_controller"
	error_controller "example.com/ref_manager/controllers/error_controller"
)

func formPostNewRef(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Form submitted.")
	var newRef db_controller.NewReference
	newRef.Title = r.FormValue("title")
	newRef.Author_last = r.FormValue("last")
	newRef.Author_first = r.FormValue("first")
	date, err := strconv.Atoi(r.FormValue("date"))
	if err != nil {
		error_controller.ErrorResponse(w, "Unable to add reference "+err.Error(), http.StatusBadRequest)
		return
	}
	newRef.Publication_date = date
	listID, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		error_controller.ErrorResponse(w, "Unable to add reference "+err.Error(), http.StatusBadRequest)
		return
	}
	newRef.List = listID
	_, err = db_controller.AddReference(newRef)
	// id, err := db_controller.AddReference(newRef)
	if err != nil {
		error_controller.ErrorResponse(w, "Unable to add reference "+err.Error(), http.StatusBadRequest)
		return
	}
	// message := fmt.Sprintf("Successfully added new reference (ID: %v)", id)
	// w.Write([]byte(message))
	http.Redirect(w, r, "/list/"+strconv.Itoa(listID), http.StatusSeeOther)
}

func showNewRefForm(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	fmt.Println(query)
	listID := string(query["listID"][0])
	// w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// html, err := os.ReadFile("./lib/templates/addRef.html")
	// if err != nil {
	// 	error_controller.ErrorResponse(w, "Unable to render this page", http.StatusInternalServerError)
	// }
	// fmt.Fprint(w, string(html))
	t, _ := template.ParseFiles("./lib/templates/layout.html", "./lib/templates/addRef.html")
	// t.ExecuteTemplate(w, "layout1", nil)
	// t, _ := template.ParseFiles("./lib/templates/addRef.html")
	t.ExecuteTemplate(w, "layout1", map[string]string{
		"ListID": listID,
	})
}

func showEditRefForm(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	fmt.Println(query)
	refId := string(query["refId"][0])
	refTitle := string(query["refTitle"][0])
	refAuthorLast := string(query["refAuthorLast"][0])
	refAuthorFirst := string(query["refAuthorFirst"][0])
	refDate := string(query["refDate"][0])
	listID := string(query["listID"][0])

	t, _ := template.ParseFiles("./lib/templates/layout.html", "./lib/templates/editRef.html")
	// t.ExecuteTemplate(w, "layout1", nil)
	// t, _ := template.ParseFiles("./lib/templates/editRef.html")
	t.ExecuteTemplate(w, "layout1", map[string]string{
		"refID":          refId,
		"refTitle":       refTitle,
		"refAuthorLast":  refAuthorLast,
		"refAuthorFirst": refAuthorFirst,
		"refDate":        refDate,
		"listID":         listID,
	})
}

func queryToColumnHeader(param string) (header string) {
	switch param {
	case "title":
		return "title"
	case "last":
		return "author_last"
	case "first":
		return "author_first"
	case "date":
		return "publication_date"
	case "listID":
		return "list"
	case "refID":
		return "id"
	}
	return "title"
}

var lastSort string
var sortAsc bool

func viewAllRefs(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sortBy")
	if lastSort == sortBy {
		sortAsc = !sortAsc
	}
	lastSort = sortBy
	sortBy = queryToColumnHeader(sortBy)

	refs, err := db_controller.SelectAllReferences(sortBy, sortAsc)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	t, _ := template.ParseFiles("./lib/templates/layout.html", "./lib/templates/allRefs.html")
	t.ExecuteTemplate(w, "layout1", refs)
}

func modifyExistingRefForm(w http.ResponseWriter, r *http.Request) {
	refTitle := r.FormValue("title")
	refAuthorLast := r.FormValue("last")
	refAuthorFirst := r.FormValue("first")
	refDate := r.FormValue("date")
	refID, err := strconv.Atoi(r.FormValue("refID"))
	if err != nil {
		fmt.Printf("Error parsing refID: %s", err)
		return
	}
	listID := r.FormValue("listID")

	data := make(map[string]interface{})
	data["Title"] = refTitle
	data["Author_last"] = refAuthorLast
	data["Author_first"] = refAuthorFirst
	data["Publication_date"] = refDate
	data["List"] = listID

	affectedID := refID

	fmt.Println(data, listID, affectedID)

	for key, val := range data {
		if key != "id" {
			column := key
			switch v := val.(type) {
			// case float64:
			// 	//value := strconv.FormatInt(int64(val.(float64)), 10)
			// 	value := v
			// 	db_controller.EditEntry("texts", column, value, affectedID)
			case string:
				value := v
				db_controller.EditEntry("texts", column, value, affectedID)
			default:
				error_controller.ErrorResponse(w, "Error updating entry: "+err.Error(), http.StatusBadRequest)
			}
		}
	}

	http.Redirect(w, r, "/list/"+listID, http.StatusSeeOther)
}

func deleteRefForm(w http.ResponseWriter, r *http.Request) {
	refID, err := strconv.Atoi(r.URL.Path[len("/reference/delete/"):])
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	references, err := db_controller.SelectAllReferences("title", true)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	var currentRef db_controller.Reference
	for _, ref := range references {
		if ref.ID == refID {
			currentRef.ID = ref.ID
			currentRef.Title = ref.Title
			currentRef.List = ref.List
		}
	}
	title, err := db_controller.DeleteRef(currentRef)
	if err != nil {
		fmt.Printf("Error deleting reference %v: %s", title, err)
		return
	}
	fmt.Println("Deleted refeference.")
	redirectToList := strconv.Itoa(currentRef.List)
	http.Redirect(w, r, "/list/"+redirectToList, http.StatusSeeOther)
}

func outputSelected(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	selected := r.Form["selected"]
	//fmt.Println(r.Form)
	// for key, value := range r.Form {
	// 	fmt.Printf("%s = %s\n", key, value)
	// }
	fmt.Println(selected)
	http.Redirect(w, r, "/references", http.StatusSeeOther)
}
