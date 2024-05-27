import React from "react";
import ReactDOM from "react-dom/client";
import "../public/style.css";

import Login from "./login.js";
import Dashboard from "./dashboard.js";

const root = ReactDOM.createRoot(document.getElementById("root"));
root.render(<Dashboard />);
