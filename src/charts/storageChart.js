import React, { useState, useEffect } from "react";
import { AgChartsReact } from "ag-charts-react";

const StorageChart = () => {
    var totalStorage = 100000;
    var usedStorage = 50000;
    var freeStorage = totalStorage - usedStorage;
    const [pieOptions, setOptions] = useState({
        data: [
        { asset: "Used Storage", amount: usedStorage },
        { asset: "Free Storage", amount: freeStorage },
        ],
        title: {
        text: "Storage Usage",
        },
        series: [
        {
            type: "pie",
            angleKey: "amount",
            calloutLabelKey: "asset",
            sectorLabelKey: "amount",
            sectorLabel: {
            color: "white",
            fontWeight: "bold",
            formatter: ({ value }) => `${(value / 1000).toFixed(0)}MB`,
            },
        },
        ],
    });
    
    return (
        <div className="h-1/3 w-1/3">
        <AgChartsReact className="" options={pieOptions} />
        </div>
    );
    };
    export default StorageChart;