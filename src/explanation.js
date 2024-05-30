import React from "react";

function Explanation() {
  return (
    <div className="bg-gradient-to-b from-green-600">
      <div className="bg-white px-4 mb-8 py-8 rounded-3xl mx-auto sm:max-w-xl md:max-w-full lg:max-w-screen-xl md:px-24 lg:px-8">
        <div className="flex flex-col items-center justify-between w-full mb-10 lg:flex-row">
          <div className="mb-16 lg:mb-0 lg:max-w-lg lg:pr-5">
            <div className="max-w-xl mb-6">
              <h2 className="font-sans text-3xl sm:mt-0 mt-6 font-medium tracking-tight text-black sm:text-4xl sm:leading-none max-w-lg mb-6">
                Gratis online omgeving
              </h2>
              <p className="text-black text-base md:text-lg">
                Heb jij een een betrouwbare online omgeving nodig? Wij bieden
                een gratis online omgeving aan voor alle studenten. Hierbij kan
                je gebruik maken van verschillende technieken zoals Node.js of
                wil je gebruik maken van een database en ben je een student?
              </p>
            </div>
            <div className="space-x-4">
              <button className="text-neutral-800 text-lg font-medium inline-flex items-center">
                <span> Vraag nu je online omgeving aan!</span>
              </button>
            </div>
          </div>
          <img
            className="rounded-3xl"
            alt="logo"
            width="420"
            height="120"
            src="https://upload.wikimedia.org/wikipedia/commons/thumb/d/d9/Node.js_logo.svg/2560px-Node.js_logo.svg.png"
          />
        </div>
        <div className="p-1 bg-green-500"></div>
        <br />
      </div>

      <div className="px-4 bg-white mb-8 py-8 rounded-3xl mx-auto sm:max-w-xl md:max-w-full lg:max-w-screen-xl md:px-24 lg:px-8">
        <div className="flex flex-col items-center justify-between w-full mb-10 lg:flex-row">
          <img
            className="rounded-3xl"
            alt="logo"
            width="420"
            height="120"
            src="https://cdn.pixabay.com/photo/2015/09/05/20/02/coding-924920_1280.jpg"
          />
          <div className="mb-16 lg:mb-0 lg:max-w-lg lg:pr-5">
            <div className="max-w-xl mb-6">
              <h2 className="font-sans text-3xl sm:mt-0 mt-6 font-medium tracking-tight text-black sm:text-4xl sm:leading-none max-w-lg mb-6">
                Voor studenten door studenten
              </h2>
              <p className="text-black text-base md:text-lg">
                Studententuin.nl is een platform voor studenten door studenten.
                Iedereen kent het wel, je hebt voor een school project een
                online omgeving nodig maar je wilt natuurlijk niet te veel
                betalen. Bij studententuin.nl kan je gratis een online omgeving
                aanvragen. Indien je meer resources nodig hebt kan je bij ons
                ook terecht. Hierbij kan je voordelig meer resources
                aanschaffen.
              </p>
            </div>
            <div className="space-x-4">
              <button className="text-neutral-800 text-lg font-medium inline-flex items-center">
                <span> Bekijk beschikbare pakketen</span>
              </button>
            </div>
          </div>
        </div>
        <div className="p-1 bg-green-500"></div>
        <br />
      </div>

      <div className="px-4 bg-white mb-8 py-8 rounded-3xl mx-auto sm:max-w-xl md:max-w-full lg:max-w-screen-xl md:px-24 lg:px-8">
        <div className="flex flex-col items-center justify-between w-full mb-10 lg:flex-row">
          <div className="mb-16 lg:mb-0 lg:max-w-lg lg:pr-5">
            <div className="max-w-xl mb-6">
              <h2 className="font-sans text-3xl sm:mt-0 mt-6 font-medium tracking-tight text-black sm:text-4xl sm:leading-none max-w-lg mb-6">
                Handleiding
              </h2>
              <p className="text-black text-base md:text-lg">
                Voorkom dat je vastloopt bij het opzetten van je online
                omgeving. Wij hebben een handleiding geschreven waarin stap voor
                stap wordt uitgelegd hoe je een online omgeving kan aanvragen.
                Hierbij wordt er ook uitgelegd hoe je gebruik kan maken van de
                verschillende technieken die wij aanbieden voor je online
                omgeving.
              </p>
            </div>
            <div className="space-x-4">
              <button className="text-neutral-800 text-lg font-medium inline-flex items-center">
                <span> Bekijk de handleiding</span>
              </button>
            </div>
          </div>
          <img
            className="rounded-3xl"
            alt="logo"
            width="420"
            height="120"
            src="https://images.unsplash.com/photo-1585829365343-ea8ed0b1cb5b?q=80&w=2670&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"
          />
        </div>
        <div className="p-1 bg-green-500"></div>
        <br />
      </div>

      <div>
        <h1 className="font-sans text-3xl pt-12 font-medium tracking-tight text-black mb-6 text-center">
          Pakketen die wij aanbieden
        </h1>
      </div>
      <div className="py-12 flex items-center justify-center px-4 sm:max-w-xl md:max-w-full lg:max-w-screen-xl md:px-24 lg:px-0 mx-auto">
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-20">
          <div className="bg-white rounded-lg overflow-hidden shadow-lg transition-transform transform hover:scale-105">
            <div className="p-1 bg-green-500"></div>
            <div className="p-8">
              <h2 className="text-3xl font-bold text-gray-800 mb-4">
                Gratis pakket
              </h2>
              <p className="text-gray-600 mb-6">
                Beste pakket voor studenten die iets kleins willen bouwen voor
                bijvoorbeeld huiswerk projecten
              </p>
              <p className="text-4xl font-bold text-gray-800 mb-6">Gratis</p>
              <ul className="text-sm text-gray-600 mb-6">
                <li className="mb-2 flex items-center">
                  <svg
                    className="w-4 h-4 mr-2 text-green-500"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M5 13l4 4L19 7"
                    ></path>
                  </svg>
                  300 MB opslag
                </li>
                <li className="mb-2 flex items-center">
                  <svg
                    className="w-4 h-4 mr-2 text-green-500"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M5 13l4 4L19 7"
                    ></path>
                  </svg>
                  Gratis hosting voor 20 weken
                </li>
                <li className="flex items-center">
                  <svg
                    className="w-4 h-4 mr-2 text-green-500"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M5 13l4 4L19 7"
                    ></path>
                  </svg>
                  1 subdomein
                </li>
              </ul>
            </div>
            <div className="p-4">
              <button className="w-full bg-green-500 text-white rounded-full px-4 py-2 hover:bg-green-400 focus:outline-none focus:shadow-outline-green active:bg-green-400">
                Kies dit pakket
              </button>
            </div>
          </div>

          <div className="bg-white rounded-lg overflow-hidden shadow-lg transition-transform transform hover:scale-105">
            <div className="p-1 bg-green-500"></div>
            <div className="p-8">
              <h2 className="text-3xl font-bold text-gray-800 mb-4">
                Basis pakket
              </h2>
              <p className="text-gray-600 mb-6">
                Als je net iets meer nodig hebt voor je projecten
              </p>
              <br />
              <p className="text-4xl font-bold text-gray-800 mb-6">
                €4.99 <span className="text-base">/Maand</span>
              </p>
              <ul className="text-sm text-gray-600 mb-6">
                <li className="mb-2 flex items-center">
                  <svg
                    className="w-4 h-4 mr-2 text-green-500"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M5 13l4 4L19 7"
                    ></path>
                  </svg>
                  10 GB opslag
                </li>
                <li className="mb-2 flex items-center">
                  <svg
                    className="w-4 h-4 mr-2 text-green-500"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M5 13l4 4L19 7"
                    ></path>
                  </svg>
                  onbeperkte hosting tijd
                </li>
                <li className="flex items-center">
                  <svg
                    className="w-4 h-4 mr-2 text-green-500"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M5 13l4 4L19 7"
                    ></path>
                  </svg>
                  2 subdomeinen
                </li>
              </ul>
            </div>
            <div className="p-4">
              <button className="w-full bg-green-500 text-white rounded-full px-4 py-2 hover:bg-green-400 focus:outline-none focus:shadow-outline-green active:bg-green-400">
                Kies dit pakket
              </button>
            </div>
          </div>

          <div className="bg-white rounded-lg overflow-hidden shadow-lg transition-transform transform hover:scale-105">
            <div className="p-1 bg-green-500"></div>
            <div className="p-8">
              <h2 className="text-3xl font-bold text-gray-800 mb-4">
                Premium pakket
              </h2>
              <p className="text-gray-600 mb-6">
                Voor meerdere hobbyprojecten die ieder een eigen subdomein nodig
                heeft
              </p>
              <br />
              <p className="text-4xl font-bold text-gray-800 mb-6">
                €9.99 <span className="text-base">/Maand</span>
              </p>

              <ul className="text-sm text-gray-600 mb-6">
                <li className="mb-2 flex items-center">
                  <svg
                    className="w-4 h-4 mr-2 text-green-500"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M5 13l4 4L19 7"
                    ></path>
                  </svg>
                  25 GB opslag
                </li>
                <li className="mb-2 flex items-center">
                  <svg
                    className="w-4 h-4 mr-2 text-green-500"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M5 13l4 4L19 7"
                    ></path>
                  </svg>
                  onbeperkte hosting tijd
                </li>
                <li className="flex items-center">
                  <svg
                    className="w-4 h-4 mr-2 text-green-500"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M5 13l4 4L19 7"
                    ></path>
                  </svg>
                  5 subdomeinen
                </li>
              </ul>
            </div>
            <div className="p-4">
              <button className="w-full bg-green-500 text-white rounded-full px-4 py-2 hover:bg-green-400 focus:outline-none focus:shadow-outline-green active:bg-green-400">
                Kies dit pakket
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Explanation;
