import { useContext } from "react";
import { useParams } from "react-router-dom";
import { WebsocketHandler } from "./WebsocketRounder";

const MainPage = () => {
  //const prop = useContext(WebsocketHandler);
  const { roomID } = useParams();
  console.log(roomID);
  return <div>Hi {roomID}</div>;
};

export default MainPage;
