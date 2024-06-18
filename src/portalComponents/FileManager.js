import React, { useRef, useState } from "react";

export default function FileManager() {
  const fileInputRef = useRef(null);
  const [log, setLog] = useState("");

  function handleFileSelect(event) {
    const files = event.target.files;
    if (files.length) {
      const entries = [...files].map(file => ({
        isFile: true,
        isDirectory: false,
        file: (callback) => callback(file),
        name: file.webkitRelativePath ,
      }));

      makedir(entries)
        .then(output)
        .catch(handleSecurityLimitation);
    } else {
      notadir();
    }
  }

  async function makedir(entries) {
    const systems = entries.map((entry) => traverse(entry, {}));
    return Promise.all(systems);

    async function traverse(entry, fs) {
      if (entry.isDirectory) {
        fs[entry.name] = {};
        let dirReader = entry.createReader();
        await new Promise((res, rej) => {
          dirReader.readEntries(async (entries) => {
            for (let e of entries) {
              await traverse(e, fs[entry.name]);
            }
            res();
          }, rej);
        });
      } else if (entry.isFile) {
        await new Promise((res, rej) => {
          entry.file((file) => {
            fs[entry.name] = file;
            res();
          }, rej);
        });
      }
      return fs;
    }
  }

  function notadir() {
    setLog("Wasn't a directory, or webkitdirectory is not supported");
  }

  function output(system_trees) {
    function checkFile(key, value) {
      if (value instanceof File) {
        return `{[File] ${value.name}, ${value.size}b}`;
      } else return value;
    }

    // Upload the files here
    uploadFiles(system_trees);
  }

  function handleSecurityLimitation(error) {
    console.error(error);
    setLog('Faced security limitations. Please go to this fiddle: https://jsfiddle.net/x85vtnef/');
  }

  function uploadFiles(system_trees) {
    const formData = new FormData();
    const files = extractFiles(system_trees);
  
    files.forEach((file) => {
      formData.append('paths', file.webkitRelativePath);
      formData.append('files', file, file.webkitRelativePath );
    });
    const json = JSON.stringify(system_trees, null, 2);


    fetch('/upload', {
      method: 'POST',
      body: formData,
    })
      .then(response => response.json())
      .then(result => {
        console.log(result);
        setLog(`Upload successful`);
        window.location.reload();
      })
      .catch(error => {
        console.error('Upload failed', error);
        setLog(`Upload failed: ${error.message}`);
      });

  }

  function extractFiles(system_trees) {
    const files = [];
    (function traverse(fs) {
      Object.keys(fs).forEach(key => {
        if (fs[key] instanceof File) {
          files.push(fs[key]);
        } else if (typeof fs[key] === 'object') {
          traverse(fs[key]);
        }
      });
    })(system_trees);
    return files;
  }

  return (
    <div>
      <input
        type="file"
        ref={fileInputRef}
        webkitdirectory="true"
        directory="true"
        multiple
        onChange={handleFileSelect}
      />
      <pre>{log}</pre>
    </div>
  );
}
