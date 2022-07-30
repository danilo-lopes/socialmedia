/*
Copyright 2022 Danilo S. Lopes.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at:

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"api/src/authentication"
	"api/src/database"
	"api/src/models"
	"api/src/prommetrics"
	"api/src/repositories"
	"api/src/responses"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// CreatePublication create a Publication in database
func CreatePublication(w http.ResponseWriter, r *http.Request) {
	prommetrics.PromRequestsCurrent.Inc()
	now := time.Now()
	httpDuration := time.Since(now)

	userID, erro := authentication.ExtractUserID(r)
	if erro != nil {
		responses.Erro(now, w, http.StatusUnauthorized, erro)
		return
	}

	body, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		responses.Erro(now, w, http.StatusUnprocessableEntity, erro)
		return
	}

	var publication models.Publication
	if erro := json.Unmarshal(body, &publication); erro != nil {
		responses.Erro(now, w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.Connect()
	if erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}

	publication.AuthorID = userID

	if erro := publication.Prepare(); erro != nil {
		responses.Erro(now, w, http.StatusBadRequest, erro)
		return
	}

	repository := repositories.NewPublicationRepository(db)
	publication.ID, erro = repository.Create(publication)
	defer db.Close()
	if erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}
	prommetrics.PromCountNewPublication.Inc()
	prommetrics.PromTimeTookToCreatePublication.WithLabelValues(fmt.Sprintf("%d", http.StatusOK)).Observe(httpDuration.Seconds())

	responses.JSON(now, w, http.StatusCreated, publication)
}

// GetPublications return publications in feed
func GetPublications(w http.ResponseWriter, r *http.Request) {
	prommetrics.PromRequestsCurrent.Inc()
	now := time.Now()

	userID, erro := authentication.ExtractUserID(r)
	if erro != nil {
		responses.Erro(now, w, http.StatusUnauthorized, erro)
		return
	}

	db, erro := database.Connect()
	if erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}

	repository := repositories.NewPublicationRepository(db)
	publications, erro := repository.Get(userID)
	defer db.Close()
	if erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(now, w, http.StatusOK, publications)
}

// GetPublication return one Publication
func GetPublication(w http.ResponseWriter, r *http.Request) {
	prommetrics.PromRequestsCurrent.Inc()
	now := time.Now()

	params := mux.Vars(r)
	publicationID, erro := strconv.ParseUint(params["publicationID"], 10, 64)
	if erro != nil {
		responses.Erro(now, w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.Connect()
	if erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}

	repository := repositories.NewPublicationRepository(db)
	publication, erro := repository.SearchByID(publicationID)
	if erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	responses.JSON(now, w, http.StatusOK, publication)
}

// UpdatePublication updates one Publication
func UpdatePublication(w http.ResponseWriter, r *http.Request) {
	prommetrics.PromRequestsCurrent.Inc()
	now := time.Now()

	userID, erro := authentication.ExtractUserID(r)
	if erro != nil {
		responses.Erro(now, w, http.StatusUnauthorized, erro)
		return
	}

	params := mux.Vars(r)
	publicationID, erro := strconv.ParseUint(params["publicationID"], 10, 64)
	if erro != nil {
		responses.Erro(now, w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.Connect()
	if erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}

	repository := repositories.NewPublicationRepository(db)

	publicationDatabase, erro := repository.SearchByID(publicationID)
	defer db.Close()
	if erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}

	if publicationDatabase.AuthorID != userID {
		responses.Erro(now, w, http.StatusForbidden, errors.New("is not possible to update publications from another user"))
		return
	}

	body, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		responses.Erro(now, w, http.StatusUnprocessableEntity, erro)
		return
	}

	var publication models.Publication
	if erro := json.Unmarshal(body, &publication); erro != nil {
		responses.Erro(now, w, http.StatusBadRequest, erro)
		return
	}

	if erro := publication.Prepare(); erro != nil {
		responses.Erro(now, w, http.StatusBadRequest, erro)
		return
	}

	if erro := repository.Update(publication.ID, publication); erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(now, w, http.StatusNoContent, nil)
}

// DeletePublication delete one Publication
func DeletePublication(w http.ResponseWriter, r *http.Request) {
	prommetrics.PromRequestsCurrent.Inc()
	now := time.Now()
	httpDuration := time.Since(now)

	userID, erro := authentication.ExtractUserID(r)
	if erro != nil {
		responses.Erro(now, w, http.StatusUnauthorized, erro)
		return
	}

	params := mux.Vars(r)
	publicationID, erro := strconv.ParseUint(params["publicationID"], 10, 64)
	if erro != nil {
		responses.Erro(now, w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.Connect()
	if erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}

	repository := repositories.NewPublicationRepository(db)
	publicationDatabase, erro := repository.SearchByID(publicationID)
	defer db.Close()
	if erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}

	if publicationDatabase.AuthorID != userID {
		responses.Erro(now, w, http.StatusForbidden, errors.New("is not possible to delete publications from another user"))
		return
	}

	if erro := repository.Delete(publicationID); erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}
	prommetrics.PromCountDeletePublication.Inc()
	prommetrics.PromTimeTookToDeletePublication.WithLabelValues(fmt.Sprintf("%d", http.StatusOK)).Observe(httpDuration.Seconds())

	responses.JSON(now, w, http.StatusNoContent, nil)
}

// GetUserPublications return all user publications
func GetUserPublications(w http.ResponseWriter, r *http.Request) {
	prommetrics.PromRequestsCurrent.Inc()
	now := time.Now()

	params := mux.Vars(r)
	userID, erro := strconv.ParseUint(params["userID"], 10, 64)
	if erro != nil {
		responses.Erro(now, w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.Connect()
	if erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}

	repository := repositories.NewPublicationRepository(db)
	publications, erro := repository.GetByUser(userID)
	defer db.Close()
	if erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(now, w, http.StatusOK, publications)
}

// LikePublication likes an publication
func LikePublication(w http.ResponseWriter, r *http.Request) {
	prommetrics.PromRequestsCurrent.Inc()
	now := time.Now()

	params := mux.Vars(r)
	publicationID, erro := strconv.ParseUint(params["publicationID"], 10, 64)
	if erro != nil {
		responses.Erro(now, w, http.StatusBadRequest, erro)
		return
	}

	likerID, erro := authentication.ExtractUserID(r)
	if erro != nil {
		responses.Erro(now, w, http.StatusUnauthorized, erro)
		return
	}

	db, erro := database.Connect()
	if erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}

	repository := repositories.NewPublicationRepository(db)
	if erro := repository.LikePublication(publicationID, likerID); erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	responses.JSON(now, w, http.StatusNoContent, nil)
}

// UnLikePublication unlikes an publication
func UnLikePublication(w http.ResponseWriter, r *http.Request) {
	prommetrics.PromRequestsCurrent.Inc()
	now := time.Now()

	params := mux.Vars(r)
	publicationID, erro := strconv.ParseUint(params["publicationID"], 10, 64)
	if erro != nil {
		responses.Erro(now, w, http.StatusBadRequest, erro)
		return
	}

	unLikerID, erro := authentication.ExtractUserID(r)
	if erro != nil {
		responses.Erro(now, w, http.StatusUnauthorized, erro)
		return
	}

	db, erro := database.Connect()
	if erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}

	repository := repositories.NewPublicationRepository(db)
	if erro := repository.UnLikePublication(publicationID, unLikerID); erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	responses.JSON(now, w, http.StatusNoContent, nil)
}

// GetLikers return all users who liked an publication
func GetLikers(w http.ResponseWriter, r *http.Request) {
	prommetrics.PromRequestsCurrent.Inc()
	now := time.Now()

	params := mux.Vars(r)
	publicationID, erro := strconv.ParseUint(params["publicationID"], 10, 64)
	if erro != nil {
		responses.Erro(now, w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.Connect()
	if erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}

	repository := repositories.NewPublicationRepository(db)
	users, erro := repository.GetLikers(publicationID)
	defer db.Close()
	if erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(now, w, http.StatusOK, users)
}
