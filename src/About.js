import React from "react";

export default function About() {
  return (
    <div className="bg-gradient-to-b from-primary-green to-neutral-warm min-h-screen">
      <div className="px-4 overflow-hidden shadow-lg mb-6 py-6 mx-auto sm:max-w-xl md:max-w-full lg:max-w-screen-xl md:px-16 lg:px-6 bg-white">
        <div className="flex flex-col items-center justify-between w-full lg:flex-row">
          <div className="mb-12 lg:mb-0 lg:max-w-lg lg:pr-4">
            <div className="max-w-xl mb-4">
              <h2 className="font-sans text-3xl sm:mt-0 mt-4 font-medium tracking-tight text-black sm:text-4xl sm:leading-none max-w-lg mb-4">
                Onze missie
              </h2>
              <div className="p-1 bg-house-green"></div>
              <br />
              <p className="text-black text-base md:text-lg">
                Bij Studententuin.nl begrijpen we dat veel studenten tegen het
                probleem aanlopen van het vinden van betrouwbare en
                toegankelijke testomgevingen. Vaak zijn de beschikbare
                omgevingen niet stabiel of blijven ze niet lang genoeg
                beschikbaar. Ons doel is om studenten te voorzien van een
                werkende omgeving die ze kunnen gebruiken voor hun projecten en
                tests, zodat ze zich volledig kunnen concentreren op hun
                leerproces en ontwikkeling.
              </p>
            </div>
          </div>
          <img
            className=""
            alt="logo"
            width="560"
            height="160"
            src="https://www.mr-online.nl/wp-content/uploads/2020/12/Depositphotos_169316092_original-scaled.jpg"
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
                Over ons
              </h2>
              <div className="p-1 bg-house-green"></div>
              <p className="text-black pt-6 text-base md:text-lg">
                Wij zijn een groep 1ste jaars IT-studenten die niet zoeken naar
                problemen maar juist oplossingen. We begrijpen de belangen van
                online omgevingen voor studenten en ook tekort voor een gratis
                omgevingen. Daarom hebben wij besloten om een platform te
                creëren waar bij studenten hun eigen omgevingen kunnen beheren
                zowel gratis of indien nodig met een extra pakket.
              </p>
            </div>
          </div>
        </div>
        <br />
      </div>
    </div>
  );
}
