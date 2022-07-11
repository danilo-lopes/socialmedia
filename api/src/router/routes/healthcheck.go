package routes

import (
	"api/src/controllers"
	"net/http"
)

var healthcheckRoutes = []Route{
	{
		URI:                    "/live",
		Method:                 http.MethodGet,
		Function:               controllers.Live,
		AuthenticationRequired: false,
	},
	{
		URI:                    "/ready",
		Method:                 http.MethodGet,
		Function:               controllers.Ready,
		AuthenticationRequired: false,
	},
}
