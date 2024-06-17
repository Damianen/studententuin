import { expect } from "chai";
import assert from "assert";

const validateSubdomainRequestChaiExpect = (req, res, next) => {
    try {
        assert(
            req.body.subdomainName,
            "sub-domain naam veld ontbreekt of is onjuist"
        );
        expect(
            req.body.subdomainName,
            "Sub-domein naam moet een string zijn"
        ).to.be.a("string");
        expect(req.body.subdomainName, "Sub-domein naam kan niet leeg zijn").to
            .not.be.empty;
        expect(
            req.body.subdomainName,
            "Sub-domein naam kan alleen bestaan uit letters en nummers en moet tussen de 1 en 15 karakters lang zijn"
        ).to.match(/^[a-zA-Z0-9]{1,15}$/);

        assert(
            req.body.emailAddress,
            "Email adres veld ontbreekt of is incorrect"
        );
        expect(
            req.body.emailAddress,
            "Email adres moet een string zijn"
        ).to.be.a("string");
        expect(req.body.emailAddress, "Email adres kan niet leeg zijn").to.not
            .be.empty;
        expect(
            req.body.emailAddress,
            "Email adres moet een valide school email zijn"
        ).to.match(
            /^[\w.%+-]+@(student\.uva\.nl|umail\.leidenuniv\.nl|students\.uu\.nl|student\.ru\.nl|eur\.nl|student\.wur\.nl|student\.maastrichtuniversity\.nl|student\.tue\.nl|student\.tilburguniversity\.edu|student\.vu\.nl|student\.rug\.nl|student\.tudelft\.nl|student\.utwente\.nl|student\.fryslan\.frl|student\.ou\.nl|auc\.nl|student\.nyenrode\.nl|student\.uvh\.nl|student\.pthu\.nl|student\.tias\.edu|studentihs\.nl|studentiss\.nl|student\.roac\.nl|student\.uwc\.nl|student\.igsr\.edu\.sr|student\.ifdp\.edu\.sr|student\.uoc\.cw|hva\.nl|student\.hr\.nl|student\.fontys\.nl|hhs\.nl|student\.hanze\.nl|student\.avans\.nl|student\.saxion\.nl|student\.windesheim\.nl|student\.han\.nl|student\.nhlstenden\.com|student\.zuyd\.nl|student\.inholland\.nl|student\.codarts\.nl|student\.aeres\.nl|student\.ahk\.nl|student\.buas\.nl|student\.has\.nl|student\.vhluniversity\.com|student\.amfi\.nl|student\.reinwardtacademie\.nl|student\.filmacademie\.nl|student\.artez\.nl|student\.rietveldacademie\.nl|student\.wittenborg\.eu|student\.akvstjoost\.nl|student\.wdka\.nl)$/
        );

        assert(req.body.password, "Wachtwoord veld ontbreekt of is onjuist");
        expect(req.body.password, "Wachtwoord moet een string zijn").to.be.a(
            "string"
        );
        expect(req.body.password, "Wachtwoord kan niet leeg zijn").to.not.be
            .empty;
        expect(
            req.body.password,
            "Wachtwoord moet een hoofdletter, een kleine letter, een cijfer, en een speciaal karakter bevatten (behalve een punt) en moet tussen de 8 en 20 karakters lang zijn"
        ).to.match(/^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])(?=.*[^a-zA-Z\d.]).{8,20}$/);
        expect(
            req.body.password,
            "Wachtwoord mag niet de naam van het subdomein bevatten"
        ).to.not.contain(req.body.subdomainName);

        assert(req.body.productPackage, "Product pakket ontbreekt");
        expect(
            req.body.productPackage,
            "Product pakket moet een string zijn"
        ).to.be.a("string");
        expect(req.body.productPackage, "Product pakket kan niet leeg zijn").to
            .not.be.empty;
        expect(
            req.body.productPackage,
            "Product pakket moet 'free', 'basic', of 'premium' zijn"
        ).to.match(/^(free|basic|premium)$/);

        next();
    } catch (error) {
        res.status(400).json({
            status: 400,
            message: error.message,
            data: {},
        });
    }
};

export default validateSubdomainRequestChaiExpect;
