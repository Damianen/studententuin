import React from "react";
import SiteDropdown from "./components/SiteDropdown.js";
import AccountDropdown from "./components/AccountDropdown.js";

function Manage() {
  const [selectedItem, setSelectedItem] = React.useState("Logs");
  return (
    <div>
      <header className="bg-green-600">
        <nav className="mx-auto flex max-w-full items-center justify-start p-6">
          <div className="flex flex-1 justify-start gap-4 text-center items-center">
            <a href="#" className="-m-1.5 p-1.5">
              <img src="logo.png" className="w-16 h-16" />
            </a>
            <button className="inline-block rounded-md border border-transparent bg-green-500 px-8 py-2 text-center font-medium text-white hover:bg-green-400">
              Upgraden
            </button>
            <a href="#" className="text-lg font-semibold leading-6 text-white">
              Handleiding
            </a>
          </div>
          <div class="relative inline-block text-left">
            <div className="mx-auto flex max-w-full items-center justify-start p-6 gap-4">
            <button className="relative block h-8 w-8 border-2 border-gray-600 rounded-full overflow-hidden">
                <img
                  className="w-full h-full"
                  src="https://images.unsplash.com/photo-1487412720507-e7ab37603c6f?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=256&q=80"
                  alt="Your avatar"
                />
              </button>
              <SiteDropdown className="inline-block text-left z-10" />
              
            </div>
          </div>
        </nav>
      </header>

      <div className="flex">
        <div className="w-1/4 bg-green-200">
          <ul className="py-4">
            <a href="#" onClick={() => setSelectedItem("Logs")}>
              <li className="px-4 py-2 cursor-pointer hover:bg-green-900 hover:text-white">
                Logs
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
          </ul>
        </div>
        <div className="w-3/4 bg-white">
          <div className="p-4 h-screen">
            <h1 className="text-2xl font-semibold">{selectedItem}</h1>
            <p>This is the {selectedItem} page content.</p>
          </div>
        </div>
      </div>
    </div>
  );
}
export default Manage;
