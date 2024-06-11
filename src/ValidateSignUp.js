import { use, expect } from 'chai'
import assert from "assert"

const validateSubdomainRequestChaiExpect = (req, res, next) => {
  try {
    assert(req.body.subdomainName, "sub-domain naam veld ontbreekt of is onjuist");
    expect(req.body.subdomainName).to.be.a("string");
    expect(req.body.subdomainName).to.not.be.empty;
    expect(req.body.subdomainName)
      .to.match(
        /^[a-zA-Z0-9]{1,15}$/,
        "Sub-domein naam kan alleen bestaan uit letters en nummers en moet tussen de 1 en 15 karakters lang zijn"
      );

    assert(req.body.emailAddress, "Email adres veld ontbreekt of is incorrect");
    expect(req.body.emailAddress, "Email adres veld ontbreekt of is incorrect")
      .to.be.a("string");
    expect(
      req.body.emailAddress,
      "Email adres veld kan niet leeg zijn"
    ).to.not.be.empty;
    expect(req.body.emailAddress)
      .to.match(
        /^[\w.%+-]+@(?:student\.uva\.nl|umail\.leidenuniv\.nl|students\.uu\.nl|student\.ru\.nl|eur\.nl|student\.wur\.nl|student\.maastrichtuniversity\.nl|student\.tue\.nl|student\.tilburguniversity\.edu|student\.vu\.nl|student\.rug\.nl|student\.tudelft\.nl|student\.utwente\.nl|student\.fryslan\.frl|student\.ou\.nl|auc\.nl|student\.nyenrode\.nl|student\.uvh\.nl|student\.pthu\.nl|student\.tias\.edu|studentihs\.nl|studentiss\.nl|student\.roac\.nl|student\.uwc\.nl|student\.igsr\.edu\.sr|student\.ifdp\.edu\.sr|student\.uoc\.cw|hva\.nl|student\.hr\.nl|student\.fontys\.nl|hhs\.nl|student\.hanze\.nl|student\.avans\.nl|student\.saxion\.nl|student\.windesheim\.nl|student\.han\.nl|student\.nhlstenden\.com|student\.zuyd\.nl|student\.inholland\.nl|student\.codarts\.nl|student\.aeres\.nl|student\.ahk\.nl|student\.buas\.nl|student\.has\.nl|student\.vhluniversity\.com|student\.amfi\.nl|student\.reinwardtacademie\.nl|student\.filmacademie\.nl|student\.artez\.nl|student\.rietveldacademie\.nl|student\.wittenborg\.eu|student\.akvstjoost\.nl|student\.wdka\.nl)$/,
        "Email adres moet een @ en . bevatten, ook moet dit een valide school email zijn"
      );

    assert(req.body.password, "Wachtwoord veld ontbreekt of in onjuist");
    expect(req.body.password).to.be.a("string");
    expect(req.body.password).to.not.be.empty;
    expect(req.body.password)
      .to.match(
        /^(?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{8,20}$/,
        "Wachtwoord moet een hoofdletter, een kleine letter en een cijfer bevatten en moet tussen de 8 en 20 karakters lang zijn"
      );

    assert(req.body.productPackage, "Product pakket ontbreekt");
    expect(req.body.productPackage).to.be.a("string");
    expect(req.body.productPackage).to.not.be.empty;
    expect(req.body.productPackage)
      .to.match(
        /^(free|basic|premium)$/,
        "Product pakket moet gratis, basis of premium zijn"
      );

    next();
  } catch (error) {

    const errorArray = error.message.split(":");
    res.status(400).json({
      status: 400,
      message: errorArray[0],
      data: {}
  });
  }
};

export default validateSubdomainRequestChaiExpect;
