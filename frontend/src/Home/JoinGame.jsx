const joingame = async ({ username, gameID }) => {
  const baseurl = import.meta.env.VITE_BACKEND_URL;
  const fetchurl = `${baseurl}/addplayer?roomID=${gameID}`;
  console.log(username);
  const jsondata = JSON.stringify({ username: username });
  try {
    const response = await fetch(fetchurl, {
      method: "POST",
      headers: {
        "Content-Type": "application-json",
      },
      body: jsondata,
    });
    return response.status;
  } catch (e) {
    console.log(e);
    return 500;
  }
};
const JoinGame = ({ setgameID, gameID, username, Navigate }) => {
  return (
    <div>
      <h3>Join a Game!</h3>
      <p>Write the GameID!</p>
      <input
        onChange={(e) => {
          setgameID(e.target.value);
        }}
        value={gameID}
      ></input>
      <button
        onClick={(e) => {
          const action = async () => {
            const data = await joingame({ username, gameID });
            if (data == 200) {
              console.log("move player to game!");
              Navigate("/login");
            }
          };
          action();
        }}
      >
        JoinGame!
      </button>
    </div>
  );
};

export default JoinGame;
