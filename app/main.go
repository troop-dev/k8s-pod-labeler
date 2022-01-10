package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/troop-dev/k8s-pod-labeler/app/mutate"
)

func main() {
	log.Println("Starting server ...")

	mux := http.NewServeMux()

	// TODO: make labels & annotations configurable
	labels := map[string]string{}
	annotations := map[string]string{
		"cluster-autoscaler.kubernetes.io/safe-to-evict": "true",
	}

	mutateHandler := &mutateHandler{
		labels:      labels,
		annotations: annotations,
	}

	mux.HandleFunc("/", handleHealthCheck)
	mux.HandleFunc("/healthz", handleHealthCheck)
	mux.HandleFunc("/mutate", mutateHandler.Handle)

	s := &http.Server{
		Addr:           ":8443",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1048576
	}
	// TODO: move to config
	sslCertFilePath := "/etc/webhook/certs/tls.crt"
	sslKeyFilePath := "/etc/webhook/certs/tls.key"
	// start listening and block until shutdown
	log.Fatal(s.ListenAndServeTLS(sslCertFilePath, sslKeyFilePath))
}

// handleHealthCheck adds a health check route
func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// mutateHandler implements the request handler for mutation requests
type mutateHandler struct {
	labels      map[string]string
	annotations map[string]string
}

// Handle handles mutation requests
func (m *mutateHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// read the body / request
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		sendError(err, w)
		return
	}

	// mutate the request
	mutated, err := mutate.Mutate(body, m.labels, m.annotations)
	if err != nil {
		sendError(err, w)
		return
	}

	// and write it back
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(mutated); err != nil {
		sendError(err, w)
	}
}

func sendError(err error, w http.ResponseWriter) {
	log.Println(err)
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "%s", err)
}
