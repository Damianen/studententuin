import React, { useState } from "react";
import FolderInput from "./FileManager/FolderInput.js";
import TreeFileComponent from "./FileManager/TreeFileComponent.js";
import { convertListFileToObjectParentTree, readTemplate } from "./FileManager/utils.js";

function FileManager() {
  const [root, setRoot] = useState("");
  const [fileTree, setFileTree] = useState(null);
  const [mappedFileItems, setMappedFileItems] = useState({});
  const [currentSelectFile, setCurrentSelectFile] = useState(null);

  return (
    <div className="w-full h-full flex">
      <div className="w-86 overflow-auto">
        <FolderInput
          onChange={(event) => {
            const files = event.target.files || {};
            const tempedMappedFiles = {};
            const mappedFiles = Object.values(files).map((item) => {
              tempedMappedFiles[item.webkitRelativePath] = item;
              return {
                name: item.name,
                webkitRelativePath: item.webkitRelativePath,
              };
            });
            const template = convertListFileToObjectParentTree(mappedFiles);
            setRoot(Object.keys(template)[0]);
            const value = readTemplate(template, {});

            setFileTree(value);
            setMappedFileItems(tempedMappedFiles);
          }}
        />
        {fileTree && (
          <TreeFileComponent
            data={fileTree}
            root={root}
            onSelectItems={(items, treeId) => {
              const file = mappedFileItems[items[0]];
              if (file) {
                setCurrentSelectFile((window.URL ? URL : webkitURL).createObjectURL(file));
                console.log(file);
              } else {
                setCurrentSelectFile(null);
              }
            }}
          />
        )}
      </div>
      <div className="flex-1">
        {currentSelectFile && (
          <embed
            className="w-full h-full"
            src={currentSelectFile}
          />
        )}
      </div>
    </div>
  );
}

export default FileManager;