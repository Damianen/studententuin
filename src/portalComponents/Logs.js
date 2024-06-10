import React, { useState, useEffect } from "react";
import axios from "axios";

const Logs = () => {
  const [stdoutLogs, setStdoutLogs] = useState("");
  const [stderrLogs, setStderrLogs] = useState("");
  const [error, setError] = useState(null);
  const [showStdout, setShowStdout] = useState(true);
  const [showStderr, setShowStderr] = useState(false);

  useEffect(() => {
    axios
      .get("https://studententuin.nl/api/logs")
      .then((response) => {
        setStdoutLogs(response.data.stdout || "");
        setStderrLogs(response.data.stderr || "");
      })
      .catch((error) => {
        setError("Error fetching logs" + error);
      });
  }, []);

  return (
    <div className="container mx-auto px-4 py-8">
      <p className="text-lg mb-4">
        Op deze pagina kunt u de stdout- en stderr-logs bekijken en bedienen.
        Klik op de knop "Show Stdout" om de stdout-logs weer te geven of te
        verbergen. Klik op de knop "Show Stderr" om de stderr-logs weer te geven
        of te verbergen.
      </p>
      {error ? (
        <p className="text-red-500">{error}</p>
      ) : (
        <div className="space-y-4">
          <button
            className="bg-primary-green hover:bg-green-400 text-white font-bold py-2 px-4 rounded mr-2"
            onClick={() => setShowStdout(!showStdout)}
          >
            {showStdout ? "Hide Stdout" : "Show Stdout"}
          </button>
          <button
            className="bg-primary-green hover:bg-green-400 text-white font-bold py-2 px-4 rounded mr-2"
            onClick={() => setShowStderr(!showStderr)}
          >
            {showStderr ? "Hide Stderr" : "Show Stderr"}
          </button>
          {showStdout && (
            <div>
              <h1 className="text-xl font-bold">Stdout Logs</h1>
              <pre className="bg-gray-200 p-4 rounded">{stdoutLogs}</pre>
            </div>
          )}
          {showStderr && (
            <div>
              <h1 className="text-xl font-bold">Stderr Logs</h1>
              <pre className="bg-gray-200 p-4 rounded">{stderrLogs}</pre>
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default Logs;
