import express from "express";
import session from "express-session";
import router from './server/router.js';
import userRoutes from './server/routes/user.routes.js';
import userService from "./server/services/user.service.js";
import fs, {readdir, stat } from "fs";
import path, { relative } from "path";
import multer from "multer";
import { promisify } from "util";
import fastFolderSize from "fast-folder-size";

const app = express();
const port = process.env.PORT || 3001;
const __dirname = path.resolve();

app.use(express.json());
app.set("view engine", "ejs");
app.use(express.static("public"));
app.use(router);
app.use(userRoutes);

const getRelativePath = (req) => {
  const userSubdomain = userService.getSubdomainByUser(req);
  let relativepath;
  if (userSubdomain) {
    let subdomain = userSubdomain.SubDomainName;
    console.log('Subdomain:', subdomain);
    directory = subdomain || 'test';
    relativepath = '../' + directory;
    console.log ('Relative path:', relativepath);
  }
  return relativepath;
}

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

const dirSize = async relativepath => {
  const fastFolderSizeASync= promisify(fastFolderSize);
  const size = await fastFolderSizeASync(relativepath);
  console.log('Directory size:', size);
  return size;
}


app.get('/filetree', async (req, res) => {
  try {
    let relativepath = getRelativePath(req);
    console.log('Relative path:', relativepath);
    if(relativepath){
    res.json(buildFileTree(relativepath));
    }else{
      res.status(400).json({ message: 'No subdomain found' });
    }
  } catch (error) {
    console.error("Error:", error);
    res.status(500).json({ message: 'Server error' });
  }
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

app.get('/dir-info', async (req, res) => {
  let relativepath = getRelativePath(req);
  console.log('Relative path:', relativepath);
  const size = await dirSize(relativepath);
  let userPackage = await userService.getUserPackage(req);
  let totalStorage;
  if (userPackage = 'free') {
    totalStorage = 314572800;
  }
  const usedStorage = size;
  const storagePercentage = (usedStorage / totalStorage) * 100;
  res.json({ size, storagePercentage });
});

app.get('/delete-file', (req, res) => {
  console.log('File deleted:', relativepath+clickedNode)
  fs.unlinkSync( relativepath + clickedNode);
  res.json({ message: 'File deleted'
   });
});



const storage = multer.diskStorage({
  destination: (req, file, cb) => {
    let folderPath = path.join(relativepath, clickedNode, req.body.relativepath || '');
    
    // Check if the path is a file or a directory
    if (fs.existsSync(folderPath)) {
      const stats = fs.statSync(folderPath);
      if (stats.isFile()) {
        // If it's a file, get the directory of the file
        folderPath = path.dirname(folderPath);
      }
    } else {
      // If the path doesn't exist, create it
      fs.mkdirSync(folderPath, { recursive: true });
    }

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


app.listen(port, () => {
  console.log(`Server listening on port ${port}`);
});
