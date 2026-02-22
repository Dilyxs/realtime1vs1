import { useState, useContext } from "react";

const CreateAccount = () => {
  const [AccountInfo, setAccountInfo] = useContext(UserAccount);
  const [userInfo, setuserInfo] = useState({});
  return (
    <div>
      <h1>Enter Your desired username!</h1>
      <input type="text" value={userInfo?.username ?? ""}></input>
      <h1>Enter Your desired Password!</h1>
      <input type="text" value={userInfo?.password ?? ""}></input>
    </div>
  );
};
export default CreateAccount;
