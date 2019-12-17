package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type Person struct {
	Name string
	Age  int
	Sex  string
}

type response struct {
	Data interface{} `json:"data"`
}

//AuthorizeRequest Middleware validates requests.
func AuthorizeRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user, pass, ok := r.BasicAuth()

		fmt.Println("username: ", user)
		fmt.Println("password: ", pass)

		if !ok || !checkUsernameAndPassword(user, pass) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Please enter your username and password for this site"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func HelloHandler(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {

	var p Person

	p = Person{
		"krishna",
		40,
		"male",
	}

	_, err := json.Marshal(p)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "Unable to json encode Person")
	}
	return p, http.StatusOK, nil
}

func ResponseHandler(h func(http.ResponseWriter, *http.Request) (interface{}, int, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, status, err := h(w, r)

		if err != nil {
			data = err.Error()
		}

		w.Header().Set("Content-Type", "application/json")
		enableCors(&w)
		w.WriteHeader(status)

		if data != nil {
			err = json.NewEncoder(w).Encode(response{Data: data})
		}
	})
}

func SetRoutes(router *mux.Router) *mux.Router {
	personRouter := mux.NewRouter()
	personRouter.Handle("/person", ResponseHandler(HelloHandler))
	router.PathPrefix("/person").Handler(AuthorizeRequest(personRouter))
	return router
}

func main() {
	router := mux.NewRouter()
	router = SetRoutes(router)
	//handler := cors.Default().Handler(router)
	http.ListenAndServe(":80", router)
}

func checkUsernameAndPassword(username, password string) bool {
	return username == "abc" && password == "123"
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET,POST, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
