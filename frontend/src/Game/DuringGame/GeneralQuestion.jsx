import { useEffect, useState } from "react";
import { answerGeneralQuestion } from "./helperfunc";

const GeneralQuestion = ({ gameState, username, roomID }) => {
  const [lastQuestion, setlastQuestion] = useState(null);
  const [hasBeenAnswered, sethasBeenAnswered] = useState(false);
  useEffect(() => {
    if (gameState[1]?.generalQuestions.length == 0) {
      return;
    }

    if (
      lastQuestion !=
      gameState[1]?.generalQuestions[gameState[1]?.generalQuestions.length - 1]
    ) {
      setlastQuestion(
        gameState[1]?.generalQuestions[
          gameState[1]?.generalQuestions.length - 1
        ],
      );
      sethasBeenAnswered(false);
    }
  }, [gameState]);
  if (lastQuestion == null || hasBeenAnswered) {
    return;
  }
  const { question, options, topic, difficulty } = lastQuestion;
  return (
    <div class="fixed top-0 left-0 w-full z-50 bg-white p-4">
      <h3>{topic}</h3>
      <p>{difficulty}</p>
      <p>{question}</p>
      {options.map((option, id) => (
        <p
          key={id}
          onClick={async (e) => {
            sethasBeenAnswered(true);
            await answerGeneralQuestion(
              roomID,
              username,
              e.key,
              lastQuestion?.questionID,
            );
          }}
        >
          {option}
        </p>
      ))}
    </div>
  );
};

export default GeneralQuestion;
