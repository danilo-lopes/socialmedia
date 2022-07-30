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

package middlewares

import (
	"api/src/authentication"
	"api/src/prommetrics"
	"api/src/responses"
	"log"
	"net/http"
	"time"
)

// Logger log the request information in terminal
func Logger(path string, nextFunction http.HandlerFunc) http.HandlerFunc {
	return handler(path, nextFunction)
}

// Authenticate validates if the User is authenticated
func Authenticate(nextFunction http.HandlerFunc) http.HandlerFunc {
	now := time.Now()
	return func(w http.ResponseWriter, r *http.Request) {
		if erro := authentication.ValidateToken(r); erro != nil {
			responses.Erro(now, w, http.StatusUnauthorized, erro)
			return
		}
		nextFunction(w, r)
	}
}

func handler(pattern string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		next(w, r)
		log.Printf("%s - - '%s %s %s' -", r.RemoteAddr, r.Method, r.RequestURI, r.Proto)
		prommetrics.PromHandlerDuration.WithLabelValues(pattern).Observe(time.Since(now).Seconds())
	})
}
