/* eslint-disable quote-props */
/* eslint-disable no-var */
/* eslint-disable quotes */
/* eslint-disable strict */
const express = require("express");
const bodyParser = require("body-parser");

const app = express();
app.use(bodyParser.json());

// Setting for Hyperledger Fabric
const { FileSystemWallet, Gateway } = require("fabric-network");
const path = require("path");

// Default to connection-org1. This is changed in the server listener statement at the bottom.
var ccpPath = path.resolve(__dirname, '..', '..', 'first-network', 'connection-org1.json');

// Changeloanstatus route - Only for Org1 (Bank) and Org2 (Notary)
app.post("/api/changeloanstatus", async function (req, res) {
    res.header("Access-Control-Allow-Origin", "*");
    try {
        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), "wallet");
        const wallet = new FileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const userExists = await wallet.exists("user1");
        if (!userExists) {
            console.log(
                'An identity for the user "user1" does not exist in the wallet'
            );
            console.log("Run the registerUser.js application before retrying");
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccpPath, {
            wallet,
            identity: "user1",
            discovery: { enabled: true, asLocalhost: true }
        });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork("mychannel");

        // Get the contract from the network.
        const contract = network.getContract("mortgageregister");

        const transientData = { loan_status: Buffer.from(JSON.stringify(req.body)) };

        const result = await contract
            .createTransaction("changeLoanStatus")
            .setTransient(transientData)
            .submit();

        console.log(
            `Transaction has been submit, result is: ${result.toString()}`
        );
        res.status(200).json({ response: result.toString() });
    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        res.status(500).json({ error: error });
        process.exit(1);
    }
});

// Issueloan route - Only for Org1 (Bank)
app.post("/api/issueloan", async function (req, res) {
    res.header("Access-Control-Allow-Origin", "*");

    try {
        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), "wallet");
        const wallet = new FileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const userExists = await wallet.exists("user1");
        if (!userExists) {
            console.log(
                'An identity for the user "user1" does not exist in the wallet'
            );
            console.log("Run the registerUser.js application before retrying");
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccpPath, {
            wallet,
            identity: "user1",
            discovery: { enabled: true, asLocalhost: true }
        });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork("mychannel");

        // Get the contract from the network.
        const contract = network.getContract("mortgageregister");
        const transientData = { loan: Buffer.from(JSON.stringify(req.body)) };

        // Pass data hard-coded:
        // const transientData = {
        //     loan: Buffer.from(
        //         JSON.stringify({
        //             loanUID: "loan6",
        //             buyer: "thomas",
        //             notary: "notary1",
        //             startDate: "testdate1",
        //             endDate: "testdate2",
        //             loanValue: 50000,
        //             currency: "EUR",
        //             interestRate: 5
        //         })
        //     )
        // };

        // via cURL:
        // curl -X POST -d '{"loanUID": "loan1", "buyer": "buyer1", "notary": "notary1", "startDate": "testdate1", "endDate": "testdate2", "loanValue": 50000, "currency": "EUR", "interestRate": 5}' -H "Content-Type: application/json" http://localhost:8080/api/issueloan

        await contract
            .createTransaction("issueLoan")
            .setTransient(transientData)
            .submit();

        console.log("Transaction has been submitted");
        console.log(transientData);
        res.send("Transaction has been submitted");

        // Disconnect from the gateway.
        await gateway.disconnect();
    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        res.status(500).json({ error: error });
        process.exit(1);
    }
});

// Readloan (public data) route - available for all organizations
app.get("/api/readloan/public/:loanuid", async function (req, res) {
    res.header("Access-Control-Allow-Origin", "*");
    try {
        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), "wallet");
        const wallet = new FileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const userExists = await wallet.exists("user1");
        if (!userExists) {
            console.log(
                'An identity for the user "user1" does not exist in the wallet'
            );
            console.log("Run the registerUser.js application before retrying");
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccpPath, {
            wallet,
            identity: "user1",
            discovery: { enabled: true, asLocalhost: true }
        });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork("mychannel");

        // Get the contract from the network.
        const contract = network.getContract("mortgageregister");

        // example with "loan1" via cURL:
        // curl http://localhost:8080/api/readloan/loan1

        const result = await contract.evaluateTransaction(
            "readLoan",
            req.params.loanuid,
            "collectionLoans"
        );

        console.log(result);

        const parsedResult = JSON.parse(result.toString());

        console.log("Query succeeded. Result:");
        console.log(parsedResult);
        res.send(parsedResult);

        // Disconnect from the gateway.
        await gateway.disconnect();
    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        res.status(500).json({ error: error });
        process.exit(1);
    }
});

// Readloan (private data) route - only for ORG1 and ORG2
app.get("/api/readloan/private/:loanuid", async function (req, res) {
    res.header("Access-Control-Allow-Origin", "*");

    try {
        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), "wallet");
        const wallet = new FileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const userExists = await wallet.exists("user1");
        if (!userExists) {
            console.log(
                'An identity for the user "user1" does not exist in the wallet'
            );
            console.log("Run the registerUser.js application before retrying");
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccpPath, {
            wallet,
            identity: "user1",
            discovery: { enabled: true, asLocalhost: true }
        });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork("mychannel");

        // Get the contract from the network.
        const contract = network.getContract("mortgageregister");

        // example with "loan1" via cURL:
        // curl http://localhost:8080/api/readloan/loan1

        const result = await contract.evaluateTransaction(
            "readLoan",
            req.params.loanuid,
            "collectionLoanPrivateInfo"
        );

        console.log(result);

        const parsedResult = JSON.parse(result.toString());

        console.log("Query succeeded. Result:");
        console.log(parsedResult);
        res.send(parsedResult);

        // Disconnect from the gateway.
        await gateway.disconnect();
    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        res.status(500).json({ error: error });
        process.exit(1);
    }
});

app.listen(8080, function () {
    console.log(`Fabric API listening on port 8080!`);
});
