package main

import (
	"flag"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/ClickHouse/clickhouse_exporter/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/log"
)

var (
	listeningAddress    = flag.String("telemetry.address", ":9116", "Address on which to expose metrics.")
	metricsEndpoint     = flag.String("telemetry.endpoint", "/metrics", "Path under which to expose metrics.")
	clickhouseScrapeURIs = flag.String("scrape_uris", "http://localhost:8123/;http://localhost2:8123/", "URIs to clickhouse http endpoint")
	clickhouseOnly      = flag.Bool("clickhouse_only", false, "Expose only Clickhouse metrics, not metrics from the exporter itself")
	insecure            = flag.Bool("insecure", true, "Ignore server certificate if using https")
	user                = os.Getenv("CLICKHOUSE_USER")
	password            = os.Getenv("CLICKHOUSE_PASSWORD")
)

func main() {
	flag.Parse()

	var uriArr []url.URL
	for _, uri := range strings.Split(*clickhouseScrapeURIs, ":") {
		uri, err := url.Parse(uri)
		if err != nil {
			log.Fatal(err)
		}
		uriArr = append(uriArr, *uri)
	}
	log.Printf("Scraping %s", *clickhouseScrapeURIs)

	registerer := prometheus.DefaultRegisterer
	gatherer := prometheus.DefaultGatherer
	if *clickhouseOnly {
		reg := prometheus.NewRegistry()
		registerer = reg
		gatherer = reg
	}

	exporters := exporter.NewExporters(uriArr, *insecure, user, password)
	for _, e := range exporters {
		registerer.MustRegister(e)
	}

	http.Handle(*metricsEndpoint, promhttp.HandlerFor(gatherer, promhttp.HandlerOpts{}))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Clickhouse Exporter</title></head>
			<body>
			<h1>Clickhouse Exporter</h1>
			<p><a href="` + *metricsEndpoint + `">Metrics</a></p>
			</body>
			</html>`))
	})

	log.Fatal(http.ListenAndServe(*listeningAddress, nil))
}
