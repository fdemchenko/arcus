package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/fdemchenko/arcus/internal/config"
)

func main() {
	cfg := config.MustLoad()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]string{
			"env": cfg.Env,
		}

		json.NewEncoder(w).Encode(data)
	})

	address := fmt.Sprintf("%s:%d", cfg.HTTPServer.Host, cfg.HTTPServer.Port)
	server := http.Server{
		Addr:    address,
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
