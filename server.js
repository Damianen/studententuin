import express from "express";
import session from "express-session";
import router from "./server/router.js";
import userRoutes from "./server/routes/user.routes.js";
import userService from "./server/services/user.service.js";
import fs, { readdir, stat, readFile } from "fs";
import path, { relative } from "path";
import multer from "multer";
import { promisify } from "util";
import fastFolderSize from "fast-folder-size";
import xml2js from "xml2js";
import bodyParser from "body-parser";

const app = express();
const port = process.env.PORT || 3001;
const __dirname = path.resolve();

app.use(express.json());
app.set("view engine", "ejs");
app.use(express.static("public"));
app.use(router);
app.use(userRoutes);
app.use(bodyParser.json());

const getRelativePath = async (req) => {
  const userSubdomain = await userService.getSubdomainByUser(req);
  let relativepath;
  if (userSubdomain) {
    let subdomain = userSubdomain.SubDomainName;
    console.log("Subdomain:", subdomain);
    let directory = subdomain || "test";
    relativepath = "../" + directory;
    console.log("Relative path:", relativepath);
  }
  return relativepath;
};

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

const dirSize = async (relativepath) => {
  const fastFolderSizeASync = promisify(fastFolderSize);
  const size = await fastFolderSizeASync(relativepath);
  console.log("Directory size:", size);
  return size;
};

app.get("/filetree", async (req, res) => {
  try {
    let relativepath = await getRelativePath(req);
    console.log("Relative path:", relativepath);
    if (relativepath) {
      res.json(buildFileTree(relativepath));
    } else {
      res.status(400).json({ message: "No subdomain found" });
    }
  } catch (error) {
    console.error("Error:", error);
    res.status(500).json({ message: "Server error" });
  }
});

app.post("/selected-node", (req, res) => {
  req.session.selectedNode = req.body.selectedNode;
  console.log(req.session.selectedNode + " in selected node");
  clickedNode = req.session.selectedNode;
  console.log(clickedNode + " in clicked node");
  req.session.save(function (err) {
    if (err) {
      console.log(err);
      res.status(500).json({ message: "Session save error" });
      return;
    } else {
      res.json({ message: "Selected node updated" });
      console.log("session saved");
    }
  });
});

app.get("/api/appsettings", (req, res) => {
  const configPath = path.join(__dirname, "web.config");

  fs.readFile(configPath, (err, data) => {
    if (err) {
      return res.status(500).send("Error reading config file");
    }

    xml2js.parseString(data, (err, result) => {
      if (err) {
        return res.status(500).send("Error parsing config file");
      }

      const appSettings =
        result.configuration?.appSettings?.[0]?.add?.map((setting) => ({
          key: setting.$.key,
          value: setting.$.value,
        })) || [];

      res.json(appSettings);
    });
  });
});

app.delete("/api/appsettings/:key", (req, res) => {
  console.log("Deleting key:", req.params.key);
  const configPath = path.join(__dirname, "web.config");

  try {
    const data = fs.readFileSync(configPath);

    xml2js.parseString(data, (err, result) => {
      if (err) {
        console.error("Error parsing XML:", err);
        return res.status(500).json({ message: "Error parsing XML" });
      }

      const appSettings = result.configuration.appSettings[0];

      const settings = appSettings?.add || [];

      const settingIndexToRemove = settings.findIndex((setting) => {
        return setting.$.key === req.params.key;
      });

      if (settingIndexToRemove !== -1) {
        settings.splice(settingIndexToRemove, 1);

        const builder = new xml2js.Builder();
        const updatedXml = builder.buildObject(result);

        try {
          fs.writeFileSync(configPath, updatedXml);
          res.json({ message: "Variable removed successfully" });
        } catch (err) {
          console.error("Error writing file:", err);
          return res
            .status(500)
            .json({ message: "Error writing appsettings file" });
        }
      } else {
        res.status(404).json({ message: "Variable not found" });
      }
    });
  } catch (err) {
    console.error("Error reading file:", err);
    return res.status(500).json({ message: "Error reading appsettings file" });
  }
});

