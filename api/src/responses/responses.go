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

package responses

import (
	"api/src/prommetrics"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

// JSON returns json response
func JSON(startedTime time.Time, w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		prommetrics.PromRequestsDuration.Observe(time.Since(startedTime).Seconds())
		prommetrics.PromRequestsCurrent.Dec()
		prommetrics.PromRequestStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
		if erro := json.NewEncoder(w).Encode(data); erro != nil {
			log.Fatal(erro)
		}
	}
}

// Erro returns an error json response
func Erro(startedTime time.Time, w http.ResponseWriter, statusCode int, erro error) {
	prommetrics.PromClientErrors.Inc()

	JSON(startedTime, w, statusCode, struct {
		Erro string `json:"erro"`
	}{
		Erro: erro.Error(),
	})
}
