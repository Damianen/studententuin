import React, { useEffect, useState } from "react";
import { Treebeard } from "react-treebeard";
import Uppy from "../Uppy.js";
import { Button } from "@material-tailwind/react";
import fs from "fs";

const style = {
  tree: {
    base: {
      listStyle: "none",
      backgroundColor: "white",
      margin: 0,
      padding: 0,
      color: "#9DA5AB",
      fontFamily: "lucida grande ,tahoma,verdana,arial,sans-serif",
      fontSize: "22px",
    },
    node: {
      base: {
        position: "relative",
      },
      link: {
        cursor: "pointer",
        position: "relative",
        padding: "5px 5px",
        display: "block",
      },
      activeLink: {
        background: "#f2f0eb",
      },
      toggle: {
        base: {
          position: "relative",
          display: "inline-block",
          verticalAlign: "top",
          marginLeft: "-5px",
          height: "24px",
          width: "24px",
        },
        wrapper: {
          position: "absolute",
          top: "50%",
          left: "50%",
          margin: "-7px 0 0 -7px",
          height: "14px",
        },
        height: 14,
        width: 14,
        arrow: {
          fill: "0f0f0f",
          strokeWidth: 0,
        },
      },
      header: {
        base: {
          display: "inline-block",
          verticalAlign: "top",
          color: "#0f0f0f",
        },
        connector: {
          width: "2px",
          height: "12px",
          borderLeft: "solid 2px black",
          borderBottom: "solid 2px black",
          position: "absolute",
          top: "0px",
          left: "-21px",
        },
        title: {
          lineHeight: "24px",
          verticalAlign: "middle",
        },
      },
      subtree: {
        listStyle: "none",
        paddingLeft: "19px",
      },
      loading: {
        color: "#E2C089",
      },
    },
  },
};

function FileTree() {
  const [fileTree, setFileTree] = useState({});
  const [cursor, setCursor] = useState(false);
  const [selectedNode, setSelectedNode] = useState("");
  const [refresh, setRefresh] = useState(false);

  useEffect(() => {
    fetch("/filetree")
      .then((response) => response.json())
      .then((data) => setFileTree(processData(data)));
  }, [refresh]);

  const handleDeleteFile = () => {
    fetch("/delete-file", {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ selectedNode: selectedNode }),
    })
    .then(response => response.json())
    .then(data => {
      if (data.success) {
        console.log('File deleted successfully');
        
      } else {
        console.log('Failed to delete file');
        setRefresh(!refresh);
      }
    });
  }

  const handleFileChange = (e) => {
    const file = e.target.files[0];
    const formData = new FormData();
    formData.append("file", e.target.files[0]);
    fetch("/upload", {
      method: "POST",
      body: formData,
    })
      .then((response) => response.json())
      .then((data) => console.log(data))
      .catch((error) => console.error(error));
  };

 const processData = (data) => {
  const processedData = { name: "root", toggled: true, children: [], path: "" };
  function processNode(node, processedNode) {
    Object.keys(node).forEach((key) => {
      const path = `${processedNode.path}/${key}`;
      if (typeof node[key] === "string") {
        processedNode.children.push({ name: node[key], path });
      } else {
        const newNode = { name: key, children: [], path };
        processedNode.children.push(newNode);
        processNode(node[key], newNode);
      }
    });
  }
  processNode(data, processedData);
  return processedData;
};

  const onToggle = (node, toggled) => {
    if (cursor) {
      cursor.active = false;
    }
    node.active = true;
    if (node.children) {
      node.toggled = toggled;
    }
    setCursor(node);
    setFileTree(Object.assign({}, fileTree));
    setSelectedNode(node.path);
    console.log(node.path);
    fetch("/selected-node", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ selectedNode: node.path }),
    })
    .then((response) => response.json())
    .then((data) => console.log(data))
    .catch((error) => console.error(error));
  };
  return (
    <div>
      <div>
        <Treebeard
          className="text-black"
          data={fileTree}
          onToggle={onToggle}
          style={style}
        />
      </div>
      <div>
        <Uppy />
      </div>
      <div>Selected Node: {selectedNode}</div>{" "}
      {/* display the selected node name */}
      <div>
      <button className="inline-block  border border-transparent bg-primary-green px-6 py-2 text-center font-medium text-white hover:bg-house-green" onClick={handleDeleteFile}>
                  Delete
                </button>
  </div>
    </div>
  
  );
}

export default FileTree;
