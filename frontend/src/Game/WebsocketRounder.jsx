import { createContext, useEffect, useRef, useState } from "react";
import { useNavigate } from "react-router-dom";
export const WebsocketHandler = createContext({});

const WebsocketRounder = ({ children, token, roomID, AccountInfo }) => {
  const ws = useRef(null);
  const [messages, setMessages] = useState([]);
  const [status, setStatus] = useState("connecting");
  const navigate = useNavigate();

  useEffect(() => {
    if (!isLoaded || !AccountInfo?.username) {
      navigate("/login");
    }
  }, [AccountInfo, isLoaded, navigate]);

  useEffect(() => {
    if (!AccountInfo?.username) return;

    const connect = async () => {
      try {
        websocket.onopen = () => setStatus("open");
        websocket.onmessage = (ev) => {
          setMessages((prev) => [...prev, ev.data]);
        };
        websocket.onclose = () => {
          if (!cancelled) setStatus("closed");
        };
        websocket.onerror = () => {
          if (!cancelled) setStatus("error");
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
  }, [AccountInfo, roomID]);

  return (
    <WebsocketHandler.Provider value={{ messages, ws, status }}>
      {children}
    </WebsocketHandler.Provider>
  );
};
export default WebsocketRounder;
