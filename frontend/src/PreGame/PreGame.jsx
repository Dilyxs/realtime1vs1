import { useContext } from "react";
import ClientPOV from "./ClientPOV";
import OwnerPOV from "./OwnerPOV";
import { WebsocketHandler } from "../Game/WebsocketRounder";
import PlayerActive from "./PlayerActive";
import { UserAccount } from "../UserLogin/AccountInfoWrapper";
import PlayerReadyOption from "./PlayerReadyOption";

const PreGame = ({ isGameMaster, roomID }) => {
  const { status, gameState, ws } = useContext(WebsocketHandler);
  const { AccountInfo } = useContext(UserAccount);
  const { username } = AccountInfo;
  if (!status) {
    return <div>Loading...</div>;
  }
  return (
    <div>
      <PlayerActive
        gameState={gameState}
        isGameMaster={isGameMaster}
        username={username}
      ></PlayerActive>
      <div>
        {isGameMaster && <OwnerPOV room_id={roomID}></OwnerPOV>}
        {!isGameMaster && <ClientPOV></ClientPOV>}
      </div>
      <PlayerReadyOption ws={ws} username={username}></PlayerReadyOption>
    </div>
  );
};

export default PreGame;
