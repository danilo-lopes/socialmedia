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

var usersRoutes = []Route{
	{
		URI:                    "/users",
		Method:                 http.MethodPost,
		Function:               controllers.CreateUser,
		AuthenticationRequired: false,
	},
	{
		URI:                    "/users",
		Method:                 http.MethodGet,
		Function:               controllers.GetUsers,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/users/{userID}",
		Method:                 http.MethodGet,
		Function:               controllers.GetUser,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/users/{userID}",
		Method:                 http.MethodPut,
		Function:               controllers.UpdateUser,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/users/{userID}",
		Method:                 http.MethodDelete,
		Function:               controllers.DeleteUser,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/users/{userID}/follow",
		Method:                 http.MethodPost,
		Function:               controllers.FollowUser,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/users/{userID}/unfollow",
		Method:                 http.MethodPost,
		Function:               controllers.UnFollowUser,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/users/{userID}/followers",
		Method:                 http.MethodGet,
		Function:               controllers.GetFollowers,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/users/{userID}/following",
		Method:                 http.MethodGet,
		Function:               controllers.GetFollowing,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/users/{userID}/updatepass",
		Method:                 http.MethodPost,
		Function:               controllers.UpdatePass,
		AuthenticationRequired: true,
	},
}
