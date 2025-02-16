package controler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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
	PGN string `json:"pgn"`
	PGNFile string `json:"pgnFile"`
}

func main() {
	config := Config{
		Debug: os.Getenv("DEBUG") != "",
		SP:    os.Getenv("SP") != "",
	}

	if !config.Debug {
		log.SetOutput(nil)
	}

	// TODO make this with react or something
	//http.HandleFunc("/", handleIndex)
	http.HandleFunc("/selfplay", handleSelfplay)
	http.HandleFunc("/spinfo", handleSPMove)
	http.HandleFunc("/info", handleCalcMove)
	http.HandleFunc("/analysis", handleAnalysis)

	fs := http.FileServer(http.Dir("templates"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	port := "5000"
	go openBrowser("http://localhost:" + port)

	fmt.Printf("server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleSelfplay(w http.ResponseWriter, r *http.Request) {
	// chess move calc
}

func handleCalcMove(w http.ResponseWriter, r *http.Request) {

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
