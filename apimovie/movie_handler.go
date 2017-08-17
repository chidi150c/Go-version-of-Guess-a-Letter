package apimovie

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"user-apiv2/data"
	"user-apiv2/tools"

	"github.com/pkg/errors"
	"github.com/pressly/chi"
	"github.com/satori/go.uuid"
)

// MovieHandler represents an HTTP API handler for movies.
type MovieHandler struct {
	mux *chi.Mux
	//redirectURL string
	Session *Session
	//movie        *data.Movie
	Logger *log.Logger
}

// NewMovieHandler returns a new instance of MovieHandler.
//mv data.MovieService
func NewMovieHandler(s *Session) *MovieHandler {
	h := &MovieHandler{
		mux:     chi.NewRouter(),
		Logger:  log.New(os.Stderr, "", log.LstdFlags),
		Session: s,
	}

	h.mux.Get("/movies", s.Validate(h.listHandler))
	h.mux.Get("/movies/add", s.Validate(h.addFormHandler))
	h.mux.Post("/movies", s.Validate(h.createHandler))
	h.mux.Get("/movies/:id", s.Validate(h.detailHandler))
	h.mux.Post("/movies/:id", s.Validate(h.updateHandler))
	h.mux.Post("/movies/delete/:id", s.Validate(h.deleteHandler))
	h.mux.Get("/movies/edit/:id", s.Validate(h.editFormHandler))

	return h
}

func (h *MovieHandler) listHandler(w http.ResponseWriter, r *http.Request) {
	user, err := h.Session.UserFromRequest(r)
	if err != nil {
		tools.Error(w, err, http.StatusBadRequest, h.Logger)
		return
	}
	movies, err := h.Session.Movieservice.ListMovies()
	if err != nil {
		tools.Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}
	dat := struct {
		Movs        data.Muu //[]*data.Movie
		AddMovieURL string
	}{
		Movs:        data.Muu(movies),
		AddMovieURL: "/movies/add?redirect=/", //the resaon is for list page to know if movie is loged in
	}
	if err := tools.ListmoviesTmpl.Execute(w, r, dat, user, true); err != nil {
		tools.Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}
	return
}

func (h *MovieHandler) addFormHandler(w http.ResponseWriter, r *http.Request) {
	user, err := h.Session.UserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	//TODO ....
	if err := tools.EditTmpl.Execute(w, r, nil, user, true); err != nil {
		tools.Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}
	return
}

