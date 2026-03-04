import { useState } from "react";
import { StartGameWithServer } from "./helperfunc";

const OwnerPOV = ({ room_id }) => {
  const [gameOption, setgameOption] = useState(0);
  return (
    <div>
      <h3>Do you want to start the game?</h3>

      <ul>
        <li>
          <div onClick={() => setgameOption(0)}>React</div>
        </li>
        <li>
          <div onClick={() => setgameOption(1)}>Backend(Go)</div>
        </li>
        <li>
          <div onClick={() => setgameOption(2)}>Embedded(C)</div>
        </li>
        <li>
          <div onClick={() => setgameOption(3)}>BlockChain(Rust)</div>
        </li>
      </ul>
      {
        //TODO: gameOption does not let the user know which mode is on!,replace by text later OR do like a coloring around selected option!
      }
      <p>Selected Mode: {gameOption}</p>

      <button
        onClick={() => {
          StartGameWithServer(gameOption, room_id);
        }}
      >
        Yes!
      </button>
    </div>
  );
};

export default OwnerPOV;
