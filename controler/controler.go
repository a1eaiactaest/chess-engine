package controler

import (
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

func main() {
	config := Config{
		Debug: os.Getenv("DEBUG") != "",
		SP:    os.Getenv("SP") != "",
	}

	if !config.Debug {
		log.SetOutput(nil)
	}

	http.HandleFunc("/", handleIndex)
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
