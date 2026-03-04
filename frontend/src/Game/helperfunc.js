export const get_Token = async ({ username, password, roomid }) => {
  const baseurl = import.meta.env.VITE_BACKEND_URL;
  const ws_url = `${baseurl}/tokenforws`;
  try {
    const response = await fetch(ws_url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ username, password, roomid: Number(roomid) }),
    });
    return await response.json();
  } catch (e) {
    console.log(e);
  }
};

export const preGameStatusFetcher = async ({ gameid, username }) => {
  const baseurl = import.meta.env.VITE_BACKEND_URL;
  const fetchurl = `${baseurl}/game/${gameid}?username=${username}`;
  try {
    const response = await fetch(fetchurl, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
    });
    return await response.json();
  } catch (e) {
    console.log(e);
    return { ErrorCode: 1, message: "Failed to fetch pre-game status" };
  }
};

export const HandleNewMessage = (ev, setgameState, seenIds, setMessages) => {
  const newMessage = JSON.parse(ev.data);
  if (newMessage.id == null) return;
  if (seenIds.current.has(newMessage?.id)) return;

  setMessages((prev) => [...prev, newMessage]);
  seenIds.current.add(newMessage?.id);

  const { gamePhase, username } = newMessage;
  if (gamePhase == 0) {
    setgameState((prev) => ({
      ...prev,
      [gamePhase]: {
        ...prev[gamePhase],
        [username]: { isReady: newMessage?.isReady ? true : false },
      },
    }));
  } else if (gamePhase == 1) {
    const { gameInfo } = newMessage;
    if (newMessage.hasStarted) {
      setgameState((prev) => ({
        ...prev,
        [gamePhase]: {
          ...prev[gamePhase],
          hasStarted: true,
          gameInfo: gameInfo,
        },
      }));
    }
  }
};