app.post("/api/appsettings", (req, res) => {
  const { key, value } = req.body;
  const configPath = path.join(__dirname, "web.config");

  fs.readFile(configPath, (err, data) => {
    if (err) {
      return res.status(500).send("Error reading config file");
    }

    xml2js.parseString(data, (err, result) => {
      if (err) {
        return res.status(500).send("Error parsing config file");
      }

      const appSettings = result.configuration.appSettings[0];
      appSettings.add.push({ $: { key, value } });

      const builder = new xml2js.Builder();
      const updatedXml = builder.buildObject(result);

      fs.writeFile(configPath, updatedXml, (err) => {
        if (err) {
          return res.status(500).send("Error writing to config file");
        }

        res.status(200).send("Variable added successfully");
      });
    });
  });
});

app.get("/dir-info", async (req, res) => {
  let relativepath = getRelativePath(req);
  console.log("Relative path:", relativepath);
  const size = await dirSize(relativepath);
  let userPackage = await userService.getUserPackage(req);
  let totalStorage;
  if ((userPackage = "free")) {
    totalStorage = 314572800;
  }
  const usedStorage = size;
  const storagePercentage = (usedStorage / totalStorage) * 100;
  res.json({ size, storagePercentage });
});

app.get("/delete-file", (req, res) => {
  console.log("File deleted:", relativepath + clickedNode);
  fs.unlinkSync(relativepath + clickedNode);
  res.json({ message: "File deleted" });
});

const storage = multer.diskStorage({
  destination: (req, file, cb) => {
    let folderPath = path.join(
      relativepath,
      clickedNode,
      req.body.relativepath || ""
    );

    if (fs.existsSync(folderPath)) {
      const stats = fs.statSync(folderPath);
      if (stats.isFile()) {
        folderPath = path.dirname(folderPath);
      }
    } else {
      fs.mkdirSync(folderPath, { recursive: true });
    }

    cb(null, folderPath);
  },
  filename: (req, file, cb) => {
    cb(null, file.originalname);
  },
});

const upload = multer({ storage: storage });

app.post("/upload", upload.array("files"), (req, res) => {
  console.log("uploading files to path:", relativepath + clickedNode);
  res.send("File uploaded successfully");
});

const readPostBuildCommandsFromFile = () => {
  try {
    const data = fs.readFileSync("postbuild.sh", "utf8");
    // Split the data by newline characters to get each line as a command
    const commands = data.split("\n").filter((line) => line.trim() !== "");
    return commands;
  } catch (err) {
    console.error("Error reading post build commands file:", err);
    return [];
  }
};

const POST_BUILD_SCRIPT_PATH = "postbuild.sh";

app.get("/api/postbuildcommands", (req, res) => {
  const commands = readPostBuildCommandsFromFile();
  res.json(commands);
});

const writePostBuildCommandsToFile = (commands) => {
  try {
    fs.writeFileSync(POST_BUILD_SCRIPT_PATH, commands.join("\n"), "utf8");
    console.log("Post build commands written to file successfully");
  } catch (err) {
    console.error("Error writing post build commands to file:", err);
  }
};

app.post("/api/postbuildcommands", (req, res) => {
  const { command } = req.body;
  if (!command) {
    return res.status(400).json({ error: "Command is required" });
  }

  const commands = readPostBuildCommandsFromFile();
  commands.push(command);
  writePostBuildCommandsToFile(commands);

  res.status(201).send("Command added successfully");
});

app.delete("/api/postbuildcommands/:index", (req, res) => {
  const index = parseInt(req.params.index);
  if (isNaN(index) || index < 1) {
    return res.status(400).json({ error: "Invalid command index" });
  }

  const commands = readPostBuildCommandsFromFile();
  if (index - 1 < 0 || index - 1 >= commands.length) {
    return res.status(404).json({ error: "Command index out of bounds" });
  }

  commands.splice(index - 1, 1);
  writePostBuildCommandsToFile(commands);

  res.send("Command removed successfully");
});

app.listen(port, () => {
  console.log(`Server listening on port ${port}`);
});
