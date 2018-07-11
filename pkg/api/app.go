package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (theApp *App) Initialize(userDB, passwordDB, nameDB string) {
	connectionString := fmt.Sprintf("%s:%s@/%s", userDB, passwordDB, nameDB)

	var err error
	theApp.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	theApp.Router = mux.NewRouter()
	theApp.initializeRoutes()
}

func (theApp *App) initializeRoutes() {
	theApp.Router.HandleFunc("/photographer", theApp.createPhotographer).Methods("POST")
	theApp.Router.HandleFunc("/photographer/{username}", theApp.getPhotographer).Methods("GET")
	theApp.Router.HandleFunc("/photographer/{username}", theApp.updatePhotographer).Methods("PUT")
	theApp.Router.HandleFunc("/photographer/{username}", theApp.deletePhotographer).Methods("DELETE")
}

func (theApp *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, theApp.Router))
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, statusCode int, errorMessage string) {
	respondWithJSON(w, statusCode, map[string]string{"error": errorMessage})
}

func (theApp *App) createPhotographer(w http.ResponseWriter, r *http.Request) {
	var p photographer

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := p.createPhotographer(theApp.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, p)
}

func matchRegexUnderDashAlphanumeric(checked string) (bool, error) {
	regex, err := regexp.Compile("^[A-Za-z0-9][A-Za-z0-9_-]*$")

	return regex.MatchString(checked), err
}

func (theApp *App) getPhotographer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	matchRegex, err := matchRegexUnderDashAlphanumeric(vars["username"])
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !matchRegex {
		respondWithError(w, http.StatusBadRequest, "User should not be exist")
		return
	}

	p := photographer{Username: vars["username"]}
	if err := p.getPhotographer(theApp.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func (theApp *App) updatePhotographer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	matchRegex, err := matchRegexUnderDashAlphanumeric(vars["username"])
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !matchRegex {
		respondWithError(w, http.StatusBadRequest, "User should not be exist")
		return
	}

	var p photographer
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	p.Username = vars["username"]
	if err := p.updatePhotographer(theApp.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, p)
}

func (theApp *App) deletePhotographer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	matchRegex, err := matchRegexUnderDashAlphanumeric(vars["username"])
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !matchRegex {
		respondWithError(w, http.StatusBadRequest, "User should not be exist")
		return
	}

	p := photographer{Username: vars["username"]}

	if err := p.deletePhotographer(theApp.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}
