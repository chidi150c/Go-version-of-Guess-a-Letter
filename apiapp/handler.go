package apiapp

import (
	"fmt"
	"net/http"
	"strings"
	"user-apiv2/apigame"
	"user-apiv2/apiuser"
)

// Handler is a collection of all the service handlers.
type Handler struct {
	UserHandler *apiuser.UserHandler
	GameHandler *apigame.GameHandler
}

//initializies the Handler struct
func NewHandler(u *apiuser.UserHandler, g *apigame.GameHandler) *Handler {
	return &Handler{
		UserHandler: u,
		GameHandler: g,
	}
}

// ServeHTTP delegates a request to the appropriate subhandler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	if strings.HasPrefix(r.URL.Path, "/tools/asset/") {
		fmt.Println()
		http.StripPrefix("/tools/asset/", http.FileServer(http.Dir("./tools/asset/"))).ServeHTTP(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/game") {
		fmt.Printf("\n\n\n ************************ In handler for game ***********************  \n\n\n\n")
		fmt.Println()
		h.GameHandler.ServeHTTP(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/") {
		fmt.Println()
		h.UserHandler.ServeHTTP(w, r)
	} else {
		fmt.Println()
		http.NotFound(w, r)
	}
}
