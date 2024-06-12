import React, { useState, useEffect } from "react";
import Uppy from "@uppy/core";
import { Dashboard } from "@uppy/react";
import XHRUpload from "@uppy/xhr-upload";
import axios from "axios";

import "@uppy/core/dist/style.min.css";
import "@uppy/dashboard/dist/style.min.css";

function Component() {
  const [usedStorage, setUsedStorage] = useState(0);
  const [totalStorage, setTotalStorage] = useState(0);
  const [user, setUser] = useState(null);
  const [isLoading, setLoading] = useState(true);

  useEffect(() => {
    const fetchUser = async () => {
      try {
        const response = await axios.get("/api/getUserByEmailFromSession");
        console.log("API response:", response.data); // Logging om de response te controleren
        setUser(response.data);
      } catch (error) {
        console.error("Error fetching user:", error);
      } finally {
        setLoading(false);
      }
    };
    const fetchUsedStorage = async () => {
      const response = await fetch("/dir-info", {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (response.ok) {
        const data = await response.json();
        console.log(data.size);
        setUsedStorage(data.size); // Update usedStorage with the fetched size
        setTotalStorage(data.totalStorage);
      }
    };
    fetchUser(); // Call the async function
    fetchUsedStorage(); // Call the async function
  }, []);
  const availableStorage = totalStorage - usedStorage;

  const uppy = new Uppy({
    autoProceed: false,
    restrictions: {},
    onBeforeFileAdded: (currentFile, files) => checkUpload(currentFile, files),
  });

  function checkUpload(currentFile, files) {
    // 10MB
    const maxTotalFileSize = availableStorage;
    let TotalFileSize = currentFile.data.size;

    for (var key in files) {
      TotalFileSize = TotalFileSize + files[key].size;
    }

    var grandTotalFileSize = currentFile.data.size + TotalFileSize;

    if (grandTotalFileSize >= maxTotalFileSize) {
      uppy.info("You have exceeded the maximum upload size");
      return false;
    }
    console.log("TotalFileSize: " + TotalFileSize);
    return true;
  }

  uppy.use(XHRUpload, {
    endpoint: user && user.subDomainName ? `https://${user.subDomainName}studententuin.nl/upload` : '',
    fieldName: "files",
    formData: true,
    allowedMetaFields: [{ id: "relativePath", name: "Relative Path" }],
  });

  uppy.on("file-added", (file) => {
    const relativePath = file.name;
    if (file.meta["relativePath"]) {
      relativePath = file.meta["relativePath"];
    }
    file.setFileMeta({ relativePath });
    console.log(file);
    console.log(file.meta);
  });

  if(isLoading){
    return <div>Loading...</div>;
  }
  return <Dashboard uppy={uppy} />;
}

export default Component;
