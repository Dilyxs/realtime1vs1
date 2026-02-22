import { BrowserRouter, Routes, Route } from "react-router-dom";
import AccountInfoWrapper from "./UserLogin/AccountInfoWrapper";
import CreateAccount from "./UserLogin/CreatAccount";
function App() {
  return (
    <AccountInfoWrapper>
      <BrowserRouter>
        <Routes>
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
