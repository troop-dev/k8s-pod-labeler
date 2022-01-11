package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/troop-dev/k8s-pod-labeler/app/mutate"
)

// Run starts the API server and blocks until shutdown
func Run(cfg *Config) {
	log.Println("configuring server ...")

	mux := http.NewServeMux()

	// TODO: make labels & annotations configurable
	labels, err := parseB64Map(cfg.LabelsB64)
	if err != nil {
		log.Fatalf("failed to load labels from base64 '%s', %v", cfg.LabelsB64, err)
	}
	annotations, err := parseB64Map(cfg.AnnotationsB64)
	if err != nil {
		log.Fatalf("failed to load annotations from base64 '%s', %v", cfg.AnnotationsB64, err)
	}

	mutateHandler := &mutateHandler{
		labels:      labels,
		annotations: annotations,
	}

	mux.HandleFunc("/", handleHealthCheck)
	mux.HandleFunc("/healthz", handleHealthCheck)
	mux.HandleFunc("/mutate", mutateHandler.Handle)

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.Port),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1048576
	}
	// TODO: move to config
	sslCertFilePath := "/etc/webhook/certs/tls.crt"
	sslKeyFilePath := "/etc/webhook/certs/tls.key"
	// start listening and block until shutdown
	log.Printf("listening on port %d\n", cfg.Port)
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

// sendError is a helper function that writes internal error codes for http
// response
func sendError(err error, w http.ResponseWriter) {
	log.Println(err)
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "%s", err)
}

// parseB64Map accepts a base64-encoded json object and
// unmarshals it to a map of strings
func parseB64Map(in string) (map[string]string, error) {
	// decode from base64
	decoded, err := base64.RawStdEncoding.DecodeString(in)
	if err != nil {
		return nil, err
	}
	out := map[string]string{}
	// unmarshal json
	if len(decoded) > 0 {
		if err := json.Unmarshal(decoded, &out); err != nil {
			return nil, err
		}
	}
	return out, nil
}
