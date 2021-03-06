package referencescontroller

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

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
	//fmt.Println(query)
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

	type refsForAllRefs struct {
		Title            string
		Author_last      string
		Author_first     string
		Publication_date string
		ID               string
		List             string
		ListName         string
	}
	var data []refsForAllRefs

	refs, err := db_controller.SelectAllReferences(sortBy, sortAsc)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	lists, err := db_controller.SelectAllLists()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	for _, ref := range refs {
		for _, list := range lists {
			if ref.List == list.ID {
				currentRef := refsForAllRefs{
					Title:            ref.Title,
					Author_last:      ref.Author_last,
					Author_first:     ref.Author_first,
					Publication_date: strconv.Itoa(ref.Publication_date),
					ID:               strconv.Itoa(ref.ID),
					List:             strconv.Itoa(ref.List),
					ListName:         list.Title,
				}
				data = append(data, currentRef)
			}
		}
	}

	t, _ := template.ParseFiles("./lib/templates/layout.html", "./lib/templates/allRefs.html")
	t.ExecuteTemplate(w, "layout1", data)
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

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

var fileOutput = regexp.MustCompile("file")

var browserOutput = regexp.MustCompile("browser")

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func outputSelected(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	selected := r.Form["selected"]
	//fmt.Println(r.Form)
	// for key, value := range r.Form {
	// 	fmt.Printf("%s = %s\n", key, value)
	// }
	//fmt.Println(selected)
	refs, err := db_controller.SelectAllReferences("title", true)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	filteredRefs := []db_controller.Reference{}
	for _, ref := range refs {
		if contains(selected, strconv.Itoa(ref.ID)) {
			filteredRefs = append(filteredRefs, ref)
		}
	}

	home, _ := os.UserHomeDir()

	var filePath string

	downloadsExists, _ := exists(home + "/downloads")

	if downloadsExists {
		filePath, _ = filepath.Abs(home + "/downloads/data.txt")
	} else {
		filePath, _ = filepath.Abs(home + "/data.txt")
	}

	fmt.Println(filePath)
	_, err = os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	for _, ref := range filteredRefs {
		title := ref.Title
		last := ref.Author_last
		first := ref.Author_first
		date := ref.Publication_date
		reference := fmt.Sprintf("%v, %v. %v, %v.\n", last, first, title, date)

		if _, err = f.Write([]byte(reference)); err != nil {
			log.Fatal(err)
		}
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	if fileOutput.FindStringSubmatch(strings.ToLower(r.FormValue("action"))) != nil {
		http.Redirect(w, r, "/references", http.StatusSeeOther)
	} else if browserOutput.FindStringSubmatch(strings.ToLower(r.FormValue("action"))) != nil {
		http.ServeFile(w, r, filePath)
	}
}
