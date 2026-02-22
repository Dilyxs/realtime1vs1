import { BrowserRouter, Routes, Route } from "react-router-dom";
import AccountInfoWrapper from "./UserLogin/AccountInfoWrapper";
import CreateAccount from "./UserLogin/CreatAccount";
import Homepage from "./Home/Homepage";
function App() {
  return (
    <AccountInfoWrapper>
      <BrowserRouter>
        <Routes>
          <Route path="/home" element={<Homepage></Homepage>}></Route>
          <Route
            path="/newaccount"
            element={<CreateAccount></CreateAccount>}
          ></Route>
        </Routes>
      </BrowserRouter>
    </AccountInfoWrapper>
  );
}

export default App;
