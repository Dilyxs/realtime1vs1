const verifyUserCanJoinGame = async () => {
  const baseurl = import.meta.env.VITE_BACKEND_URL;
  const fetchurl = `${baseurl}/`;
  try {
    fetch(fetchurl);
  } catch (e) {}
};

const SaveToLocalStorage = (obj) => {
  const keys = Object.keys(obj);
  keys.forEach((key) => {
    localStorage.setItem(key, obj[key]);
  });
};

const GetKeyFromLocalStorage = (key) => {
  return localStorage.getItem(key);
};
