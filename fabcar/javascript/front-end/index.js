const express = require("express");
const app = express();
const port = 3000;
const bodyParser = require("body-parser");
const ejsLint = require("ejs-lint");
const request = require("request");

app.set("view engine", "ejs");
app.use(express.static(__dirname + "/public"));
app.use(bodyParser.urlencoded({ extended: true }));

// Issue loan routes
app.get("/issueloan", (req, res) => res.render("issueloan"));

app.post("/issueloan", function (req, res) {
    try {
        request.post('http://localhost:8080/api/issueloan', {
            json: {
                loanUID: req.body.loanUID,
                issuer: req.body.issuer,
                buyer: req.body.buyer,
                notary: req.body.notary,
                startDate: req.body.startDate,
                endDate: req.body.endDate,
                loanValue: parseInt(req.body.loanValue),
                currency: req.body.currency,
                interestRate: parseFloat(req.body.interestRate)
            }
        }, (error, response, body) => {
            if (error) {
                console.error(error);
            }
        });
        res.render("success");
    }
    catch (error) {
        console.log(error);
        res.render("error");
    }
});


// Public readloan routes
app.get("/readloan/public", (req, res) => res.render("readloanpublic"));

app.get("/readloan/public/:loanuid", function (req, res) {
    try {
        request("http://localhost:8080/api/readloan/public/" + req.params.loanuid, function (error, response, body) {
            if (error) {
                console.log(error);
            }

            res.render("readloanpublic", { data: JSON.parse(body) });
        });
    }
    catch (error) {
        console.log(error);
        res.render("error");
    }
});

// Private readloan routes
app.get("/readloan/private", (req, res) => res.render("readloanprivate"));

app.get("/readloan/private/:loanuid", function (req, res) {
    try {
        request("http://localhost:8080/api/readloan/private/" + req.params.loanuid, function (error, response, body) {
            if (error) {
                console.log(error);
            }

            res.render("readloanprivate", { data: JSON.parse(body) });
        });
    }
    catch (error) {
        console.log(error);
        res.render("error");
    }
})


// Deactivate loan routes
app.get("/deactivateloan", (req, res) => res.render("deactivateloan"));

app.post("/deactivateloan", function (req, res) {
    try {
        request.post('http://localhost:8080/api/changeloanstatus', {
            json: req.body
        }, (error, res, body) => {
            if (error) {
                console.error(error)
                return
            }
            console.log(`statusCode: ${res.statusCode}`)
        });
    }
    catch (error) {
        console.log(error);
        res.render("error");
    }
    res.render("success");
});

app.get("/", (req, res) => res.render("issueloan"));

app.listen(port, () => console.log(`Bank app listening on port ${port}!`));
