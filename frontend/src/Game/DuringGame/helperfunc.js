//TODO: Later Check on if they are null or undefined before proceeding!
export const answerGeneralQuestion = async (
  roomID,
  username,
  optionID,
  questionID,
) => {
  const baseurl = import.meta.env.VITE_BACKEND_URL;
  const fetchurl = `${baseurl}/answergeneralquestion`;
  data = {
    roomID,
    username,
    questionID,
    option: optionID,
  };
  const jsonData = JSON.stringify(data);
  try {
    const response = await fetch(fetchurl, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: jsonData,
    });
    const jsonResponse = await response.json();
    return jsonResponse;
  } catch (e) {
    console.log(e);
    return { Error: e };
  }
};
export const answerNicheQuestion = async (roomID, username, answer) => {
  const baseurl = import.meta.env.VITE_BACKEND_URL;
  const fetchurl = `${baseurl}/answernichequestion`;
  const data = {
    roomID,
    username,
    answer,
  };
  const jsonData = JSON.stringify(data);
  try {
    const response = await fetch(fetchurl, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: jsonData,
    });
    const jsonResponse = await response.json();
    return jsonResponse;
  } catch (e) {
    console.log(e);
    return { Error: e };
  }
};
