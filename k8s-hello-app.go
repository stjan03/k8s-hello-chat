// A simple webserver that written in golang that publish...
//      /               Default page the publish "Hello world...
//      /env            Publish the pod specific enviroment data in json format
//      /health         Healthcheck "OK"
//      /metrics        Prometheus metrics


// Enviroment varibles that can be set
//    APP_* is manually set
//    K8S_* is automaticlly set by the pod deployment

// Created by Stefan Jansson
// 2021-07-17

// Compile with...
// go build k8s-hello-app.go


// ----------- The magic code ------------------
package main

import (
        "fmt"
        "net/http"
        "os"
        "log"
        "strings"
        "encoding/json"
        "github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promhttp"
)


func defaultHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "Hello kubernetes, world of containers...")
  fmt.Fprintln(w, " - Appversion: " + os.Getenv("APP_VERSION"))
  fmt.Fprintln(w, " - Hostname: " + os.Getenv("HOSTNAME"))
  podInfo.Inc()
  //fmt.Fprintln(w, "PodInfo: " + podInfo.get())
}


func healthHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "ok")
}


func envHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  //w.Header().Set("Access-Control-Allow-Origin", "*")
  m := make(map[string]string)
  for _, e := range os.Environ() {
    if i := strings.Index(e, "="); i >= 0 {
      m[e[:i]] = e[i+1:]
    }
  }
  env_json, _ := json.MarshalIndent(m, "", "\t")
  fmt.Fprintf(w, string(env_json))
}


var (
  podInfo = prometheus.NewCounter(prometheus.CounterOpts{
    Name: "pod_info",
    Help: "Number of pod hits",
    ConstLabels: prometheus.Labels{
      "hostname":  os.Getenv("HOSTNAME"),
      "appversion":  os.Getenv("APP_VERSION"),
    },
  })
)


func init() {
  prometheus.MustRegister(podInfo)
}


func main() {
  http.HandleFunc("/", defaultHandler)
  http.HandleFunc("/env", envHandler)
  http.HandleFunc("/health", healthHandler)
  http.Handle("/metrics", promhttp.Handler())

  log.Println("Listening on port 8080")
  log.Fatal(http.ListenAndServe(":8080", nil))
}
