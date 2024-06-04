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

app.use(express.json());
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
  console.log(req.session.selectedNode + " in selected node");
  req.session.save(function (err) {
    if (err) {
      console.log(err);
      res.status(500).json({ message: 'Session save error' });
      return;
    }else{
      res.json({ message: 'Selected node updated' });
      console.log('session saved');
    }
});
});


function checkSelectedNode(req, res, next) {
  if (req.session && req.session.selectedNode) {
    next();
    console.log(selectedNode = req.session.selectedNode);
  } else {
    res.status(400).json({ message: 'Selected node is not defined' });
    console.log('selected node is not defined');
  }
}


const storage = multer.diskStorage({
  destination: function (req, file, cb) {
   if (req.session) {
    console.log(req.session.selectedNode + ' in storage');
      cb(null, req.session.selectedNode);
    } else {
      cb(new Error('Session is not defined'));
    }
  },
  filename: function (req, file, cb) {
    cb(null, file.originalname);
  }
});

app.post('/upload', checkSelectedNode, multer({ storage: storage }).single('file'), (req, res) => {
  res.json({ message: 'File uploaded' });
  console.log('file uploaded');
});

app.listen(port, () => {
  console.log(`Server listening on port ${port}`);
});
