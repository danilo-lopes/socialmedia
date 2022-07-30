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
	"api/src/security"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// CreateUser creates a new "User" in database
func CreateUser(w http.ResponseWriter, r *http.Request) {
	prommetrics.PromRequestsCurrent.Inc()
	now := time.Now()
	httpDuration := time.Since(now)

	body, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		responses.Erro(now, w, http.StatusUnprocessableEntity, erro)
		return
	}

	var user models.User
	if erro = json.Unmarshal(body, &user); erro != nil {
		responses.Erro(now, w, http.StatusBadRequest, erro)
		return
	}

	if erro := user.Prepare("registration"); erro != nil {
		responses.Erro(now, w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.Connect()
	if erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}

	repository := repositories.NewUsersRepository(db)
	user.ID, erro = repository.Create(user)
	defer db.Close()
	if erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}
	prommetrics.PromCountCreatedUsers.Inc()
	prommetrics.PromTimeTookToCreateUser.WithLabelValues(fmt.Sprintf("%d", http.StatusOK)).Observe(httpDuration.Seconds())

	responses.JSON(now, w, http.StatusCreated, user)
}

// GetUsers with given "nick" or "email", will return all "Users" from database
func GetUsers(w http.ResponseWriter, r *http.Request) {
	prommetrics.PromRequestsCurrent.Inc()
	now := time.Now()

	nameOrNick := strings.ToLower(
		r.URL.Query().Get("user"),
	)

	db, erro := database.Connect()
	if erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}

	repository := repositories.NewUsersRepository(db)
	users, erro := repository.Search(nameOrNick)
	defer db.Close()
	if erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(now, w, http.StatusOK, users)
}

// GetUser return specific "User" from database
func GetUser(w http.ResponseWriter, r *http.Request) {
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

	repository := repositories.NewUsersRepository(db)
	user, erro := repository.SearchByID(userID)
	defer db.Close()
	if erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(now, w, http.StatusOK, user)
}

// UpdateUser upadate "User" attributes in database
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	prommetrics.PromRequestsCurrent.Inc()
	now := time.Now()

	params := mux.Vars(r)
	userID, erro := strconv.ParseUint(params["userID"], 10, 64)
	if erro != nil {
		responses.Erro(now, w, http.StatusBadRequest, erro)
		return
	}

	userIDInsideToken, erro := authentication.ExtractUserID(r)
	if erro != nil {
		responses.Erro(now, w, http.StatusUnauthorized, erro)
		return
	}

	if userID != userIDInsideToken {
		responses.Erro(now, w, http.StatusForbidden, errors.New("permission denied"))
		return
	}

	body, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		responses.Erro(now, w, http.StatusUnprocessableEntity, erro)
		return
	}

	var user models.User
	if erro := json.Unmarshal(body, &user); erro != nil {
		responses.Erro(now, w, http.StatusBadRequest, erro)
		return
	}

	if erro := user.Prepare("edit"); erro != nil {
		responses.Erro(now, w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.Connect()
	if erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}

	repository := repositories.NewUsersRepository(db)
	if erro := repository.Update(userID, user); erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	responses.JSON(now, w, http.StatusNoContent, nil)
}

// DeleteUser deletes a "User" in database
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	prommetrics.PromRequestsCurrent.Inc()
	now := time.Now()
	httpDuration := time.Since(now)

	params := mux.Vars(r)
	userID, erro := strconv.ParseUint(params["userID"], 10, 64)
	if erro != nil {
		responses.Erro(now, w, http.StatusBadRequest, erro)
		return
	}

	userIDInsideToken, erro := authentication.ExtractUserID(r)
	if erro != nil {
		responses.Erro(now, w, http.StatusUnauthorized, erro)
		return
	}

	if userID != userIDInsideToken {
		responses.Erro(now, w, http.StatusForbidden, errors.New("permission denied"))
	}

	db, erro := database.Connect()
	if erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}

	repository := repositories.NewUsersRepository(db)
	if erro := repository.Delete(userID); erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()
	prommetrics.PromCountDeletedUsers.Inc()
	prommetrics.PromTimeTookToDeleteUser.WithLabelValues(fmt.Sprintf("%d", http.StatusOK)).Observe(httpDuration.Seconds())

	responses.JSON(now, w, http.StatusNoContent, nil)
}

