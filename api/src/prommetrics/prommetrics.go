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

package prommetrics

import "github.com/prometheus/client_golang/prometheus"

// Prometheus Metrics Instantiation
var (
	PromOpenedConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "sm_open_connections",
			Help: "The current number of open connections",
		},
	)

	PromHandlerDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "sm_handlers_duration_seconds",
			Help: "Handlers request duration in seconds",
		}, []string{"path"},
	)

	PromRequestsDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name: "sm_request_duration_seconds",
			Help: "The duration of the requests to the sm service",
		},
	)

	PromRequestsCurrent = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "sm_requests_current",
			Help: "The current number of requests to the sm service",
		},
	)

	PromRequestStatus = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "sm_requests_total",
			Help: "The total number of requests to the sm service by status",
		}, []string{"status"},
	)

	PromClientErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "sm_errors",
			Help: "The total number of sm client errors",
		},
	)

	PromTimeTookToCreateUser = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "sm_time_to_create_user",
			Help:    "Time took to create a new user",
			Buckets: []float64{1, 2, 5, 6, 10},
		}, []string{"status"},
	)

	PromTimeTookToDeleteUser = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "sm_time_to_delete_user",
			Help:    "Time took to delete a new user",
			Buckets: []float64{1, 2, 5, 6, 10},
		}, []string{"status"},
	)

	PromCountCreatedUsers = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "sm_created_users_total",
			Help: "Quantity of users created",
		},
	)

	PromCountDeletedUsers = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "sm_deleted_users_total",
			Help: "Quantity of users deleted",
		},
	)

	PromTimeTookToCreatePublication = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "sm_time_to_create_publication",
			Help:    "Time took to create a new publication",
			Buckets: []float64{1, 2, 5, 6, 10},
		}, []string{"status"},
	)

	PromTimeTookToDeletePublication = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "sm_time_to_delete_publication",
			Help:    "Time took to delete a publication",
			Buckets: []float64{1, 2, 5, 6, 10},
		}, []string{"status"},
	)

	PromCountNewPublication = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "sm_created_publications_total",
			Help: "Quantity of publications created",
		},
	)

	PromCountDeletePublication = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "sm_deleted_publications_total",
			Help: "Quantity of publications deleted",
		},
	)
)

var Metrics []prometheus.Collector

func Load() {
	Metrics = append(Metrics, PromOpenedConnections)
	Metrics = append(Metrics, PromHandlerDuration)
	Metrics = append(Metrics, PromRequestsDuration)
	Metrics = append(Metrics, PromRequestsCurrent)
	Metrics = append(Metrics, PromRequestStatus)
	Metrics = append(Metrics, PromClientErrors)
	Metrics = append(Metrics, PromTimeTookToCreateUser)
	Metrics = append(Metrics, PromTimeTookToDeleteUser)
	Metrics = append(Metrics, PromCountCreatedUsers)
	Metrics = append(Metrics, PromCountDeletedUsers)
	Metrics = append(Metrics, PromTimeTookToCreatePublication)
	Metrics = append(Metrics, PromTimeTookToDeletePublication)
	Metrics = append(Metrics, PromCountNewPublication)
	Metrics = append(Metrics, PromCountDeletePublication)
}
