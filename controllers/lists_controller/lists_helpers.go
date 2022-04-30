package listscontroller

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	db_controller "example.com/ref_manager/controllers/db_controller"
	error_controller "example.com/ref_manager/controllers/error_controller"
)

func formPostNewList(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Form submitted.")
	var newList db_controller.NewList
	newList.Title = r.FormValue("title")
	var err error
	_, err = db_controller.AddList(newList)
	if err != nil {
		error_controller.ErrorResponse(w, "Unable to add list "+err.Error(), http.StatusBadRequest)
		return
	}
	// message := fmt.Sprintf("Successfully added new list (ID: %v)", id)
	// w.Write([]byte(message))
	http.Redirect(w, r, "/lists", http.StatusSeeOther)
}

func showNewListForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// html, err := os.ReadFile("./lib/templates/addList.html")
	// if err != nil {
	// 	error_controller.ErrorResponse(w, "Unable to render this page", http.StatusInternalServerError)
	// }
	// fmt.Fprint(w, string(html))
	t, _ := template.ParseFiles("./lib/templates/layout.html", "./lib/templates/addList.html")
	t.ExecuteTemplate(w, "layout1", nil)
}

func renderAllLists(w http.ResponseWriter, r *http.Request) {
	lists, err := db_controller.SelectAllLists()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	t, _ := template.ParseFiles("./lib/templates/layout.html", "./lib/templates/allLists.html")
	//t, _ := template.ParseFiles("./lib/templates/layout.html")
	//t, _ := template.ParseFiles("./lib/templates/allLists.html")
	//t.Execute(w, lists)
	t.ExecuteTemplate(w, "layout1", lists)
	//fmt.Println(t)
	//t.ExecuteTemplate(w, "layout", lists)
}

func showSingleList(w http.ResponseWriter, r *http.Request) {
	// type refsAllStrings struct {
	// 	Title            string
	// 	Author_last      string
	// 	Author_first     string
	// 	Publication_date string
	// 	ID               string
	// 	List             string
	// }

	type Data struct {
		Title  string
		ListID string
		Refs   []db_controller.Reference
	}

	listID := r.URL.Path[len("/list/"):]
	intID, err := strconv.Atoi(listID)
	if err != nil {
		error_controller.ErrorResponse(w, "Unable to render this page", http.StatusInternalServerError)
	}
	refs, err := db_controller.SelectAllReferences("title", true)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	filteredRefs := []db_controller.Reference{}
	for _, ref := range refs {
		if ref.List == intID {
			filteredRefs = append(filteredRefs, ref)
		}
	}
	lists, err := db_controller.SelectAllLists()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	var currentList string
	for _, list := range lists {
		if list.ID == intID {
			currentList = list.Title
		}
	}

	data := Data{currentList, listID, filteredRefs}
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	t, _ := template.ParseFiles("./lib/templates/layout.html", "./lib/templates/refsForList.html")
	// t.ExecuteTemplate(w, "layout1", nil)
	// t, _ := template.ParseFiles("./lib/templates/refsForList.html")
	t.ExecuteTemplate(w, "layout1", data)
}

func showEditListForm(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	// fmt.Println(query)
	listID := string(query["listID"][0])
	listTitle := string(query["listTitle"][0])
	//fmt.Println(listTitle)
	t, _ := template.ParseFiles("./lib/templates/layout.html", "./lib/templates/editList.html")
	//t, _ := template.ParseFiles("./lib/templates/editList.html")
	t.ExecuteTemplate(w, "layout1", map[string]string{
		"listID":   listID,
		"listName": listTitle,
	})
}

func modifyExistingListForm(w http.ResponseWriter, r *http.Request) {
	listTitle := r.FormValue("listTitle")
	listID, err := strconv.Atoi(r.FormValue("listID"))
	fmt.Println(listTitle, listID)
	if err != nil {
		error_controller.ErrorResponse(w, "Unable to render this page", http.StatusInternalServerError)
	}

	db_controller.EditEntry("lists", "title", listTitle, listID)

	http.Redirect(w, r, "/lists", http.StatusSeeOther)
	// message := fmt.Sprintf("Successfully updated list (ID: %v)", affectedID)
	// w.Write([]byte(message))
}

func deleteListForm(w http.ResponseWriter, r *http.Request) {
	listID, err := strconv.Atoi(r.URL.Path[len("/lists/delete"):])
	fmt.Println(listID)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	lists, err := db_controller.SelectAllLists()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	var currentList db_controller.List
	for _, list := range lists {
		if list.ID == listID {
			currentList.ID = list.ID
			currentList.Title = list.Title
		}
	}
	title, err := db_controller.DeleteList(currentList)
	if err != nil {
		fmt.Printf("Error deleting list %v: %s", title, err)
		return
	}
	http.Redirect(w, r, "/lists", http.StatusSeeOther)
}
