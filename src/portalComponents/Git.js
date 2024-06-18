import React, { useState, useEffect } from "react";
import axios from "axios";

const Git = () => {
    const [selectedBranch, setSelectedBranch] = useState("master");
    const [showExplanation, setShowExplanation] = useState(false);
    const [sshKey, setSshKey] = useState("");
    const [repositoryName, setRepositoryName] = useState("");
    const [selectedRepo, setSelectedRepo] = useState("");
    const [showRepoChange, setShowRepoChange] = useState(false);
    const [subdomain, setSubdomain] = useState("");
    const [branchFromDb, setBranchFromDb] = useState("master");
    const [customBranchName, setCustomBranchName] = useState("");
    const [user, setUser] = useState({});

    useEffect(() => {
        const fetchUser = async () => {
            try {
                const response = await axios.get(
                    "/api/getUserByEmailFromSession"
                );
                console.log("API response:", response.data); // Logging to check the response
                setUser(response.data);
            } catch (error) {
                console.error("Error fetching user:", error);
            }
        };

        fetchUser();
    }, []);

    useEffect(() => {
        const fetchRepo = async () => {
            if (user && user.subDomainName) {
                try {
                    const response = await axios.get(
                        `/api/getRepoFromSubdomain/${encodeURIComponent(
                            user.subDomainName
                        )}`
                    );
                    setSelectedRepo(response.data.github);
                } catch (error) {
                    console.error("Error fetching repo:", error);
                }
            }
        };

        const fetchBranch = async () => {
            if (user && user.subDomainName) {
                try {
                    const response = await axios.get(
                        `/api/getBranch/${encodeURIComponent(
                            user.subDomainName
                        )}`
                    );
                    if (response.data.githubBranch) {
                        setBranchFromDb(response.data.githubBranch);
                    } else {
                        setBranchFromDb("master");
                    }
                } catch (error) {
                    console.error("Error fetching branch:", error);
                }
            }
        };

        if (user) {
            fetchRepo();
            fetchBranch();
        }
    }, [user]);

    const showExplanationHandler = () => {
        setShowExplanation(!showExplanation);
    };

    const runDeploy = () => {
        // Run the deploy script
    };

    const showRepoChangeHandler = () => {
        setShowRepoChange(!showRepoChange);
    };

    const generateSshKey = () => {
        fetch(`/getSSHKey/${user.subDomainName}`, {
            method: "GET",
            headers: new Headers({
                "Content-Type": "application/json",
            }),
        })
            .then((response) => {
                if (!response.ok) {
                    throw new Error("Network response was not ok");
                }
                return response.json();
            })
            .then((data) => {
                setSshKey(data.key);
            })
            .catch((error) => {
                console.error("Error fetching SSH key:", error);
            });
    };
    const handleBranchChange = (e) => {
        e.preventDefault();
        const branchToSet =
            selectedBranch === "custom" ? customBranchName : selectedBranch;

        fetch(`/api/setBranch/`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                branch: branchToSet,
                subDomainName: user.subDomainName,
            }),
        })
            .then((response) => response.json())
            .then((data) => {
                console.log(data);
            });
        setBranchFromDb(branchToSet);
    };
    const handleNewRepo = (e) => {
        e.preventDefault();
        setSelectedRepo(e.target.value);
        fetch(`/newRepo`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                repo: repositoryName,
                subdomain: user.subDomainName,
                branch: branchFromDb,
            }),
        })
            .then((response) => response.json())
            .then((data) => {
                console.log(data);
            });
    };

    return (
        <div className="text-black text-base md:text-lg p-4">
            <h1>Git integratie</h1>
            <p> Verbind hier je eigen repository met je domein</p>
            <hr className="h-px my-5 bg-gray-200 border-0 dark:bg-gray-700"></hr>
            <h2>
                Stap 1 kopieer de SSH key en voeg hem toe aan je repository als
                je al een key hebt gegenereerd wordt de oude key overschreven
            </h2>
            <button
                className="inline-block rounded-md border border-transparent bg-primary-green px-8 py-2 text-center font-medium text-white hover:bg-green-400"
                onClick={generateSshKey}
            >
                Genereer SSH key
            </button>
            <p className="text-center">{sshKey} </p>
            <hr className="h-px my-5 bg-gray-200 border-0 dark:bg-gray-700"></hr>

            <div>
                {selectedRepo ? (
                    <div>
                        <h2>Stap 2 check de gekoppelde github Repo</h2>
                        <p>Repository: {selectedRepo}</p>
                    </div>
                ) : (
                    <div>
                        <h2>
                            Stap 2 vul de naam van je repository in met SSH url
                        </h2>
                        <form onSubmit={handleNewRepo}>
                            <input
                                className="p-4 px-8 py-2"
                                type="text"
                                placeholder="Enter repository name"
                                value={repositoryName}
                                onChange={(e) =>
                                    setRepositoryName(e.target.value)
                                }
                            />
                            <button
                                className="inline-block rounded-md border border-transparent bg-primary-green px-8 py-2 text-center font-medium text-white hover:bg-green-400"
                                type="submit"
                            >
                                Voeg repository
                            </button>
                        </form>
                    </div>
                )}
            </div>

            <hr className="h-px my-5 bg-gray-200 border-0 dark:bg-gray-700"></hr>
            <h2>Stap 3 kies de branch die je wilt deployen</h2>
            <p>Huidige branch: {branchFromDb}</p>
            <form onSubmit={handleBranchChange}>
                <select
                    className="p-4 px-8 py-2 mx-2"
                    value={selectedBranch}
                    onChange={(e) => setSelectedBranch(e.target.value)}
                >
                    <option value="release">Release</option>
                    <option value="develop">Develop</option>
                    <option value="main">Main</option>
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
                    Verander branch
                </button>
            </form>
            <hr className="h-px my-5 bg-gray-200 border-0 dark:bg-gray-700"></hr>
            <h2>
                Je kunt nu handmatig deployen door middel van de knop hieronder
            </h2>
            <p>
                Let op! Dit zal de huidige branch pullen en daarna al je post
                build commands uitvoeren (deze kun je in het tabje Omgevings
                variablen/post build commands veranderen)
            </p>
            <button
                className="inline-block rounded-md border border-transparent bg-primary-green px-8 py-2 text-center font-medium text-white hover:bg-green-400"
                type="button"
                onClick={runDeploy}
            >
                Deploy
            </button>
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
                        Kopieer de onderstaande webhook en voeg hem toe aan je
                        repository
                    </p>
                    <p className="text-center p-8">
                        https://webhook.studententuin.nl/
                    </p>
                    <hr className="h-px my-5 bg-gray-200 border-0 dark:bg-gray-700"></hr>
                    <p>
                        {" "}
                        Configureer het zo dat de webhook getriggerd wordt als
                        je een push doet. Dit resulteert in dat elke keer als er
                        een bestand wordt gepushed naar de repository
                        Studententuin op de hoogte wordt gesteld van de
                        verandering en automatisch pulled, er zal dan dus niet
                        fysiek gepulled hoeven worden de bestanden worden
                        automatisch gedeployed. Voor meer informatie over
                        webhooks kijk bij de documentatie van{" "}
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
