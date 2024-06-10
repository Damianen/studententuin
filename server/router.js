import { Router } from "express";
import bodyParser from "body-parser";
import bcrypt from "bcryptjs";
import userService from "../server/services/user.service.js";
import subDomainService from "../server/services/subdomain.service.js";
import cors from "cors";
import session from "express-session";

const router = Router();

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

router.get("/api/getRepoFromSubdomain/:subdomain", async (req, res) => {
  try {
    console.log("Requesting repo for subdomain:", req.params.subdomain);
    const repo = await subDomainService.getRepo(req.params.subdomain);
    console.log("Repository:", repo);
    if (repo) {
      res.status(200).json(repo);
    } else {
      res.status(404).json({ message: "Repository not found" });
    }
  } catch (error) {
    console.error("Error:", error);
    res.status(500).json({ message: "An error occurred" });
  }
});

router.get("/api/getSubdomainFromUser", async (req, res) => {
  try {
    const subdomain = await userService.getSubdomainByUser(req);
    if (subdomain) {
      res.status(200).json(subdomain);
    } else {
      res.status(404).json({ message: "Subdomain not found" });
    }
  } catch (error) {
    console.error("Error:", error);
    res.status(500).json({ message: "An error occurred" });
  }
});

export default router;
