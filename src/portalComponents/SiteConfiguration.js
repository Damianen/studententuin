import React, { useState, useEffect } from "react";

const SiteConfiguration = () => {
  const [variables, setVariables] = useState([]);
  const [showVariableInput, setShowVariableInput] = useState(false);
  const [postBuildCommands, setPostBuildCommands] = useState([]);
  const [showPostBuildCommands, setShowPostBuildCommands] = useState(false);

  const showPostBuildCommandsHandler = () => {
    setShowPostBuildCommands(!showPostBuildCommands);
  };

  var postBuildCommandsList = [
    { id: 1, content: "npm run build" },
    { id: 2, content: "npm run build2" },
    { id: 3, content: "npm run test" },
  ];
  const addPostBuildCommandHandler = () => {
    const command = document.getElementById("postBuildCommand").value;
    const newCommand = { id: postBuildCommands.length + 1, content: command };
    if (command === "") {
      return;
    } else {
      setPostBuildCommands([...postBuildCommands, newCommand]);
      showPostBuildCommandsHandler();
    }
  };
  useEffect(() => {
    setPostBuildCommands([...postBuildCommandsList]);
  }, []);

  var variableList = [
    { id: 1, content: "DB-HOST" },
    { id: 2, content: "DB-USER" },
    { id: 3, content: "DB-PASSWORD" },
  ];

  const showVariableInputHandler = () => {
    setShowVariableInput(!showVariableInput);
  };
  const addVariableHandler = () => {
    const variableName = document.getElementById("variableName").value;
    const variableValue = document.getElementById("variableValue").value;
    const newVariable = { id: variables.length + 1, content: variableName };
    if (variableName === "" || variableValue === "") {
      return;
    } else {
      setVariables([...variables, newVariable]);
      showVariableInputHandler();
    }
  };
  useEffect(() => {
    setVariables([...variableList]);
  }, []);

  return (
    <div>
      <div>
        <hr className="h-px my-5 bg-gray-200 border-0 dark:bg-gray-700"></hr>
        <h1 className="font-sans text-xl font-medium">
          Voeg hier je environmental variables toe
        </h1>
        <ul className="shadow-md rounded-lg">
          {variables.map((log, index) => (
            <div
              key={index}
              className={`p-3 ${
                index % 2 === 0 ? "bg-white" : "bg-neutral-warm"
              }`}
            >
              {log.content}
            </div>
          ))}
        </ul>
          <button
            className="bg-house-green hover:bg-green-600 text-white font-bold py-2 px-4 float-right rounded"
            onClick={showVariableInputHandler}
          >
            Toggle adding a variable
          </button>
        {showVariableInput && (
          <div>
            <form>
              <input
                id="variableName"
                className="p-4 px-8 py-2"
                type="text"
                placeholder="Enter variable name"
                required
              />
              <input
                id="variableValue"
                className="p-4 px-8 py-2"
                type="text"
                placeholder="Enter variable value"
                required
              />
              <button
                type="submit"
                className="inline-block rounded-md border border-transparent bg-primary-green px-8 py-2 text-center font-medium text-white hover:bg-green-400"
                onClick={() => {
                  addVariableHandler();
                }}
              >
                Add Variable
              </button>
            </form>
          </div>
        )}
      </div>
      <hr className="h-px my-5 bg-gray-200 border-0 dark:bg-gray-700"></hr>
      <h1 className="font-sans text-xl font-medium">Post build commands</h1>
      <div>
        <ul className="shadow-md rounded-lg">
          {postBuildCommands.map((log, index) => (
            <div
              key={index}
              className={`p-3 ${
                index % 2 === 0 ? "bg-white" : "bg-neutral-warm"
              }`}
            >
              {log.content}
            </div>
          ))}
        </ul>
        <button
          className="bg-house-green hover:bg-green-600 text-white font-bold py-2 px-4 float-right rounded"
          onClick={showPostBuildCommandsHandler}
        >
          Toggle Post Build Commands
        </button>
      </div>
      {showPostBuildCommands && (
        <div>
          <form>
            <input
              id="postBuildCommand"
              className="p-4 px-8 py-2"
              type="text"
              placeholder="Enter post build command"
              required
            />
            <button
              type="submit"
              className="inline-block rounded-md border border-transparent bg-primary-green px-8 py-2 text-center font-medium text-white hover:bg-green-400"
            >
              Add Post Build Command
            </button>
          </form>
        </div>
      )}
    </div>
  );
};

export default SiteConfiguration;
