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

package routes

import (
	"api/src/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

// Route represents an API route
type Route struct {
	URI                    string
	Method                 string
	Function               func(http.ResponseWriter, *http.Request)
	AuthenticationRequired bool
}

type PromRoute struct {
	URI                    string
	Method                 string
	Function               http.Handler
	AuthenticationRequired bool
}

// Configure instanciate all API routes into mux router
func Configure(r *mux.Router) *mux.Router {
	apiRoutes := usersRoutes
	apiRoutes = append(apiRoutes, loginRoute)
	apiRoutes = append(apiRoutes, publicationsRoutes...)
	apiRoutes = append(apiRoutes, healthcheckRoutes...)

	for _, apiRoute := range apiRoutes {
		if apiRoute.AuthenticationRequired {
			r.HandleFunc(apiRoute.URI,
				middlewares.Logger(apiRoute.URI, middlewares.Authenticate(apiRoute.Function)),
			).Methods(apiRoute.Method)
		} else {
			r.HandleFunc(apiRoute.URI,
				middlewares.Logger(apiRoute.URI, apiRoute.Function),
			).Methods(apiRoute.Method)
		}
	}

	// Prometheus api route
	r.Handle(
		metricsRoute.URI,
		metricsRoute.Function,
	).Methods(metricsRoute.Method)

	return r
}
