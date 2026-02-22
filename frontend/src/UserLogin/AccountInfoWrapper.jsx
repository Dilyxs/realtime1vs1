import { createContext, useState } from "react";
const UserAccount = createContext(null);
const AccountInfoWrapper = ({ children }) => {
  const [AccountInfo, setAccountInfo] = useState({});
  return (
    <UserAccount.Provider value={{ AccountInfo, setAccountInfo }}>
      {children}
    </UserAccount.Provider>
  );
};
export default AccountInfoWrapper;
