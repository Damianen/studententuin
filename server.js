import express from "express";

import router from "./server/router.js";
import fs from "fs";
import path from "path";

const app = express();
const port = process.env.PORT || 3001;
const subdomain = ''
const directory = subdomain || 'studententuin'
const relativepath = '../' + directory

app.set("view engine", "ejs");
app.use(express.static("public"));
app.use(router);

function buildFileTree(dirPath) {
  const tree = {};
  const files = fs.readdirSync(dirPath);
  files.forEach((file) => {
    const filePath = path.join(dirPath, file);
    if (fs.statSync(filePath).isDirectory()) {
      tree[file] = buildFileTree(filePath);
    } else {
      tree[file] = file;
    }
  });
  return tree;
}

app.get('/filetree', (req, res) => {
  res.json(buildFileTree(relativepath));
});


app.listen(port, () => {
  console.log(`Server listening on port ${port}`);
});