// FollowUser permits an "User" to "Follow" another "User"
func FollowUser(w http.ResponseWriter, r *http.Request) {
	prommetrics.PromRequestsCurrent.Inc()
	now := time.Now()

	followerID, erro := authentication.ExtractUserID(r)
	if erro != nil {
		responses.Erro(now, w, http.StatusUnauthorized, erro)
		return
	}

	params := mux.Vars(r)
	userID, erro := strconv.ParseUint(params["userID"], 10, 64)
	if erro != nil {
		responses.Erro(now, w, http.StatusBadRequest, erro)
		return
	}

	if followerID == userID {
		responses.Erro(now, w, http.StatusForbidden, errors.New("is not possible to follow itself"))
		return
	}

	db, erro := database.Connect()
	if erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}

	repository := repositories.NewUsersRepository(db)
	if erro := repository.Follow(userID, followerID); erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	responses.JSON(now, w, http.StatusNoContent, nil)
}

// UnFollowUser permits an "User" to "Unfollow" another "User"
func UnFollowUser(w http.ResponseWriter, r *http.Request) {
	prommetrics.PromRequestsCurrent.Inc()
	now := time.Now()

	followerID, erro := authentication.ExtractUserID(r)
	if erro != nil {
		responses.Erro(now, w, http.StatusUnauthorized, erro)
		return
	}

	params := mux.Vars(r)
	userID, erro := strconv.ParseUint(params["userID"], 10, 64)
	if erro != nil {
		responses.Erro(now, w, http.StatusBadRequest, erro)
		return
	}

	if followerID == userID {
		responses.Erro(now, w, http.StatusForbidden, errors.New("is not possible to unfollow itself"))
		return
	}

	db, erro := database.Connect()
	if erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}

	repository := repositories.NewUsersRepository(db)
	if erro := repository.UnFollow(userID, followerID); erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	responses.JSON(now, w, http.StatusNoContent, nil)
}

// GetFollowers return all "Followers" from specific "User"
func GetFollowers(w http.ResponseWriter, r *http.Request) {
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

	repository := repositories.NewUsersRepository(db)
	followers, erro := repository.GetFollowers(userID)
	defer db.Close()
	if erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(now, w, http.StatusOK, followers)
}

// GetFollowing return all "Users" a "User" is "Following"
func GetFollowing(w http.ResponseWriter, r *http.Request) {
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

	repository := repositories.NewUsersRepository(db)
	users, erro := repository.GetFollowing(userID)
	defer db.Close()
	if erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(now, w, http.StatusOK, users)
}

// UpdatePass update an "User" password
func UpdatePass(w http.ResponseWriter, r *http.Request) {
	prommetrics.PromRequestsCurrent.Inc()
	now := time.Now()

	userIDInsideToken, erro := authentication.ExtractUserID(r)
	if erro != nil {
		responses.Erro(now, w, http.StatusUnauthorized, erro)
		return
	}

	params := mux.Vars(r)
	userID, erro := strconv.ParseUint(params["userID"], 10, 64)
	if erro != nil {
		responses.Erro(now, w, http.StatusBadRequest, erro)
		return
	}

	if userIDInsideToken != userID {
		responses.Erro(now, w, http.StatusForbidden, errors.New("is only allowed to update your own password"))
		return
	}

	body, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		responses.Erro(now, w, http.StatusBadRequest, erro)
		return
	}

	var pass models.Pass
	if erro := json.Unmarshal(body, &pass); erro != nil {
		responses.Erro(now, w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.Connect()
	if erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}

	repository := repositories.NewUsersRepository(db)
	userPassHash, erro := repository.GetUserPass(userID)
	if erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
	}

	if erro := security.ValidatePass(userPassHash, pass.Current); erro != nil {
		responses.Erro(now, w, http.StatusUnauthorized, errors.New("the password is incorrect"))
		return
	}

	hashedPass, erro := security.Hash(pass.New)
	if erro != nil {
		responses.Erro(now, w, http.StatusBadRequest, erro)
		return
	}

	if erro := repository.UpadateUserPass(userID, string(hashedPass)); erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	responses.JSON(now, w, http.StatusNoContent, nil)
}

// LikedPublications return all publications a user liked
func LikedPublications(w http.ResponseWriter, r *http.Request) {
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

	repository := repositories.NewUsersRepository(db)
	publications, erro := repository.LikedPublications(userID)
	defer db.Close()
	if erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(now, w, http.StatusOK, publications)
}
