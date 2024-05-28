import React from "react";
import { Link } from "react-router-dom";

export default function Navigation() {
  const [dialog, setDialog] = React.useState(false);

  const toggleDialog = () => {
    setDialog(!dialog);
  };

  return (
    <header className="bg-primary-green">
      <nav className="lg:mx-auto flex max-w-full lg:items-center justify-between p-6 ">
        <div className="flex lg:flex-shrink-0 lg:flex-grow-0 lg:justify-start lg:gap-4">
          <Link to="/" className="-m-1.5 p-1.5">
            <img src="logo.png" className="w-16 h-16" />
          </Link>
        </div>
        <div
          className={"flex " + (dialog ? "hidden" : "lg:hidden")}
          onClick={toggleDialog}
        >
          <button
            type="button"
            className="-m-2.5 inline-flex items-center justify-center rounded-md p-2.5 text-primary-green"
          >
            <svg
              className="h-16 w-16"
              fill="none"
              viewBox="0 0 24 24"
              strokeWidth="1.5"
              stroke="currentColor"
              aria-hidden="true"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5"
              />
            </svg>
          </button>
        </div>
        <div className="lg:flex lg:flex-shrink-1 lg:flex-grow lg:justify-center lg:items-center lg:gap-4 lg:px-2 hidden">
          <Link to="/" className="text-lg font-semibold leading-6 text-white">
            Technologien
          </Link>
          <Link to="/requestForm" className="text-lg font-semibold leading-6 text-white">
            Aanvragen
          </Link>
          <Link to="/" className="text-lg font-semibold leading-6 text-white">
            Prijzen
          </Link>
          <Link to="/handleiding" className="text-lg font-semibold leading-6 text-white">
            Handleiding
          </Link>
        </div>
        <div className="lg:flex lg:flex-shrink-0 lg:flex-grow lg:justify-end lg:items-center lg:gap-4 hidden">
          <Link to="/login" className="text-lg font-semibold leading-6 text-white">
            Login
          </Link>
          <Link to="/requestform">
            <button className="inline-block  border border-transparent bg-primary-green px-8 py-2 text-center font-medium text-white hover:bg-house-green">
              Account aanmaken
            </button>
          </Link>
          <Link to="/manage">
            <button className="inline-block  border border-transparent bg-primary-green px-8 py-2 text-center font-medium text-white hover:bg-house-green">
              Jouw omgeving
            </button>
          </Link>
        </div>
        <div
          className={dialog ? "lg:hidden" : "hidden"}
          role="dialog"
          aria-modal="true"
        >
          <div className="fixed inset-0 z-10"></div>
          <div className="fixed inset-y-0 right-0 z-10 w-full overflow-y-auto bg-white px-6 py-6 sm:max-w-sm sm:ring-1 sm:ring-gray-900/10">
            <div className="flex items-center justify-between">
              <Link to="/" className="-m-1.5 p-1.5">
                <img className="h-8 w-auto" src="logo.png" alt="" />
              </Link>
              <button
                type="button"
                className="-m-2.5 rounded-md p-2.5 text-gray-700"
                onClick={toggleDialog}
              >
                <svg
                  className="h-6 w-6"
                  fill="none"
                  viewBox="0 0 24 24"
                  strokeWidth="1.5"
                  stroke="currentColor"
                  aria-hidden="true"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M6 18L18 6M6 6l12 12"
                  />
                </svg>
              </button>
            </div>
            <div className="mt-6 flow-root">
              <div className="-my-6 divide-y divide-gray-500/10">
                <div className="space-y-2 py-6">
                  <Link
                    to="/"
                    className="-mx-3 block rounded-lg px-3 py-2 text-base font-semibold leading-7 text-gray-900 hover:bg-gray-50"
                  >
                    Technologien
                  </Link>
                  <Link
                    to="/requestForm"
                    className="-mx-3 block rounded-lg px-3 py-2 text-base font-semibold leading-7 text-gray-900 hover:bg-gray-50"
                  >
                    Aanvragen
                  </Link>
                  <Link
                    to="/"
                    className="-mx-3 block rounded-lg px-3 py-2 text-base font-semibold leading-7 text-gray-900 hover:bg-gray-50"
                  >
                    Prijzen
                  </Link>
                  <Link
                    to="/handleiding"
                    className="-mx-3 block rounded-lg px-3 py-2 text-base font-semibold leading-7 text-gray-900 hover:bg-gray-50"
                  >
                    Handleiding
                  </Link>
                </div>
                <div className="py-6">
                  <Link
                    to="/manage"
                    className="-mx-3 block rounded-lg px-3 py-2.5 text-base font-semibold leading-7 text-gray-900 hover:bg-gray-50"
                  >
                    Jouw Omgeving
                  </Link>
                  <Link
                    to="/requestForm"
                    className="-mx-3 block rounded-lg px-3 py-2.5 text-base font-semibold leading-7 text-gray-900 hover:bg-gray-50"
                  >
                    Account Aanmaken
                  </Link>
                  <Link
                    to="/login"
                    className="-mx-3 block rounded-lg px-3 py-2.5 text-base font-semibold leading-7 text-gray-900 hover:bg-gray-50"
                  >
                    Log in
                  </Link>
                </div>
              </div>
            </div>
          </div>
        </div>
      </nav>
    </header>
  );
}
