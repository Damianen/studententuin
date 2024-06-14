import React, { useState, useEffect } from "react";

const SiteConfiguration = () => {
    const [variables, setVariables] = useState([]);
    const [showVariableInput, setShowVariableInput] = useState(true);
    const [postBuildCommands, setPostBuildCommands] = useState([]);
    const [showPostBuildCommands, setShowPostBuildCommands] = useState(true);

    useEffect(() => {
        fetch("/api/appsettings")
            .then((response) => response.json())
            .then((data) => {
                if (Array.isArray(data)) {
                    const initialVariables = data.map((setting, index) => ({
                        id: index + 1,
                        key: setting.key,
                        value: setting.value,
                    }));
                    setVariables(initialVariables);
                } else {
                    setVariables([]);
                }
            })
            .catch((error) => {
                console.error("Error fetching app settings:", error);
                setVariables([]);
            });
    }, []);

    useEffect(() => {
        fetch("/api/postbuildcommands")
            .then((response) => response.json())
            .then((data) => {
                if (Array.isArray(data)) {
                    const initialCommands = data.map((command, index) => ({
                        id: index + 1,
                        content: command,
                    }));
                    setPostBuildCommands(initialCommands);
                } else {
                    setPostBuildCommands([]);
                }
            })
            .catch((error) => {
                console.error("Error fetching post build commands:", error);
                setPostBuildCommands([]);
            });
    }, []);

    const showPostBuildCommandsHandler = () => {
        setShowPostBuildCommands(!showPostBuildCommands);
    };

    const addPostBuildCommandHandler = (e) => {
        e.preventDefault();
        const command = document
            .getElementById("postBuildCommand")
            .value.trim();
        if (!command) return;

        const newCommand = {
            id: postBuildCommands.length + 1,
            content: command,
        };
        setPostBuildCommands([...postBuildCommands, newCommand]);

        fetch("/api/postbuildcommands", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ command }),
        })
            .then((response) => response.text())
            .then((data) => {
                console.log("Command added successfully:", data);
            })
            .catch((error) => {
                console.error("Error adding command:", error);
            });
    };

    const removeCommandHandler = (commandId) => {
        console.log("Removing command with ID:", commandId);
        fetch(`/api/postbuildcommands/${commandId}`, {
            method: "DELETE",
        })
            .then((response) => {
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                return response.text();
            })
            .then((data) => {
                console.log("Command removed successfully:", data);

                setPostBuildCommands(
                    postBuildCommands.filter(
                        (command) => command.id !== commandId
                    )
                );
            })
            .catch((error) => {
                console.error("Error removing command:", error);
            });
    };

    const showVariableInputHandler = () => {
        setShowVariableInput(!showVariableInput);
    };

    const addVariableHandler = (e) => {
        e.preventDefault();
        const variableName = document
            .getElementById("variableName")
            .value.trim();
        const variableValue = document
            .getElementById("variableValue")
            .value.trim();
        if (!variableName || !variableValue) return;

        if (
            variables.some(
                (variable) =>
                    variable.key.toLowerCase() === variableName.toLowerCase()
            )
        ) {
            alert("Variable name already exists. Please use a different name.");
            return;
        }

        const newVariable = {
            id: variables.length + 1,
            key: variableName,
            value: variableValue,
        };
        setVariables([...variables, newVariable]);
        showVariableInputHandler();

        fetch("/api/appsettings", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ key: variableName, value: variableValue }),
        })
            .then((response) => response.text())
            .then((data) => {
                console.log("Variable added successfully:", data);
            })
            .catch((error) => {
                console.error("Error adding variable:", error);
            });
    };

    const removeVariableHandler = (variableKey) => {
        console.log("Removing variable:", variableKey);
        fetch(`/api/appsettings/${variableKey}`, {
            method: "DELETE",
        })
            .then((response) => {
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                return response.text();
            })
            .then((data) => {
                console.log("Variable removed successfully:", data);

                setVariables(
                    variables.filter((variable) => variable.key !== variableKey)
                );
            })
            .catch((error) => {
                console.error("Error removing variable:", error);
            });
    };

    return (
        <div>
            <div>
                <hr className="h-px my-5 bg-gray-200 border-0 dark:bg-gray-700" />
                <h1 className="font-sans text-xl font-medium">
                    Voeg hier je environmental variables toe
                </h1>
                <ul className="shadow-md rounded-lg">
                    {variables.map((variable) => (
                        <div
                            key={variable.id}
                            className="flex justify-between items-center p-3 bg-white"
                        >
                            <span>
                                {variable.key}: {variable.value}
                            </span>
                            <button
                                className="text-red-500"
                                onClick={() =>
                                    removeVariableHandler(variable.key)
                                }
                            >
                                Verwijder
                            </button>
                        </div>
                    ))}
                </ul>
                <button
                    className="bg-house-green hover:bg-green-600 text-white font-bold mt-1 py-2 px-4 float-right rounded"
                    onClick={showVariableInputHandler}
                >
                    Variable toevoegen
                </button>
                {showVariableInput && (
                    <div>
                        <form onSubmit={addVariableHandler}>
                            <input
                                id="variableName"
                                className="px-8 py-2"
                                type="text"
                                placeholder="Enter variable name"
                                required
                            />
                            <input
                                id="variableValue"
                                className="py-2 px-8"
                                type="text"
                                placeholder="Enter variable value"
                                required
                            />
                            <button
                                type="submit"
                                className="inline-block rounded-md border border-transparent bg-primary-green px-8 py-2 text-center font-medium text-white hover:bg-green-400"
                            >
                                Voeg variable toe
                            </button>
                        </form>
                    </div>
                )}
            </div>
            <hr className="h-px my-5 bg-gray-200 border-0 dark:bg-gray-700" />
            <h1 className="font-sans text-xl font-medium">
                Post build commands
            </h1>
            <div>
                <ul className="shadow-md rounded-lg">
                    {postBuildCommands.map((command) => (
                        <div
                            key={command.id}
                            className="flex justify-between items-center p-3 bg-white"
                        >
                            <span>{command.content}</span>
                            <button
                                className="text-red-500"
                                onClick={() => removeCommandHandler(command.id)}
                            >
                                Verwijder
                            </button>
                        </div>
                    ))}
                </ul>
                <button
                    className="bg-house-green hover:bg-green-600 text-white font-bold py-2 px-4 float-right rounded"
                    onClick={showPostBuildCommandsHandler}
                >
                    Post Build Command toevoegen
                </button>
            </div>
            {showPostBuildCommands && (
                <div>
                    <form onSubmit={addPostBuildCommandHandler}>
                        <input
                            id="postBuildCommand"
                            className="p-4 px-8 py-2"
                            type="text"
                            placeholder="Enter command"
                            required
                        />
                        <button
                            type="submit"
                            className="inline-block rounded-md border border-transparent bg-primary-green px-8 py-2 text-center font-medium text-white hover:bg-green-400"
                        >
                            Voeg Post Build Commands
                        </button>
                    </form>
                </div>
            )}
        </div>
    );
};

export default SiteConfiguration;
