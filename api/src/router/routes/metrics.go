package routes

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var metricsRoute = PromRoute{
	URI:                    "/metrics",
	Method:                 http.MethodGet,
	Function:               promhttp.Handler(),
	AuthenticationRequired: false,
}
