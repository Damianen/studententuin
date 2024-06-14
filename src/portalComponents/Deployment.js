import React, { useState, useEffect } from "react";

const Deployment = () => {
    const [logs, setLogs] = useState([]);

    useEffect(() => {
        // Fetch logs data from your API or database
        // Replace the API_URL with your actual API endpoint
        // fetch(API_URL)
        //     .then(response => response.json())
        //     .then(data => setLogs(data))
        //     .catch(error => console.error(error));
    }, []);

    return (
        <div>
            <h1>Dit is de Deployment page</h1>
            <ul>
                {logs.map((log) => (
                    <li key={log.id}>
                        {}
                        {log.content}
                    </li>
                ))}
            </ul>
        </div>
    );
};

export default Deployment;
