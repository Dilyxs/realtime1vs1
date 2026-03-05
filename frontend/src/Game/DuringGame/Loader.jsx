import { useEffect, useState } from "react";

//:TODO: Put this at the foregrond/ in-frint at all time
const Loader = ({ GameNicheInfo }) => {
  const [counter, setcounter] = useState(10);
  useEffect(() => {
    const intervalID = setInterval(() => {
      setcounter((prev) => (prev > 0 ? prev - 1 : 0));
    }, 1000);
    return () => clearInterval(intervalID);
  });
  const {
    problemTopic,
    problemTimeRequired,
    problemDifficulty,
    problemDescription,
  } = GameNicheInfo;
  return (
    <div className="fixed z-50 bg-white p-4">
      <h1>{problemTopic} contest...</h1>
      <div>
        <p>{problemDescription}</p>
      </div>
      <div>
        <p>Time Estimate: {problemTimeRequired}</p>
        <p>Difficulty:{problemDifficulty}</p>
      </div>
      <h2>Ready To Start? </h2>
      <h1>{counter}</h1>
    </div>
  );
};

export default Loader;
