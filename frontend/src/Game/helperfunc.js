export const get_Token = ({ username, password, roomid }) => {
  const baseurl = import.meta.env.VITE_BACKEND_URL;
  const ws_url = `${baseurl}/tokenforws`;
  try {
    const response = fetch(ws_url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ username, password, roomid: Number(roomid) }),
    });
    return response.json();
  } catch (e) {
    console.log(e);
  }
};
