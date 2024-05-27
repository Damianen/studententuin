import React from "react";

function Homepage() {
  return (
    <div className="bg-gradient-to-b from-green-600">
      <div className="px-4 bg-white rounded-lg overflow-hidden shadow-lg mb-6 py-6 mx-auto sm:max-w-xl md:max-w-full lg:max-w-screen-xl md:px-16 lg:px-6">
        <div className="flex flex-col items-center justify-between w-full  lg:flex-row">
          <div className="mb-12 lg:mb-0 lg:max-w-lg lg:pr-4">
            <div className="max-w-xl mb-4">
              <h2 className="font-sans text-3xl sm:mt-0 mt-4 font-medium tracking-tight text-black sm:text-4xl sm:leading-none max-w-lg mb-4">
                Gratis online omgeving
              </h2>
              <div className="p-1 bg-green-500"></div>
              <br />
              <p className="text-black text-base md:text-lg">
                Heb jij een een betrouwbare online omgeving nodig? Wij bieden
                een gratis online omgeving aan voor alle studenten. Hierbij kan
                je gebruik maken van verschillende technieken zoals Node.js of
                wil je gebruik maken van een database en ben je een student?
                <br />
                <br />
              </p>
            </div>
            <div className="flex justify-center pt-8">
              <button className="inline-block rounded-md border border-transparent bg-green-500 px-6 py-2 text-center font-medium text-white hover:bg-green-400">
                Vraag een online omgeving aan
              </button>
            </div>
          </div>
          <img
            className="rounded-3xl"
            alt="logo"
            width="560"
            height="160"
            src="https://upload.wikimedia.org/wikipedia/commons/thumb/d/d9/Node.js_logo.svg/2560px-Node.js_logo.svg.png"
          />
        </div>
        <br />
      </div>

      <div className="px-4 bg-white rounded-lg overflow-hidden shadow-lg mb-6 py-6 mx-auto sm:max-w-xl md:max-w-full lg:max-w-screen-xl md:px-16 lg:px-6">
        <div className="flex flex-col items-center justify-between w-full lg:flex-row">
          <img
            className="rounded-3xl"
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
              <div className="p-1 bg-green-500"></div>
              <p className="text-black pt-6 text-base md:text-lg">
                Studententuin.nl is een platform voor studenten door studenten.
                Iedereen kent het wel, je hebt voor een school project een
                online omgeving nodig maar je wilt natuurlijk niet te veel
                betalen. Bij studententuin.nl kan je gratis een online omgeving
                aanvragen. Indien je meer resources nodig hebt kan je bij ons
                ook terecht.
                <br />
              </p>
            </div>
            <div className="flex justify-center pt-12">
              <button className="inline-block rounded-md border border-transparent bg-green-500 px-6 py-2 text-center font-medium text-white hover:bg-green-400">
                Bekijk beschikbare pakketen
              </button>
            </div>
          </div>
        </div>
        <br />
      </div>

      <div className="px-4 bg-white rounded-lg overflow-hidden shadow-lg mb-6 py-6 mx-auto sm:max-w-xl md:max-w-full lg:max-w-screen-xl md:px-16 lg:px-6">
        <div className="flex flex-col items-center justify-between w-full lg:flex-row">
          <div className="mb-12 lg:mb-0 lg:max-w-lg lg:pr-4">
            <div className="max-w-xl mb-4">
              <h2 className="font-sans text-3xl sm:mt-0 mt-4 font-medium tracking-tight text-black sm:text-4xl sm:leading-none max-w-lg mb-4">
                Handleiding
              </h2>
              <div className="p-1 bg-green-500"></div>
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
              <button className="inline-block rounded-md border border-transparent bg-green-500 px-6 py-2 text-center font-medium text-white hover:bg-green-400">
                Bekijk de handleiding
              </button>
            </div>
          </div>
          <img
            className="rounded-3xl"
            alt="logo"
            width="560"
            height="160"
            src="https://images.unsplash.com/photo-1585829365343-ea8ed0b1cb5b?q=80&w=2670&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"
          />
        </div>
        <br />
      </div>

      <div>
        <h1 className="font-sans text-3xl pt-8 font-medium tracking-tight text-black mb-4 text-center">
          Pakketen die wij aanbieden
        </h1>
      </div>
      <div className="py-8 flex items-center justify-center px-4 sm:max-w-xl md:max-w-full lg:max-w-screen-xl md:px-16 lg:px-4 mx-auto">
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-12">
          <div className="bg-white rounded-lg overflow-hidden shadow-lg transition-transform transform hover:scale-105">
            <div className="p-1 bg-green-500"></div>
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
              <button className="inline-block rounded-md border border-transparent bg-green-500 px-6 py-2 text-center font-medium text-white hover:bg-green-400">
                Kies dit pakket
              </button>
            </div>
          </div>

          <div className="bg-white rounded-lg overflow-hidden shadow-lg transition-transform transform hover:scale-105">
            <div className="p-1 bg-green-500"></div>
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
              <button className="inline-block rounded-md border border-transparent bg-green-500 px-6 py-2 text-center font-medium text-white hover:bg-green-400">
                Kies dit pakket
              </button>
            </div>
          </div>

          <div className="bg-white rounded-lg overflow-hidden shadow-lg transition-transform transform hover:scale-105">
            <div className="p-1 bg-green-500"></div>
            <div className="p-6">
              <h2 className="text-3xl font-bold text-gray-800 mb-4">
                Premium pakket
              </h2>
              <p className="text-gray-600 mb-4">
                Voor meerdere hobbyprojecten die ieder een eigen subdomein nodig
                heeft
              </p>
              <br />
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
              <button className="inline-block rounded-md border border-transparent bg-green-500 px-6 py-2 text-center font-medium text-white hover:bg-green-400">
                Kies dit pakket
              </button>
            </div>
          </div>
        </div>
      </div>
      <div>
        <h1 className="font-sans text-3xl pt-8 font-medium tracking-tight text-black mb-4 text-center mt-5">
          Technologieën die wij ondersteunen
        </h1>
      </div>
      <div class="container mx-auto py-8 px-4 sm:px-16 lg:px-4">
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-20 mb-6 py-6 rounded-3xl mx-auto sm:max-w-xl md:max-w-full lg:max-w-screen-xl md:px-16 lg:px-6">
          <div class="flex flex-col items-center justify-center h-full">
            <img
              class="rounded-3xl"
              alt="Node.js logo"
              width="140"
              height="40"
              src="https://miro.medium.com/v2/resize:fit:800/1*v2vdfKqD4MtmTSgNP0o5cg.png"
            />
            <h1 class="text-center mt-4">Node.js</h1>
            <p class="text-center ">Here's some text about Node.js</p>
          </div>
          <div class="flex flex-col items-center justify-center h-full">
            <img
              class="rounded-3xl"
              alt=".NET logo"
              width="140"
              height="40"
              src="https://upload.wikimedia.org/wikipedia/commons/thumb/7/7d/Microsoft_.NET_logo.svg/1280px-Microsoft_.NET_logo.svg.png"
            />
            <h1 class="text-center mt-4">.NET</h1>
            <p class="text-center">Here's some text about .NET</p>
          </div>
          <div class="flex flex-col items-center justify-center h-full">
            <img
              class="rounded-3xl"
              alt="SQL Server logo"
              width="140"
              height="40"
              src="https://www.derdack.com/wp-content/uploads/sites/2/2020/12/Microsoft_SQL_Server_Logo.png"
            />
            <h1 class="text-center mt-11">SQL Server</h1>
            <p class="text-center">Here's some text about SQL Server</p>
          </div>
          <div class="flex flex-col items-center justify-center h-full">
            <img
              class="rounded-3xl"
              alt="MySQL logo"
              width="140"
              height="40"
              src="https://futuresolutionsonline.co.uk/wp-content/uploads/2023/04/mySQL-logo.png"
            />
            <h1 class="text-center mt-4">MySQL</h1>
            <p class="text-center">Here's some text about MySQL</p>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Homepage;
