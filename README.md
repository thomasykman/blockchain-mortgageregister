This repository was created as the practical part of our (myself & Thomas Ykman) official exam paper for the Blockchain Developer & Architect course @ Howest. Read our full exam paper at https://www.kareldesmet.be/assets/files/BC-project-ykman-desmet.pdf to learn more about the business logic and how to interpret and run this code.

In short, this repository represents a modified [*fabric-samples*](https://github.com/hyperledger/fabric-samples) folder, which is used in its original form as a basis for running sample applications that interact with a Hyperledger Fabric network. We have used bits and pieces of different examples and strung them together with some of our own findings to eventually get to the chaincode that you can find in the folder chaincode/mortgageregister. 

Besides the actual chaincode, there are 2 Node.js processes you can run to start the web application and the API that connects to the Hyperledger Fabric network. 

Front-end: /fabcar/javascript/front-end/index.js
Back-end code: /fabcar/javascript/app.js

We found it very dangerous to limit the contents of this repository to just our own example, as it relies on the files from other examples. Believe me, we've tried. I can't remember how many times I've started the network, installed and instantiated the chaincode, only to find out that something wasn't working and to start all over again. That's why serving this folder as a whole was a safer bet for the submission of our final paper. Demo'ing this was quite nerve-wrecking :)

## Warning
Be warned: the exam paper above is over 30 pages long and contains significant detail with regards to how you can interact with a blockchain network. However, as we have noticed, Hyperledger Fabric is a project which is in constant development and requires a significant amount of fine-tuning. The slightest version bump might render the code and/or terminal commands from our paper useless. Therefore, it is extremely important that you use Hyperledger Fabric v1.4 to run the network with which the application will interact. 

## Disclaimer
At the time of writing, the latest version of Hyperledger Fabric is already v2.2 so please do not use this repository, nor an outdated Hyperledger Fabric version, as a basis for anything which you might want to use in production.
