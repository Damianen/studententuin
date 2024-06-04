import React from "react";
import { Link, useLocation } from "react-router-dom";
import Joyride, { ACTIONS, EVENTS, ORIGIN, STATUS } from 'react-joyride';

export default function Navigation() {
  const [dialog, setDialog] = React.useState(false);
  const location = useLocation();

  const toggleDialog = () => {
    setDialog(!dialog);
  };

  const handleScrollToTech = () => {
    if (location.pathname === "/") {
      const techSection = document.getElementById("technologieen");
      techSection.scrollIntoView({ behavior: "smooth" });
    } else {
      window.location.href = "/#technologieen";
    }
  };

  const handleScrollToPrice = () => {
    if (location.pathname === "/") {
      const techSection = document.getElementById("pakketen");
      techSection.scrollIntoView({ behavior: "smooth" });
    } else {
      window.location.href = "/#pakketen";
    }
  };

  const [joyrideRun, setRun] = React.useState(false);

  const handleClickStart = () => {
    setRun(true);
  };

const joyrideSteps = [
  {
    target: '.joyride-step-1',
    content: 'Vraag hier makkelijk je sub-domein aan!',
    disableBeacon: true
  },
  {
    target: '.joyride-step-2',
    content: 'Een makkelijke handleiding om te weten hoe jij je sub-domein kan beheren!',
  },
  {
    target: '.joyride-step-3',
    content: 'Kies een pakket dat bij jou past!',
  },
  {
    target: '.joyride-step-4',
    content: 'Meerdere technologieën die wij ondersteunen!',
  },
];

const handleJoyrideCallback = (data = null) => {
  const { action, index, origin, status, type } = data;
  
  if ([STATUS.FINISHED, STATUS.SKIPPED].includes(status)) {
    // You need to set our running state to false, so we can restart if we click start again.
    setRun(false);
  }
};

const joyrideStyling = {
  arrowColor: '#fff',
  backgroundColor: '#fff',
  beaconSize: 36,
  overlayColor: 'rgba(0, 0, 0, 0.5)',
  primaryColor: '#f04',
  spotlightShadow: '0 0 15px rgba(0, 0, 0, 0.5)',
  textColor: '#333',
  width: undefined,
  zIndex: 100,
};

  return (
    <header className="bg-primary-green">
      <Joyride callback={handleJoyrideCallback} steps={joyrideSteps} continuous={true} run={joyrideRun} showSkipButton={true} 
      styles={{
          options: {
            primaryColor: '#006241'
          },
        }}
        locale={{ back: 'Terug', close: 'Afsluiten', last: 'Afsluiten', next: 'Volgende', skip: 'Overslaan' }}
        
        />
      <nav className="lg:mx-auto flex max-w-full lg:items-center justify-between p-6">
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
            className="-m-2.5 inline-flex items-center justify-center rounded-md p-2.5 text-house-green"
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
          <button
            onClick={handleScrollToPrice}
            className="text-lg font-semibold leading-6 text-white"
          >
            Pakketen
          </button>

          <button
            onClick={handleScrollToTech}
            className="text-lg font-semibold leading-6 text-white"
          >
            Technologien
          </button>
          <Link
            to="/requestForm"
            className="text-lg font-semibold leading-6 text-white"
          >
            Aanvragen
          </Link>
          <Link
            to="/handleiding"
            className="text-lg font-semibold leading-6 text-white"
          >
            Handleiding
          </Link>
          <button onClick={handleClickStart} className="inline-block border border-transparent bg-house-green px-8 py-2 text-center font-medium text-white hover:bg-light-green hover:text-black">
            Introductie
          </button>
        </div>
        <div className="lg:flex lg:flex-shrink-0 lg:flex-grow lg:justify-end lg:items-center lg:gap-4 hidden">
          <Link
            to="/login"
            className="text-lg font-semibold leading-6 text-white"
          >
            Login
          </Link>
          <Link to="/requestform">
            <button className="inline-block border border-transparent bg-house-green px-8 py-2 text-center font-medium text-white hover:bg-light-green hover:text-black">
              Account aanmaken
            </button>
          </Link>
          <Link to="/manage">
            <button className="inline-block border border-transparent bg-house-green px-8 py-2 text-center font-medium text-white hover:bg-light-green hover:text-black">
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
                  <button
                    onClick={handleScrollToPrice}
                    className="-mx-3 block rounded-lg px-3 py-2 text-base font-semibold leading-7 text-gray-900 hover:bg-gray-50"
                  >
                    Pakketen
                  </button>
                  <button
                    onClick={handleScrollToTech}
                    className="-mx-3 block rounded-lg px-3 py-2 text-base font-semibold leading-7 text-gray-900 hover:bg-gray-50"
                  >
                    Technologien
                  </button>
                  <Link
                    to="/requestForm"
                    className="-mx-3 block rounded-lg px-3 py-2 text-base font-semibold leading-7 text-gray-900 hover:bg-gray-50"
                  >
                    Aanvragen
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
