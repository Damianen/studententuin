import React, { useState, useEffect } from 'react';

const Analytics = () => {
    const [logs, setLogs] = useState([]);

    useEffect(() => {
    }, []);

    return (
        <div>
            <h1>Dit zijn de Analytics</h1>
            <ul>
                {logs.map(log => (
                    <li key={log.id}>
                        {}
                        {log.content}
                    </li>
                ))}
            </ul>
        </div>
    );
};

export default Analytics;