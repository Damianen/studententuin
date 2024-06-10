import React from "react";
import Uppy from "@uppy/core";
import { Dashboard } from "@uppy/react";
import XHRUpload from "@uppy/xhr-upload";

import "@uppy/core/dist/style.min.css";
import "@uppy/dashboard/dist/style.min.css";

function Component() {
  const uppy = new Uppy({
    autoProceed: true,
    restrictions: {},
  });

  uppy.use(XHRUpload, {
    endpoint: "https://studententuin.nl/upload",
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

  return <Dashboard uppy={uppy} />;
}

export default Component;
