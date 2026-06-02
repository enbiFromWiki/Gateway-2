import { useState, useEffect } from "react";
import "./App.css";

function App() {
  const [ip, setIp] = useState<string>();

  useEffect(() => {
    async function getIp() {
      const res = await fetch("http://127.0.0.1:8080/api/ip");
      const data = await res.text();
      setIp(data);
    }
    getIp();
  }, []);

  return <div className="ip-holder">{ip}</div>;
}

export default App;
