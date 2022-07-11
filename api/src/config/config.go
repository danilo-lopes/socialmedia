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

package config

import (
	"fmt"
	"os"
	"strconv"
)

var (
	// Database information
	DatabaseStringConnection string = ""
	DatabaseHost             string = ""
	DatabasePort             string = ""

	// API Service Port
	APIPort int = 0

	// Used to assign the token
	SecretKey []byte
)

// Load inicialize environment variables, configure database string connection and set which port the api will run.
func Load() {
	var erro error

	APIPort, erro = strconv.Atoi(os.Getenv("API_PORT"))
	if erro != nil {
		APIPort = 8080
	}

	DatabaseHost = os.Getenv("DB_HOST")
	DatabasePort = os.Getenv("DB_PORT")
	DatabaseStringConnection = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?tls=skip-verify&charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		DatabaseHost,
		DatabasePort,
		os.Getenv("DB_NAME"),
	)

	SecretKey = []byte(os.Getenv("SECRET_KEY"))
}
