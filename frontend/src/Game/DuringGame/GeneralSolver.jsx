const GeneralSolver = ({
  GameNicheInfo,
  UserWrittenSolution,
  setUserWrittenSolution,
}) => {
  const { problemDescription } = GameNicheInfo;
  return (
    <div>
      <div className="flex flex-col">
        <div>
          <h1>Description:</h1>
          <p>{problemDescription}</p>
        </div>
        <div>
          {
            //TODO: Gonna need a package to ahve a much prettier UI possibly with LSP!
          }
          <h1>Write your code here:</h1>
          <input
            type="text"
            value={UserWrittenSolution}
            onChange={(e) => {
              setUserWrittenSolution(e.target.value);
            }}
          ></input>
          {
            //NOTE: for testing purpose let's make a 15 min counter ideally the chosenquestion amount
            //<Timer></Timer>
          }
          <button
            onClick={() => {
              console.log("submit result!");
            }}
          >
            <h1>Submit!</h1>
          </button>
        </div>
      </div>
    </div>
  );
};

export default GeneralSolver;
