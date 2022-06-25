package controllers

import (
	"api/src/authentication"
	"api/src/database"
	"api/src/models"
	"api/src/repositories"
	"api/src/responses"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CreatePublication create a Publication in database
func CreatePublication(w http.ResponseWriter, r *http.Request) {

	userID, erro := authentication.ExtractUserID(r)
	if erro != nil {
		responses.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	body, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		responses.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var publication models.Publication
	if erro := json.Unmarshal(body, &publication); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.Connect()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	publication.AuthorID = userID

	if erro := publication.Prepare(); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	repository := repositories.NewPublicationRepository(db)
	publication.ID, erro = repository.Create(publication)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	responses.JSON(w, http.StatusCreated, publication)
}

// GetPublications return publications in feed
func GetPublications(w http.ResponseWriter, r *http.Request) {

}

// GetPublication return one Publication
func GetPublication(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	publicationID, erro := strconv.ParseUint(params["publicationID"], 10, 64)
	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.Connect()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	repository := repositories.NewPublicationRepository(db)
	publication, erro := repository.SearchByID(publicationID)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	responses.JSON(w, http.StatusOK, publication)
}

// UpdatePublication updates one Publication
func UpdatePublication(w http.ResponseWriter, r *http.Request) {

}

// DeletePublication delete one Publication
func DeletePublication(w http.ResponseWriter, r *http.Request) {

}
