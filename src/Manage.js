import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import SiteDropdown from "./components/SiteDropdown.js";
import SiteDropdownHamburger from "./components/SiteDropdownHamburger.js";
import Logs from "./portalComponents/Logs.js";
import Analytics from "./portalComponents/Analytics.js";
import Domain from "./portalComponents/Domain.js";
import Deployment from "./portalComponents/Deployment.js";
import Forms from "./portalComponents/Forms.js";
import SiteConfiguration from "./portalComponents/SiteConfiguration.js";
import Git from "./portalComponents/Git.js";
import FileManager from "./portalComponents/FileManager.js";
import FileTree from "./portalComponents/Filetree.js";
import Cookies from "js-cookie";

function Manage() {
  const [selectedItem, setSelectedItem] = useState("Logs");
  const navigate = useNavigate();

  useEffect(() => {
    // Maak een request naar de backend om te controleren of de gebruiker is ingelogd
    fetch("/api/check-session", { credentials: "include" })
      .then((res) => res.json())
      .then((data) => {
        if (!data.isAuthenticated) {
          navigate("/login");
        }
      })
      .catch((error) => {
        console.error("Error:", error);
        navigate("/login");
      });
  }, [navigate]);

  return (
    <div>
      <header className="bg-house-green">
        <nav className="mx-auto flex max-w-full items-center justify-start p-6">
          <div className="flex flex-1 gap-4 text-center items-center">
            <div className="flex justify-start">
              <a href="#" className="-m-1.5 p-1.5">
                <img src="logo.png" className="w-16 h-16" />
              </a>
            </div>
            <div className="ml-auto flex items-center">
              <SiteDropdownHamburger className="inline-block text-left z-10" />

              <div
                className="hidden w-full md:block md:w-auto"
                id="navbar-default"
              >
                <ul className="font-medium flex flex-col p-4 md:p-0 mt-4 border border-gray-100 rounded-lg md:flex-row md:space-x-8 rtl:space-x-reverse md:mt-0 md:border-0 md:bg-white dark:bg-house-green md:dark:bg-house-green dark:border-gray-700">
                  <li>
                    <button className="inline-block rounded-md border border-transparent bg-primary-green px-8 py-2 text-center font-medium text-white hover:bg-green-400">
                      Upgraden
                    </button>
                  </li>
                  <li>
                    <a
                      href="#"
                      className="inline-block px-8 py-2 text-center font-medium text-white"
                    >
                      Handleiding
                    </a>
                  </li>
                  <li>
                    <div className="flex items-center gap-1">
                      <button className="relative block h-8 w-8 border-2 border-gray-600 rounded-full overflow-hidden">
                        <img
                          className="w-full h-full"
                          src="https://images.unsplash.com/photo-1487412720507-e7ab37603c6f?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=256&q=80"
                          alt="Your avatar"
                        />
                      </button>
                      <SiteDropdown className="inline-block text-left z-10" />
                    </div>
                  </li>
                </ul>
              </div>
            </div>
          </div>
        </nav>
      </header>

      <div className="flex">
        <div className="w-4/12 bg-light-green">
          <ul className="py-4">
            <a href="#" onClick={() => setSelectedItem("Logs")}>
              <li className="px-4 py-2 cursor-pointer hover:bg-green-900 hover:text-white">
                Logs email: {user.email}
              </li>
            </a>
            <a href="#" onClick={() => setSelectedItem("Analytics")}>
              <li className="px-4 py-2 cursor-pointer hover:bg-green-900 hover:text-white">
                Analytics
              </li>
            </a>
            <a href="#" onClick={() => setSelectedItem("Domain")}>
              <li className="px-4 py-2 cursor-pointer hover:bg-green-900 hover:text-white">
                Domain
              </li>
            </a>
            <a href="#" onClick={() => setSelectedItem("Deployment")}>
              <li className="px-4 py-2 cursor-pointer hover:bg-green-900 hover:text-white">
                Deployment
              </li>
            </a>
            <a href="#" onClick={() => setSelectedItem("Forms")}>
              <li className="px-4 py-2 cursor-pointer hover:bg-green-900 hover:text-white">
                Forms
              </li>
            </a>
            <a href="#" onClick={() => setSelectedItem("Site Configuration")}>
              <li className="px-4 py-2 cursor-pointer hover:bg-green-900 hover:text-white">
                Site Configuration
              </li>
            </a>
            <a href="#" onClick={() => setSelectedItem("Git")}>
              <li className="px-4 py-2 cursor-pointer hover:bg-green-900 hover:text-white">
                Git
              </li>
            </a>
            <a href="#" onClick={() => setSelectedItem("Files")}>
              <li className="px-4 py-2 cursor-pointer hover:bg-green-900 hover:text-white">
                Files
              </li>
            </a>
          </ul>
        </div>
        <div
          className="w-10/12 bg-white
        "
        >
          <div className="p-4 h-screen overflow-auto">
            <h1 className="text-2xl font-semibold">{selectedItem}</h1>
            {selectedItem === "Logs" ? (
              <Logs />
            ) : selectedItem === "Analytics" ? (
              <Analytics />
            ) : selectedItem === "Domain" ? (
              <Domain />
            ) : selectedItem === "Deployment" ? (
              <Deployment />
            ) : selectedItem === "Forms" ? (
              <Forms />
            ) : selectedItem === "Site Configuration" ? (
              <SiteConfiguration />
            ) : selectedItem === "Git" ? (
              <Git />
            ) : (
              <FileTree />
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
export default Manage;
