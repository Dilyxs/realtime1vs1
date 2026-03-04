import { LetServerKnowAboutReadiness } from "./helperfunc";

const PlayerReadyOption = ({ ws, username }) => {
  return (
    <div>
      <h1>Are you ready?</h1>
      <button onClick={() => LetServerKnowAboutReadiness(ws, false, username)}>
        No
      </button>
      <button onClick={() => LetServerKnowAboutReadiness(ws, true, username)}>
        Yes
      </button>
    </div>
  );
};

export default PlayerReadyOption;
