import { Router } from "express";
import bodyParser from "body-parser";
import bcrypt from "bcrypt";
import userService from "../server/services/user.service.js";
import cors from "cors";

const router = Router();

router.use(bodyParser.urlencoded({ extended: false }));
router.use(bodyParser.json());
router.use(cors());

router.get("/", (req, res, next) => {
  res.render("index.ejs");
});

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
    req.session = user;

    // Geef een succesbericht terug
    res.json({ message: "Login successful" });
  } catch (error) {
    console.error("Error:", error);
    res.status(500).json({ message: "An error occurred" });
  }
});

export default router;
