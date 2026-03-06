import { useState } from "react";
import { answerNicheQuestion } from "./helperfunc";

const GeneralSolver = ({
  GameNicheInfo,
  UserWrittenSolution,
  setUserWrittenSolution,
  roomID,
  username,
}) => {
  const { problemDescription } = GameNicheInfo;
  const [questionAnswered, setquestionAnswered] = useState(false);
  return (
    <div>
      <div className="flex flex-col">
        <div>
          <h1>Description:</h1>
          <p>{problemDescription}</p>
        </div>
        <div>
          {
            //TODO: Gonna need a package to ahve a much prettier UI possibly with LSP!
          }
          <h1>Write your code here:</h1>
          <input
            rows={10}
            cols={50}
            type="text"
            value={UserWrittenSolution}
            onChange={(e) => {
              setUserWrittenSolution(e.target.value);
            }}
          ></input>
          {
            //:NOTE: for testing purpose let's make a 15 min counter ideally the chosenquestion amount
            //<Timer></Timer>
          }
          {!questionAnswered && (
            <div>
              <button
                onClick={async () => {
                  const response = await answerNicheQuestion(
                    roomID,
                    username,
                    UserWrittenSolution,
                  );
                  if (response?.succesful) {
                    setquestionAnswered(true);
                  }
                }}
              ></button>
              <h1>Submit!</h1>
            </div>
          )}
          {questionAnswered && (
            <div>
              <p>Wait for timer to finish or everyone to submit</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default GeneralSolver;
