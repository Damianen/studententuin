import React from "react";

export default function About() {
    return (
        <div className="bg-gradient-to-b from-primary-green to-neutral-warm min-h-screen">
            <div className="overflow-hidden shadow-lg mb-6 py-1 mx-auto sm:max-w-xl md:max-w-full lg:max-w-screen-xl md:px-16 lg:px-6 bg-white">
                <div className="flex flex-col items-center justify-between w-full lg:flex-row">
                    <div className="lg:w-full">
                        <img
                            className="w-full h-64 object-cover mt-12"
                            alt="logo"
                            src="https://www.mr-online.nl/wp-content/uploads/2020/12/Depositphotos_169316092_original-scaled.jpg"
                        />
                    </div>
                    <div className="px-4 py-6 lg:py-3 lg:px-8 lg:w-full">
                        <div className="max-w-xl mb-4">
                            <h2 className="font-sans text-3xl sm:mt-0 mt-4 font-medium tracking-tight text-black sm:text-4xl sm:leading-none max-w-lg mb-4">
                                Onze missie
                            </h2>
                            <div className="p-1 bg-house-green"></div>
                            <br />
                            <p className="text-black text-base md:text-lg">
                                Bij Studententuin.nl begrijpen we dat veel
                                studenten tegen het probleem aanlopen van het
                                vinden van betrouwbare en toegankelijke
                                testomgevingen. Vaak zijn de beschikbare
                                omgevingen niet stabiel of blijven ze niet lang
                                genoeg beschikbaar. Ons doel is om studenten te
                                voorzien van een werkende omgeving die ze kunnen
                                gebruiken voor hun projecten en om gewoon in
                                rond te kunnen spelen, zodat ze zich volledig
                                kunnen concentreren op hun leerproces en
                                ontwikkeling. Hierbij willen vooral de nadruk
                                leggen op een soort zandbak waarbij studenten al
                                hun eigen applicaties op kunnen hosten en
                                natuurlijk kunnen testen.
                            </p>
                        </div>
                    </div>
                </div>

                <div className="flex flex-col items-center justify-between w-full lg:flex-row mt-6">
                    <div className="px-4 py-6 lg:py-3 lg:px-8 lg:w-full">
                        <div className="max-w-xl mb-4">
                            <h2 className="font-sans text-3xl sm:mt-0 mt-4 font-medium tracking-tight text-black sm:text-4xl sm:leading-none max-w-lg mb-4">
                                Over ons
                            </h2>
                            <div className="p-1 bg-house-green"></div>
                            <br />
                            <p className="text-black text-base md:text-lg">
                                Wij zijn een groep eerstejaars IT-studenten die
                                niet zoeken naar problemen maar juist naar
                                oplossingen. We begrijpen de belangen van online
                                omgevingen voor studenten en het tekort aan
                                gratis omgevingen. Daarom hebben wij besloten
                                met behulp van school om een platform te creëren
                                waar studenten hun eigen omgevingen kunnen
                                beheren, zowel gratis als met een extra pakket
                                indien nodig. Ons verdienmodel is erop gericht
                                om de serverkosten te dekken. We hebben geen
                                winstoogmerk maar een educatief
                                ondersteuningsdoel voor zowel gebruikers als
                                voor onszelf.
                            </p>
                        </div>
                    </div>
                    <div className="lg:w-full mt-12">
                        <img
                            className="w-full h-64 object-cover"
                            alt="coding"
                            src="https://cdn.pixabay.com/photo/2015/09/05/20/02/coding-924920_1280.jpg"
                        />
                    </div>
                </div>
            </div>
        </div>
    );
}
