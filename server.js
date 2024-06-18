import express from "express";
import session from "express-session";
import router from "./server/router.js";
import userRoutes from "./server/routes/user.routes.js";
import userService from "./server/services/user.service.js";
import subDomainService from "./server/services/subdomain.service.js";
import getLatestLogs from "./server/routes/log.routes.js";
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

app.get("/api/appsettings", async (req, res) => {
    let relativepath = await getRelativePath(req);
    const configPath = path.join(relativepath, "web.config");

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
    let relativepath = await getRelativePath(req);
    const configPath = path.join(relativepath, "web.config");

    fs.readFile(configPath, (err, data) => {
        if (err) {
            console.error("Error reading config file:", err);
            return res
                .status(500)
                .json({ message: "Error reading config file" });
        }

        xml2js.parseString(data, (err, result) => {
            if (err) {
                console.error("Error parsing config file:", err);
                return res
                    .status(500)
                    .json({ message: "Error parsing config file" });
            }

            const appSettings = result.configuration.appSettings[0];
            const settings = appSettings?.add || [];

            const settingIndexToRemove = settings.findIndex(
                (setting) => setting.$.key === req.params.key
            );

            if (settingIndexToRemove !== -1) {
                settings.splice(settingIndexToRemove, 1);

                const builder = new xml2js.Builder();
                const updatedXml = builder.buildObject(result);

                fs.writeFile(configPath, updatedXml, (err) => {
                    if (err) {
                        console.error("Error writing to config file:", err);
                        return res
                            .status(500)
                            .json({ message: "Error writing to config file" });
                    }

                    res.json({ message: "Variable removed successfully" });
                });
            } else {
                res.status(404).json({ message: "Variable not found" });
            }
        });
    });
});

app.post("/api/appsettings", async (req, res) => {
    const { key, value } = req.body;
    let relativepath = await getRelativePath(req);
    const configPath = path.join(relativepath, "web.config");

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

app.post("/newRepo", async (req, res) => {
    const { subdomain, repo, branch } = req.body;

    subDomainService.insertNewRepo(
        subdomain,
        repo,
        branch,
        (success, error) => {
            if (success) {
                if (success.status === 200) {
                    console.log("New repository added to db");
                }
            } else {
                res.status(error.status).json({ message: error.message });
            }
        }
    );

    try {
        const response = await fetch(
            "https://webhook.studententuin.nl/newRepo",
            {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    subdomain: subdomain,
                    repo: repo,
                    branch: branch,
                }),
            }
        );
        if (!response.ok) {
            throw new Error("Failed to add new repository");
        }
        res.status(200).json({ message: "Repository added" });
    } catch (error) {
        console.error("Error adding new repository:", error);
        res.status(500).json({ error: "Failed to add new repository" });
    }
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

const readPostBuildCommandsFromFile = (path) => {
    try {
        return fs.readFileSync(path, "utf8").split("\n").filter(Boolean);
    } catch (err) {
        console.error("Error reading post build commands from file:", err);
        return [];
    }
};

const POST_BUILD_SCRIPT_PATH = `postbuild.sh`;

app.get("/api/postbuildcommands", async (req, res) => {
    const relativepath = await getRelativePath(req);
    const path = `${relativepath}/${POST_BUILD_SCRIPT_PATH}`;
    const commands = readPostBuildCommandsFromFile(path);
    res.json(commands);
});

const writePostBuildCommandsToFile = (path, commands) => {
    try {
        fs.writeFileSync(path, commands.join("\n"), "utf8");
        console.log("Post build commands written to file successfully");
    } catch (err) {
        console.error("Error writing post build commands to file:", err);
    }
};

app.post("/api/postbuildcommands", async (req, res) => {
    const { command } = req.body;
    if (!command) {
        return res.status(400).json({ error: "Command is required" });
    }

    const relativepath = await getRelativePath(req);
    const path = `${relativepath}/${POST_BUILD_SCRIPT_PATH}`;
    const commands = readPostBuildCommandsFromFile(path);
    commands.push(command);
    writePostBuildCommandsToFile(path, commands);

    res.status(201).send("Command added successfully");
});

app.delete("/api/postbuildcommands/:index", async (req, res) => {
    const index = parseInt(req.params.index);
    if (isNaN(index) || index < 1) {
        return res.status(400).json({ error: "Invalid command index" });
    }

    const relativepath = await getRelativePath(req);
    const path = `${relativepath}/${POST_BUILD_SCRIPT_PATH}`;
    const commands = readPostBuildCommandsFromFile(path);
    if (index - 1 < 0 || index - 1 >= commands.length) {
        return res.status(404).json({ error: "Command index out of bounds" });
    }

    commands.splice(index - 1, 1);
    writePostBuildCommandsToFile(path, commands);

    res.send("Command removed successfully");
});

httpServer.listen(port, () => {
    console.log(`Server listening on port ${port}`);
});

const io = new Server(httpServer, {
    cors: {
        origin: "http://localhost:3001",
        methods: ["GET", "POST"],
    },
});
io.on("connection", (socket) => {
    console.log("connected");
    socket.on("getLatestLogs", async (data) => {
        const logs = await getLatestLogs(data.subdomainName);
        socket.emit("logs", logs);
    });
});
