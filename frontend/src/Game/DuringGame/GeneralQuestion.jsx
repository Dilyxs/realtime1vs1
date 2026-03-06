import { useEffect, useState } from "react";
import { answerGeneralQuestion } from "./helperfunc";

const GeneralQuestion = ({ gameState, username, roomID }) => {
  const [lastQuestion, setLastQuestion] = useState(null);
  const [hasBeenAnswered, setHasBeenAnswered] = useState(false);

  useEffect(() => {
    let timeoutID;
    if (!gameState[1]?.generalQuestions.length) return;

    const latestQuestion =
      gameState[1].generalQuestions[gameState[1].generalQuestions.length - 1];

    if (lastQuestion !== latestQuestion) {
      setLastQuestion(latestQuestion);
      setHasBeenAnswered(false);
      timeoutID = setTimeout(() => setHasBeenAnswered(true), 5000);
    }
    return () => clearTimeout(timeoutID);
  }, [gameState, lastQuestion]);

  if (!lastQuestion || hasBeenAnswered) return null;

  const { question, options, topic, difficulty, questionID } = lastQuestion;

  return (
    <div className="fixed top-0 left-0 w-full z-50 bg-white p-4">
      <h3>{topic}</h3>
      <p>{difficulty}</p>
      <p>{question}</p>
      {options.map((option, id) => (
        <p
          key={id}
          onClick={async () => {
            setHasBeenAnswered(true);
            await answerGeneralQuestion(roomID, username, option, questionID);
          }}
        >
          {option}
        </p>
      ))}
    </div>
  );
};

export default GeneralQuestion;
