package controller

import (
	"chess-engine/engine"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	Debug bool
	SP    bool
}

type MoveRequest struct {
	Depth int    `json:"depth"`
	FEN   string `json:"fen"`
}

type AnalysisRequest struct {
	Content []string `json:"content"`
}

type LichessAnalysisParams struct {
	PGN     string `json:"pgn"`
	PGNFile string `json:"pgnFile"`
}

// static
var config Config

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func Main() {
	config = Config{
		Debug: os.Getenv("DEBUG") != "",
	}

	fmt.Printf("DEBUG: %v\n", config.Debug)
	if !config.Debug {
		log.SetOutput(io.Discard)
	}

	// TODO make this with react or something
	http.HandleFunc("/", corsMiddleware(handleIndex))
	http.HandleFunc("/info", corsMiddleware(handleCalcMove))
	http.HandleFunc("/analysis", corsMiddleware(handleAnalysis))

	fs := http.FileServer(http.Dir("templates"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	port := "2828"
	fmt.Printf("server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("engine controller says hello!"))
}

func handleCalcMove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req MoveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	game, err := engine.NewGame(req.FEN)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	start := time.Now()

	move := game.IDS(req.Depth, true)
	score := game.Evaluate()

	if config.Debug {
		fmt.Printf("\nDepth: %d\n", req.Depth)
		fmt.Printf("Move: %s\n", move)
		fmt.Printf("Nodes explored: %d\n", game.GetLeaves())
		fmt.Printf("Eval: %d\n", score)
		fmt.Printf("FEN: %s\n", req.FEN)
		fmt.Printf("Time elapsed: %v\n\n", time.Since(start))
	}

	w.Write([]byte(move)) // move is a string i suppose
}

func handleAnalysis(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AnalysisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
