import { useState, useEffect, useContext } from "react";
import { preGameStatusFetcher } from "./helperfunc";
import PreGame from "../PreGame/PreGame";
import MidGame from "./DuringGame/MidGame";
import EndGame from "./PostGame/EndGame";
import { WebsocketHandler } from "./WebsocketRounder";

const MainPage = ({ roomID, AccountInfo, isLoaded, messages }) => {
  const [status, setStatus] = useState({});
  const [HasLoaded, setHasLoaded] = useState(false);
  const [GamePhase, setGamePhase] = useState("preGame");
  const { gameState } = useContext(WebsocketHandler);
  const [GameNicheInfo, setGameNicheInfo] = useState({});
  useEffect(() => {
    if (gameState[1]?.hasStarted) {
      setGamePhase("duringGame");
      setGameNicheInfo(gameState[1]?.gameInfo);
    }
  }, [gameState]);
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
  useEffect(() => {});
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
        roomID={roomID}
      ></PreGame>
    );
  }
  if (GamePhase === "duringGame") {
    return (
      <MidGame messages={messages} GameNicheInfo={GameNicheInfo}></MidGame>
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
