import React, { useState, useEffect } from "react";

const Analytics = () => {
  const [logs, setLogs] = useState([]);
  const usedStorage = 150;
  const totalStorage = 300;
  const storagePercentage = (usedStorage / totalStorage) * 100 
  const storagePercentageString = Math.round(storagePercentage) + "%"    
  useEffect(() => {}, []);

  return (
    <div>
<hr className="h-px my-5 bg-gray-200 border-0 dark:bg-gray-700"></hr>
      <h2>Storage usage</h2>
      
        <div className="w-full flex">
      <div className="w-1/2 bg-gray-200 rounded-full h-8 mb-4 dark:bg-gray-700 flex-1">
        <div
          className="bg-gray-600 h-8 rounded-full dark:bg-gray-300 text-center text-black font-bold py-1"
          style={{ width: `${storagePercentage}%`}}
        ><div>
        {storagePercentageString}
      </div></div>
      </div>
    </div>
    <hr className="h-px my-5 bg-gray-200 border-0 dark:bg-gray-700"></hr>
    </div>
  );
};

export default Analytics;
