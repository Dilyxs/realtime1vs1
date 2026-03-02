import React, { createContext, useEffect, useRef, useState } from "react";
import { useNavigate } from "react-router-dom";
export const WebsocketHandler = createContext({});
import { HandleNewMessage } from "./helperfunc.js";

const WebsocketRounder = ({
  children,
  token,
  roomID,
  AccountInfo,
  isLoaded,
}) => {
  const ws = useRef(null);
  const [messages, setMessages] = useState([]);
  const [status, setStatus] = useState("connecting");
  const seenIds = useRef(new Set());
  const [gameState, setgameState] = useState({ 0: {}, 1: {}, 2: {} });

  useEffect(() => {
    if (!AccountInfo?.username || !isLoaded) return;

    const connect = async () => {
      const baseurl = import.meta.env.VITE_BACKEND_WS;
      const ws_url = `${baseurl}/websocketconn?token=${token}&roomid=${roomID}`;
      try {
        const websocket = new WebSocket(ws_url);
        websocket.onopen = () => {
          setStatus("open");
          console.log("connected!");
        };
        websocket.onmessage = (ev) => {
          HandleNewMessage(ev, setgameState, seenIds, setMessages);
        };
        websocket.onclose = () => {
          setStatus("closed");
        };
        websocket.onerror = () => {
          setStatus("error");
        };

        ws.current = websocket;
      } catch (e) {
        console.error(e);
        if (!cancelled) setStatus("error");
      }
    };

    connect();

    return () => {
      ws.current?.close();
    };
  }, [AccountInfo, roomID, token, isLoaded]);

  return (
    <WebsocketHandler.Provider value={{ messages, ws, status, gameState }}>
      {React.cloneElement(children, {
        roomID,
        AccountInfo,
        isLoaded,
        messages,
      })}
    </WebsocketHandler.Provider>
  );
};
export default WebsocketRounder;
