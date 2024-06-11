import { Router } from "express";
import bodyParser from "body-parser";
import bcrypt from "bcryptjs";
import userService from "../server/services/user.service.js";
import subDomainService from "../server/services/subdomain.service.js";
import cors from "cors";
import path from "path";
import { fileURLToPath } from "url";
import session from "express-session";
import fs from "fs/promises"; // gebruik fs/promises voor asynchrone file operaties

const router = Router();

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Middleware
router.use(bodyParser.urlencoded({ extended: false }));
router.use(bodyParser.json());
router.use(cors());

router.use(
  session({
    secret: "your_secret_key", // Verander dit naar een sterke geheime sleutel
    resave: false,
    saveUninitialized: true,
    cookie: { secure: false, maxAge: 30 * 60 * 1000 }, // Gebruik true als je HTTPS hebt
  })
);

// Routes
router.get("/", (req, res) => {
  res.render("index.ejs");
});

router.get("/requestForm", (req, res) => {
  res.render("index.ejs");
});

router.get("/login", (req, res) => {
  res.render("index.ejs");
});

router.get("/handleiding", (req, res) => {
  res.render("index.ejs");
});

router.get("/manage", (req, res) => {
  console.log("Accessing /manage route");
  console.log("Session on /manage:", req.session);

  if (req.session.user) {
    res.render("index.ejs", { user: req.session.user });
  } else {
    res.redirect("/");
  }
});

router.get("/api/check-session", (req, res) => {
  if (req.session.user) {
    res.json({ isAuthenticated: true });
  } else {
    res.json({ isAuthenticated: false });
  }
});

router.get("/api/getUserByEmailFromSession", async (req, res) => {
  try {
    console.log("Session data:", req.session); // Voeg logging toe voor sessiegegevens

    if (!req.session.user) {
      return res.status(404).json({ message: "User not authenticated" });
    }

    const user = await userService.getUserByEmailFromSession(req);
    if (user) {
      console.log("User found:", user);
      res.status(200).json(user);
    } else {
      console.log("User not found");
      res.status(404).json({ message: "User not found" });
    }
  } catch (error) {
    console.error("Error in getUserByEmailFromSession route:", error);
    res
      .status(500)
      .json({ message: "An error occurred", details: error.message });
  }
});

router.get("/about", (req, res) => {
  res.render("index.ejs");
});

router.get("/contact", (req, res) => {
  res.render("index.ejs");
});

// Helper function to find the newest log files
const getNewestLogFiles = async (dir) => {
  const files = await fs.readdir(dir);

  const stdoutFiles = files.filter(
    (file) =>
      file.startsWith("STUDENTENTUIN01-") &&
      file.includes("-stdout-") &&
      file.endsWith(".txt")
  );

  const stderrFiles = files.filter(
    (file) =>
      file.startsWith("STUDENTENTUIN01-") &&
      file.includes("-stderr-") &&
      file.endsWith(".txt")
  );

  stdoutFiles.sort((a, b) => {
    const timestampA = parseInt(a.split("-")[3].replace(".txt", ""));
    const timestampB = parseInt(b.split("-")[3].replace(".txt", ""));
    return timestampB - timestampA;
  });

  stderrFiles.sort((a, b) => {
    const timestampA = parseInt(a.split("-")[3].replace(".txt", ""));
    const timestampB = parseInt(b.split("-")[3].replace(".txt", ""));
    return timestampB - timestampA;
  });

  const newestStdout =
    stdoutFiles.length > 0 ? path.join(dir, stdoutFiles[0]) : null;
  const newestStderr =
    stderrFiles.length > 0 ? path.join(dir, stderrFiles[0]) : null;

  return { newestStdout, newestStderr };
};

router.get("/api/logs", async (req, res) => {
  try {
    const logsDir = path.resolve(__dirname, "../logs");
    console.log(`Looking for logs in: ${logsDir}`);

    try {
      await fs.access(logsDir);
    } catch (err) {
      console.error(`Directory does not exist: ${logsDir}`, err);
      return res.status(404).json({ error: "Logs directory does not exist" });
    }

    const { newestStdout, newestStderr } = await getNewestLogFiles(logsDir);

    if (!newestStdout && !newestStderr) {
      return res.status(404).json({ error: "No log files found" });
    }

    const stdoutData = newestStdout
      ? await fs.readFile(newestStdout, "utf8")
      : null;
    const stderrData = newestStderr
      ? await fs.readFile(newestStderr, "utf8")
      : null;

    res.json({ stdout: stdoutData, stderr: stderrData });
  } catch (err) {
    console.error("Error reading log files:", err);
    res
      .status(500)
      .json({ error: "Error reading log files", details: err.message });
  }
});

router.post("/register", async (req, res) => {
  const { email, password, userPackage } = req.body;

  try {
    // Controleer of de gebruiker al bestaat
    const userExists = await userService.getUserByEmail(email);

    if (userExists) {
      return res.status(400).json({ message: "User already exists" });
    }

    // Versleutel het wachtwoord
    const hashedPassword = await bcrypt.hash(password, 10);

    // Sla de gebruiker op in de database
    const user = await userService.createUser(email, hashedPassword, userPackage);

    // Geef een succesbericht terug met de redirect-url
    res.json({ redirectUrl: "/manage", sessionId: req.sessionID });
  } catch (error) {
    console.error("Error:", error);
    res.status(500).json({ message: "An error occurred" });
  }
});

router.post("/login", async (req, res) => {
  const { email, password } = req.body;

  try {
    const user = await userService.getUserByEmail(email);
    if (!user) {
      return res.status(401).json({ message: "Invalid email or password" });
    }

    const passwordMatch = await bcrypt.compare(password, user.password);
    if (!passwordMatch) {
      return res.status(401).json({ message: "Invalid email or password" });
    }

    req.session.user = user;
    console.log("Session created:", req.session);
    res.json({ redirectUrl: "/manage", sessionId: req.sessionID });
  } catch (error) {
    console.error("Error:", error);
    res.status(500).json({ message: "An error occurred" });
  }
});

router.post("/api/setBranch", async (req, res) => {
  const { branch, subDomainName } = req.body;
  try {
    await subDomainService.setBranch(subDomainName, branch);
    res.status(200).json({ message: "Branch set" });
  } catch (error) {
    console.error("Error:", error);
    res.status(500).json({ message: "An error occurred" });
  }
});

router.get("/api/getBranch/:subDomainName", async (req, res) => {
  const subDomainName = req.params.subDomainName;
  try {
    const branch = await subDomainService.getBranch(subDomainName);
    if (branch) {
      res.status(200).json(branch);
    } else {
      res.status(404).json({ message: "Branch not found" });
    }
  } catch (error) {
    console.error("Error:", error);
    res.status(500).json({ message: "An error occurred" });
  }
});

router.get("/logout", (req, res) => {
  req.session.destroy((err) => {
    if (err) {
      console.error("Error destroying session:", err);
      res.status(500).json({ message: "Internal server error" });
    } else {
      res.clearCookie("session"); // Verwijder het sessie-cookie
      res.redirect("/");
    }
  });
});

export default router;
