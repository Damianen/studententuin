import { Button } from "@material-tailwind/react";
import React from "react";
import { useNavigate } from "react-router-dom";

export default function AccountDropdownHamburger() {
    const [isOpen, setIsOpen] = React.useState(false);

    const navigate = useNavigate();

    const handleLogout = () => {
        localStorage.removeItem("token");
        console.log("Token removed");
        navigate("/");
        console.log("Navigated to login");
    };

    const toggleDropdown = () => {
        setIsOpen(!isOpen);
    };

    const closeDropdown = () => {
        setIsOpen(false);
    };

    return (
        <div>
            <div className="relative">
                <button
                    onClick={toggleDropdown}
                    data-collapse-toggle="navbar-default"
                    type="button"
                    className="inline-flex items-center p-2 w-10 h-10 justify-center text-sm text-gray-500 rounded-lg md:hidden hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-gray-200 dark:text-gray-400 dark:hover:bg-gray-700 dark:focus:ring-gray-600"
                    aria-controls="navbar-default"
                    aria-expanded="false"
                >
                    <span className="sr-only">Open hoofdmenu</span>
                    <svg
                        className="w-5 h-5"
                        aria-hidden="true"
                        xmlns="http://www.w3.org/2000/svg"
                        fill="none"
                        viewBox="0 0 17 14"
                    >
                        <path
                            stroke="currentColor"
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth="2"
                            d="M1 1h15M1 7h15M1 13h15"
                        />
                    </svg>
                </button>
                <button
                    className={
                        isOpen
                            ? " cursor-default opacity-50 fixed inset-0 w-full h-full"
                            : "hidden"
                    }
                    onClick={closeDropdown}
                    tabIndex="-1"
                />
                <div
                    className={
                        isOpen
                            ? "absolute right-0 mt-2 w-48 bg-white rounded-lg py-2 shadow-xl"
                            : "hidden"
                    }
                >
                    <span>
                        <a
                            href="/logout"
                            className="block px-4 py-2 text-gray-800 hover:bg-green-400 hover:text-white"
                            onClick={handleLogout}
                        >
                            Logout
                        </a>
                    </span>
                </div>
            </div>
        </div>
    );
}
