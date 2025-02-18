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
	"bufio"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

const (
	namespace = "poolcounter"

	totalAcquired     = "total_acquired"
	totalReleases     = "total_releases"
	hashTableEntries  = "hashtable_entries"
	processingWorkers = "processing_workers"
	waitingWorkers    = "waiting_workers"
	connectErrors     = "connect_errors"
	fullQueues        = "full_queues"
	lockMismatch      = "lock_mismatch"
	releaseMismatch   = "release_mismatch"
	processedCount    = "processed_count"
)

// PoolCounterCollector handles scraping metrics from a poolcounter instance.
type PoolCounterCollector struct {
	poolCounterAddress      string
	collectorTimeoutSeconds int

	up                            *prometheus.Desc
	totalProcessingTimeSeconds    *prometheus.Desc
	averageProcessingTimeSeconds  *prometheus.Desc
	totalGainedTimeSeconds        *prometheus.Desc
	totalExclusiveWaitTimeSeconds *prometheus.Desc
	totalSharedWaitTimeSeconds    *prometheus.Desc
	totalAcquired                 *prometheus.Desc
	totalReleases                 *prometheus.Desc
	hashTableEntries              *prometheus.Desc
	processingWorkers             *prometheus.Desc
	waitingWorkers                *prometheus.Desc
	connectErrors                 *prometheus.Desc
	failedSends                   *prometheus.Desc
	fullQueues                    *prometheus.Desc
	lockMismatch                  *prometheus.Desc
	releaseMismatch               *prometheus.Desc
	processedCount                *prometheus.Desc
}

// newPoolCounterCollector initializes and returns a new PoolCounterCollector instance based on scraper configuration.
func newPoolCounterCollector(configuration PrometheusExporterConfiguration) *PoolCounterCollector {
	return &PoolCounterCollector{
		poolCounterAddress:      configuration.PoolCounterAddress,
		collectorTimeoutSeconds: configuration.CollectorTimeoutSeconds,
		up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "up"),
			"Whether poolcounter is up and responding to the exporter",
			nil,
			nil,
		),
		totalProcessingTimeSeconds: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "total_processing_time_seconds"),
			"Total processing time in seconds",
			nil,
			nil,
		),
		averageProcessingTimeSeconds: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "avg_processing_time_seconds"),
			"Average processing time in seconds",
			nil,
			nil,
		),
		totalGainedTimeSeconds: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "total_gained_time_seconds"),
			"Total processing time saved by the use of PoolCounter in seconds",
			nil,
			nil,
		),
		totalExclusiveWaitTimeSeconds: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "total_excl_wait_time_seconds"),
			"Total waiting time for exclusive locks in seconds",
			nil,
			nil,
		),
		totalSharedWaitTimeSeconds: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "total_shared_wait_time_seconds"),
			"Total waiting time for shared locks in seconds",
			nil,
			nil,
		),
		totalAcquired: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", totalAcquired),
			"Total acquired locks count",
			nil,
			nil,
		),
		totalReleases: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", totalReleases),
			"Total released locks count",
			nil,
			nil,
		),
		hashTableEntries: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", hashTableEntries),
			"Number of entries in poolcounter hash table",
			nil,
			nil,
		),
		processingWorkers: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", processingWorkers),
			"Number of workers busy processing tasks",
			nil,
			nil,
		),
		waitingWorkers: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", waitingWorkers),
			"Number of workers waiting for tasks to be completed",
			nil,
			nil,
		),
		connectErrors: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", connectErrors),
			"Total count of client connection errors",
			nil,
			nil,
		),
		fullQueues: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", fullQueues),
			"Number of queues full of waiting workers",
			nil,
			nil,
		),
		lockMismatch: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", lockMismatch),
			"Total count of mismatched lock requests",
			nil,
			nil,
		),
		releaseMismatch: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", releaseMismatch),
			"Total count of mismatched release requests",
			nil,
			nil,
		),
		processedCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", processedCount),
			"Total count of processed tasks",
			nil,
			nil,
		),
	}
}

// parseTimeDescription takes a poolcounter time string (e.g. 389 days 9343h 3m 28.000000s)
// and returns its duration in seconds.
func parseTimeDescription(description string) float64 {
	segments := strings.Split(description, " days ")

	if len(segments) > 1 {
		days, _ := strconv.ParseInt(segments[0], 10, 0)
		duration, _ := time.ParseDuration(strings.Replace(segments[1], " ", "", -1))

		return float64(24*60*60*days) + duration.Seconds()
	}

	duration, _ := time.ParseDuration(strings.Replace(segments[0], " ", "", -1))

	return duration.Seconds()
}

