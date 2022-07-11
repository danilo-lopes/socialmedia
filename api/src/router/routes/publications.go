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
	"api/src/controllers"
	"net/http"
)

var publicationsRoutes = []Route{
	{
		URI:                    "/publications",
		Method:                 http.MethodPost,
		Function:               controllers.CreatePublication,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/publications",
		Method:                 http.MethodGet,
		Function:               controllers.GetPublications,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/publications/{publicationID}",
		Method:                 http.MethodGet,
		Function:               controllers.GetPublication,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/publications/{publicationID}",
		Method:                 http.MethodPut,
		Function:               controllers.UpdatePublication,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/publications/{publicationID}",
		Method:                 http.MethodDelete,
		Function:               controllers.DeletePublication,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/users/{userID}/publications",
		Method:                 http.MethodGet,
		Function:               controllers.GetUserPublications,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/publications/{publicationID}/like",
		Method:                 http.MethodPost,
		Function:               controllers.LikePublication,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/publications/{publicationID}/unlike",
		Method:                 http.MethodPost,
		Function:               controllers.UnLikePublication,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/publications/{publicationID}/likers",
		Method:                 http.MethodGet,
		Function:               controllers.GetLikers,
		AuthenticationRequired: true,
	},
}
