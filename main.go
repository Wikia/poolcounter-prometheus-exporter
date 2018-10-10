package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