// handleGetMovie handles requests to create a new movie.
func (h *MovieHandler) createHandler(w http.ResponseWriter, r *http.Request) {
	user, err := h.Session.UserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	movie, err := h.movieFromForm(w, r)
	if err != nil {
		tools.Error(w, errors.Wrapf(err, "could not parse movie from form: %v", err), http.StatusInternalServerError, h.Logger)
		return
	}
	ctx := context.WithValue(context.Background(), AddMovieKey, user)
	_, err = h.Session.Movieservice.AddMovie(ctx, movie)
	if err != nil {
		tools.Error(w, errors.Wrapf(err, "could not save movie: %v", err), http.StatusInternalServerError, h.Logger)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/movies/%d", movie.ID), http.StatusFound)
	return
}

func (h *MovieHandler) detailHandler(w http.ResponseWriter, r *http.Request) {
	//check if user is authenticate before add movie service
	user, err := h.Session.UserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	ids := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(ids)
	dMov, err := h.Session.Movieservice.GetMovie(data.MovieID(id))
	d := struct {
		Mov    *data.Movie
		IsUser bool
	}{
		Mov:    dMov,
		IsUser: user.Username == dMov.OwnerID,
	}

	if err := tools.DetailTmpl.Execute(w, r, d, user, true); err != nil {
		tools.Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}
	return
}

func (h *MovieHandler) updateHandler(w http.ResponseWriter, r *http.Request) {
	//check if user is authenticate before add movie service
	user, err := h.Session.UserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	ids := chi.URLParam(r, "id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		tools.Error(w, errors.Wrapf(err, "bad movie id: %v", err), http.StatusInternalServerError, h.Logger)
		return
	}
	movie, err := h.movieFromForm(w, r)
	if err != nil {
		tools.Error(w, errors.Wrapf(err, "could not parse movie from form: %v", err), http.StatusInternalServerError, h.Logger)
		return
	}
	movie.ID = data.MovieID(id)
	ctx := context.WithValue(context.Background(), AddMovieKey, user)
	err = h.Session.Movieservice.UpdateMovie(ctx, movie)
	if err != nil {
		tools.Error(w, errors.Wrapf(err, "could not save movie: %v", err), http.StatusInternalServerError, h.Logger)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/movies/%d", movie.ID), http.StatusFound)
	return
}

// deleteHandler deletes a given movie.
func (h *MovieHandler) deleteHandler(w http.ResponseWriter, r *http.Request) {
	//check if user is authenticate before add movie service
	user, err := h.Session.UserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	ids := chi.URLParam(r, "id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		tools.Error(w, errors.Wrapf(err, "bad movie id: %v", err), http.StatusInternalServerError, h.Logger)
		return
	}
	ctx := context.WithValue(context.Background(), AddMovieKey, user)
	err = h.Session.Movieservice.DeleteMovie(ctx, data.MovieID(id))
	if err != nil {
		tools.Error(w, errors.Wrapf(err, "could not delete movie: %v", err), http.StatusInternalServerError, h.Logger)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
	return
}

func (h *MovieHandler) editFormHandler(w http.ResponseWriter, r *http.Request) {
	//check if user is authenticate before add movie service
	user, err := h.Session.UserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	// User is Logged in
	movie, err := h.movieFromRequest(r)
	if err != nil {
		tools.Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}
	if movie.OwnerID == user.Username {
		if err := tools.EditTmpl.Execute(w, r, movie, user, true); err != nil {
			tools.Error(w, err, http.StatusInternalServerError, h.Logger)
			return
		}
	}
	tools.Error(w, errors.New("User Attempting to Edit not Owned movie is prohibited"), 403, h.Logger)
	return
}

func (h *MovieHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

// handleGetMovie handles requests to create a new movie.

// movieFromRequest retrieves a movie from the database given a movie ID in the
// URL's path.
func (h *MovieHandler) movieFromRequest(r *http.Request) (*data.Movie, error) {
	ids := chi.URLParam(r, "id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		return nil, errors.Wrapf(err, "bad movie id: %v", err)
	}
	movie, err := h.Session.Movieservice.GetMovie(data.MovieID(id))
	if err != nil {
		return nil, errors.Wrapf(err, "could not find movie: %v", err)
	}
	return movie, nil
}

func uploadFileFromForm(w http.ResponseWriter, r *http.Request) (url string, err error) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("image")
	if err == http.ErrMissingFile {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	// random filename, retaining existing extension.
	name := uuid.NewV4().String() + path.Ext(handler.Filename)
	defer file.Close()
	f, err := os.OpenFile("./asset/public_images/"+name, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer f.Close()
	io.Copy(f, file)
	const publicURL = "/%s/%s"
	return fmt.Sprintf(publicURL, "asset/public_images", name), nil
}

func (h *MovieHandler) movieFromForm(w http.ResponseWriter, r *http.Request) (*data.Movie, error) {

	imageURL, err := uploadFileFromForm(w, r)
	if err != nil {
		return nil, errors.Wrapf(err, "could not upload file: %v", err)
	}
	if imageURL == "" {
		imageURL = r.FormValue("imageURL")
	}
	movie := &data.Movie{
		Name:           r.FormValue("name"),
		OwnerName:      r.FormValue("ownername"),
		Author:         r.FormValue("author"),
		OwnerID:        r.FormValue("ownerid"),
		Token:          r.FormValue("token"),
		ProductionDate: r.FormValue("productiondate"),
		Description:    r.FormValue("description"),
		ImageURL:       imageURL,
		Category:       r.FormValue("category"),
	}
	// If the form didn't carry the user information for the creator, populate it
	// from the currently logged in user (or mark as anonymous).
	if movie.OwnerName == "" {
		user, err := h.Session.UserFromRequest(r)
		if err != nil {
			return nil, err
		}
		if user != nil {
			// Logged in.
			movie.OwnerName = user.DisplayName
			movie.OwnerID = string(user.ID)
		} else {
			// Not logged in.
			movie.SetCreatorAnonymous()
		}
	}
	return movie, nil
}
