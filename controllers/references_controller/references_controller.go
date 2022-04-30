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
var formOutputRoute = regexp.MustCompile("output")

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
		} else if formOutputRoute.FindStringSubmatch(r.URL.Path) != nil {
			outputSelected(w, r)
		} else {
			viewAllRefs(w, r)
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET method is supported for this route.")
	}
}
