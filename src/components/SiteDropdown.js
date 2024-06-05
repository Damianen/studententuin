import React, { useEffect, useState } from "react";
import axios from "axios";

export default function AccountDropdown() {
  const [isOpen, setIsOpen] = React.useState(false);
  const [user, setUser] = useState(null);

  useEffect(() => {
    // Maak een HTTP-verzoek om de gebruikersgegevens op te halen
    const fetchUser = async () => {
      try {
        const response = await axios.get("/api/getUserByEmailFromSession");
        setUser(response.data);
      } catch (error) {
        console.error("Error fetching user:", error);
      }
    };

    fetchUser(); // Roep de functie voor het ophalen van gebruikersgegevens aan wanneer het component gemonteerd is
  }, []); // De lege array als tweede argument zorgt ervoor dat useEffect alleen wordt uitgevoerd bij de eerste render

  const toggleDropdown = () => {
    setIsOpen(!isOpen);
  };

  const closeDropdown = () => {
    setIsOpen(false);
  };

  console.log(user);

  return (
    <div>
      <div className="relative">
        <button
          onClick={toggleDropdown}
          className="relative overflow-hidden focus:outline-none focus:border-white text-lg font-semibold leading-6 text-white"
        >
          <a>
            {user ? user.SubDomainName + ".studententuin.nl" : "Loading..."}
          </a>
        </button>
        <button
          className={
            isOpen
              ? " cursor-default bg-black opacity-50 fixed inset-0 w-full h-full"
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
          <a
            href="#"
            className="block px-4 py-2 text-gray-800 hover:bg-indigo-500 hover:text-white"
          >
            Account settings
          </a>
          <a
            href="#"
            className="block px-4 py-2 text-gray-800 hover:bg-indigo-500 hover:text-white"
          >
            Support
          </a>
          <a
            href="/logout"
            className="block px-4 py-2 text-gray-800 hover:bg-indigo-500 hover:text-white"
          >
            Sign Out
          </a>
        </div>
      </div>
    </div>
  );
}
