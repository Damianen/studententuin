import React from "react";
import { Link } from "react-router-dom";

function Homepage() {
  return (
    <div className="bg-gradient-to-b from-primary-green to-neutral-warm ">
      <div className="px-4 bg-white  overflow-hidden shadow-lg mb-6 py-6 mx-auto sm:max-w-xl md:max-w-full lg:max-w-screen-xl md:px-16 lg:px-6">
        <div className="flex flex-col items-center justify-between w-full  lg:flex-row">
          <div className="mb-12 lg:mb-0 lg:max-w-lg lg:pr-4">
            <div className="max-w-xl mb-4">
              <h2 className="font-sans text-3xl sm:mt-0 mt-4 font-medium tracking-tight text-black sm:text-4xl sm:leading-none max-w-lg mb-4">
                Gratis online omgeving
              </h2>
              <div className="p-1 bg-house-green"></div>
              <br />
              <p className="text-black text-base md:text-lg">
                Heb jij een een betrouwbare online omgeving nodig? Wij bieden
                een gratis online omgeving aan voor alle studenten. Hierbij kan
                je gebruik maken van verschillende technologien zoals Node.js of
                wil je gebruik maken van een database en ben je een student? Vul
                de aanvraagformulier in en krijg je eigen online omgeving!
              </p>
            </div>
            <div className="flex justify-center pt-8">
              <Link to="requestfrom">
                <button className="inline-block  border border-transparent bg-primary-green px-6 py-2 text-center font-medium text-white hover:bg-house-green">
                  Vraag nu een online omgeving aan!
                </button>
              </Link>
            </div>
          </div>
          <img
            className=""
            alt="logo"
            width="560"
            height="160"
            src="https://upload.wikimedia.org/wikipedia/commons/thumb/d/d9/Node.js_logo.svg/2560px-Node.js_logo.svg.png"
          />
        </div>
        <br />
      </div>

      <div className="px-4 bg-white overflow-hidden shadow-lg mb-6 py-6 mx-auto sm:max-w-xl md:max-w-full lg:max-w-screen-xl md:px-16 lg:px-6">
        <div className="flex flex-col items-center justify-between w-full lg:flex-row">
          <img
            className=""
            alt="logo"
            width="560"
            height="160"
            src="https://cdn.pixabay.com/photo/2015/09/05/20/02/coding-924920_1280.jpg"
          />
          <div className="mb-12 lg:mb-0 lg:max-w-lg lg:pr-4">
            <div className="max-w-xl mb-4">
              <h2 className="font-sans text-3xl sm:mt-0 mt-4 font-medium tracking-tight text-black sm:text-4xl sm:leading-none max-w-lg mb-4">
                Voor studenten door studenten
              </h2>
              <div className="p-1 bg-house-green"></div>
              <p className="text-black pt-6 text-base md:text-lg">
                Studententuin.nl is een platform voor studenten door studenten.
                Iedereen kent het wel, je hebt voor een school project een
                online omgeving nodig maar je wilt natuurlijk niet te veel
                betalen. Bij studententuin.nl kan je gratis een online omgeving
                aanvragen. Indien je meer resources nodig hebt kan je bij ons
                ook terecht.
              </p>
            </div>
            <div className="flex justify-center pt-12">
              <Link>
                <button className="inline-block border border-transparent bg-primary-green px-6 py-2 text-center font-medium text-white hover:bg-house-green">
                  Bekijk beschikbare pakketen
                </button>
              </Link>
            </div>
          </div>
        </div>
        <br />
      </div>

      <div className="px-4 bg-white overflow-hidden shadow-lg mb-6 py-6 mx-auto sm:max-w-xl md:max-w-full lg:max-w-screen-xl md:px-16 lg:px-6">
        <div className="flex flex-col items-center justify-between w-full lg:flex-row">
          <div className="mb-12 lg:mb-0 lg:max-w-lg lg:pr-4">
            <div className="max-w-xl mb-4">
              <h2 className="font-sans text-3xl sm:mt-0 mt-4 font-medium tracking-tight text-black sm:text-4xl sm:leading-none max-w-lg mb-4">
                Handleiding
              </h2>
              <div className="p-1 bg-house-green"></div>
              <br />
              <p className="text- text-base md:text-lg">
                Voorkom dat je vastloopt bij het opzetten van je online
                omgeving. Wij hebben een handleiding geschreven waarin stap voor
                stap wordt uitgelegd hoe je een online omgeving kan aanvragen.
                Hierbij wordt er ook uitgelegd hoe je gebruik kan maken van de
                verschillende technieken die wij aanbieden voor je online
                omgeving.
                <br />
                <br />
              </p>
            </div>
            <div className="flex justify-center">
              <Link to="handleiding">
                <button className="inline-block border border-transparent bg-primary-green px-6 py-2 text-center font-medium text-white hover:bg-house-green">
                  Bekijk de handleiding
                </button>
              </Link>
            </div>
          </div>
          <img
            className=""
            alt="logo"
            width="560"
            height="160"
            src="https://images.unsplash.com/photo-1585829365343-ea8ed0b1cb5b?q=80&w=2670&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"
          />
        </div>
        <br />
      </div>
      <div className="text-center">
        <h1 className="font-sans text-3xl pt-8 font-medium tracking-tight text-black mb-4">
          Pakketen die wij aanbieden
        </h1>
        <div className="flex justify-center">
          <div className="inline-block">
            <div className="w-72 h-2 bg-primary-green"></div>
          </div>
        </div>
      </div>
      <div className="py-8 flex items-center justify-center sm:max-w-xl md:max-w-full lg:max-w-screen-xl md:px-16 lg:px-0 mx-auto">
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-12">
          <div className="bg-white overflow-hidden shadow-lg transition-transform transform hover:scale-105">
            <div className="p-1 bg-primary-green"></div>
            <div className="p-6">
              <h2 className="text-3xl font-bold text-gray-800 mb-4">
                Gratis pakket
              </h2>
              <p className="text-gray-600 mb-4">
                Beste pakket voor studenten die iets kleins willen bouwen voor
                bijvoorbeeld huiswerk projecten
              </p>
              <p className="text-4xl font-bold text-gray-800 mb-4">Gratis</p>
              <ul className="text-sm text-gray-600 mb-4">
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
            <div className="p-4 flex justify-center">
              <Link to="requestForm">
                <button className="inline-block border border-transparent bg-primary-green px-6 py-2 text-center font-medium text-white hover:bg-house-green">
                  Kies dit pakket
                </button>
              </Link>
            </div>
          </div>

          <div className="bg-white overflow-hidden shadow-lg transition-transform transform hover:scale-105">
            <div className="p-1 bg-primary-green"></div>
            <div className="p-6">
              <h2 className="text-3xl font-bold text-gray-800 mb-4">
                Basis pakket
              </h2>
              <p className="text-gray-600 mb-4">
                Als je net iets meer nodig hebt voor je projecten
              </p>
              <br />
              <p className="text-4xl font-bold text-gray-800 mb-4">
                €4.99 <span className="text-base">/Maand</span>
              </p>
              <ul className="text-sm text-gray-600 mb-4">
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
                  5 GB opslag
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
            <div className="p-4 flex justify-center">
              <Link to="requestForm">
                <button className="inline-block border border-transparent bg-primary-green px-6 py-2 text-center font-medium text-white hover:bg-house-green">
                  kies dit pakket
                </button>
              </Link>
            </div>
          </div>

          <div className="bg-white overflow-hidden shadow-lg transition-transform transform hover:scale-105">
            <div className="p-1 bg-primary-green"></div>
            <div className="p-6">
              <h2 className="text-3xl font-bold text-gray-800 mb-4">
                Premium pakket
              </h2>
              <p className="text-gray-600 mb-4">
                Voor meerdere hobbyprojecten die ieder een eigen subdomein nodig
                heeft
              </p>
              <p className="text-4xl font-bold text-gray-800 mb-4">
                €9.99 <span className="text-base">/Maand</span>
              </p>

              <ul className="text-sm text-gray-600 mb-4">
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
                  5 subdomeinen
                </li>
              </ul>
            </div>
            <div className="p-4 flex justify-center">
              <Link to="requestForm">
                <button className="inline-block border border-transparent bg-primary-green px-6 py-2 text-center font-medium text-white hover:bg-house-green">
                  Kies dit pakket
                </button>
              </Link>
            </div>
          </div>
        </div>
      </div>
      <div className="text-center">
        <h1 className="font-sans text-3xl pt-8 font-medium tracking-tight text-black mb-4">
          Technologieën die wij ondersteunen
        </h1>
        <div className="flex justify-center">
          <div className="inline-block">
            <div className="w-72 h-2 bg-primary-green"></div>
          </div>
        </div>
      </div>

      <div className="container mx-auto py-3">
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-20 mb-6 py-6 mx-auto sm:max-w-xl md:max-w-full lg:max-w-screen-xl md:px-16 lg:px-0">
          <div className="flex flex-col items-center justify-center h-full text-center border bg-white shadow-xl transition-transform transform hover:scale-105">
            <div className="w-full p-1 bg-primary-green"></div>
            <div className="flex flex-col items-center justify-center h-full text-center p-4">
              <img
                className="w-32 h-32 object-contain mb-4"
                alt="Node.js logo"
                src="https://miro.medium.com/v2/resize:fit:800/1*v2vdfKqD4MtmTSgNP0o5cg.png"
              />
              <h1 className="text-lg font-bold">Node.js</h1>
              <p className="mt-2">
                Node.js is een open-source, cross-platform JavaScript runtime
                environment en library voor het uitvoeren van JavaScript-code
              </p>
            </div>
          </div>
          <div className="flex flex-col items-center justify-center h-full text-center border bg-white shadow-xl transition-transform transform hover:scale-105">
            <div className="w-full p-1 bg-primary-green"></div>
            <div className="flex flex-col items-center justify-center h-full text-center p-4">
              <img
                className="w-32 h-32 object-contain mb-4"
                alt=".NET logo"
                src="https://upload.wikimedia.org/wikipedia/commons/thumb/7/7d/Microsoft_.NET_logo.svg/1280px-Microsoft_.NET_logo.svg.png"
              />
              <h1 className="text-lg font-bold">.NET</h1>
              <p className="mt-2">
                .NET is een veelzijdig framework voor het bouwen van moderne,
                schaalbare applicaties voor web, mobiel en desktop.
              </p>
            </div>
          </div>
          <div className="flex flex-col items-center justify-center h-full text-center border bg-white shadow-xl transition-transform transform hover:scale-105">
            <div className="w-full p-1 bg-primary-green"></div>
            <div className="flex flex-col items-center justify-center h-full text-center p-4">
              <img
                className="w-32 h-32 object-contain mb-4"
                alt="SQL Server logo"
                src="https://www.svgrepo.com/show/303229/microsoft-sql-server-logo.svg"
              />
              <h1 className="text-lg font-bold">Microsoft SQL</h1>
              <p className="mt-2">
                Microsoft SQL Server is een krachtige relationele
                databasebeheersysteem voor het opslaan en beheren van gegevens.
              </p>
            </div>
          </div>
          <div className="flex flex-col items-center justify-center h-full text-center border bg-white shadow-xl transition-transform transform hover:scale-105">
            <div className="w-full p-1 bg-primary-green"></div>
            <div className="flex flex-col items-center justify-center h-full text-center p-4">
              <img
                className="w-32 h-32 object-contain mb-4"
                alt="MySQL logo"
                src="https://futuresolutionsonline.co.uk/wp-content/uploads/2023/04/mySQL-logo.png"
              />
              <h1 className="text-lg font-bold">MySQL</h1>
              <p className="mt-2">
                MySQL is een populaire open-source database die vaak wordt
                gebruikt voor webapplicaties vanwege zijn snelheid en
                flexibiliteit.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Homepage;
