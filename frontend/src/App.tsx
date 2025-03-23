import React from 'react';
import { ChessGame } from './components/ChessGame';
import './App.css';

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <h1>Chess Engine</h1>
      </header>
      <main>
        <ChessGame />
      </main>
    </div>
  );
}

export default App;
