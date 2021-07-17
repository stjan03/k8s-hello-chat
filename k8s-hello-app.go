// A simple webserver that written in golang that publish...
//      /               Default page the publish "Hello world...
//      /json           Publish the pod specific enviroment data in json format
//      /health         Healthcheck "OK"
//      /metrics        Prometheus metrics


// Enviroment varibles that can be set
//    APP_* is manually set
//    K8S_* is automaticlly set by the pod deployment

// Created by Stefan
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
        //"encoding/json"
        "github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promhttp"
)


func defaultHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "Hello world..." + os.Getenv("APP_TEXT"))
  podInfo.Inc()
}


func healthHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "ok")
}


func jsonHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  w.Header().Set("Access-Control-Allow-Origin", "*")
//      response := "{\n  \"Version\": \"" + os.Getenv("APP_VERSION") +
//                    "\",\n  \"Text\": \"" + os.Getenv("APP_TEXT") +
//                    "\",\n  \"Nodename\": \"" + os.Getenv("K8S_NODE_NAME") +
//                    "\",\n  \"Podname\": \"" + os.Getenv("K8S_POD_NAME") +
//                    "\",\n  \"Namespace\": \"" + os.Getenv("K8S_POD_NAMESPACE") +
//                    "\",\n  \"PodIP\": \"" + os.Getenv("K8S_POD_IP") +
//                    "\",\n  \"HostIP\": \"" + os.Getenv("K8S_HOST_IP") +
//                    "\",\n  \"Serviceaccount\": \"" + os.Getenv("K8S_POD_SERVICE_ACCOUNT") +
//                    "\",\n  \"Hostname\": \"" + os.Getenv("HOSTNAME") +
//                    "\"\n}\n"
  response := os.Environ()
  fmt.Fprintf(w, response)
}


var (
  podInfo = prometheus.NewCounter(prometheus.CounterOpts{
    Name: "pod_info",
    Help: "Pod information",
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
  http.HandleFunc("/json", jsonHandler)
  http.HandleFunc("/health", healthHandler)
  http.Handle("/metrics", promhttp.Handler())

  log.Println("Listening on port 8080")
  //log.Println(json.MarshalIndent(os.Environ()))
  log.Fatal(http.ListenAndServe(":8080", nil))
}
