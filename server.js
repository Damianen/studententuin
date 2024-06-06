import express from "express";
import session from "express-session";
import router from './server/router.js';
import userRoutes from './server/routes/user.routes.js';
import fs from "fs";
import path from "path";
import multer from "multer";

const app = express();
const port = process.env.PORT || 3001;
const __dirname = path.resolve();
const subdomain = '';
const directory = subdomain || 'test';
const relativepath = '../' + directory + '/';
let clickedNode = '' || '/src/test';

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
  clickedNode = req.session.selectedNode
  console.log(clickedNode + " in clicked node");
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

app.delete('/delete-file', (req, res) => {
  fs.unlinkSync( __dirname+ clickedNode);
  console.log('File deleted:', clickedNode)
  res.json({ message: 'File deleted'
   });
});


function checkSelectedNode(req, res, next) {
  console.log(clickedNode + ' in checkSelectedNode');
}


const storage = multer.diskStorage({
 destination: (req, file, cb) => {
    const folderPath = path.join(relativepath, clickedNode,
    req.body.relativepath || '')
  fs.mkdirSync(folderPath, { recursive: true });
  cb(null, folderPath);
},
filename: (req, file, cb) => {
  cb(null, file.originalname);
},
});

const upload = multer({ storage: storage });

app.post('/upload', upload.array("files") ,(req, res) => {
  console.log('uploading files to path:', relativepath + clickedNode);
  res.send('File uploaded successfully');
});
app.use(userRoutes);

app.listen(port, () => {
  console.log(`Server listening on port ${port}`);
});
