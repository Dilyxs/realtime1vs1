import { v7 as uuidv7 } from "uuid";
export const LetServerKnowAboutReadiness = (ws, isReady, username) => {
  const id = uuidv7().toString();
  const data = {
    type: "userIsReady",
    main: { id, username, isReady, gamePhase: 0 },
  };
  const jsonData = JSON.stringify(data);
  ws.current.send(jsonData);
};

export const StartGameWithServer = async (gameoption, room_id) => {
  const data = {
    room_id: Number(room_id),
    question_topic: gameoption,
  };
  const baseurl = import.meta.env.VITE_BACKEND_URL;
  const fetchurl = `${baseurl}/startgame`;
  const jsonData = JSON.stringify(data);
  try {
    const response = await fetch(fetchurl, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: jsonData,
    });
    if (!response.ok) {
      console.log("Failed to start the game");
      return { Err: "Failed to start the game" };
    }
    const responseData = await response.json();
    return responseData;
  } catch (e) {
    console.log(e);
    return { Err: e };
  }
};
