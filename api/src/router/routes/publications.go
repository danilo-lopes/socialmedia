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
		URI:                    "/publications/{postID}",
		Method:                 http.MethodGet,
		Function:               controllers.GetPublication,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/publications/{postID}",
		Method:                 http.MethodPut,
		Function:               controllers.UpdatePublication,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/publications/{postID}",
		Method:                 http.MethodDelete,
		Function:               controllers.DeletePublication,
		AuthenticationRequired: true,
	},
}
