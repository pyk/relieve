package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	// "time"

	"github.com/gorilla/mux"
	"github.com/pyk/relieve/database"
)

var (
	PORT = os.Getenv("PORT")
)

var (
	db *database.Database
)

// apiError define structure of API error
type apiError struct {
	Tag     string `json:"-"`
	Error   error  `json:"-"`
	Message string `json:"error"`
	Code    int    `json:"code"`
}

// ApiHandler global API mux
type ApiHandler func(w http.ResponseWriter, r *http.Request) *apiError

func (api ApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// add header on every response
	w.Header().Add("Server", "Relieve by Sunday Code")
	w.Header().Add("X-Wisdom-Media-Type", "relieve.v0")
	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	// if handler return an &apiError
	err := api(w, r)
	if err != nil {
		// http log
		log.Printf("%s %s %s [%s] %s", r.RemoteAddr, r.Method, r.URL, err.Tag, err.Error)

		// response proper http status code
		w.WriteHeader(err.Code)

		// response JSON
		resp := json.NewEncoder(w)
		err_json := resp.Encode(err)
		if err_json != nil {
			log.Println("Encode JSON for error response was failed.")

			return
		}

		return
	}

	// http log
	log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
}

// psikologHandler handle psikolog endpoint
func psikologHandler(w http.ResponseWriter, r *http.Request) *apiError {
	if r.Method != "POST" {
		http.Redirect(w, r, "https://sundaycode.co", 302)
		return nil
	}

	// read data from POST request and decode data to User type
	var p *database.Psikolog
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		return &apiError{
			"psikologHandler Decode",
			err,
			"Internal server error",
			http.StatusInternalServerError,
		}
	}

	// insert data to database
	err = db.InsertPsikolog(p)
	if err != nil {
		return &apiError{
			"psikologHandler db.InsertPsikolog",
			err,
			"Internal server error",
			http.StatusInternalServerError,
		}
	}

	return nil
}

// userHandler handle user endpoint
func userHandler(w http.ResponseWriter, r *http.Request) *apiError {
	if r.Method != "POST" {
		http.Redirect(w, r, "https://sundaycode.co", 302)
		return nil
	}

	// read data from POST request and decode data to User type
	var user *database.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		return &apiError{
			"usersHandler Decode",
			err,
			"Internal server error",
			http.StatusInternalServerError,
		}
	}

	// insert data to database
	err = db.InsertUser(user)
	if err != nil {
		return &apiError{
			"usersHandler db.InsertUser",
			err,
			"Internal server error",
			http.StatusInternalServerError,
		}
	}

	return nil
}

// notFoundHandler handle a not found response
func notFoundHandler(w http.ResponseWriter, r *http.Request) *apiError {
	return &apiError{
		"notFoundHandler",
		errors.New("Not Found"),
		"Not Found",
		http.StatusNotFound,
	}
}

// indexHandler handle request to '/'
// redirect to github pages
func indexHandler(w http.ResponseWriter, r *http.Request) *apiError {
	http.Redirect(w, r, "https://sundaycode.co", 302)
	return nil
}

func main() {
	var err error
	db, err = database.New()
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	// index handler doesn't need database utils
	r.Handle("/", ApiHandler(indexHandler))

	// not found handler
	r.NotFoundHandler = ApiHandler(notFoundHandler)

	// insert data to users table
	// POST /v0/users
	r.Handle("/v0/users", ApiHandler(userHandler))

	// insert data to psikologs table
	// POST /v0/psikolog
	r.Handle("/v0/psikologs", ApiHandler(psikologHandler))

	// server listener
	http.Handle("/", r)
	log.Printf("Listening on :%s", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
