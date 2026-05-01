const DisplayScore = ({ username, code, score }) => {
  return (
    <div>
      <h1>User: {username}</h1>
      <div>
        <h3>User Written Code:</h3>
        <p>{code}</p>
      </div>
      <h1>Score: {score}</h1>
      <div className="flex flex-col"></div>
    </div>
  );
};

export default DisplayScore;
