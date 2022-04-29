package listscontroller

import (
	"fmt"
	"net/http"
	"regexp"
)

var apiRoute = regexp.MustCompile("api")
var formEditRoute = regexp.MustCompile("edit")
var formNewRoute = regexp.MustCompile("new")
var formDeleteRoute = regexp.MustCompile("delete")

func ListHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if formNewRoute.FindStringSubmatch(r.URL.Path) != nil {
			showNewListForm(w, r)
		} else if formEditRoute.FindStringSubmatch(r.URL.Path) != nil {
			showEditListForm(w, r)
		} else if formDeleteRoute.FindStringSubmatch(r.URL.Path) != nil {
			fmt.Println("its a get (delete)")
			deleteListForm(w, r)
		} else {
			showSingleList(w, r)
		}
		return
	case "POST":
		if apiRoute.FindStringSubmatch(r.URL.Path) != nil {
			addListAPI(w, r)
		} else if formNewRoute.FindStringSubmatch(r.URL.Path) != nil {
			formPostNewList(w, r)
			// can't use PUT in HTML form
		} else if formEditRoute.FindStringSubmatch(r.URL.Path) != nil {
			fmt.Println("attempting edit")
			modifyExistingListForm(w, r)
		}
		return
	case "PUT":
		if apiRoute.FindStringSubmatch(r.URL.Path) != nil {
			modifyExistingListAPI(w, r)
		}
		return
	case "DELETE":
		deleteListAPI(w, r)
	default:
		fmt.Fprintf(w, "Sorry, only GET and DELETE methods supported for this route.")
	}
}

func ListsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if apiRoute.FindStringSubmatch(r.URL.Path) != nil {
			sendAllListsJSON(w, r)
		} else {
			renderAllLists(w, r)
		}
		return
	default:
		fmt.Fprintf(w, "Sorry, only GET method is supported for this route.")
	}
}
