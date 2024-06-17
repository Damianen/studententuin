import React, { useState, useEffect } from "react";
import { io } from "socket.io-client";
import axios from "axios";

const socket = io("https://webhook.studententuin.nl", { autoConnect: false });

const Deployment = () => {
    const [isConnected, setIsConnected] = useState(socket.connected);
    const [user, setUser] = useState(undefined);
    const [logs, setLogs] = useState("");
    const [state, setState] = useState(false);

    useEffect(() => {
        const fetchUser = async () => {
            try {
                const response = await axios.get(
                    "/api/getUserByEmailFromSession"
                );
                console.log("API response:", response.data); // Logging om de response te controleren
                setUser(response.data);
                socket.connect();
            } catch (error) {
                console.error("Error fetching user:", error);
            }
        };

        fetchUser();

        return () => {
            socket.disconnect();
        }
    }, []);

    useEffect(() => {
        function onLogs(logs) {
            setLogs(logs.logs);
        }

        function onConnect() {
            console.log("connected");
            setIsConnected(true);
        }

        function onDisconnect() {
            setIsConnected(false);
        }

        socket.on("logs", onLogs);
        socket.on("connect", onConnect);
        socket.on("disconnect", onDisconnect);

        return () => {
            socket.off("connect", onConnect);
            socket.off("disconnect", onDisconnect);
            socket.off("logs", onLogs);
        };
    }, []);

    useEffect(() => {
        const interval = setInterval(() => {
            setState(!state);
            if (user) {
                socket.emit("logRequest", {
                    subdomainName: user.subDomainName,
                });
            }
        }, 1000);

        return () => clearInterval(interval);
    }, [state]);


    return (
        <div className="container mx-auto px-4 py-8">
            {!isConnected ? (
                <p>Laden...</p>
            ) : (
                    <div>
                        <pre className="bg-gray-200 p-4 rounded overflow-x-auto max-h-96 whitespace-pre-wrap">
                            {logs}
                        </pre>
                    </div>
                )}
        </div>
    );
}

export default Deployment;
