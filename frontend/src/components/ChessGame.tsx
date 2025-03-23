import React, { useState } from 'react';
import { Chessboard } from 'react-chessboard';
import { Chess, Square } from 'chess.js';
import axios from 'axios';
import { GameState } from '../types/chess';

const API_URL = 'http://localhost:2828';

export const ChessGame: React.FC = () => {
  const [gameState, setGameState] = useState<GameState>({
    game: new Chess(),
    isThinking: false,
    evaluation: null,
    lastMove: null,
  });

  const makeMove = async (from: Square, to: Square) => {
    const newGame = new Chess(gameState.game.fen());
    try {
      // Make player's move
      const moveResult = newGame.move({
        from,
        to,
        promotion: 'q' // Always promote to queen for simplicity
      });

      if (!moveResult) {
        console.error('Invalid player move:', from, to);
        return;
      }

      setGameState(prev => ({ ...prev, game: newGame, lastMove: `${from}${to}` }));
      
      // Get computer's move
      setGameState(prev => ({ ...prev, isThinking: true }));
      
      const response = await axios.post<string>(`${API_URL}/info`, {
        depth: 4,
        fen: newGame.fen(),
      });

      const computerMove = response.data;
      const newGameAfterComputer = new Chess(newGame.fen());
      
      // Convert the backend's move format to chess.js format
      // The backend returns moves like "e2e4" (from-to format)
      const computerFrom = computerMove.slice(0, 2) as Square;
      const computerTo = computerMove.slice(2, 4) as Square;
      
      const computerMoveResult = newGameAfterComputer.move({
        from: computerFrom,
        to: computerTo,
        promotion: 'q' // Always promote to queen for simplicity
      });
      
      if (computerMoveResult) {
        setGameState(prev => ({
          ...prev,
          game: newGameAfterComputer,
          isThinking: false,
          lastMove: computerMove,
        }));
      } else {
        console.error('Invalid computer move:', computerMove);
      }
    } catch (error) {
      console.error('Error making move:', error);
    }
  };

  const onDrop = (sourceSquare: Square, targetSquare: Square) => {
    makeMove(sourceSquare, targetSquare);
    return true;
  };

  return (
    <div className="chess-game">
      <div className="game-info">
        <p>Current Turn: {gameState.game.turn() === 'w' ? 'White' : 'Black'}</p>
        {gameState.isThinking && <p>Computer is thinking...</p>}
        {gameState.lastMove && <p>Last Move: {gameState.lastMove}</p>}
      </div>
      <div className="board-container" style={{ width: '500px', height: '500px' }}>
        <Chessboard
          position={gameState.game.fen()}
          onPieceDrop={onDrop}
          boardOrientation="white"
          customBoardStyle={{
            borderRadius: '4px',
            boxShadow: '0 2px 10px rgba(0, 0, 0, 0.3)',
          }}
          boardWidth={500}
        />
      </div>
    </div>
  );
}; 