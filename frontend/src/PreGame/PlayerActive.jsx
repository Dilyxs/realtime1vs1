const PlayerActive = ({ gameState, isGameMaster, username }) => {
  const players = Object.keys(gameState[0]);
  return (
    <div>
      <h2>Players in the room:</h2>
      {players.map((player, idx) => (
        <div key={idx}>
          {
            //:TODO: add more robust owner management, this is legit only show who the lobby master is ONLY FOR THE lobby master's
            //:CONSIDER maybe adding when addding the owner, right away drop off a HubMessage With who the owner is
            //:OR like right, now after detecting it via this detection, drop off a message
          }
          {isGameMaster && username == player && <div>Lobby Master!</div>}
          <h1>{player}</h1>
          <p>is Ready?:{gameState[0][player]?.isReady ? "yes" : "no"}</p>
        </div>
      ))}
    </div>
  );
};

export default PlayerActive;
