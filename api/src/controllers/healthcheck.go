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
	"api/src/config"
	"api/src/database"
	"api/src/prommetrics"
	"api/src/repositories"
	"api/src/responses"
	"errors"
	"net/http"
	"time"
)

// Ready validates if our API is live and can process the requests received
func Live(w http.ResponseWriter, r *http.Request) {
	prommetrics.PromRequestsCurrent.Inc()
	now := time.Now()

	db, erro := database.Connect()
	if erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}
	repository := repositories.NewHealthcheckRepository(db)
	if erro := repository.SimulateDatabaseInsert(); erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	responses.JSON(now, w, http.StatusOK, nil)
}

// Ready validates if our API is ready to receive network connection and provide his main functionality
func Ready(w http.ResponseWriter, r *http.Request) {
	prommetrics.PromRequestsCurrent.Inc()
	now := time.Now()

	var hosts = []string{
		config.DatabaseHost,
	}

	db, erro := database.Connect()
	if erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}
	repository := repositories.NewHealthcheckRepository(db)

	for _, host := range hosts {
		if erro := repository.DNSResolver(host); erro != nil {
			responses.Erro(now, w, http.StatusInternalServerError, errors.New(host))
			return
		}
	}

	if erro := repository.PingDatabase(); erro != nil {
		responses.Erro(now, w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	responses.JSON(now, w, http.StatusOK, nil)
}
