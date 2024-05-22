package main

import (
	"flag"
	"fmt"
	"net/http"
)

func redirectHandler(targetURL string, w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, targetURL+r.RequestURI, http.StatusFound)
}

func main() {
	var (
		port        int
		targetIP    string
		targetPort  int
	)

	flag.IntVar(&port, "port", 3000, "Porta de entrada do servidor")
	flag.StringVar(&targetIP, "target-ip", "127.0.0.1", "Endereço IP para redirecionar (sem porta)")
	flag.IntVar(&targetPort, "target-port", 80, "Porta para redirecionar")

	flag.Parse()

	targetURL := fmt.Sprintf("http://%s:%d", targetIP, targetPort)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		redirectHandler(targetURL, w, r)
	})

	fmt.Printf("Servidor iniciado na porta %d. Redirecionando para %s...\n", port, targetURL)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
