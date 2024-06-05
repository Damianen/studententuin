import { use, expect } from 'chai'
import assert from "assert"

const validateSubdomainRequestChaiExpect = (req, res, next) => {
  try {
    assert(req.body.subdomainName, "sub-domain name field is missing or incorrect");
    expect(req.body.subdomainName).to.be.a("string");
    expect(req.body.subdomainName).to.not.be.empty;
    expect(req.body.subdomainName)
      .to.match(
        /^[a-zA-Z0-9]{1,15}$/,
        "sub-domain name can only contain letters and numbers and be between 1 and 15 characters long"
      );

    assert(req.body.emailAddress, "emailAddress field is missing or incorrect");
    expect(req.body.emailAddress, "emailAddress field is missing or incorrect")
      .to.be.a("string");
    expect(
      req.body.emailAddress,
      "emailAddress field is missing or incorrect"
    ).to.not.be.empty;
    expect(req.body.emailAddress)
      .to.match(
        /^[\w.%+-]+@(?:student\.uva\.nl|umail\.leidenuniv\.nl|students\.uu\.nl|student\.ru\.nl|eur\.nl|student\.wur\.nl|student\.maastrichtuniversity\.nl|student\.tue\.nl|student\.tilburguniversity\.edu|student\.vu\.nl|student\.rug\.nl|student\.tudelft\.nl|student\.utwente\.nl|student\.fryslan\.frl|student\.ou\.nl|auc\.nl|student\.nyenrode\.nl|student\.uvh\.nl|student\.pthu\.nl|student\.tias\.edu|studentihs\.nl|studentiss\.nl|student\.roac\.nl|student\.uwc\.nl|student\.igsr\.edu\.sr|student\.ifdp\.edu\.sr|student\.uoc\.cw|hva\.nl|student\.hr\.nl|student\.fontys\.nl|hhs\.nl|student\.hanze\.nl|student\.avans\.nl|student\.saxion\.nl|student\.windesheim\.nl|student\.han\.nl|student\.nhlstenden\.com|student\.zuyd\.nl|student\.inholland\.nl|student\.codarts\.nl|student\.aeres\.nl|student\.ahk\.nl|student\.buas\.nl|student\.has\.nl|student\.vhluniversity\.com|student\.amfi\.nl|student\.reinwardtacademie\.nl|student\.filmacademie\.nl|student\.artez\.nl|student\.rietveldacademie\.nl|student\.wittenborg\.eu|student\.akvstjoost\.nl|student\.wdka\.nl)$/,
        "email adress must contain @ and . and be a valid school email adress"
      );

    assert(req.body.password, "password field is missing or incorrect");
    expect(req.body.password).to.be.a("string");
    expect(req.body.password).to.not.be.empty;
    expect(req.body.password)
      .to.match(
        /^(?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{8,20}$/,
        "Password must contain at least one number and one uppercase and lowercase letter, and at least 8 and max 20 characters long"
      );

    assert(req.body.productPackage, "Product package field is missing or incorrect");
    expect(req.body.productPackage).to.be.a("string");
    expect(req.body.productPackage).to.not.be.empty;
    expect(req.body.productPackage)
      .to.match(
        /^(free|basic|premium)$/,
        "Product package must be either free, basic or premium"
      );

    next();
  } catch (error) {
    res.status(400).json({
      status: 400,
      message: error.message,
      data: {}
  });
  }
};

export default validateSubdomainRequestChaiExpect;
