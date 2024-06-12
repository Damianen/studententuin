import React from "react";
import { Tabs, Tab } from "./components/Tabs.js";

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
                        <Tab label="Files">
                            {" "}
                            <div className="py-4">
                                <div className="flex min-h-full flex-col justify-start px-6 py-12 lg:px-8 space-y-7">
                                    <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                        <h2 className="text-4xl font-bold leading-9 tracking-tight text-gray-900">
                                            Domein beheren: Files
                                        </h2>
                                        <p className="mt-3">
                                            Deze pagina dient als handleiding
                                            voor het tabje files op jouw
                                            omgeving. Hierin wordt er uitgelegd
                                            wat precies de files tab is en waar
                                            je het kan vinden.
                                        </p>
                                    </div>
                                    <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                                    <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                        <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                                            Wat kan ik met het files tabje?
                                        </h2>
                                        <p className="mt-3">
                                            Natuurlijk wil je als gebruiker ook
                                            je bestanden uploaden, delete of
                                            zels folders uploaden. Met het files
                                            tabje kan je alle je bestanden
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
                                            Hoe werkt ik met het files tabje?
                                        </h2>
                                        <p className="mt-3">
                                            Om files te deleten klik je op de
                                            file die je wilt verwijderen,
                                            Vervolgens klik je op het knopje
                                            "Delete" hiermee verwijder je het
                                            geselecteerde bestand. Om een
                                            bestand te uploaden kan je kiezen
                                            tussen een losse file of een
                                            directory of te wel een folder.{" "}
                                            <br></br>
                                            <br></br>Om losse files te uploaden
                                            klik je op het knopje "upload single
                                            file", Vervolgens verschijnt een
                                            drag en drop box op je scherm hier
                                            mee kan je of je bestanden slepen
                                            naar het vakje of indien je dit niet
                                            wilt ook klikken op "browse files"
                                            om zelf te kiezen in windows
                                            explorer. <br></br>
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
                        <Tab label="Analytics">
                            <div className="py-4">
                                <div className="flex min-h-full flex-col justify-start px-6 py-12 lg:px-8 space-y-7">
                                    <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                        <h2 className="text-4xl font-bold leading-9 tracking-tight text-gray-900">
                                            Domein beheren: Analytics
                                        </h2>
                                        <p className="mt-3">
                                            Deze pagina dient als handleiding
                                            voor het tabje analytics op jouw
                                            omgeving. Hierin wordt er uitgelegd
                                            wat precies de analytics zijn en
                                            waar je ze kan vinden.
                                        </p>
                                    </div>
                                    <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                                    <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                        <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                                            Wat zijn analytics?
                                        </h2>
                                        <p className="mt-3">
                                            analytics zijn een handige tool om
                                            te zien hoeveel "resources" jouw
                                            site gebruikt. hiermee kan je zien
                                            of je te veel opslag gebruikt of
                                            hoeveel opslag je nog over hebt. Met
                                            deze informatie kan je misschien
                                            besluiten om je site te
                                            optimaliseren of om een upgrade te
                                            doen.
                                        </p>
                                    </div>
                                    <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                                    <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                        <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                                            Waar kan ik mijn analytics tab
                                            vinden?
                                        </h2>
                                        <p className="mt-3"></p>De analytics tab
                                        kan je vinden onder "jouw omgeving". Om
                                        naar jouw omgeving te gaan kan je
                                        klikken op de knop jouw omgeving rechts
                                        bovenin. Hierna kan je op de tab
                                        analytics klikken om de analytics te
                                        bekijken. hierbij zie je verschillende
                                        grafieken en tabellen die je een goed
                                        beeld geven van hoeveel resources je
                                        site gebruikt zoals je opslag gebruik.
                                    </div>
                                    <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                                </div>
                            </div>
                        </Tab>
                        <Tab label="Deployment">
                            <div className="py-4">
                                <div className="flex min-h-full flex-col items-center justify-center object-center text-center px-6 py-12 lg:px-8 space-y-7">
                                    <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                        <h1 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                                            Nog onder constructie, kom later
                                            terug voor meer informatie
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

                        <Tab label="Site Configuration">
                            {" "}
                            <div className="py-4">
                                <div className="flex min-h-full flex-col justify-start px-6 py-12 lg:px-8 space-y-7">
                                    <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                        <h2 className="text-4xl font-bold leading-9 tracking-tight text-gray-900">
                                            Domein beheren: Site Configuration
                                        </h2>
                                        <p className="mt-3">
                                            Deze pagina dient als handleiding
                                            voor het tabje Site Configuration op
                                            jouw omgeving. Hierin wordt precies
                                            uitgelegd wat je kan doen met de
                                            Site Configuration tab en waar je
                                            het kan vinden.
                                        </p>
                                    </div>
                                    <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                                    <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                        <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                                            Wat is Site Configuration?
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
                                            toevoegen die je krijgt van ons.
                                            Hierna kan je de repository
                                            toevoegen aan jouw omgeving.
                                            Vervolgens kies je de branch die je
                                            wilt toevoegen aan jouw omgeving.
                                            Ten slotte klik je op automatisch
                                            deployen.
                                        </p>
                                    </div>
                                </div>
                            </div>
                        </Tab>
                        <Tab label="Logs">
                            <div className="py-4">
                                <div className="flex min-h-full flex-col justify-start px-6 py-12 lg:px-8 space-y-7">
                                    <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                                        <h2 className="text-4xl font-bold leading-9 tracking-tight text-gray-900">
                                            Domein beheren: Logs
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
                <Tab label=".NET">
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
                <Tab label="MySQL">
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
