package routes

import (
	"api/src/controllers"
	"net/http"
)

var metricRoutes = Route{
	URI:                    "/live",
	Method:                 http.MethodGet,
	Function:               controllers.Metrics,
	AuthenticationRequired: false,
}
