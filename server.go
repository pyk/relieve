package main

import (
	"encoding/json"
	"errors"
	"fmt"
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

type StatusRequest struct {
	Status bool `json:"status"`
}

var successReq = StatusRequest{true}
var failureReq = StatusRequest{false}

// ApiHandler global API mux
type ApiHandler func(w http.ResponseWriter, r *http.Request) *apiError

func (api ApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// for response need an array. requested by @fachrian
	var errs []*apiError

	// add header on every response
	w.Header().Add("Server", "Relieve by Sunday Code")
	w.Header().Add("X-Wisdom-Media-Type", "relieve.v0")
	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	// if handler return an &apiError
	err := api(w, r)
	errs = append(errs, err)
	if err != nil {
		// http log
		log.Printf("%s %s %s [%s] %s", r.RemoteAddr, r.Method, r.URL, err.Tag, err.Error)

		// response proper http status code
		w.WriteHeader(err.Code)

		// response JSON
		resp := json.NewEncoder(w)
		err_json := resp.Encode(errs)
		if err_json != nil {
			log.Println("Encode JSON for error response was failed.")

			return
		}

		return
	}

	// http log
	log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
}

// get information about reliever with specified ID
func relieverHandler(w http.ResponseWriter, r *http.Request) *apiError {
	// get reliever ID, if not specified then return an bad request status
	// GET /v0/reliever?reliever_id=ID
	if r.Method == "GET" {
		relieverID := r.FormValue("reliever_id")
		if relieverID == "" {
			return &apiError{
				"relieverHandler GET",
				errors.New("relieverHandler reliever_id not specified"),
				"reliever_id not specified.",
				http.StatusBadRequest,
			}
		}
		rl, err := db.GetPsikologByID(relieverID)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				return &apiError{
					"relieverHandler GET",
					err,
					"reliever_id not exists",
					http.StatusBadRequest,
				}
			}
			return &apiError{
				"relieverHandler GET",
				err,
				"Bad request",
				http.StatusBadRequest,
			}
		}
		var rls []database.Reliever
		enc := json.NewEncoder(w)
		rls = append(rls, rl)
		err = enc.Encode(rls)
		if err != nil {
			return &apiError{
				"relieverHandler Encode",
				err,
				"Encoding JSON failed",
				http.StatusInternalServerError,
			}
		}

		return nil
	}
	return nil
}

// checkWisdomHandler handle a GET request to check status of wisdom.
// if user already give a piskolog wisdom point then return
// {"wisdom_point_status":"false"}
// otherwise return
// {"wisdom_point_status":"true"}
// GET /v0/checkwisdom?psikolog_id=12&user_id=1
func checkWisdomHandler(w http.ResponseWriter, r *http.Request) *apiError {
	// response should be an array.
	var status []database.WisdomPointStatus

	psikologID := r.FormValue("psikolog_id")
	userID := r.FormValue("user_id")
	if psikologID != "" && userID != "" {
		s, err := db.CheckWisdomPoint(userID, psikologID)
		if err != nil {
			return &apiError{
				"checkWisdomHandler CheckWisdomPoint",
				err,
				"Internal server error",
				http.StatusInternalServerError,
			}
		}

		enc := json.NewEncoder(w)
		status = append(status, s)
		err = enc.Encode(status)
		if err != nil {
			return &apiError{
				"checkWisdomHandler CheckWisdomPoint encode JSON",
				err,
				"Internal server error",
				http.StatusInternalServerError,
			}
		}
		return nil
	}

	fmt.Fprintln(w, "200 OK")
	return nil
}

// wisdomHandler handle a GET & POST request to get wisdom point of psikolog
// and insert new wisdom point.
// GET /v0/wisdom?psikolog_id=12
// POST /v0/wisdom ; with data: {"user_id": 1, "psikolog_id": 1}
func wisdomHandler(w http.ResponseWriter, r *http.Request) *apiError {
	// insert new wisdom point
	if r.Method == "POST" {
		var wp *database.WisdomPoint
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&wp)
		if err != nil {
			return &apiError{
				"wisdomHandler Decode",
				err,
				"Internal server error",
				http.StatusInternalServerError,
			}
		}

		// insert data to database
		err = db.InsertWisdomPoint(wp)
		if err != nil {
			if err.Error() == "pq: duplicate key value violates unique constraint \"wisdom_points_wisdom_user_id_wisdom_psikolog_id_key\"" {
				return &apiError{
					"wisdomHandler db.InsertWisdomPoint",
					err,
					"Bad request. Record exists.",
					http.StatusBadRequest,
				}
			}
			return &apiError{
				"wisdomHandler db.InsertWisdomPoint",
				err,
				"Internal server error",
				http.StatusInternalServerError,
			}
		}
		return nil
	}

	// get psikolog wisdom point
	psikologID := r.FormValue("psikolog_id")
	// response should be an array
	var psikolog_points []database.PsikologPoint
	if psikologID != "" {
		wp, err := db.GetWisdomPointByID(psikologID)
		if err != nil {
			return &apiError{
				"wisdomHandler GetWisdomPointById",
				err,
				"Internal server error",
				http.StatusInternalServerError,
			}
		}

		psikolog_points = append(psikolog_points, wp)
		enc := json.NewEncoder(w)
		err = enc.Encode(psikolog_points)
		if err != nil {
			return &apiError{
				"wisdomHandler GetWisdomPointById encode JSON",
				err,
				"Internal server error",
				http.StatusInternalServerError,
			}
		}
		return nil
	}

	return nil
}

