import { useEffect, useState } from "react";
import { answerGeneralQuestion } from "./helperfunc";

const GeneralQuestion = ({ gameState, username, roomID }) => {
  const [lastQuestion, setLastQuestion] = useState(null);
  const [hasBeenAnswered, setHasBeenAnswered] = useState(false);
  const [timer, settimer] = useState(5);

  useEffect(() => {
    if (!gameState[1]?.generalQuestions.length) return;

    const latestQuestion =
      gameState[1].generalQuestions[gameState[1].generalQuestions.length - 1];

    if (lastQuestion !== latestQuestion) {
      setLastQuestion(latestQuestion);
      setHasBeenAnswered(false);
      settimer(5);
    }
  }, [gameState, lastQuestion]);
  useEffect(() => {
    if (!lastQuestion || hasBeenAnswered) return;
    const timeoutID = setTimeout(() => setHasBeenAnswered(true), 5000);
    const intervalID = setInterval(() => {
      settimer((prev) => (prev > 0 ? prev - 1 : 0));
    }, 1000);

    return () => {
      clearTimeout(timeoutID);
      clearInterval(intervalID);
    };
  }, [hasBeenAnswered, lastQuestion]);

  if (!lastQuestion || hasBeenAnswered) return null;

  const { question, options, topic, difficulty, questionID } = lastQuestion;

  return (
    <div className="fixed top-0 left-0 w-full z-50 bg-green-900 p-4">
      <h3>{topic}</h3>
      <p>{difficulty}</p>
      <p>{question}</p>
      {options.map((option, id) => (
        <p
          className="cursor-pointer hover:bg-yellow-50"
          key={id}
          onClick={async () => {
            setHasBeenAnswered(true);
            await answerGeneralQuestion(roomID, username, id, questionID);
          }}
        >
          {option}
        </p>
      ))}
      <h3>{timer}</h3>
    </div>
  );
};

export default GeneralQuestion;
