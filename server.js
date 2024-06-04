import express from "express";
import session from "express-session";
import router from "./server/router.js";
import fs from "fs";
import path from "path";
import multer from "multer";

const app = express();
const port = process.env.PORT || 3001;
const subdomain = ''
const directory = subdomain || 'studententuin'
const relativepath = '../' + directory

app.set("view engine", "ejs");
app.use(express.static("public"));
app.use(router);
app.use(session({
  secret: 'secret',
  resave: false,
  saveUninitialized: true,
  cookie: { secure: true }
}));


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

app.post('/selected-node', (req, res) => {
  req.session.selectedNode = req.body.selectedNode;
  res.json({ message: 'Selected node updated' });
  console.log(req.session.selectedNode);
});

const storage = multer.diskStorage({
  destination: function (req, file, cb) {
   if (req.session) {
      cb(null, req.session.selectedNode);
    } else {
      cb(new Error('Session is not defined'));
    }
  },
  filename: function (req, file, cb) {
    cb(null, file.originalname);
  }
});

app.post('/upload', multer({ storage: storage }).single('file'), (req, res) => {
  res.json({ message: 'File uploaded' });
});

app.listen(port, () => {
  console.log(`Server listening on port ${port}`);
});
