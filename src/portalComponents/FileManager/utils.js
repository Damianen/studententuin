import fs from 'fs';
import path from 'path';

export function convertListFileToObjectParentTree(list) {
  const obj = {};
  
  list.forEach((item) => {
    const splitPath = item.webkitRelativePath.split('/');
		
    if (splitPath.length === 1) {
      obj[splitPath[0]] = splitPath[0];
    } else {
      let tempPointer = obj;
      let i = 0;
      while(i < splitPath.length - 1) {
      	if (tempPointer[splitPath[i]]) {
						tempPointer = tempPointer[splitPath[i]];
				} else {
          tempPointer[splitPath[i]] = {};
          tempPointer = tempPointer[splitPath[i]];
        }
        i++;
      }
      tempPointer[item.webkitRelativePath] = null;
    }
  });
  return obj;
};

function createDirectoriesAndFilesFromObject(obj, currentPath = '.') {
  for (const [key, value] of Object.entries(obj)) {
    const newPath = path.join(currentPath, key);
    if (value === null) {
      fs.writeFileSync(newPath, ''); // Create a new file
    } else {
      fs.mkdirSync(newPath, { recursive: true }); // Create a new directory
      createDirectoriesAndFilesFromObject(value, newPath); // Recursively create subdirectories and files
    }
  }
}

export const readTemplate = (template, data) => {
  for (const [key, value] of Object.entries(template)) {
    data[key] = {
      index: key,
      canMove: true,
      isFolder: value !== null,
      children: value !== null ? Object.keys(value || {}) : undefined,
      data: key,
      canRename: true
    };

    if (value !== null) {
      readTemplate(value, data);
    }
  }
  return data;
};