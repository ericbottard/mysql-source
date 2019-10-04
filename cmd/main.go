/*
 * Copyright 2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package main

import (
	"context"
	"fmt"
	"github.com/projectriff-samples/mysql-source/pkg/source"
	client2 "github.com/projectriff/http-gateway/pkg/client"
	"os"
)

func main() {
	query := os.Getenv("QUERY")
	update := os.Getenv("UPDATE")
	dataSourceName := os.Getenv("DATASOURCE")
	gateway := os.Getenv("GATEWAY")
	topics := os.Getenv("TOPICS")
	if query == "" || update == "" || dataSourceName == "" || gateway == "" || topics == ""{
		panic("Expected all of the following ENV VARS to be set: QUERY, UPDATE, DATASOURCE, GATEWAY and TOPICS")
	}


	var client *client2.StreamClient
	var err error
	if client, err = client2.NewStreamClient(gateway, topics, "text/plain") ; err != nil {
		panic(err)
	}

	s, err := source.NewSource(query, update, dataSourceName, client)
	if err != nil {
		panic(err)
	}

	if n, err := s.Run(context.Background()) ; err != nil {
		panic(err)
	} else {
		fmt.Printf("Emitted %d messages on %s\n", n, topics)
	}

	s.Close()

}
