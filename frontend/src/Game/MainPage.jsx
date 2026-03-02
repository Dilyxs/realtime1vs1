import { useState, useEffect } from "react";
import { preGameStatusFetcher } from "./helperfunc";
import PreGame from "../PreGame/PreGame";
import MidGame from "./DuringGame/MidGame";
import EndGame from "./PostGame/EndGame";

const MainPage = ({ roomID, AccountInfo, isLoaded, messages }) => {
  const [status, setStatus] = useState({});
  const [HasLoaded, setHasLoaded] = useState(false);
  const [GamePhase, setGamePhase] = useState("preGame");
  useEffect(() => {
    if (!isLoaded) return;
    const action = async () => {
      const response = await preGameStatusFetcher({
        gameid: roomID,
        username: AccountInfo?.username,
      });
      setStatus(response);
      setHasLoaded(true);
    };
    action();
  }, [isLoaded]);
  if (!HasLoaded) {
    return <div>Loading....</div>;
  }
  if (status?.ErrorCode) {
    return <div>Oops Internal Error for ya!</div>;
  }
  if (!status?.isAllowedToJoin) {
    return <div>You do not have permission!</div>;
  }
  if (GamePhase === "preGame") {
    return (
      <PreGame
        isGameMaster={status?.isGameMaster}
        messages={messages}
      ></PreGame>
    );
  }
  if (GamePhase === "duringGame") {
    return (
      <MidGame
        isGameMaster={status?.isGameMaster}
        messages={messages}
      ></MidGame>
    );
  }
  if (GamePhase === "postGame") {
    return (
      <EndGame
        isGameMaster={status?.isGameMaster}
        messages={messages}
      ></EndGame>
    );
  }
};

export default MainPage;
