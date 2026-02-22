import { useState, useContext } from "react";
import { UserAccount } from "./AccountInfoWrapper";
import { useNavigate } from "react-router-dom";

const submitNewUser = async (logindata) => {
  const baseurl = import.meta.env.VITE_BACKEND_URL;
  const fetchurl = `${baseurl}/createuser`;
  const jsonData = JSON.stringify(logindata);
  try {
    const response = await fetch(fetchurl, {
      method: "POST",
      headers: { "Content-Type": "application-json" },
      body: jsonData,
    });
    return response.status;
  } catch (e) {
    console.log("can't reach server!");
    return 500;
  }
};

const CreateAccount = () => {
  const Navigate = useNavigate();
  const { AccountInfo, setAccountInfo } = useContext(UserAccount);
  const [userInfo, setuserInfo] = useState({});
  return (
    <div>
      <h1>Enter Your desired username!</h1>
      <input
        type="text"
        value={userInfo.username}
        onChange={(e) => {
          setuserInfo(e.target.value);
        }}
      ></input>
      <h1>Enter Your desired Password!</h1>
      <input
        type="text"
        value={userInfo?.password}
        onChange={(e) => {
          setuserInfo(e.target.value);
        }}
      ></input>
      <button
        onClick={() => {
          const action = async () => {
            const status = await submitNewUser({ userInfo });
            if (status >= 200 && status < 300) {
              setAccountInfo(userInfo);
              Navigate("/home");
            } else {
              console.log(status);
            }
          };
          action();
        }}
      >
        Create Account!
      </button>
    </div>
  );
};
export default CreateAccount;