func (collector *PoolCounterCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.totalProcessingTimeSeconds
	ch <- collector.averageProcessingTimeSeconds

	ch <- collector.up
	ch <- collector.totalAcquired
	ch <- collector.totalReleases
	ch <- collector.hashTableEntries
	ch <- collector.processingWorkers
	ch <- collector.waitingWorkers
	ch <- collector.connectErrors
	ch <- collector.fullQueues
	ch <- collector.lockMismatch
	ch <- collector.releaseMismatch
	ch <- collector.processedCount
}

func (collector *PoolCounterCollector) Collect(ch chan<- prometheus.Metric) {
	var finalErr error = nil
	var upValue float64 = 1 // 1 or 0

	defer func() {
		ch <- prometheus.MustNewConstMetric(collector.up, prometheus.GaugeValue, upValue)
		if finalErr != nil {
			log.Error(finalErr)
		}
	}()

	conn, finalErr := net.Dial("tcp", collector.poolCounterAddress)
	if finalErr != nil {
		upValue = 0
		return
	}

	defer conn.Close()

	conn.SetDeadline(time.Now().Add(time.Duration(collector.collectorTimeoutSeconds) * time.Second))

	_, finalErr = conn.Write([]byte("STATS FULL\n"))
	if finalErr != nil {
		upValue = 0
		return
	}

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ": ")

		if len(parts) < 2 {
			break // sanity
		}

		name, value := parts[0], parts[1]

		switch name {
		case "total processing time":
			value := parseTimeDescription(value)
			ch <- prometheus.MustNewConstMetric(collector.totalProcessingTimeSeconds, prometheus.CounterValue, value)
		case "average processing time":
			value := parseTimeDescription(value)
			ch <- prometheus.MustNewConstMetric(collector.averageProcessingTimeSeconds, prometheus.GaugeValue, value)
		case "gained time":
			value := parseTimeDescription(value)
			ch <- prometheus.MustNewConstMetric(collector.totalGainedTimeSeconds, prometheus.CounterValue, value)
		case "waiting time for me":
			value := parseTimeDescription(value)
			ch <- prometheus.MustNewConstMetric(collector.totalExclusiveWaitTimeSeconds, prometheus.CounterValue, value)
		case "waiting time for anyone":
			value := parseTimeDescription(value)
			ch <- prometheus.MustNewConstMetric(collector.totalSharedWaitTimeSeconds, prometheus.CounterValue, value)
		case totalAcquired:
			value, _ := strconv.ParseFloat(value, 0)
			ch <- prometheus.MustNewConstMetric(collector.totalAcquired, prometheus.CounterValue, value)
		case totalReleases:
			value, _ := strconv.ParseFloat(value, 0)
			ch <- prometheus.MustNewConstMetric(collector.totalReleases, prometheus.CounterValue, value)
		case hashTableEntries:
			value, _ := strconv.ParseFloat(value, 0)
			ch <- prometheus.MustNewConstMetric(collector.hashTableEntries, prometheus.GaugeValue, value)
		case processingWorkers:
			value, _ := strconv.ParseFloat(value, 0)
			ch <- prometheus.MustNewConstMetric(collector.processingWorkers, prometheus.GaugeValue, value)
		case waitingWorkers:
			value, _ := strconv.ParseFloat(value, 0)
			ch <- prometheus.MustNewConstMetric(collector.waitingWorkers, prometheus.GaugeValue, value)
		case connectErrors:
			value, _ := strconv.ParseFloat(value, 0)
			ch <- prometheus.MustNewConstMetric(collector.connectErrors, prometheus.CounterValue, value)
		case fullQueues:
			value, _ := strconv.ParseFloat(value, 0)
			ch <- prometheus.MustNewConstMetric(collector.fullQueues, prometheus.CounterValue, value)
		case lockMismatch:
			value, _ := strconv.ParseFloat(value, 0)
			ch <- prometheus.MustNewConstMetric(collector.lockMismatch, prometheus.CounterValue, value)
		case releaseMismatch:
			value, _ := strconv.ParseFloat(value, 0)
			ch <- prometheus.MustNewConstMetric(collector.releaseMismatch, prometheus.CounterValue, value)
		case processedCount:
			value, _ := strconv.ParseFloat(value, 0)
			ch <- prometheus.MustNewConstMetric(collector.processedCount, prometheus.CounterValue, value)
		}
	}

	if finalErr = scanner.Err(); finalErr != nil {
		upValue = 0
	}
}
