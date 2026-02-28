import { BrowserRouter, Routes, Route } from "react-router-dom";
import AccountInfoWrapper from "./UserLogin/AccountInfoWrapper";
import CreateAccount from "./UserLogin/CreatAccount";
import Homepage from "./Home/Homepage";
import LoginAccount from "./UserLogin/LoginAccount";
import MainPage from "./Game/MainPage";
import WebsocketRounder from "./Game/WebsocketRounder";
import TokenFetcher from "./Game/TokenFetcher";
function App() {
  return (
    <AccountInfoWrapper>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<LoginAccount></LoginAccount>}></Route>
          <Route path="/home" element={<Homepage></Homepage>}></Route>
          <Route
            path="/game/:roomID"
            element={
              //<WebsocketRounder>
              <TokenFetcher>
                <MainPage></MainPage>
              </TokenFetcher>
              //</WebsocketRounder>
            }
          ></Route>
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
