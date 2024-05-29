import React from "react";
import SiteDropdown from "./components/SiteDropdown.js";
import Logs from "./portalComponents/Logs.js";
import Analytics from "./portalComponents/Analytics.js";
import Domain from "./portalComponents/Domain.js";
import Deployment from "./portalComponents/Deployment.js";
import Forms from "./portalComponents/Forms.js";
import SiteConfiguration from "./portalComponents/SiteConfiguration.js";
import Git from "./portalComponents/Git.js";

function Manage() {
  const [selectedItem, setSelectedItem] = React.useState("Logs");
  return (
    <div>
      <header className="bg-house-green">
        {/* <nav className="mx-auto flex max-w-full items-center justify-start p-6">
          <div className="flex flex-1 justify-start gap-4 text-center items-center">
            <a href="#" className="-m-1.5 p-1.5">
              <img src="logo.png" className="w-16 h-16" />
            </a>
            <button className="inline-block rounded-md border border-transparent bg-primary-green px-8 py-2 text-center font-medium text-white hover:bg-green-400">
              Upgraden
            </button>
            <a href="#" className="text-lg font-semibold leading-6 text-white">
              Handleiding
            </a>
          </div>
          <div className="relative inline-block text-left">
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
        </nav> */}
        <nav class="bg-white border-gray-200 dark:bg-house-green">
          <div class="max-w-screen-xl flex flex-wrap items-center justify-between mx-auto p-4">
            <a href="#" className="-m-1.5 p-1.5">
              <img src="logo.png" className="w-16 h-16" />
            </a>
            <button
              data-collapse-toggle="navbar-default"
              type="button"
              class="inline-flex items-center p-2 w-10 h-10 justify-center text-sm text-gray-500 rounded-lg md:hidden hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-gray-200 dark:text-gray-400 dark:hover:bg-gray-700 dark:focus:ring-gray-600"
              aria-controls="navbar-default"
              aria-expanded="false"
            >
              <span class="sr-only">Open main menu</span>
              <svg
                class="w-5 h-5"
                aria-hidden="true"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 17 14"
              >
                <path
                  stroke="currentColor"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M1 1h15M1 7h15M1 13h15"
                />
              </svg>
            </button>
            <div class="hidden w-full md:block md:w-auto" id="navbar-default">
              <ul class="font-medium flex flex-col p-4 md:p-0 mt-4 border border-gray-100 rounded-lg md:flex-row md:space-x-8 rtl:space-x-reverse md:mt-0 md:border-0 md:bg-white dark:bg-house-green md:dark:bg-house-green dark:border-gray-700">
                <li>
                  <button className="inline-block rounded-md border border-transparent bg-primary-green px-8 py-2 text-center font-medium text-white hover:bg-green-400">
                    Upgraden
                  </button>
                </li>
                <li>
                  <a
                    href="#"
                    class="inline-block px-8 py-2 text-center font-medium text-white"
                  >
                    Handleiding
                  </a>
                </li>
                <li>
                  <div className="mx-auto flex max-w-full items-center justify-start p-1 gap-1">
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
        </nav>
      </header>

      <div className="flex">
        <div className="w-4/12 bg-light-green">
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
              ""
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
export default Manage;
