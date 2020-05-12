// Copyright (c) 2020 Fandom, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type PrometheusExporterConfiguration struct {
	PoolCounterAddress      string `default:"localhost:7531" split_words:"true"`
	ListenAddress           string `default:"localhost:8000" split_words:"true"`
	LogsAsJson              bool   `default:"false" split_words:"true"`
	CollectorTimeoutSeconds int    `default:"5" split_words:"true"`
	ServerTimeoutSeconds    int    `default:"3" split_words:"true"`
}

func main() {
	var configuration PrometheusExporterConfiguration
	err := envconfig.Process("exporter", &configuration)

	if err != nil {
		log.Fatal(err.Error())
	}

	if configuration.LogsAsJson {
		log.SetFormatter(&log.JSONFormatter{
			FieldMap: log.FieldMap{
				log.FieldKeyTime: "@timestamp",
				log.FieldKeyMsg:  "@message",
			},
		})
	}

	srv := &http.Server{
		Addr:         configuration.ListenAddress,
		WriteTimeout: time.Second * time.Duration(configuration.ServerTimeoutSeconds),
		ReadTimeout:  time.Second * time.Duration(configuration.ServerTimeoutSeconds),
		IdleTimeout:  time.Second * time.Duration(configuration.ServerTimeoutSeconds),
	}

	log.Infof("Configuring collector to scrape from %s", configuration.PoolCounterAddress)

	collector := newPoolCounterCollector(configuration)
	prometheus.MustRegister(collector)

	http.Handle("/prometheus", promhttp.Handler())

	http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		response.Write([]byte("EXPORTER ALIVE"))
	})

	log.Infof("Starting Prometheus collector on %s", configuration.ListenAddress)

	if err := srv.ListenAndServe(); err != nil {
		log.Error(err)
	}
}
