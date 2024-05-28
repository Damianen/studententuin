import React, { useState, useEffect } from "react";

const Git = () => {
  const [selectedBranch, setSelectedBranch] = useState("master");
  const [showExplanation, setShowExplanation] = useState(false);

  const showExplanationHandler = () => {
    setShowExplanation(!showExplanation);
  };

  return (
    <div className="text-black text-base md:text-lg p-4">
      <h1>Git intergratie</h1>
      <p> Verbind hier je eigen repository met je domein</p>
      <hr className="h-px my-5 bg-gray-200 border-0 dark:bg-gray-700"></hr>
      <h2>Stap 1 vul de naam van je repository in met SSH encriptie</h2>
      <form>
        <input
          className="p-4 px-8 py-2"
          type="text"
          placeholder="Enter repository name"
        />
        <button
          className="inline-block rounded-md border border-transparent bg-primary-green px-8 py-2 text-center font-medium text-white hover:bg-green-400"
          type="submit"
        >
          Add Repository
        </button>
      </form>

      <hr className="h-px my-5 bg-gray-200 border-0 dark:bg-gray-700"></hr>

      <h2>Stap 2 kopieer de SSH key en voeg hem toe aan je repository</h2>
      <p className="text-center">
        ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQD2Xxv9N2X4Jv
      </p>

      <hr className="h-px my-5 bg-gray-200 border-0 dark:bg-gray-700"></hr>
      <h2>Stap 3 kies de branch die je wilt deployen</h2>
      <form>
        <select
          className="p-4 px-8 py-2 mx-2"
          value={selectedBranch}
          onChange={(e) => setSelectedBranch(e.target.value)}
        >
          <option value="master">Master</option>
          <option value="develop">Develop</option>
          <option value="custom">Custom</option>
        </select>
        {selectedBranch === "custom" && (
          <input
            className="p-4 px-8 py-2 mx-2"
            type="text"
            placeholder="Enter custom branch name"
          />
        )}
        <button
          className="inline-block rounded-md border border-transparent bg-primary-green px-8 py-2 text-center font-medium text-white hover:bg-green-400"
          type="submit"
        >
          Deploy
        </button>
      </form>
      <hr className="h-px my-5 bg-gray-200 border-0 dark:bg-gray-700"></hr>
      <h2>Wil je automatisch deployen klik dan hier</h2>
      <button
        className="inline-block rounded-md border border-transparent bg-primary-green px-8 py-2 text-center font-medium text-white hover:bg-green-400"
        type="submit"
        onClick={showExplanationHandler}
      >
        Automatisch deployen
      </button>
      {showExplanation && (
        <div>
          <p className="">
            {" "}
            Kopieer de onderstaande webhook en voeg hem toe aan je repository
          </p>
          <p className="text-center p-8">
            https://studententuin.nl/api/webhook
          </p>
          <hr className="h-px my-5 bg-gray-200 border-0 dark:bg-gray-700"></hr>
          <p>
            {" "}
            Configureer het zo dat de webhook getriggerd wordt als je een push
            doet. Dit resulteert in dat elke keer als er een bestand wordt
            gepushed naar de repository Studententuin op de hoogte wordt gesteld
            van de verandering en automatisch pulled, er zal dan dus niet fysiek
            gepulled hoeven worden de bestanden worden automatisch gedeployed. Voor meer informatie over webhooks kijk bij de documentatie van <a className="text-blue-600 hover:underline hover:text-blue-800 visited:text-purple-600" href="https://docs.github.com/en/webhooks">GitHub</a>
          </p>
        </div>
      )}
    </div>
  );
};

export default Git;
