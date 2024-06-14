import React, { useState, useEffect } from "react";
import { NavLink, useNavigate } from "react-router-dom";
import SiteDropdown from "./components/SiteDropdown.js";
import SiteDropdownHamburger from "./components/SiteDropdownHamburger.js";
import Logs from "./portalComponents/Logs.js";
import Analytics from "./portalComponents/Analytics.js";
import Domain from "./portalComponents/Domain.js";
import Deployment from "./portalComponents/Deployment.js";
import SiteConfiguration from "./portalComponents/SiteConfiguration.js";
import Git from "./portalComponents/Git.js";
import Cookies from "js-cookie";
import FileManager from "./portalComponents/FileManager.js";
import FileTree from "./portalComponents/Filetree.js";

function Manage() {
    const [selectedItem, setSelectedItem] = useState("Bestanden");
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
            <header className="bg-primary-green">
                <nav className="mx-auto flex max-w-full items-center justify-start p-6">
                    <div className="flex flex-1 gap-4 text-center items-center">
                        <div className="flex justify-start">
                            <NavLink to="/" className="-m-1.5 p-1.5">
                                <img src="logo.png" className="w-16 h-16" />
                            </NavLink>
                        </div>
                        <div className="ml-auto flex items-center">
                            <SiteDropdownHamburger className="inline-block text-left z-10" />

                            <div
                                className="hidden w-full md:block md:w-auto"
                                id="navbar-default"
                            >
                                <ul className="font-medium flex flex-col p-4 md:p-0 mt-4 border border-gray-100 rounded-lg md:flex-row md:space-x-8 rtl:space-x-reverse md:mt-0 md:border-0">
                                    <li>
                                        <NavLink
                                            to="/handleiding"
                                            className="inline-block px-8 py-3 text-center font-medium text-white"
                                        >
                                            Handleiding
                                        </NavLink>
                                    </li>

                                    <li>
                                        <div className="flex items-center gap-1">
                                            <div className="inline-block px-8 py-2 text-center font-medium text-white z-10">
                                                <SiteDropdown />
                                            </div>
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
                        <a
                            href="#"
                            onClick={() => setSelectedItem("Bestanden")}
                        >
                            <li className="px-4 py-2 cursor-pointer hover:bg-green-900 hover:text-white">
                                Bestanden
                            </li>
                        </a>
                        <a href="#" onClick={() => setSelectedItem("Opslag")}>
                            <li className="px-4 py-2 cursor-pointer hover:bg-green-900 hover:text-white">
                                Opslag
                            </li>
                        </a>
                        <a href="#" onClick={() => setSelectedItem("Domein")}>
                            <li className="px-4 py-2 cursor-pointer hover:bg-green-900 hover:text-white">
                                Domein
                            </li>
                        </a>
                        <a
                            href="#"
                            onClick={() => setSelectedItem("Deployment")}
                        >
                            <li className="px-4 py-2 cursor-pointer hover:bg-green-900 hover:text-white">
                                Deployment
                            </li>
                        </a>
                        <a
                            href="#"
                            onClick={() =>
                                setSelectedItem(
                                    "Omgevings variablen/post build commands"
                                )
                            }
                        >
                            <li className="px-4 py-2 cursor-pointer hover:bg-green-900 hover:text-white">
                                Omgevings variablen/post build commands
                            </li>
                        </a>
                        <a href="#" onClick={() => setSelectedItem("Git")}>
                            <li className="px-4 py-2 cursor-pointer hover:bg-green-900 hover:text-white">
                                Git
                            </li>
                        </a>
                        <a href="#" onClick={() => setSelectedItem("Logboek")}>
                            <li className="px-4 py-2 cursor-pointer hover:bg-green-900 hover:text-white">
                                Logboek
                            </li>
                        </a>
                    </ul>
                </div>
                <div
                    className="w-10/12 bg-white
        "
                >
                    <div className="p-4 h-screen overflow-auto z-5">
                        <h1 className="text-2xl font-semibold">
                            {selectedItem}
                        </h1>
                        {selectedItem === "Bestanden" ? (
                            <FileTree />
                        ) : selectedItem === "Opslag" ? (
                            <Analytics />
                        ) : selectedItem === "Domein" ? (
                            <Domain />
                        ) : selectedItem === "Deployment" ? (
                            <Deployment />
                        ) : selectedItem ===
                          "Omgevings variablen/post build commands" ? (
                            <SiteConfiguration />
                        ) : selectedItem === "Git" ? (
                            <Git />
                        ) : (
                            <div className="z-0">
                                <Logs />
                            </div>
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
}
export default Manage;
