import React, { useState, useEffect } from "react";
import axios from "axios";

const Logs = () => {
  const [stdoutLogs, setStdoutLogs] = useState("");
  const [stderrLogs, setStderrLogs] = useState("");
  const [error, setError] = useState(null);

  useEffect(() => {
    axios
      .get("http://localhost:3001/api/logs")
      .then((response) => {
        setStdoutLogs(response.data.stdout || "");
        setStderrLogs(response.data.stderr || "");
      })
      .catch((error) => {
        setError("Error fetching logs");
      });
  }, []);

  return (
    <div>
      {error ? (
        <p>{error}</p>
      ) : (
        <div>
          <h2>Stdout Logs</h2>
          <pre>{stdoutLogs}</pre>
          <h2>Stderr Logs</h2>
          <pre>{stderrLogs}</pre>
        </div>
      )}
    </div>
  );
};

export default Logs;
