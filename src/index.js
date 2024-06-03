import React from "react";
import ReactDOM from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import "../public/style.css";

import Main from "./Main.js";
import HomePage from "./Homepage.js";
import Login from "./login.js";
import RequestForm from "./RequestForm.js";
import Manage from "./Manage.js";
import About from "./About.js";
import Contact from "./contact.js";

const router = createBrowserRouter([
  {
    path: "/",
    element: (
      <Main>
        <HomePage />
      </Main>
    ),
  },
  {
    path: "login",
    element: (
      <Main>
        <Login />
      </Main>
    ),
  },
  {
    path: "requestForm",
    element: (
      <Main>
        <RequestForm />
      </Main>
    ),
  },
  {
    path: "handleiding",
    element: <Main></Main>,
  },
  {
    path: "manage",
    element: (
      <Main manage={true}>
        <Manage />
      </Main>
    ),
  },
  {
    path: "about",
    element: (
      <Main>
        <About></About>
      </Main>
    ),
  },
  {
    path: "contact",
    element: (
      <Main>
        <Contact />
      </Main>
    ),
  },
  {
    path: "wachtwoord",
    element: <Main></Main>,
  },
]);

const root = ReactDOM.createRoot(document.getElementById("root"));
root.render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
);
