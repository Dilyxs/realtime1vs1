import { createContext, useState } from "react";
export const UserAccount = createContext({});
const AccountInfoWrapper = ({ children }) => {
  const [AccountInfo, setAccountInfo] = useState({});
  return (
    <UserAccount.Provider value={{ AccountInfo, setAccountInfo }}>
      {children}
    </UserAccount.Provider>
  );
};
export default AccountInfoWrapper;
