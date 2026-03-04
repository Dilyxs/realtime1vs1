import { useState } from "react";

const MidGame = ({ GameNicheInfo }) => {
  const [problemLoader, setproblemLoader] = useState(false);
  return <div>{GameNicheInfo}</div>;
};

export default MidGame;
