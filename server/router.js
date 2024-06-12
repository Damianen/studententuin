import { Router } from "express";
import bodyParser from "body-parser";
import bcrypt from "bcryptjs";
import userService from "../server/services/user.service.js";
import cors from "cors";
import path from "path";
import { fileURLToPath } from "url";
import session from "express-session";
import fs from "fs/promises"; // gebruik fs/promises voor asynchrone file operaties

const router = Router();

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

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

router.get("/", (req, res, next) => {
  res.render("index.ejs");
});

router.get("/requestForm", (req, res, next) => {
  res.render("index.ejs");
});

router.get("/login", (req, res, next) => {
  res.render("index.ejs");
});

router.get("/handleiding", (req, res, next) => {
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
    // Als er een gebruiker is opgeslagen in de sessie, retourneer een JSON-object met isAuthenticated:true
    res.json({ isAuthenticated: true });
  } else {
    // Als er geen gebruiker is opgeslagen in de sessie, retourneer een JSON-object met isAuthenticated:false
    res.json({ isAuthenticated: false });
  }
});

router.get("/api/getUserByEmailFromSession", async (req, res) => {
  try {
    const user = await userService.getUserByEmailFromSession(req);
    if (user) {
      res.status(200).json(user);
    } else {
      res.status(404).json({ message: "User not found" });
    }
  } catch (error) {
    console.error("Error:", error);
    res.status(500).json({ message: "An error occurred" });
  }
});

router.get("/about", (req, res, next) => {
  res.render("index.ejs");
});

router.get("/contact", (req, res, next) => {
  res.render("index.ejs");
});

// router.get("/check-session", (req, res) => {
//   if (req.session.user) {
//     res.json({ session: req.session.user });
//   } else {
//     res.status(401).json({ message: "No active session" });
//   }
// });

router.post("/login", async (req, res) => {
  const { email, password } = req.body;

  try {
    // Zoek de gebruiker in de database op basis van het e-mailadres
    const user = await userService.getUserByEmail(email);

    // Als de gebruiker niet gevonden is, geef een foutmelding terug
    if (!user) {
      return res.status(401).json({ message: "Invalid email or password" });
    }

    // Vergelijk het ingediende wachtwoord met het versleutelde wachtwoord in de database
    const passwordMatch = await bcrypt.compare(password, user.password);

    // Als de wachtwoorden niet overeenkomen, geef een foutmelding terug
    if (!passwordMatch) {
      return res.status(401).json({ message: "Invalid email or password" });
    }

    // Als de inloggegevens correct zijn, sla de gebruiker op in de sessie
    req.session.user = user;

    console.log("Session created:", req.session);

    // Geef een succesbericht terug met de redirect-url
    res.json({ redirectUrl: "/manage", sessionId: req.sessionID });
  } catch (error) {
    console.error("Error:", error);
    res.status(500).json({ message: "An error occurred" });
  }
});

router.get("/logout", (req, res) => {
  // Vernietig de sessie van de gebruiker
  req.session.destroy((err) => {
    if (err) {
      console.error("Error destroying session:", err);
      res.status(500).json({ message: "Internal server error" });
    } else {
      // Stuur een succesvolle respons terug
      res.clearCookie("session"); // Verwijder het sessie-cookie
      res.redirect("/");
    }
  });
});

// Helper functie om het nieuwste logbestand te vinden per type (stdout/stderr)
const getNewestLogFiles = async (dir) => {
  const files = await fs.readdir(dir);

  // Filter logbestanden voor stdout en stderr
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

  // Sorteer de logbestanden op datum
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

  // Return de nieuwste logbestanden voor stdout en stderr
  const newestStdout =
    stdoutFiles.length > 0 ? path.join(dir, stdoutFiles[0]) : null;
  const newestStderr =
    stderrFiles.length > 0 ? path.join(dir, stderrFiles[0]) : null;

  return { newestStdout, newestStderr };
};

router.get("/api/logs", async (req, res) => {
  try {
    const logsDir = path.join(__dirname, "../logs");
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
    res
      .status(500)
      .json({ error: "Error reading log files", details: err.message });
  }
});

export default router;
