import React from "react";

export default function Manuel() {
    return (
        <div className="flex min-h-full flex-col justify-start px-6 py-12 lg:px-8 space-y-7">
            <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                <h2 className="text-4xl font-bold leading-9 tracking-tight text-gray-900">
                    Hoe zet ik mijn sub-domein op?
                </h2>
                <p className="mt-3">
                    Deze pagina dient als handleiding voor het aanvragen van een
                    subdomein. Hierin worden de stappen uigelegd die nodig zijn
                    om succesvol een subdomein aan te vragen en beheren. Volg de
                    instructies om ervoor te zorgen dat je aanvraag en setup
                    soepel verloopt.
                </p>
            </div>
            <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-black opacity-20" />
            <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                    Stap 1: Maak een account/vraag een domein aan
                </h2>
                <p className="mt-3">
                    Vraag via het{" "}
                    <a className="text-primary-green" href="#">
                        aanvraag formulier
                    </a>{" "}
                    je eerste subdomein aan.{" "}
                </p>
            </div>
            <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-black opacity-20" />
            <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                    Stap 2: Log in om de status van je aanvraag te bekijken
                </h2>
                <p className="mt-3">
                    Log via het{" "}
                    <a className="text-primary-green" href="#">
                        inlog formulier
                    </a>{" "}
                    in om de status van je sub-domein aanvraag te bekijken. Als
                    je aanvraag wordt goedgekeurt/afgekeurt zal je dit ook
                    ontvangen over email.
                </p>
            </div>
            <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-black opacity-20" />
            <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                    Stap 3: Omgeving beheren
                </h2>
                <p className="mt-3">
                    Log via het{" "}
                    <a className="text-primary-green" href="#">
                        inlog formulier
                    </a>{" "}
                    in om je online omgeving te beheren. Hier kan je meerdere
                    settings vinden zoals node.js support, console bekijken,
                    logs bekijken en nog veel meer.
                </p>
            </div>
        </div>
    );
}
