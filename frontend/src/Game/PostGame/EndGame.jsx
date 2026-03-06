import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";

import DisplayScore from "./DisplayScore";
const EndGame = ({ finalAnswer, isGameMaster, messages }) => {
  if (finalAnswer == {}) {
    return;
  }
  const [currentIndex, setcurrentIndex] = useState(0);
  const [finished, setfinished] = useState(false);
  const Navigate = useNavigate();
  const { score, code } = finalAnswer;
  const numElements = Object.keys(code).length;
  useEffect(() => {
    const internvalID = setInterval(() => {
      if (currentIndex == numElements - 1) {
        setfinished(true);
        return;
      }
      setcurrentIndex((prev) => prev + 1);
    }, 15000);
    return () => clearInterval(internvalID);
  }, [finalAnswer]);

  return (
    <div>
      {finished && (
        <div>
          <h1>Game Finished!</h1>
          <button
            onClick={() => {
              Navigate("/home");
            }}
          >
            Go back Home
          </button>
        </div>
      )}
      {!finished && (
        <div>
          <DisplayScore
            username={Object.keys(score)[currentIndex]}
            code={code[Object.keys(score)[currentIndex]]}
            score={score[Object.keys(score)[currentIndex]]}
          ></DisplayScore>
        </div>
      )}
    </div>
  );
};

export default EndGame;
