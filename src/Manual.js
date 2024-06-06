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
                  Deze pagina dient als handleiding voor het aanvragen van een
                  subdomein. Hierin worden de stappen uigelegd die nodig zijn om
                  succesvol een subdomein aan te vragen. Volg de instructies om
                  ervoor te zorgen dat je aanvraag en setup soepel verloopt.
                </p>
              </div>
              <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
              <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                  Stap 1: Klik op het knopje account aanmaken rechts boven in of
                  op de "aanvraag formulier" tekst hieronder
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
                  Stap 2: Vul de gegevens in de aanvraagformulier in.
                </h2>
                <p className="mt-3">
                  <p className="font-bold">
                    Je hebt wel een geldige schoolmail nodig om je in te
                    schrijven!
                  </p>
                  <br></br>
                  Vul de gegevens in het aanvraagformulier in en klik op
                  verzenden. Je krijgt een email met een bevesteging dat je
                  account is aangemaakt.
                  <p className="font-bold"></p>
                  Tip: gebruik een password manager om je wachtwoord te
                  onthouden.
                </p>
              </div>
              <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
              <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                  Stap 3: Omgeving beheren
                </h2>
                <p className="mt-3">
                  Na het aanvragen van je subdomein wil je natuurlijk je domein
                  beheren om meteen aan de gang te gaan met je applicatie. om je
                  online omgeving te beheren kan je makkelijk inloggen via het{" "}
                  <a className="text-primary-green" href="#">
                    inlog formulier
                  </a>{" "}
                  of via je de knop jouw omgeving rechts bovenin, Hier heb je
                  alle tools om je online omgeving te beheren. Hier kan je
                  meerdere settings vinden zoals node.js support, console
                  bekijken, logs bekijken en nog veel meer.
                  <br></br>
                  <br></br>
                  Weet je niet precies waar je terecht moet? klik dan op de tab
                  domein beheren om meer informatie te krijgen over het beheren
                  van je domein.
                </p>
              </div>
            </div>
          </div>
        </Tab>
        <Tab label="Domein beheren">
          <Tabs>
            <Tab label="Logs">
              <div className="py-4">
                <div className="flex min-h-full flex-col justify-start px-6 py-12 lg:px-8 space-y-7">
                  <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                    <h2 className="text-4xl font-bold leading-9 tracking-tight text-gray-900">
                      Domein beheren: Logs
                    </h2>
                    <p className="mt-3">
                      Deze pagina dient als handleiding voor het beheren van een
                      subdomein. Hierin wordt er uitgelegd wat precies de logs
                      zijn en waar je ze kan vinden.
                    </p>
                  </div>
                  <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                  <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                    <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                      Wat zijn logs?
                    </h2>
                    <p className="mt-3">
                      Logs zijn bestanden waarin alle activiteiten van je domein
                      worden bijgehouden. Hierin kan je zien welke acties er
                      zijn uitgevoerd en of er fouten zijn opgetreden. Logs zijn
                      handig om te debuggen en om te zien wat er precies is
                      gebeurd in je domein. Je kan logs vinden onder de knop je
                      "jouw omgeving" rechts bovenin. Hierna kan je op de jouw
                      omgeving pagina op logs klikken om de logs te bekijken.
                    </p>
                  </div>
                  <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                  <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                    <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                      Waar kan ik mijn logs vinden?
                    </h2>
                    <p className="mt-3">
                      Je kan logs vinden onder de knop je "jouw omgeving" rechts
                      bovenin. Hierna kan je op de jouw omgeving pagina op logs
                      klikken om de logs te bekijken. Hier kan je alle logs zien
                      die betrekking hebben op jouw applicatie. Ze staan op
                      chronologische volgorde geordend. Zoals eerder benoemd kan
                      je de logs weer gebruiken om je code te debuggen en errors
                      op te sporen.
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
                      Deze pagina dient als handleiding voor het tabje analytics
                      op jouw omgeving. Hierin wordt er uitgelegd wat precies de
                      analytics zijn en waar je ze kan vinden.
                    </p>
                  </div>
                  <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                  <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                    <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                      Wat zijn analytics?
                    </h2>
                    <p className="mt-3">
                      analytics zijn een handige tool om te zien hoeveel
                      "resources" jouw site gebruikt. hiermee kan je zien of je
                      site te zwaar is en of je te veel opslag gebruikt. onder
                      het tabje analytics kan je zien hoeveel geheugen je al
                      hebt gebruikt en hoeveel je nog over hebt. Ook kan je zien
                      of je site te veel CPU gebruikt. Met deze informatie kan
                      je misschien besluiten om je site te optimaliseren of om
                      een upgrade te doen.
                    </p>
                  </div>
                  <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                  <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                    <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                      Waar kan ik mijn analytics tab vinden?
                    </h2>
                    <p className="mt-3"></p>De analytics tab kan je vinden onder
                    "jouw omgeving". Om naar jouw omgeving te gaan kan je
                    klikken op de knop jouw omgeving rechts bovenin. Hierna kan
                    je op de tab analytics klikken om de analytics te bekijken.
                    hierbij zie je verschillende grafieken en tabellen die je
                    een goed beeld geven van hoeveel resources je site gebruikt
                    zoals je opslag gebruik.
                  </div>
                  <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                </div>
              </div>
            </Tab>

            <Tab label="Domein">
              {" "}
              <div className="py-4">
                <div className="flex min-h-full flex-col justify-start px-6 py-12 lg:px-8 space-y-7">
                  <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                    <h2 className="text-4xl font-bold leading-9 tracking-tight text-gray-900">
                      Domein beheren: Domein
                    </h2>
                    <p className="mt-3">
                      Deze pagina dient als handleiding voor het tabje domein op
                      jouw omgeving. Hierin wordt er uitgelegd wat precies de
                      domein tab is en waar je het kan vinden.
                    </p>
                  </div>
                  <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                  <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                    <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                      Wat is een domein?
                    </h2>
                    <p className="mt-3">
                      Een domein is een unieke naam die je kan gebruiken om je
                      site te vinden. Een domein is een soort adres dat je kan
                      gebruiken om je site te bezoeken. Een domein is uniek en
                      kan maar een keer gebruikt worden. Een domein kan je kopen
                      bij verschillende bedrijven zoals google of hostnet. Een
                      domein kan je koppelen aan je site zodat mensen je site
                      kunnen vinden. Bij studententuin hebben wij al een domein
                      die je kan gebruiken waarbij je een "subdomein" bij ons
                      krijgt. met behulp van dit subdomein kunnen mensen jouw
                      site vinden.
                    </p>
                  </div>
                  <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                  <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                    <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                      Waar kan ik mijn domein tab vinden?
                    </h2>
                    <p className="mt-3">
                      De domein tab kan je vinden onder "jouw omgeving". Om naar
                      jouw omgeving te gaan kan je klikken op de knop jouw
                      omgeving rechts bovenin. Hierna kan je op de tab domein
                      klikken om jouw domein te bekijken. Hier kan je jouw
                      domein inzien en indien je een basis of premium pakket
                      hebt meerdere domeinen aanmaken.
                    </p>
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
                      Nog onder constructie, kom later terug voor meer
                      informatie
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

            <Tab label="Forms">
              <div className="py-4">
                <div className="flex min-h-full flex-col items-center justify-center object-center text-center px-6 py-12 lg:px-8 space-y-7">
                  <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                    <h1 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                      Nog onder constructie, kom later terug voor meer
                      informatie
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
            <Tab label="Git">
              <div className="py-4">
                <div className="flex min-h-full flex-col justify-start px-6 py-12 lg:px-8 space-y-7">
                  <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                    <h2 className="text-4xl font-bold leading-9 tracking-tight text-gray-900">
                      Domein beheren: Git
                    </h2>
                    <p className="mt-3">
                      Deze pagina dient als handleiding voor het tabje git op
                      jouw omgeving. Hierin wordt er uitgelegd wat git is en wat
                      precies je kan doen met de git tab en waar je het kan
                      vinden.
                    </p>
                  </div>
                  <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                  <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                    <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                      Wat is git?
                    </h2>
                    <p className="mt-3">
                      Git is een versiebeheersysteem dat wordt gebruikt om
                      wijzigingen in de code bij te houden. Met git kan je
                      verschillende versies van je code bijhouden en makkelijk
                      terug gaan naar een vorige versie. Git is handig om samen
                      in een team te werken en om je code te backuppen. Git is
                      een krachtige tool die je kan helpen om je code te beheren
                      en te organiseren. Een site waar je gratis je code kan
                      bewaren is github.com. Hier kan je makkelijk je code
                      uploaden en beheren.
                      <br></br>
                      <p className="font-bold my-2">
                        Tip: Als student kan je ook gratis in aanmerking komen
                        voor github pro via de education pack
                      </p>
                      Github kan ook worden gebruikt om je code te deployen op
                      onze servers. Indien je jouw github repository hebt
                      toegevoegd aan jouw omgeving wordt automatisch je code van
                      je git repository opgehaald en op je domein gezet. hierbij
                      wordt de branch "main", "develop" of een custom branch
                      gebruikt waarbij wij automatisch jouw code van de branch
                      voor jouw online zetten op onze server.
                    </p>
                  </div>
                  <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                  <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                    <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                      Waar kan ik mijn git tab vinden?
                    </h2>
                    <p className="mt-3">
                      De git tab kan je vinden onder "jouw omgeving". Om naar
                      jouw omgeving te gaan kan je klikken op de knop jouw
                      omgeving rechts bovenin. Hierna kan je op de tab git
                      klikken om jouw git repository te bekijken of toe te
                      voegen.
                    </p>
                  </div>
                  <hr className="h-0.5 my-8 w-9/12 mx-auto border-0 bg-primary-green opacity-20" />
                  <div className="sm:mx-auto sm:w-full sm:max-w-5xl">
                    <h2 className="text-2xl font-bold leading-9 tracking-tight text-gray-900">
                      Hoe voeg ik mijn git repository toe?
                    </h2>
                    <p className="mt-3">
                      Om je git repository toe te voegen kan je klikken op de
                      knop "add repository". Hier moet je de SSH key toevoegen
                      aan je repository op github. dit kan door op github naar
                      je repository te gaan en naar settings te gaan. Hier kan
                      je naar deploy keys gaan en een nieuwe key toevoegen. Hier
                      kan je de key toevoegen die je krijgt van ons. Hierna kan
                      je de repository toevoegen aan jouw omgeving. Vervolgens
                      kies je de branch die je wilt toevoegen aan jouw omgeving.
                      Ten slotte klik je op automatisch deployen.
                    </p>
                  </div>
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
                  Nog onder constructie, kom later terug voor meer informatie
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
                  Nog onder constructie, kom later terug voor meer informatie
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
                  Nog onder constructie, kom later terug voor meer informatie
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
                  Nog onder constructie, kom later terug voor meer informatie
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
