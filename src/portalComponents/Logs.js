import React, { useState, useEffect } from "react";

const Logs = () => {
  const [logs, setLogs] = useState([]);

  useEffect(() => {
    // Fetch logs data from your API or database
    // Replace the API_URL with your actual API endpoint
    // fetch(API_URL)
    //     .then(response => response.json())
    //     .then(data => setLogs(data))
    //     .catch(error => console.error(error));
    setLogs([
      { id: 1, content: "Log 1" },
      { id: 2, content: "Log 2" },
      { id: 3, content: "Log 3" },
    ]);
  }, []);

  return (
    <div>
      <h1>Dit zijn de logs</h1>
      <ul className="shadow-md rounded-lg">
        {logs.map((log, index) => (
          <div
            key={index}
            className={`p-3 ${
              index % 2 === 0 ? "bg-with" : "bg-neutral-warm"
            }`}
          >
            {log.content}
          </div>
        ))}
      </ul>
    </div>
  );
};

export default Logs;
