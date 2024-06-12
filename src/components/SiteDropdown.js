import React, { useEffect, useState } from "react";
import axios from "axios";

export default function AccountDropdown() {
  const [isOpen, setIsOpen] = useState(false);
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchUser = async () => {
      try {
        const response = await axios.get("/api/getUserByEmailFromSession");
        console.log("API response:", response.data); // Logging om de response te controleren
        setUser(response.data);
      } catch (error) {
        console.error("Error fetching user:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchUser();
  }, []);

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
          className="relative overflow-hidden focus:outline-none focus:border-white text-lg font-semibold leading-6 text-white"
        >
          <a>
            {loading
              ? "Loading..."
              : user && user
              ? `${user.subDomainName}.studententuin.nl`
              : "No user data"}
          </a>
        </button>
        <button
          className={
            isOpen
              ? " cursor-default bg-black opacity-50 fixed inset-0 w-full h-full z-40"
              : "hidden"
          }
          onClick={closeDropdown}
          tabIndex="-1"
        />
        <div
          className={
            isOpen
              ? "absolute right-0 mt-2 w-48 bg-white rounded-lg py-2 shadow-xl z-50"
              : "hidden"
          }
        >
          <a
            href="/logout"
            className="block px-4 py-2 text-gray-800 hover:bg-green-400 hover:text-white"
          >
            Sign Out
          </a>
        </div>
      </div>
    </div>
  );
}
