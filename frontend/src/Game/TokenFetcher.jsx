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
  useEffect(() => {
    const action = async () => {
      const result = await get_Token({ username, password, roomid: roomID });
      settoken(result?.token);
    };
    if (isLoaded && username && password && roomID) {
      action();
    }
  }, [isLoaded, username, password, roomID]);
  if (!isLoaded) {
    return <div>Loading...</div>;
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
    <div>{React.cloneElement(children, { token, roomID, AccountInfo })}</div>
  );
};

export default TokenFetcher;
