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
}
