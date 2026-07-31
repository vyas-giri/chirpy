package main

import (
	"net/http"
	"sync/atomic"
	"log"
	"encoding/json"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

type validateRequest struct {
	Body string `json:"body"`
}

func handlerValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	req := validateRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Body == "" {
		respondWithError(w, http.StatusBadRequest, "Body field is required")
		return
	}

	if len(req.Body) > 280 {
		respondWithError(w, http.StatusBadRequest, "Body exceeds 280 characters")
		return
	}

	response := map[string]bool{"valid": true}
	respondWithJSON(w, http.StatusOK, response)
}

func main() {
	const port = "8080"
	const filePathRoot = "."

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot))))
	mux.Handle("/app/", fsHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidate)

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}
	
	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(server.ListenAndServe())
}