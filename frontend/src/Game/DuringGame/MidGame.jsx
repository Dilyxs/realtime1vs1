import { useEffect, useState } from "react";
import Loader from "./Loader";
import GeneralQuestion from "./GeneralQuestion";
import GeneralSolver from "./GeneralSolver";

const MidGame = ({
  GameNicheInfo,
  gameState,
  roomID,
  username,
  UserWrittenSolution,
  setUserWrittenSolution,
}) => {
  const [problemLoader, setproblemLoader] = useState(true);
  useEffect(() => {
    const timeoutID = setTimeout(() => {
      setproblemLoader(false);
    }, 10000);

    return () => clearTimeout(timeoutID);
  }, []);
  if (problemLoader) {
    return (
      <Loader
        GameNicheInfo={GameNicheInfo}
        roomID={roomID}
        username={username}
      ></Loader>
    );
  }

  return (
    <div>
      <GeneralQuestion
        gameState={gameState}
        username={username}
        roomID={roomID}
      ></GeneralQuestion>
      <GeneralSolver
        GameNicheInfo={GameNicheInfo}
        UserWrittenSolution={UserWrittenSolution}
        setUserWrittenSolution={setUserWrittenSolution}
      ></GeneralSolver>
    </div>
  );
};

export default MidGame;
