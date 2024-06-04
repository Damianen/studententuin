import React from "react";

export default function AccountDropdown() {
  const [isOpen, setIsOpen] = React.useState(false);

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
          <a>voorbeeld.studententuin.nl</a>
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
            href="#"
            className="block px-4 py-2 text-gray-800 hover:bg-indigo-500 hover:text-white"
          >
            Sign Out
          </a>
        </div>
      </div>
    </div>
  );
}
