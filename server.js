import express from "express";
import session from "express-session";
import router from "./server/router.js";
import userRoutes from "./server/routes/user.routes.js";
import userService from "./server/services/user.service.js";
import getLatestLogs from "./server/routes/log.routes.js"
import fs, { readdir, stat, readFile } from "fs";
import path, { relative } from "path";
import multer from "multer";
import { promisify } from "util";
import fastFolderSize from "fast-folder-size";
import xml2js from "xml2js";
import bodyParser from "body-parser";
import { createServer } from "http";
import { Server } from "socket.io";

const app = express();
const httpServer = createServer(app);
const port = process.env.PORT || 3001;
const __dirname = path.resolve();

app.use(express.json());
app.set("view engine", "ejs");
app.use(express.static("public"));
app.use(router);
app.use(
    session({
        secret: "secret",
        resave: false,
        saveUninitialized: true,
        cookie: { secure: true },
    })
);

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
    const clickedNode = req.session.selectedNode;
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

app.delete("/api/appsettings/:key", async (req, res) => {
    console.log("Deleting key:", req.params.key);
    let configPath = await getRelativePath(req);

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
        return res
            .status(500)
            .json({ message: "Error reading appsettings file" });
    }
});

app.post("/api/appsettings", async (req, res) => {
    const { key, value } = req.body;
    const configPath = await getRelativePath(req);

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
    let relativepath = await getRelativePath(req);
    console.log("Relative path:", relativepath);
    const size = await dirSize(relativepath);
    let userPackage = await userService.getUserPackage(req);
    let totalStorage;
    if ((userPackage = "free")) {
        totalStorage = 314572800;
    }
    const usedStorage = size;
    const storagePercentage = (usedStorage / totalStorage) * 100;
    res.json({ size, storagePercentage, totalStorage });
});
const deleteFolderRecursive = function (directoryPath) {
    if (fs.existsSync(directoryPath)) {
        fs.readdirSync(directoryPath).forEach((file, index) => {
            const curPath = path.join(directoryPath, file);
            if (fs.lstatSync(curPath).isDirectory()) {
                // recurse
                deleteFolderRecursive(curPath);
            } else {
                // delete file
                fs.unlinkSync(curPath);
            }
        });
        fs.rmdirSync(directoryPath);
    }
};

app.get("/delete-file", async (req, res) => {
    let relativepath = await getRelativePath(req);
    let clickedNode = req.session.selectedNode;
    console.log("File/dir to delete:", clickedNode);
    let fullPath = relativepath + clickedNode;
    console.log("File/dir deleted:", fullPath);

    if (clickedNode === "") {
        return res
            .status(400)
            .json({ message: "No file or directory selected" });
    } else {
        if (fs.lstatSync(fullPath).isDirectory()) {
            deleteFolderRecursive(fullPath);
        } else {
            fs.unlinkSync(fullPath);
        }
    }
    console.log("File deleted:", fullPath);
    res.status(200).json({ message: "File deleted" });
});

const storage = multer.diskStorage({
    destination: async (req, file, cb) => {
        let filePath;
        if (Array.isArray(req.body.paths)) {
            req.body.paths.forEach((p) => {
                if (p.includes(file.originalname)) {
                    filePath = p;
                }
            });
        } else {
            filePath = req.body.paths;
        }
        let filePath2 = filePath.substring(0, filePath.lastIndexOf(`/`));
        let relativePath = await getRelativePath(req);
        let clickedNode = req.session.selectedNode;
        let stat = fs.statSync(path.join(relativePath, clickedNode));

        if (stat.isFile()) {
            // If clickedNode is a file, get the directory above it
            clickedNode = path.dirname(clickedNode);
        }
        const uploadPath = path.join(relativePath, clickedNode, filePath2);
        // Create the directory if it doesn't exist
        if (!fs.existsSync(uploadPath)) {
            fs.mkdirSync(uploadPath, { recursive: true });
        }

        cb(null, uploadPath);
    },
    filename: (req, file, cb) => {
        cb(null, path.basename(file.originalname || file.webkitRelativePath));
    },
});

const upload = multer({ storage: storage });

app.post("/upload", upload.array("files"), (req, res) => {
    res.status(200).json({ message: "Files uploaded successfully" });
});

const Uppystorage = multer.diskStorage({
    destination: async (req, file, cb) => {
        let relativePath = await getRelativePath(req);
        let clickedNode = req.session.selectedNode;

        let stat = fs.statSync(path.join(relativePath, clickedNode));

        if (stat.isFile()) {
            // If clickedNode is a file, get the directory above it
            clickedNode = path.dirname(clickedNode);
        }
        let folderPath = path.join(relativePath, clickedNode);
        cb(null, folderPath);
    },
    filename: (req, file, cb) => {
        cb(null, file.originalname);
    },
});
const uploadUppy = multer({ storage: Uppystorage });

app.post("/upload/uppy", uploadUppy.array("files"), (req, res) => {
    res.status(200).json({ message: "Files uploaded successfully" });
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

let relativePath;

(async () => {
    relativePath = await getRelativePath(req);
})();

const POST_BUILD_SCRIPT_PATH = `${relativePath}/postbuild.sh`;

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

httpServer.listen(port, () => {
  console.log(`Server listening on port ${port}`);
});

const io = new Server(httpServer, {});
io.on("connection", (socket) => {
    console.log('connected');
    socket.on("getLatestLogs", async () => {
        console.log("on get latest")
        const data = await getLatestLogs();
        socket.emit("logs", data);
    });
});
