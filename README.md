# Chess

Python implementation of a chess engine based on Minimax with Alpha Beta Pruning.

## Description

This is a chess engine implementation in Go that uses the Minimax algorithm with Alpha-Beta pruning for move calculation. It includes both a command-line interface and a web server component.

## Prerequisites

- Go 1.23.2 or higher
- [notnil/chess](https://github.com/notnil/chess) package

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd chess-engine
```

2. Install dependencies:
```bash
go mod download
```

## Usage

### Command Line Interface

Run the engine with a specific FEN position:

```bash
go run main.go "your-fen-string"
```

If no FEN string is provided, it will use a default position.

### Web Server

The engine also provides a web server interface that runs on port 2828.

1. Start the server:
```bash
go run main.go
```

2. Available endpoints:

- `GET /` - Basic health check
- `POST /info` - Calculate best move for a given position
  ```json
  {
    "depth": 5,
    "fen": "your-fen-string"
  }
  ```
- `POST /analysis` - Analyze chess positions (WIP)

### Environment Variables

- `DEBUG`: Set to any value to enable debug output

## Features

- Minimax algorithm with Alpha-Beta pruning
- Iterative Deepening Search (IDS)
- Position evaluation using piece values and piece-square tables
- Web API for move calculation
- Support for FEN position input

## Project Structure

- `main.go` - Entry point
- `engine/` - Core chess engine implementation
- `controller/` - Web server and API handlers
- `test/` - Test utilities

## License

[Add your license information here]
