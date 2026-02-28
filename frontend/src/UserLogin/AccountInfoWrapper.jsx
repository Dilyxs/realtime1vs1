import { createContext, useEffect, useState } from "react";
export const UserAccount = createContext({});
const AccountInfoWrapper = ({ children }) => {
  const [AccountInfo, setAccountInfo] = useState({});
  const [isLoaded, setisLoaded] = useState(false);
  useEffect(() => {
    const data = localStorage.getItem("AccountInfo");

    if (data) {
      setAccountInfo(JSON.parse(data));
    }
    setisLoaded(true);
  }, []);

  useEffect(() => {
    if (isLoaded) {
      localStorage.setItem("AccountInfo", JSON.stringify(AccountInfo));
    }
  }, [AccountInfo, isLoaded]);
  return (
    <UserAccount.Provider value={{ AccountInfo, setAccountInfo, isLoaded }}>
      {children}
    </UserAccount.Provider>
  );
};
export default AccountInfoWrapper;
