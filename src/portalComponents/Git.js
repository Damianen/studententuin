import React, { useState, useEffect } from "react";
import axios from "axios";

const Git = () => {
  const [selectedBranch, setSelectedBranch] = useState("master");
  const [showExplanation, setShowExplanation] = useState(false);
  const [sshKey, setSshKey] = useState("");
  const [repositoryName, setRepositoryName] = useState("");
  const [showRepoChange, setShowRepoChange] = useState(false);
  const [subdomain, setSubdomain] = useState("");
  const [branchFromDb, setBranchFromDb] = useState("master");
  const [customBranchName, setCustomBranchName] = useState("");

  useEffect(() => {
    const fetchSubdomain = async () => {
      const response = await axios.get("/api/getUserByEmailFromSession");
      const subdomain = response.data.SubDomainName;
      setSubdomain(subdomain);
      console.log(subdomain);
      return subdomain;
    };

    fetchSubdomain().then((subdomain) => {
      fetch(`/api/getRepoFromSubdomain/${encodeURIComponent(subdomain)}`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      })
        .then((response) => response.json())
        .then((data) => {
          setRepositoryName(data.github);
        });

      fetch(`/api/getBranch/${encodeURIComponent(subdomain)}`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      })
        .then((response) => response.json())
        .then((data) => {
          setBranchFromDb(data.githubBranch);
        });
    });
  }, []);

  const showExplanationHandler = () => {
    setShowExplanation(!showExplanation);
  };

  const showRepoChangeHandler = () => {
    setShowRepoChange(!showRepoChange);
  };

  const generateSshKey = () => {
    // Get the generated ssh key from the script
  };

  const handleBranchChange = (e) => {
    e.preventDefault();
    const branchToSet = selectedBranch === "custom" ? customBranchName : selectedBranch;

    fetch(`/api/setBranch/`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ branch: branchToSet, subDomainName: subdomain }),
    })
      .then((response) => response.json())
      .then((data) => {
        console.log(data);
      });
      setBranchFromDb(branchToSet);
  };

  return (
    <div className="text-black text-base md:text-lg p-4">
      <h1>Git integratie</h1>
      <p> Verbind hier je eigen repository met je domein</p>
      <hr className="h-px my-5 bg-gray-200 border-0 dark:bg-gray-700"></hr>
      <h2>Stap 1 check de gekoppelde github Repo</h2>
      <div>
        {repositoryName !== "" ? (
          <p>Repository: {repositoryName}</p>
        ) : (
          <div>
            <h2>Stap 1 vul de naam van je repository in met SSH encriptie</h2>
            <form>
              <input
                className="p-4 px-8 py-2"
                type="text"
                placeholder="Enter repository name"
                value={customBranchName}
                onChange={(e) => setCustomBranchName(e.target.value)}
              />
              <button
                className="inline-block rounded-md border border-transparent bg-primary-green px-8 py-2 text-center font-medium text-white hover:bg-green-400"
                type="submit"
              >
                Add Repository
              </button>
            </form>
          </div>
        )}
      </div>

      <hr className="h-px my-5 bg-gray-200 border-0 dark:bg-gray-700"></hr>

      <h2>Stap 2 kopieer de SSH key en voeg hem toe aan je repository</h2>
      <p className="text-center">{sshKey}</p>

      <hr className="h-px my-5 bg-gray-200 border-0 dark:bg-gray-700"></hr>
      <h2>Stap 3 kies de branch die je wilt deployen</h2>
      <p>Current branch: {branchFromDb}</p>
      <form onSubmit={handleBranchChange}>
        <select
          className="p-4 px-8 py-2 mx-2"
          value={selectedBranch}
          onChange={(e) => setSelectedBranch(e.target.value)}
        >
          <option value="main">Main</option>
          <option value="develop">Develop</option>
          <option value="release">Release</option>
          <option value="custom">Custom</option>
        </select>
        {selectedBranch === "custom" && (
          <input
            className="p-4 px-8 py-2 mx-2"
            type="text"
            placeholder="Enter custom branch name"
            value={customBranchName}
            onChange={(e) => setCustomBranchName(e.target.value)}
          />
        )}
        <button
          className="inline-block rounded-md border border-transparent bg-primary-green px-8 py-2 text-center font-medium text-white hover:bg-green-400"
          type="submit"
        >
          Change branch
        </button>
      </form>
      <hr className="h-px my-5 bg-gray-200 border-0 dark:bg-gray-700"></hr>
      <h2>Wil je automatisch deployen klik dan hier</h2>
      <button
        className="inline-block rounded-md border border-transparent bg-primary-green px-8 py-2 text-center font-medium text-white hover:bg-green-400"
        type="button"
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
          <p className="text-center p-8">https://studententuin.nl/api/webhook</p>
          <hr className="h-px my-5 bg-gray-200 border-0 dark:bg-gray-700"></hr>
          <p>
            {" "}
            Configureer het zo dat de webhook getriggerd wordt als je een push
            doet. Dit resulteert in dat elke keer als er een bestand wordt
            gepushed naar de repository Studententuin op de hoogte wordt gesteld
            van de verandering en automatisch pulled, er zal dan dus niet fysiek
            gepulled hoeven worden de bestanden worden automatisch gedeployed.
            Voor meer informatie over webhooks kijk bij de documentatie van{" "}
            <a
              className="text-blue-600 hover:underline hover:text-blue-800 visited:text-purple-600"
              href="https://docs.github.com/en/webhooks"
            >
              GitHub
            </a>
          </p>
        </div>
      )}
    </div>
  );
};

export default Git;
