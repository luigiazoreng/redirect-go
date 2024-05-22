package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	var (
		port       int
		targetIP   string
		targetPort int
	)

	flag.IntVar(&port, "port", 3000, "Porta de entrada do servidor")
	flag.StringVar(&targetIP, "target-ip", "127.0.0.1", "Endereço IP para redirecionar (sem porta)")
	flag.IntVar(&targetPort, "target-port", 32767, "Porta para redirecionar")
	port = 3000
	flag.Parse()

	targetURL := fmt.Sprintf("http://%s:%d", targetIP, targetPort)

	sm := http.NewServeMux()
	sm.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		redirectHandler(targetURL, w, r)
	})
	server := &http.Server{
		Handler:      sm,
		Addr:         "localhost:" + strconv.Itoa(port),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
	}

	go func() {
		fmt.Printf("Servidor iniciado na porta %d. Redirecionando para %s...\n", port, targetURL)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Erro ao iniciar o servidor: %v\n", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)

	sig := <-sigChan
	fmt.Println("Received terminate, graceful shutdown", sig)

	d := 30 * time.Second
	tc, cancel := context.WithTimeout(context.Background(), d)
	server.Shutdown(tc)
	defer cancel()

}

func redirectHandler(targetURL string, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request received:", r.Method, r.URL.Path)
	resp, err := http.Get(targetURL + r.URL.Path)
	if err != nil {
		fmt.Printf("Error processing request: %v\n", err)
		http.Error(w, "Error processing request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Request not successful: %s\n", resp.Status)
		http.Error(w, "Request not successful", resp.StatusCode)
		return
	}

	fmt.Printf("Request successful: %s\n", resp.Status)
	fmt.Println("Redirecting to:", resp.Request.URL)
	http.Redirect(w, r, resp.Request.URL.String(), http.StatusTemporaryRedirect)
}