// reportHandler handle report endpoint
func reportHandler(w http.ResponseWriter, r *http.Request) *apiError {
	if r.Method != "POST" {
		http.Redirect(w, r, "https://sundaycode.co", 302)
		return nil
	}

	// read data from POST request and decode data to *database.Report type
	var rp *database.Report
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&rp)
	if err != nil {
		return &apiError{
			"commentHandler Decode",
			err,
			"Internal server error",
			http.StatusInternalServerError,
		}
	}

	// insert data to database
	err = db.InsertReport(rp)
	if err != nil {
		return &apiError{
			"commentHandler db.InsertComment",
			err,
			"Internal server error",
			http.StatusInternalServerError,
		}
	}

	return nil
}

// commentHandler handle comment endpoint
func commentHandler(w http.ResponseWriter, r *http.Request) *apiError {
	if r.Method != "POST" {
		http.Redirect(w, r, "https://sundaycode.co", 302)
		return nil
	}

	// read data from POST request and decode data to *database.Comment type
	var c *database.Comment
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&c)
	if err != nil {
		return &apiError{
			"commentHandler Decode",
			err,
			"Internal server error",
			http.StatusInternalServerError,
		}
	}

	// insert data to database
	err = db.InsertComment(c)
	if err != nil {
		return &apiError{
			"commentHandler db.InsertComment",
			err,
			"Internal server error",
			http.StatusInternalServerError,
		}
	}

	return nil
}

// postHandler handle a post endpoint.
// * get list of posts by user with specified ID
// * post a curhat
func postHandler(w http.ResponseWriter, r *http.Request) *apiError {
	var posts []database.Post
	var err error

	// GET /v0/posts?user_id=ID
	// TODO: use regex to handle params. make sure input integer.
	// if inpit is plain string it make a server panic.
	if r.Method == "GET" {
		userID := r.FormValue("user_id")
		if userID == "" {
			return &apiError{
				"postHandler GET",
				errors.New("postHandler user_id not specified"),
				"user_id not specified.",
				http.StatusBadRequest,
			}
		}
		posts, err = db.GetAllPostsByUserID(userID)
		if err != nil {
			if err.Error() == "cannot found a list of posts" {
				return &apiError{
					"postHandler GET",
					err,
					err.Error(),
					http.StatusNotFound,
				}
			}
			return &apiError{
				"postHandler GET",
				err,
				"Internal server error",
				http.StatusInternalServerError,
			}
		}
		enc := json.NewEncoder(w)
		err = enc.Encode(posts)
		if err != nil {
			return &apiError{
				"postHandler GetAllPosts encode JSON",
				err,
				"Internal server error",
				http.StatusInternalServerError,
			}
		}

		return nil
	}

	// TODO: nyelesain endpoint post for POST
	// POST /v0/posts
	// data: user_id, psikolog_id, title, category, content
	if r.Method == "POST" {
		// save data from params
		userID := r.FormValue("user_id")
		psikologID := r.FormValue("psikolog_id")
		title := r.FormValue("title")
		category := r.FormValue("category")
		content := r.FormValue("content")

		// all params should not empty
		if userID == "" || psikologID == "" || title == "" || category == "" || content == "" {
			return &apiError{
				"postHandler POST",
				errors.New("data incomplete"),
				"POST data incomplete",
				http.StatusNotAcceptable,
			}
		}
		p := database.Post{
			UserId:     userID,
			PsikologId: psikologID,
			Title:      title,
			Category:   category,
			Content:    content,
		}

		// insert data to database
		// TODO: cari tahu kemungkinan error nya apa aja
		err = db.InsertPost(&p)
		if err != nil {
			return &apiError{
				"postHandler db.InsertPost",
				err,
				"Internal server error. Cannot insert data to database",
				http.StatusInternalServerError,
			}
		}

		// send a success message
		enc := json.NewEncoder(w)
		err = enc.Encode(successReq)
		if err != nil {
			return &apiError{
				"postHandler POST encode JSON",
				err,
				"Internal server error",
				http.StatusInternalServerError,
			}
		}

		return nil
	}

	return nil
}

// psikologHandler handle psikolog endpoint
func psikologHandler(w http.ResponseWriter, r *http.Request) *apiError {
	if r.Method != "POST" {
		http.Redirect(w, r, "https://sundaycode.co", 302)
		return nil
	}

	// read data from POST request and decode data to *database.Psikolog type
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
	var user *database.User
	if r.Method == "GET" {
		fmt.Fprintln(w, "200 OK")
		return nil
	}
	if r.Method == "POST" {
		// read data from POST request and decode data to User type

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

	// reliever handler
	r.Handle("/v0/reliever", ApiHandler(relieverHandler))

	// insert data to psikologs table
	// POST /v0/psikolog
	r.Handle("/v0/psikologs", ApiHandler(psikologHandler))

	// get & post a wisdom points
	r.Handle("/v0/wisdom", ApiHandler(wisdomHandler))
	r.Handle("/v0/checkwisdom", ApiHandler(checkWisdomHandler))

	// insert data to posts table
	// POST /v0/posts
	r.Handle("/v0/posts", ApiHandler(postHandler))

	// insert data to comments table
	// POST /v0/comments
	r.Handle("/v0/comments", ApiHandler(commentHandler))

	// insert data to reports table
	// POST /v0/reports
	r.Handle("/v0/reports", ApiHandler(reportHandler))

	// server listener
	http.Handle("/", r)
	log.Printf("Listening on :%s", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
