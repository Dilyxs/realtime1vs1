import { useParams } from "react-router-dom";

const MainPage = () => {
  const { roomID } = useParams();
  console.log(roomID);
  return <div>Hi</div>;
};

export default MainPage;
