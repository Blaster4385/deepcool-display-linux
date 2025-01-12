import styles from "./CustomGrid.module.css";
import {
  GenerateHexCommand,
  ExportGrid,
  SendPattern,
} from "../../../wailsjs/go/main/App";
import { useState } from "preact/hooks";

const CustomGrid = () => {
  const gridSize = 14;
  const [grid, setGrid] = useState(
    Array.from({ length: gridSize }, () => Array(gridSize).fill(false)),
  );
  const [history, setHistory] = useState([]);
  const [isDragging, setIsDragging] = useState(false);
  const [currentState, setCurrentState] = useState(null);

  const toggleCell = (row, col, activate) => {
    setGrid((prevGrid) => {
      const newGrid = prevGrid.map((rowArray, rowIndex) =>
        rowArray.map((cell, colIndex) =>
          rowIndex === row && colIndex === col ? activate : cell,
        ),
      );
      return newGrid;
    });
  };

  const handleMouseDown = (row, col) => {
    setHistory((prevHistory) => [...prevHistory, grid.map((row) => [...row])]);
    setIsDragging(true);
    setCurrentState(!grid[row][col]);
    toggleCell(row, col, !grid[row][col]);
  };

  const handleMouseOver = (row, col) => {
    if (isDragging) {
      toggleCell(row, col, currentState);
    }
  };

  const handleMouseUp = () => {
    setIsDragging(false);
  };

  const undo = () => {
    if (history.length > 0) {
      setGrid(history[history.length - 1]);
      setHistory((prevHistory) => prevHistory.slice(0, -1));
    }
  };

  const clearGrid = () => {
    setHistory((prevHistory) => [...prevHistory, grid.map((row) => [...row])]);
    setGrid(
      Array.from({ length: gridSize }, () => Array(gridSize).fill(false)),
    );
  };

  const sendGrid = async () => {
    try {
      console.log("Grid layout:", grid);
      const hex = await SendPattern(grid);
      console.log("Hex:", hex);
    } catch (e) {
      console.log(e);
    }
  };

  const exportGrid = async () => {
    try {
      console.log("Grid layout:", grid);
      const fileName = prompt("Enter a filename for your layout:");
      const hex = await ExportGrid(fileName, grid).then((result) => {
        alert("Layout exported successfully!");
      });
      console.log("Hex:", hex);
    } catch (e) {
      console.log(e);
    }
  };

  return (
    <div onMouseUp={handleMouseUp} onMouseLeave={handleMouseUp}>
      <div class={styles.container}>
        {grid.map((row, rowIndex) => (
          <div key={rowIndex} class={styles.row}>
            {row.map((cell, colIndex) => (
              <div
                key={colIndex}
                class={`${styles.cell} ${cell ? styles.active : ""}`}
                onMouseDown={() => handleMouseDown(rowIndex, colIndex)}
                onMouseOver={() => handleMouseOver(rowIndex, colIndex)}
              ></div>
            ))}
          </div>
        ))}
      </div>
      <button class={styles.button} onClick={undo}>
        Undo
      </button>
      <button class={styles.button} onClick={clearGrid}>
        Clear
      </button>
      <button class={styles.button} onClick={sendGrid}>
        Send Layout
      </button>
      <button class={styles.button} onClick={exportGrid}>
        Export Layout
      </button>
    </div>
  );
};

export default CustomGrid;
