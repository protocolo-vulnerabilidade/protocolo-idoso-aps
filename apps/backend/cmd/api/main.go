package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"piaps/backend/internal/armazenamento"
	"piaps/backend/internal/httpapi"
)

func main() {
	url := getenv("DATABASE_URL", "postgresql://piaps_dev:piaps_dev@localhost:5432/protocolo_idoso_dev")
	porta := getenv("PORTA", "8080")
	origens := separarCSV(getenv("CORS_ORIGENS", "http://localhost:5173,http://localhost:4173"))

	ctx := context.Background()
	pool, err := conectar(ctx, url)
	if err != nil {
		log.Fatalf("conectar ao banco: %v", err)
	}
	defer pool.Close()

	if err := armazenamento.Migrar(ctx, pool); err != nil {
		log.Fatalf("aplicar migrações: %v", err)
	}
	log.Print("migrações verificadas")

	arm := armazenamento.NovoPostgres(pool)
	handler := httpapi.NovoHandler(arm, origens)

	srv := &http.Server{
		Addr:              ":" + porta,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}

	parar := make(chan os.Signal, 1)
	signal.Notify(parar, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("API ouvindo em http://localhost:%s", porta)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("servidor: %v", err)
		}
	}()

	<-parar
	log.Print("encerrando…")
	ctxParar, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctxParar); err != nil {
		log.Printf("shutdown: %v", err)
	}
}

func conectar(ctx context.Context, url string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}
	cfg.MaxConns = 5
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	var ultimoErro error
	for i := 0; i < 30; i++ {
		if err := pool.Ping(ctx); err == nil {
			return pool, nil
		} else {
			ultimoErro = err
		}
		time.Sleep(time.Second)
	}
	pool.Close()
	return nil, ultimoErro
}

func getenv(chave, padrao string) string {
	if v := os.Getenv(chave); v != "" {
		return v
	}
	return padrao
}

func separarCSV(v string) []string {
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
