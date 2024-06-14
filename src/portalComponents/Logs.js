import React, { useState, useEffect } from "react";
import axios from "axios";
import { io } from "socket.io-client";
import Cookies from "js-cookie";

    const socket = io();
const Logs = () => {

    const [stdoutLogs, setStdoutLogs] = useState("");
    const [stderrLogs, setStderrLogs] = useState("");
    const [error, setError] = useState(null);
    const [showStdout, setShowStdout] = useState(true);
    const [showStderr, setShowStderr] = useState(false);
    const [isConnected, setIsConnected] = useState(socket.connected);
    const [state, setState] = useState(false);
    const [user, setUser] = useState(undefined);

    useEffect(() => {
        const fetchUser = async () => {
            try {
                const response = await axios.get("/api/getUserByEmailFromSession");
                console.log("API response:", response.data); // Logging om de response te controleren
                setUser(response.data);
            } catch (error) {
                console.error("Error fetching user:", error);
            }
        };

        fetchUser();
    }, []);

    useEffect(() => {
        function onLogs(logs) {
            setStdoutLogs(logs.stdout);
            setStderrLogs(logs.stderr);
            console.log("logs received");
        }

        function onConnect() {
            console.log("connected");
            setIsConnected(true);
        }

        function onDisconnect() {
            setIsConnected(false);
        }

        socket.on('logs', onLogs);
        socket.on('connect', onConnect);
        socket.on('disconnect', onDisconnect);

        return () => {
            socket.off('connect', onConnect);
            socket.off('disconnect', onDisconnect);
            socket.off('logs', onLogs);
        }
    }, []);

    useEffect(() => {
        const interval = setInterval(() => {
            setState(!state);
            if (user) {
                socket.emit('getLatestLogs', { session_id: user.subDomainName});
            }
        }, 1000);

        return () => clearInterval(interval);
    }, [state])

  const handleToggleStdout = () => {
    setShowStdout((prev) => !prev);
    if (!showStdout) {
      setShowStderr(false); // Hide stderr when toggling stdout
    }
  };

  const handleToggleStderr = () => {
    setShowStderr((prev) => !prev);
    if (!showStderr) {
      setShowStdout(false); // Hide stdout when toggling stderr
    }
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <p className="text-lg mb-4">
        Op deze pagina kunt u de stdout- en stderr-logs bekijken en bedienen.
        Klik op de knop "Show Stdout" om de stdout-logs weer te geven of te
        verbergen. Klik op de knop "Show Stderr" om de stderr-logs weer te geven
        of te verbergen.
      </p>
      {user == undefined ? (
        <p>Loading...</p>
      ) : (
        <div className="space-y-4">
          <button
            className="bg-primary-green hover:bg-green-400 text-white font-bold py-2 px-4 rounded mr-2"
            onClick={handleToggleStdout}
          >
            {showStdout ? "Hide Stdout" : "Show Stdout"}
          </button>
          <button
            className="bg-primary-green hover:bg-green-400 text-white font-bold py-2 px-4 rounded mr-2"
            onClick={handleToggleStderr}
          >
            {showStderr ? "Hide Stderr" : "Show Stderr"}
          </button>
          {showStdout && (
            <div>
              <h1 className="text-xl font-bold">Stdout Logs</h1>
              <pre className="bg-gray-200 p-4 rounded overflow-x-auto max-h-96 whitespace-pre-wrap">
                {stdoutLogs}
              </pre>
            </div>
          )}
          {showStderr && (
            <div>
              <h1 className="text-xl font-bold">Stderr Logs</h1>
              <pre className="bg-gray-200 p-4 rounded overflow-x-auto max-h-96 whitespace-pre-wrap">
                {stderrLogs}
              </pre>
            </div>
          )}
        </div>
      )}
      </div>
    );
};

export default Logs;
