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

package main

import (
	"api/src/config"
	"api/src/prommetrics"
	"api/src/router"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	r := router.Generate()
	config.Load()

	prommetrics.Load()
	for _, metric := range prommetrics.Metrics {
		prometheus.MustRegister(metric)
	}

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.APIPort),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		ConnState: func(c net.Conn, s http.ConnState) {
			switch s {
			case http.StateNew:
				prommetrics.PromOpenedConnections.Inc()
			case http.StateHijacked:
				prommetrics.PromOpenedConnections.Dec()
			case http.StateClosed:
				prommetrics.PromOpenedConnections.Dec()
			}
		},
	}

	fmt.Printf("Serving on Port %d\n", config.APIPort)
	log.Fatal(s.ListenAndServe())
}
