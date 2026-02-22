import { useState } from "react";

const requestANewgame = async ({ username }) => {
  const baseurl = import.meta.env.VITE_BACKEND_URL;
  const fetchurl = `${baseurl}/newroom`;
  const jsondata = JSON.stringify({ username: username });
  try {
    const response = await fetch(fetchurl, {
      method: "POST",
      headers: {
        "Content-Type": "application-json",
      },
      body: jsondata,
    });
    const data = await response.json();
    return { ...data, Err: null };
  } catch (e) {
    console.log(e);
    return { Err: e };
  }
};
const CreateGame = ({ username, Navigate }) => {
  const [gameID, setgameID] = useState(-1);
  const [displayGameID, setdisplayGameID] = useState(false);
  return (
    <div>
      <h3>Create A Game!</h3>
      <button
        onClick={() => {
          const action = async () => {
            const data = await requestANewgame(username);
            console.log(data);
            if (data.Err == null) {
              setgameID(data?.id);
              setdisplayGameID(true);
            }
          };
          action();
        }}
      >
        Make A Game!
      </button>
      {displayGameID && (
        <div>
          <p>Share this Game ID With Your friends!: {gameID}</p>
          <p>Ready to Join the Lobby? </p>
          <button
            onClick={(e) => {
              Navigate("/login");
            }}
          >
            Join The Lobby!
          </button>
        </div>
      )}
    </div>
  );
};
export default CreateGame;
