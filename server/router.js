import { Router } from "express";
import bodyParser from "body-parser";
import cors from "cors";

const router = Router();

router.use(bodyParser.urlencoded({ extended: false }));
router.use(bodyParser.json());
router.use(cors());

const users = [
  { email: "admin@example.com", password: "secret" }, // Voorbeeldgebruikers
];

router.get("/", (req, res, next) => {
  res.render("index.ejs");
});

router.post("/", (req, res) => {
  const { email, password } = req.body;
  const user = users.find((u) => u.email === email && u.password === password);

  if (user) {
    res.json({ message: "Login successful" });
  } else {
    res.status(401).json({ message: "Invalid email or password" });
  }
});

export default router;
