import React from "react";
import ReactDOM from "react-dom/client";
import "../public/style.css";

import Login from "./login.js";
import Dashboard from "./dashboard.js";
import Home from "./home.js";
import Navigation from "./Navigation.js";
import Footer from "./Footer.js";
import Explanation from "./explanation.js";
import Manage from "./Manage.js";

const root = ReactDOM.createRoot(document.getElementById("root"));
root.render(
  <>
  
    <Manage />
    <Footer />
  </>
);
