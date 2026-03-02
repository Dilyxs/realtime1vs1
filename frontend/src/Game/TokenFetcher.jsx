import React, { useContext, useEffect, useState } from "react";
import { UserAccount } from "../UserLogin/AccountInfoWrapper";
import { useNavigate, useParams } from "react-router-dom";
import { get_Token } from "./helperfunc";

const TokenFetcher = ({ children }) => {
  const Navigate = useNavigate();
  const { roomID } = useParams();
  const { AccountInfo, isLoaded } = useContext(UserAccount);
  const { username, password } = AccountInfo;
  const [token, settoken] = useState("");
  const [tokenReady, settokenReady] = useState(false);
  const [error, seterror] = useState("");
  useEffect(() => {
    const action = async () => {
      const result = await get_Token({ username, password, roomid: roomID });
      seterror(result?.ErrorCode);
      settoken(result?.token);
      settokenReady(true);
    };
    if (isLoaded && username && password && roomID) {
      action();
    }
  }, [isLoaded, username, password, roomID, error]);
  if (!isLoaded) {
    return <div>Loading...</div>;
  }
  if (isLoaded && error) {
    console.log(error);
    return (
      <div>
        <h1>{error == 2 ? "Game Not Found!" : "Unauthorized!"}</h1>
        <button
          onClick={() => {
            Navigate("/home");
          }}
        >
          Go Back Home!
        </button>
      </div>
    );
  }

  if (isLoaded && (!username || !password)) {
    return (
      <div>
        <h1>Pls Login In!</h1>
        <button
          onClick={() => {
            Navigate("/login");
          }}
        >
          Login!
        </button>
      </div>
    );
  }
  return (
    <div>
      {React.cloneElement(children, {
        token,
        roomID,
        AccountInfo,
        isLoaded: tokenReady,
      })}
    </div>
  );
};

export default TokenFetcher;
