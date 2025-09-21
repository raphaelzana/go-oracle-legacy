package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/sijms/go-ora/v2" // driver go-ora
	"log"
	"net/http"
	"net/url"
	"time"
)

type Response struct {
	Message string `json:"message"`
	DBTime  string `json:"db_time,omitempty"`
}

var db *sql.DB

func newOracleDBFromEnv() (*sql.DB, error) {
	user := ""    // Oracle user
	pass := ""    // Oracle password
	host := ""    // Oracle host
	port := ""    // Oracle port
	service := "" // Oracle service

	// DSN do go-ora (driver puro)
	dsn := fmt.Sprintf("oracle://%s:%s@%s:%s/%s",
		url.QueryEscape(user),
		url.QueryEscape(pass),
		host,
		port,
		service,
	)

	db, err := sql.Open("oracle", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	// Ajustes do pool de conexões
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	// Teste de conectividade com timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("falha no Ping ao Oracle: %w", err)
	}

	return db, nil
}

// Handler function for root path "/"
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Busca a hora do banco para comprovar que a conexão está OK
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var dbTime string
	err := db.QueryRowContext(ctx, "SELECT TO_CHAR(SYSDATE, 'YYYY-MM-DD HH24:MI:SS') FROM dual").Scan(&dbTime)
	if err != nil {
		http.Error(w, "Erro ao consultar o banco", http.StatusInternalServerError)
		return
	}

	response := Response{
		Message: "Hello, Buddy!",
		DBTime:  dbTime,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Erro ao serializar resposta", http.StatusInternalServerError)
		return
	}
}

func main() {
	var err error
	db, err = newOracleDBFromEnv()
	if err != nil {
		log.Fatalf("Não foi possível conectar ao Oracle: %v", err)
	}
	defer db.Close()

	http.HandleFunc("/", helloHandler) // Route setup

	log.Println("Servidor escutando em :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
