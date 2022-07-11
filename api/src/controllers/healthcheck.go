package controllers

import (
	"api/src/config"
	"api/src/database"
	"api/src/repositories"
	"api/src/responses"
	"errors"
	"net/http"
)

// Ready validates if our API is live and can process the requests received
func Live(w http.ResponseWriter, r *http.Request) {
	db, erro := database.Connect()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	repository := repositories.NewHealthcheckRepository(db)
	if erro := repository.SimulateDatabaseInsert(); erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	responses.JSON(w, http.StatusOK, nil)
}

// Ready validates if our API is ready to receive network connection and provide his main functionality
func Ready(w http.ResponseWriter, r *http.Request) {
	var hosts = []string{
		config.DatabaseHost,
	}

	db, erro := database.Connect()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	repository := repositories.NewHealthcheckRepository(db)

	for _, host := range hosts {
		if erro := repository.DNSResolver(host); erro != nil {
			responses.Erro(w, http.StatusInternalServerError, errors.New(host))
			return
		}
	}

	if erro := repository.PingDatabase(); erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	responses.JSON(w, http.StatusOK, nil)
}
