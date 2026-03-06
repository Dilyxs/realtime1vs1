import { useContext, useEffect, useState } from "react";
import { UserAccount } from "../UserLogin/AccountInfoWrapper";
import { useNavigate } from "react-router-dom";
import JoinGame from "./JoinGame";
import CreateGame from "./CreateGame";
const Homepage = () => {
  const Navigate = useNavigate();
  const { AccountInfo, setAccountInfo } = useContext(UserAccount);
  const [isLogedIn, setisLogedin] = useState(false);
  const [gameID, setgameID] = useState("");
  useEffect(() => {
    if (
      AccountInfo?.username != undefined &&
      AccountInfo?.password != undefined
    ) {
      setisLogedin(true);
    }
  }, [AccountInfo]);

  if (!isLogedIn) {
    return (
      <div>
        <div>Pls Login In!</div>
        <button
          onClick={() => {
            Navigate("/login");
          }}
        >
          Login Page
        </button>
      </div>
    );
  }
  return (
    <div className="flex flex-col">
      <div>
        <h1>Hi {AccountInfo.username}</h1>
        <CreateGame
          username={AccountInfo.username}
          Navigate={Navigate}
        ></CreateGame>
        <JoinGame
          setgameID={setgameID}
          gameID={gameID}
          username={AccountInfo.username}
          Navigate={Navigate}
        ></JoinGame>
      </div>
      <div>
        <h1>Leaderboard! : Coming Soon!</h1>
      </div>
    </div>
  );
};

export default Homepage;
