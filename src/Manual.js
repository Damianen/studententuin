import React from "react";
import { Tabs, Tab } from "./components/Tabs.js";
import Popup from "reactjs-popup";
import "reactjs-popup/dist/index.css";

export default function Manual() {
    return (
        <div>
            <Tabs>
                <Tab label="Domein aanvragen">
                    <div className="py-4">
                        <div className="flex min-h-full flex-col justify-start px-6 py-12 lg:px-8 space-y-7">
                            <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                <h2 className="text-4xl font-bold leading-9 tracking-tight text-gray-900">
                                    Hoe vraag ik een domein aan?
                                </h2>
                                <p className="mt-3">
                                    Deze pagina dient als handleiding voor het
                                    aanvragen van een subdomein. Hierin worden
                                    de stappen uigelegd die nodig zijn om
                                    succesvol een subdomein aan te vragen. Volg
                                    de instructies om ervoor te zorgen dat je
                                    aanvraag en setup soepel verloopt.
                                </p>
                            </div>
                            <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                            <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                                    Stap 1: Klik op het knopje account aanmaken
                                    rechts boven in of op de "aanvraag
                                    formulier" tekst hieronder
                                </h2>
                                <p className="mt-3">
                                    Klik hier om naar de{" "}
                                    <a className="text-primary-green" href="#">
                                        aanvraag formulier
                                    </a>{" "}
                                    te gaan.{" "}
                                </p>
                            </div>
                            <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                            <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                                    Stap 2: Vul de gegevens in de
                                    aanvraagformulier in.
                                </h2>
                                <p className="mt-3">
                                    <p className="font-bold">
                                        Je hebt wel een geldige schoolmail nodig
                                        om je in te schrijven!
                                    </p>
                                    <br></br>
                                    Vul de gegevens in het aanvraagformulier in
                                    en klik op verzenden. Je krijgt een email
                                    met een bevesteging dat je account is
                                    aangemaakt.
                                    <p className="font-bold"></p>
                                    Tip: gebruik een password manager om je
                                    wachtwoord te onthouden.
                                </p>
                            </div>
                            <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                            <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                                    Stap 3: Omgeving beheren
                                </h2>
                                <p className="mt-3">
                                    Na het aanvragen van je subdomein wil je
                                    natuurlijk je domein beheren om meteen aan
                                    de gang te gaan met je applicatie. om je
                                    online omgeving te beheren kan je makkelijk
                                    inloggen via het{" "}
                                    <a className="text-primary-green" href="#">
                                        inlog formulier
                                    </a>{" "}
                                    of via je de knop jouw omgeving rechts
                                    bovenin, Hier heb je alle tools om je online
                                    omgeving te beheren. Hier kan je meerdere
                                    settings vinden zoals node.js support,
                                    console bekijken, logs bekijken en nog veel
                                    meer.
                                    <br></br>
                                    <br></br>
                                    Weet je niet precies waar je terecht moet?
                                    klik dan op de tab domein beheren om meer
                                    informatie te krijgen over het beheren van
                                    je domein.
                                </p>
                            </div>
                        </div>
                    </div>
                </Tab>
                <Tab label="Domein beheren">
                    <Tabs>
                        <Tab label="Bestanden">
                            {" "}
                            <div className="py-4">
                                <div className="flex min-h-full flex-col justify-start px-6 py-12 lg:px-8 space-y-7">
                                    <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                        <h2 className="text-4xl font-bold leading-9 tracking-tight text-gray-900">
                                            Domein beheren: Bestanden
                                        </h2>
                                        <p className="mt-3">
                                            Deze pagina dient als handleiding
                                            voor het tabje Bestanden op jouw
                                            omgeving. Hierin wordt er uitgelegd
                                            wat precies de Bestanden tab is en
                                            waar je het kan vinden.
                                        </p>
                                    </div>
                                    <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                                    <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                        <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                                            Wat kan ik met het Bestanden tabje?
                                        </h2>
                                        <p className="mt-3">
                                            Natuurlijk wil je als gebruiker ook
                                            je bestanden uploaden, verwijderen
                                            of zels folders uploaden. Met het
                                            files tabje kan je alle je bestanden
                                            uploaden en deleten en ook je
                                            folders. Hiermee kan je makkelijk je
                                            code op jouw domein zetten en
                                            beheren. Dit is handig als je
                                            bijvoorbeeld een website wilt hosten
                                            of een applicatie wilt deployen maar
                                            je bijvoorbeeld geen github account
                                            hebt of je code niet wilt delen met
                                            anderen of als je gewoon snel je
                                            aanpassingen wilt testen.
                                        </p>
                                    </div>
                                    <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                                    <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                        <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                                            Hoe werkt ik met het bestanden
                                            tabje?
                                        </h2>
                                        <p className="mt-3">
                                            Om bestanden te verwijderen klik je
                                            op de file die je wilt verwijderen,
                                            Vervolgens klik je op het knopje
                                            "Delete" hiermee verwijder je het
                                            geselecteerde bestand. Om een
                                            bestand te uploaden kan je kiezen
                                            tussen een losse file of een
                                            directory of te wel een folder.{" "}
                                            <br></br>
                                            <br></br>Om losse bestanden te
                                            uploaden klik je op het knopje
                                            "upload single file", Vervolgens
                                            verschijnt een drag en drop box op
                                            je scherm hier mee kan je of je
                                            bestanden slepen naar het vakje of
                                            indien je dit niet wilt ook klikken
                                            op "browse files" om zelf te kiezen
                                            in windows explorer. <br></br>
                                            <br></br>Om een folder te uploaden
                                            moet je op het knopje "Upload a
                                            directory klikken" Vervolgens komt
                                            een knop te voorschijn. Na het
                                            klikken van de knop choose files kan
                                            je een folder selecteren die je wilt
                                            uploaden. Na het selecteren van de
                                            folder klik je op de knop "upload"
                                            om de folder te uploaden.
                                        </p>
                                    </div>
                                    <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                                    <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                        <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                                            Waar kan ik mijn files tab vinden?
                                        </h2>
                                        <p className="mt-3">
                                            De files tab kan je vinden onder
                                            "jouw omgeving". Om naar jouw
                                            omgeving te gaan kan je klikken op
                                            de knop jouw omgeving rechts
                                            bovenin. vervolgens wordt je
                                            omgeving geladen en is files ook
                                            meteen de eerste pagina waar je
                                            terecht komt. Hier kan je makkelijk
                                            je bestanden uploaden en beheren.
                                            Als je op een andere beheer pagina
                                            bent kun je op files klikken op het
                                            menu links van je scherm om naar de
                                            files pagina te gaan.
                                        </p>
                                    </div>
                                    <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                                </div>
                            </div>
                        </Tab>
                        <Tab label="Opslag">
                            <div className="py-4">
                                <div className="flex min-h-full flex-col justify-start px-6 py-12 lg:px-8 space-y-7">
                                    <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                        <h2 className="text-4xl font-bold leading-9 tracking-tight text-gray-900">
                                            Domein beheren: Opslag
                                        </h2>
                                        <p className="mt-3">
                                            Deze pagina dient als handleiding
                                            voor het tabje opslag op jouw
                                            omgeving. Hierin wordt er uitgelegd
                                            wat precies de opslag zijn en waar
                                            je ze kan vinden.
                                        </p>
                                    </div>
                                    <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                                    <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                        <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                                            Wat is opslag?
                                        </h2>
                                        <p className="mt-3">
                                            Opslag is hoeveel data jouw site
                                            gebruikt op de site. Voor gratis
                                            gebruikers is de opslag beperkt tot
                                            300MB. Dit is vooral zodat zoveel
                                            mogelijk mensen hun code kunnen
                                            hosten op onze server en ze veel
                                            kunnen leren. Voor betaalde
                                            gebruikers zijn er binnenkort opties
                                            om meer opslag aan te vragen als ze
                                            dat nodig hebben. Opslag is een
                                            belangrijk onderdeel van het hosten
                                            van je site. Als je site te veel
                                            opslag gebruikt kan het erook voor
                                            zorgen dat je site niet goed is
                                            geoptamiliseerd en het dus langzaam
                                            kan laden voor gebruikers wat je
                                            natuurlijk niet wilt. Opslag is dus
                                            een cruciaal onderdeel van het
                                            hosten van je site/applicatie.
                                        </p>
                                    </div>
                                    <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                                    <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                        <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                                            Waar kan ik mijn opslag tab vinden?
                                        </h2>
                                        <p className="mt-3"></p>De opslag tab
                                        kan je vinden onder "jouw omgeving". Om
                                        naar jouw omgeving te gaan kan je
                                        klikken op de knop jouw omgeving rechts
                                        bovenin. Hierna kan je op de tab opslag
                                        klikken om de opslag te bekijken.
                                        hierbij zie je verschillende grafieken
                                        en tabellen die je een goed beeld geven
                                        van hoeveel resources je site gebruikt
                                        zoals je opslag gebruik.
                                    </div>
                                    <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                                </div>
                            </div>
                        </Tab>
                        <Tab label="Deployment">
                            <div className="py-4">
                                <div className="flex min-h-full flex-col justify-start px-6 py-12 lg:px-8 space-y-7">
                                    <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                        <h2 className="text-4xl font-bold leading-9 tracking-tight text-gray-900">
                                            Domein beheren: Deployment
                                        </h2>
                                        <p className="mt-3">
                                            Deze pagina dient als handleiding
                                            voor het tabje Deployment op jouw
                                            omgeving. Hierin wordt er uitgelegd
                                            waarvoor de Deployment tabje voor
                                            dient en waar je het kan vinden.
                                        </p>
                                    </div>
                                    <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                                    <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                        <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                                            Wat is deployment
                                        </h2>
                                        <p className="mt-3">
                                            Deployment is het proces waarbij je
                                            code van je lokale machine naar een
                                            server wordt geupload en ook wordt
                                            opgezet. Dit is handig als je
                                            bijvoorbeeld een website of een
                                            applicatie wilt deployen. Met
                                            deployment kan je makkelijk inzien
                                            of je code correct is geupload en
                                            ook of er geen problemen zijn
                                            onstaan zoals bij het opzetten van
                                            de server. Deployment is een
                                            belangrijk onderdeel van het
                                            ontwikkelen van een website of
                                            applicatie. Hiermee kan je dus
                                            realtime volgen of er problemen zijn
                                            met je code of met het opzetten van
                                            je server. Het is dus een handig
                                            tool om zowel te debuggen en
                                            natuurlijk om te leren.
                                        </p>
                                    </div>
                                    <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                                    <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                        <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                                            Waar kan ik mijn Deployment tab
                                            vinden?
                                        </h2>
                                        <p className="mt-3"></p>De Deployment
                                        tab kan je vinden onder "jouw omgeving".
                                        Om naar jouw omgeving te gaan kan je
                                        klikken op de knop jouw omgeving rechts
                                        bovenin. Hierna kan je op de tab
                                        Deployment klikken om de Deployment
                                        status te bekijken. Hierbij zie je de
                                        volledige logboek van jouw server
                                        tijdens het opzetten van je website of
                                        applicatie. Hierbij kan je ook zien of
                                        er problemen zijn opgetreden tijdens het
                                        deployen van je code.
                                    </div>
                                    <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                                </div>
                            </div>
                        </Tab>

                        <Tab label="Omgevings variablen/post build commands">
                            {" "}
                            <div className="py-4">
                                <div className="flex min-h-full flex-col justify-start px-6 py-12 lg:px-8 space-y-7">
                                    <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                        <h2 className="text-4xl font-bold leading-9 tracking-tight text-gray-900">
                                            Domein beheren: Omgevings
                                            variablen/post build commands
                                        </h2>
                                        <p className="mt-3">
                                            Deze pagina dient als handleiding
                                            voor het tabje Omgevings
                                            variablen/post build commands op
                                            jouw omgeving. Hierin wordt precies
                                            uitgelegd wat je kan doen met de
                                            Omgevings variablen/post build
                                            commands tab en waar je het kan
                                            vinden.
                                        </p>
                                    </div>
                                    <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                                    <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                        <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                                            Wat zijn Omgevings variablen en post
                                            build commands?
                                        </h2>
                                        <p className="mt-3">
                                            Site Configuration is een tab waar
                                            mee je zowel enviroments variables
                                            kan toevoegen, verwijderen of aan
                                            kan passen. Enviroment variables
                                            zijn variables om bijvoorbeeld een
                                            database te koppelen aan je
                                            applicatie indien nodig. Met deze
                                            tab kan je makkeliijk dus je
                                            database toevoegen zodat je
                                            gevoelige informatie zoals je
                                            database host en database wachtwoord
                                            niet in je code hoeft te zetten.
                                            Hiermee kan je makkelijk je code
                                            deployen zonder dat je je zorgen
                                            hoeft te maken over je gevoelige
                                            informatie. Misschien ben je al
                                            bekent met enviroment variables in
                                            bijvoorbeeld een .env file. Dit is
                                            een soortgelijk concept. Ook checken
                                            we tijdens het toevoegen dat je niet
                                            perongeluk duplicate keys toevoegd.
                                            Naast enviroment variables kan je
                                            ook post build commands toevoegen
                                            voor bijvoorbeeld het bouwen van je
                                            site.
                                        </p>
                                    </div>
                                    <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                                    <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                        <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                                            Waar kan ik mijn Site Configuration
                                            tab vinden?
                                        </h2>
                                        <p className="mt-3">
                                            De Site Configuration tab kan je
                                            vinden onder "jouw omgeving". Om
                                            naar jouw omgeving te gaan kan je
                                            klikken op de knop jouw omgeving
                                            rechts bovenin. Hierna kan je op de
                                            tab Site Configuration klikken om
                                            Enviroment vairables te bekijken,
                                            bewerken of te verwijderen. Buiten
                                            je enviroment variables kan je ook
                                            post build commands Bekijken,
                                            bewerken of te verwijderen.
                                        </p>
                                    </div>
                                    <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                                </div>
                            </div>
                        </Tab>
                        <Tab label="Git">
                            <div className="py-4">
                                <div className="flex min-h-full flex-col justify-start px-6 py-12 lg:px-8 space-y-7">
                                    <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                        <h2 className="text-4xl font-bold leading-9 tracking-tight text-gray-900">
                                            Domein beheren: Git
                                        </h2>
                                        <p className="mt-3">
                                            Deze pagina dient als handleiding
                                            voor het tabje git op jouw omgeving.
                                            Hierin wordt er uitgelegd wat git is
                                            en wat precies je kan doen met de
                                            git tab en waar je het kan vinden.
                                        </p>
                                    </div>
                                    <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                                    <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                        <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                                            Wat is git?
                                        </h2>
                                        <p className="mt-3">
                                            Git is een versiebeheersysteem dat
                                            wordt gebruikt om wijzigingen in de
                                            code bij te houden. Met git kan je
                                            verschillende versies van je code
                                            bijhouden en makkelijk terug gaan
                                            naar een vorige versie. Git is
                                            handig om samen in een team te
                                            werken en om je code te backuppen.
                                            Git is een krachtige tool die je kan
                                            helpen om je code te beheren en te
                                            organiseren. Een site waar je gratis
                                            je code kan bewaren is github.com.
                                            Hier kan je makkelijk je code
                                            uploaden en beheren.
                                            <br></br>
                                            <p className="font-bold my-2">
                                                Tip: Als student kan je ook
                                                gratis in aanmerking komen voor
                                                github pro via de education pack
                                            </p>
                                            Github kan ook worden gebruikt om je
                                            code te deployen op onze servers.
                                            Indien je jouw github repository
                                            hebt toegevoegd aan jouw omgeving
                                            wordt automatisch je code van je git
                                            repository opgehaald en op je domein
                                            gezet. hierbij wordt de branch
                                            "main", "develop" of een custom
                                            branch gebruikt waarbij wij
                                            automatisch jouw code van de branch
                                            voor jouw online zetten op onze
                                            server.
                                        </p>
                                    </div>
                                    <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                                    <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                        <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                                            Waar kan ik mijn git tab vinden?
                                        </h2>
                                        <p className="mt-3">
                                            De git tab kan je vinden onder "jouw
                                            omgeving". Om naar jouw omgeving te
                                            gaan kan je klikken op de knop jouw
                                            omgeving rechts bovenin. Hierna kan
                                            je op de tab git klikken om jouw git
                                            repository te bekijken of toe te
                                            voegen.
                                        </p>
                                    </div>
                                    <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                                    <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                        <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                                            Hoe voeg ik mijn git repository toe?
                                        </h2>
                                        <p className="mt-3">
                                            Om je git repository toe te voegen
                                            kan je klikken op de knop "add
                                            repository". Hier moet je de SSH key
                                            toevoegen aan je repository op
                                            github. dit kan door op github naar
                                            je repository te gaan en naar
                                            settings te gaan. Hier kan je naar
                                            deploy keys gaan en een nieuwe key
                                            toevoegen. Hier kan je de key
                                            toevoegen die je krijgt van ons{" "}
                                            <Popup
                                                trigger={
                                                    <button>
                                                        <img
                                                            src="questionMark.png"
                                                            style={{
                                                                width: "10px",
                                                                height: "10px",
                                                            }}
                                                        />
                                                    </button>
                                                }
                                                modal
                                                nested
                                            >
                                                <div
                                                    style={{
                                                        width: "950px",
                                                        height: "400px",
                                                    }}
                                                >
                                                    <img
                                                        src="AddShh.jpg"
                                                        style={{
                                                            width: "100%",
                                                            height: "100%",
                                                        }}
                                                    />
                                                </div>
                                            </Popup>{" "}
                                            . Hierna kan je de repository
                                            toevoegen aan jouw omgeving.
                                            Vervolgens kies je de branch die je
                                            wilt toevoegen aan jouw omgeving.
                                            Ten slotte klik je op automatisch
                                            deployen. Dan kopieer je de link
                                            (https://webhook.studententuin.nl/)
                                            en die kun je toevoegen aan je
                                            github repository bij setting en dan
                                            webhooks{" "}
                                            <Popup
                                                trigger={
                                                    <button>
                                                        <img
                                                            src="questionMark.png"
                                                            style={{
                                                                width: "10px",
                                                                height: "10px",
                                                            }}
                                                        />
                                                    </button>
                                                }
                                                modal
                                                nested
                                            >
                                                <div
                                                    style={{
                                                        width: "950px",
                                                        height: "400px",
                                                    }}
                                                >
                                                    <img
                                                        src="Webhook-add.jpg"
                                                        style={{
                                                            width: "100%",
                                                            height: "100%",
                                                        }}
                                                    />
                                                </div>
                                            </Popup>{" "}
                                            . Vergeet niet de content type te
                                            veranderen naar application/json
                                            anders zal het niet werken{" "}
                                            <Popup
                                                trigger={
                                                    <button>
                                                        <img
                                                            src="questionMark.png"
                                                            style={{
                                                                width: "10px",
                                                                height: "10px",
                                                            }}
                                                        />
                                                    </button>
                                                }
                                                modal
                                                nested
                                            >
                                                <div
                                                    style={{
                                                        width: "950px",
                                                        height: "400px",
                                                    }}
                                                >
                                                    <img
                                                        src="Webhook-contentType.jpg"
                                                        style={{
                                                            width: "100%",
                                                            height: "100%",
                                                        }}
                                                    />
                                                </div>
                                            </Popup>{" "}
                                            !
                                        </p>
                                    </div>
                                </div>
                            </div>
                        </Tab>
                        <Tab label="Logboek">
                            <div className="py-4">
                                <div className="flex min-h-full flex-col justify-start px-6 py-12 lg:px-8 space-y-7">
                                    <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                        <h2 className="text-4xl font-bold leading-9 tracking-tight text-gray-900">
                                            Domein beheren: Logboek
                                        </h2>
                                        <p className="mt-3">
                                            Deze pagina dient als handleiding
                                            voor het beheren van een subdomein.
                                            Hierin wordt er uitgelegd wat
                                            precies de logs zijn en waar je ze
                                            kan vinden.
                                        </p>
                                    </div>
                                    <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                                    <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                        <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                                            Wat zijn logs?
                                        </h2>
                                        <p className="mt-3">
                                            Logs zijn bestanden waarin alle
                                            activiteiten van je domein worden
                                            bijgehouden. Hierin kan je zien
                                            welke acties er zijn uitgevoerd en
                                            of er fouten zijn opgetreden. Logs
                                            zijn handig om te debuggen en om te
                                            zien wat er precies is gebeurd in je
                                            domein. Je kan logs vinden onder de
                                            knop je "jouw omgeving" rechts
                                            bovenin. Hierna kan je op de jouw
                                            omgeving pagina op logs klikken om
                                            de logs te bekijken.
                                        </p>
                                    </div>
                                    <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                                    <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                        <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                                            Waar kan ik mijn logs vinden?
                                        </h2>
                                        <p className="mt-3">
                                            Je kan logs vinden onder de knop je
                                            "jouw omgeving" rechts bovenin.
                                            Hierna kan je op de jouw omgeving
                                            pagina op logs klikken om de logs te
                                            bekijken. Hier kan je alle logs zien
                                            die betrekking hebben op jouw
                                            applicatie. Ze staan op
                                            chronologische volgorde geordend.
                                            Zoals eerder benoemd kan je de logs
                                            weer gebruiken om je code te
                                            debuggen en errors op te sporen.
                                        </p>
                                    </div>
                                    <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                                </div>
                            </div>
                        </Tab>
                    </Tabs>
                </Tab>
                <Tab label="Node.js">
                    <div className="py-4">
                        <div className="flex min-h-full flex-col items-center justify-center object-center text-center px-6 py-12 lg:px-8 space-y-7">
                            <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                <h1 className="text-2xl font-bold leading-9 tracking-tight text-gray-900 mt-24">
                                    Nog onder constructie, kom later terug voor
                                    meer informatie
                                </h1>
                                <div class="flex justify-center items-center pt-12">
                                    <img
                                        className=" content-center"
                                        alt="under construction"
                                        src="https://t4.ftcdn.net/jpg/00/89/02/67/360_F_89026793_eyw5a7WCQE0y1RHsizu41uhj7YStgvAA.jpg"
                                    />
                                </div>
                            </div>
                        </div>
                    </div>
                </Tab>
                <Tab label="HTML">
                    <div className="py-4">
                        <div className="flex min-h-full flex-col justify-start px-6 py-12 lg:px-8 space-y-7">
                            <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                <h2 className="text-4xl font-bold leading-9 tracking-tight text-gray-900">
                                    HTML website
                                </h2>
                                <p className="mt-3">
                                    Deze pagina dient als handleiding voor het
                                    hosten van een HTML website. Hierin wordt er
                                    uitgelegd wat precies HTML hosting is en hoe
                                    je een HTML website kan hosten op onze
                                    servers.
                                </p>
                            </div>
                            <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                            <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                                    Hoe host ik een HTML website?
                                </h2>
                                <p className="mt-3">
                                    Een HTML website hosten op studententuin is
                                    zo makkelijk mogelijk gemaakt. Als studenten
                                    weten we helaas hoe lastig het soms kan zijn
                                    om makkelijk en snel een website te hosten
                                    is. om je website te hosten. Iedere student
                                    kent het wel Je hoeft bij onze niks af te
                                    weten van hoe en wat. Alles wat je moet doen
                                    is of je html files up te loaden naar de
                                    public folder of kan je ook zelf je eigen
                                    public folder up te loaden met behulp van de
                                    bestanden tabje in jouw omgeving. Hiermee
                                    kan je makkelijk je website hosten op onze
                                    servers en kan je makkelijk je website delen
                                    met anderen. Dit is handig als je
                                    bijvoorbeeld een assignment hebt voor
                                    bijvoorbeeld vakken waarbij je een website
                                    moet hosten. Naast dat studententuin sneller
                                    is dan andere gratis concurrent is het
                                    uploaden en testen van je website veel
                                    makkelijker.
                                </p>
                            </div>
                            <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                            <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                                    Waar kan ik mijn logs vinden?
                                </h2>
                                <p className="mt-3">
                                    Je kan logs vinden onder de knop je "jouw
                                    omgeving" rechts bovenin. Hierna kan je op
                                    de jouw omgeving pagina op logs klikken om
                                    de logs te bekijken. Hier kan je alle logs
                                    zien die betrekking hebben op jouw
                                    applicatie. Ze staan op chronologische
                                    volgorde geordend. Zoals eerder benoemd kan
                                    je de logs weer gebruiken om je code te
                                    debuggen en errors op te sporen.
                                </p>
                            </div>
                            <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                        </div>
                    </div>
                </Tab>
                <Tab label="Microsoft SQL">
                    <div className="py-4">
                        <div className="flex min-h-full flex-col items-center justify-center object-center text-center px-6 py-12 lg:px-8 space-y-7">
                            <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                <h1 className="text-2xl font-bold leading-9 tracking-tight text-gray-900 mt-24">
                                    Nog onder constructie, kom later terug voor
                                    meer informatie
                                </h1>
                                <div class="flex justify-center items-center pt-12">
                                    <img
                                        className=" content-center"
                                        alt="under construction"
                                        src="https://t4.ftcdn.net/jpg/00/89/02/67/360_F_89026793_eyw5a7WCQE0y1RHsizu41uhj7YStgvAA.jpg"
                                    />
                                </div>
                            </div>
                        </div>
                    </div>
                </Tab>
            </Tabs>
        </div>
    );
}
