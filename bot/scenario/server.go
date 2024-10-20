package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Diarkis/diarkis/util"
	"handson/bot/scenario/lib/report"
	"io"
	"net/http"
	"strings"
)

type String string

func (s String) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, s)
}
func handleRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		logger.Error("Got invalid http method: %v ", r.Method)
	}
	body := r.Body
	defer body.Close()
	buf := new(bytes.Buffer)
	io.Copy(buf, body)
	json.Unmarshal(buf.Bytes(), &ss)
	json.Unmarshal(buf.Bytes(), &gp.Raw.ParamsFromAPI)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Scenario Started\n")
	logger.Info("Starting Scenario [%s] with parameters [%s] ", ss.ScenarioName, ss.ScenarioPattern)
	go run()
}
func handleGetMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := report.GetPrometheusMetrics()
	fmt.Fprint(w, metrics)
	logger.Verbose("Get Metrics Called...")
}
func listen() error {
	address := util.GetEnv("BOT_ADDRESS")
	if address == "" {
		address = "localhost"
	}
	port := util.GetEnv("BOT_PORT")
	if port == "" {
		port = "9500"
	}
	host := strings.Join([]string{address, port}, ":")
	http.Handle("/", String("hello"))
	http.HandleFunc("/run/", handleRun)
	http.HandleFunc("/metrics/", handleGetMetrics)
	logger.Info("Bot server started. listening %s ...", host)
	http.ListenAndServe(host, nil)
	return nil
}
