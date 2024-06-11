import { Router } from "express";
import bodyParser from "body-parser";
import bcrypt from "bcrypt";
import userService from "../server/services/user.service.js";
import cors from "cors";
import jwt from "jsonwebtoken";
import Cookies from "js-cookie";
import validateSubdomainRequestChaiExpect from "./../src/ValidateSignUp.js";
import subDomainService from "../server/services/subdomain.service.js"
import UserSubDomainService from "../server/services/user.subdomain.service.js"

const router = Router();

router.use(bodyParser.urlencoded({ extended: false }));
router.use(bodyParser.json());
router.use(cors());

router.get("/", (req, res, next) => {
  res.render("index.ejs");
});

router.get('/requestForm', (req, res, next) => {
    res.render("index.ejs");
});

router.get('/login', (req, res, next) => {
    res.render("index.ejs");
});

router.get('/handleiding', (req, res, next) => {
    res.render("index.ejs");
});

router.get('/manage', (req, res, next) => {
    res.render("index.ejs");
});

router.get('/about', (req, res, next) => {
    res.render("index.ejs");
});

router.get('/contact', (req, res, next) => {
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

    const token = jwt.sign({ email: user.email }, process.env.JWT_SECRET
    );

    // Als de wachtwoorden niet overeenkomen, geef een foutmelding terug

    if (!passwordMatch) {
      return res.status(401).json({ message: "Invalid email or password" });
    }

    

    // Als de inloggegevens correct zijn, sla de gebruiker op in de sessie
    // req.session = user;

    // Geef een succesbericht terug
    console.log("Token:", token);
    res.cookie("token", token, { httpOnly: true });
    
    res.redirect("/manage");
  } catch (error) {
    console.error("Error:", error);
    res.status(500).json({ message: "An error occurred" });
  }
});

router.post("/requestForm", validateSubdomainRequestChaiExpect, async (req, res) => {
    const {subdomainName, emailAddress, password, productPackage} = req.body;
  

    await userService.getByEmail(emailAddress, (success, error) => {
        if (success) {
            if(success.status === 200){
                res.status(409).json({
                    status: 409,
                    message: "De gebruikte email staat al geregistreerd in het systeem!",
                });
            }
        }else{
            
            subDomainService.getByName(subdomainName, (success, error) => {
                if(success){
                    if(success.status === 200){
                        res.status(409).json({
                            status: 409,
                            message: "Het opgegeven sub-domein naam is al in gebruik!",
                        });
                    }
                }else{
                    userService.insert(emailAddress, password, productPackage, (success, error) => {
                        if(success){
                            if(success.status === 200){
                                subDomainService.insertDomainName(subdomainName, (success, error) => {
                                    if(success){
                                        if(success.status === 200){
                                            UserSubDomainService.insert(emailAddress, subdomainName, (success, error) => {
                                                if(success){
                                                    if(success.status === 200){

                                                        const data = {
                                                            email: emailAddress,
                                                            subdomain: subdomainName
                                                        }

                                                        // fetch('https://webhook.studententuin.nl/new', {
                                                        //     method: 'POST',
                                                        //     headers: {
                                                        //         'Content-Type': 'application/json'
                                                        //     },
                                                        //     body: JSON.stringify(data)
                                                        // });

                                                        res.status(200).json({
                                                            status: 200,
                                                            message: "Gebruiker/Sub-domein succesvol geregistreerd!",
                                                        });
                                                    }
                                                }else{
                                                    res.status(error.status).json({
                                                        status: error.status,
                                                        message: error.message
                                                    });
                                                }
                                            })
                                            
                                        }
                                    }else{
                                        res.status(error.status).json({
                                            status: error.status,
                                            message: error.message
                                        });
                                    }
                                })
                            }
                        }else{
                            res.status(error.status).json({
                                status: error.status,
                                message: error.message
                            });
                        }
                    })
                }
            })
            
        }
    });
  });

export default router;
