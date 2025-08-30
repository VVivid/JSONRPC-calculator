package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
	// Create JSON-RPC server
	rpcServer := NewJSONRPCServer()
	
	// HTTP handler for JSON-RPC
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers for web testing
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		// Handle OPTIONS request for CORS preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		// Only accept POST requests
		if r.Method != "POST" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`{"error": "Only POST method is allowed for JSON-RPC"}`))
			return
		}
		
		// Check content type
		contentType := r.Header.Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Content-Type must be application/json"}`))
			return
		}
		
		// Read request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Cannot read request body"}`))
			return
		}
		defer r.Body.Close()
		
		// Process JSON-RPC request
		response, err := rpcServer.HandleRequest(body)
		if err != nil {
			log.Printf("Error processing request: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Internal server error"}`))
			return
		}
		
		// Set response headers
		w.Header().Set("Content-Type", "application/json")
		
		// Handle notifications (no response)
		if response == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		
		// Send JSON-RPC response
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	})
	
	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy", "service": "JSON-RPC Calculator"}`))
	})
	
	// Start server
	port := 8090
	log.Printf("JSON-RPC Calculator Server starting on port %d", port)
	log.Printf("Health check available at: http://localhost:%d/health", port)
	log.Printf("JSON-RPC endpoint at: http://localhost:%d/", port)
	log.Println("")
	log.Println("Example curl commands:")
	log.Printf(`  curl -X POST -H "Content-Type: application/json" -d '{"jsonrpc":"2.0","method":"add","params":{"a":10,"b":20},"id":1}' http://localhost:%d/`, port)
	log.Printf(`  curl -X POST -H "Content-Type: application/json" -d '{"jsonrpc":"2.0","method":"log","params":{"message":"Hello from curl!"}}' http://localhost:%d/`, port)
	
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}