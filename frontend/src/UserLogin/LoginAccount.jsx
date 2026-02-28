import { useContext, useEffect, useState } from "react";
import { UserAccount } from "./AccountInfoWrapper";
import { useNavigate } from "react-router-dom";

const logIn = async (logindata) => {
  const baseurl = import.meta.env.VITE_BACKEND_URL;
  const fetchurl = `${baseurl}/login`;
  const jsonData = JSON.stringify(logindata);
  try {
    const response = await fetch(fetchurl, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: jsonData,
    });
    return await response.json();
  } catch (e) {
    console.log(`ran into error: ${e}`);
    return {};
  }
};

//:TODO: check with localStorage if we already rememeber if user had logged in!
const LoginAccount = () => {
  const Navigate = useNavigate();
  const [loginInfo, setloginInfo] = useState({});
  const { AccountInfo, setAccountInfo, isLoaded } = useContext(UserAccount);
  const [ErrorMessage, setErrorMessage] = useState({});
  useEffect(() => {
    if (isLoaded && AccountInfo?.username) {
      console.log("detected");
      Navigate("/home");
    }
  }, [AccountInfo, isLoaded]);
  return (
    <div>
      <h1>Enter Your Username</h1>
      <input
        onChange={(e) => {
          setloginInfo((prev) => ({
            ...prev,
            username: e.target.value,
          }));
        }}
      ></input>

      <h1>Enter Your Password</h1>
      <input
        onChange={(e) => {
          setloginInfo((prev) => ({
            ...prev,
            password: e.target.value,
          }));
        }}
      ></input>
      <button
        onClick={() => {
          const action = async () => {
            const response = await logIn(loginInfo);
            if (response?.valid == true) {
              setAccountInfo(loginInfo);
              setloginInfo({});
              Navigate("/home");
            } else if (response?.valid == false) {
              setErrorMessage({ message: "Wrong Password!" });
            } else if (response?.error_code == 1) {
              setErrorMessage({ message: "User Does Not Exist!" });
            }
          };
          action();
        }}
      >
        Login!
      </button>
      {Object.keys(ErrorMessage).length > 0 && (
        <div>{ErrorMessage?.message}</div>
      )}
    </div>
  );
};

export default LoginAccount;
