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
	"io/ioutil"
	"net/http"
	"time"
)

// Login authenticate an User
func Login(w http.ResponseWriter, r *http.Request) {
	prommetrics.PromRequestsCurrent.Inc()
	now := time.Now()

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

	db, erro := database.Connect()
	if erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}

	repository := repositories.NewUsersRepository(db)
	userFromDB, erro := repository.SearchByEmail(user.Email)
	defer db.Close()
	if erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}

	if erro := security.ValidatePass(userFromDB.Pass, user.Pass); erro != nil {
		responses.Erro(now, w, http.StatusUnauthorized, errors.New("incorrect password"))
		return
	}

	token, erro := authentication.GenerateToken(userFromDB.ID)
	if erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}

	w.Write([]byte(token))
}
