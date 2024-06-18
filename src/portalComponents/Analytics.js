import React, { useState, useEffect } from "react";

const Analytics = () => {
    const [logs, setLogs] = useState([]);
    const totalStorage = 10000; // get this from database

    const [usedStorage, setUsedStorage] = useState(0);
    const [storagePercentage, setStoragePercentage] = useState(0);
    const [storagePercentageString, setStoragePercentageString] =
        useState("0%");
    const [isLoaded, setIsLoaded] = useState(false);

    useEffect(() => {
        const fetchUsedStorage = async () => {
            const response = await fetch("/dir-info", {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                },
            });

            if (response.ok) {
                const data = await response.json();
                console.log("storage percentage:" + data.storagePercentage);
                setUsedStorage(data.size); // Update usedStorage with the fetched size
                setStoragePercentage(data.storagePercentage);
                setStoragePercentageString(
                    Math.round(data.storagePercentage) + "%"
                );
                setIsLoaded(true);
            }
        };

        fetchUsedStorage(); // Call the async function
    }, []);
    return (
        <div>
            <hr className="h-px my-5 bg-gray-200 border-0 dark:bg-gray-700"></hr>
            <h2>Opslag in gebruik</h2>
            {!isLoaded ? (
                <p>Laden...</p>
            ) : (
                <div>
                    <div className="w-1/2 bg-gray-200 rounded-full h-8 mb-4 dark:bg-gray-700 flex-1">
                        <div
                            className="bg-gray-600 h-8 rounded-full dark:bg-gray-300 text-center text-black font-bold py-1"
                            style={{ width: `${storagePercentage}%` }}
                        >
                            <div>{storagePercentageString}</div>
                        </div>
                    </div>
                    <hr className="h-px my-5 bg-gray-200 border-0 dark:bg-gray-700"></hr>
                </div>
            )}
        </div>
    );
};

export default Analytics;
