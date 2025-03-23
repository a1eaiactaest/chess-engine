import { Chess } from 'chess.js';

export interface GameState {
  game: Chess;
  isThinking: boolean;
  evaluation: number | null;
  lastMove: string | null;
}

export interface MoveRequest {
  depth: number;
  fen: string;
}

export interface MoveResponse {
  move: string;
  evaluation: number;
} 