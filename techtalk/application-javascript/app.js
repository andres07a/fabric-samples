'use strict';

const { Gateway, Wallets } = require('fabric-network');
const FabricCAServices = require('fabric-ca-client');
const path = require('path');
const { buildCAClient, registerAndEnrollUser, enrollAdmin } = require('../../test-application/javascript/CAUtil.js');
const { buildCCPOrg1, buildWallet } = require('../../test-application/javascript/AppUtil.js');

const channelName = 'canal-tt';
const chaincodeVaccine = 'vaccine';
const chaincodeUser = 'user';
const mspOrg1 = 'Org1MSP';
const walletPath = path.join(__dirname, 'wallet');
const org1UserId = 'appUser';

function prettyJSONString(inputString) {
	return JSON.stringify(JSON.parse(inputString), null, 2);
}

// pre-requisites:
// - fabric-sample two organization test-network setup with two peers, ordering service,
//   and 2 certificate authorities
//         ===> from directory /fabric-samples/test-network
//         ./network.sh up createChannel -ca
// - Use any of the techtalk chaincodes deployed on the channel "canal-tt"
//   with the chaincode name of "vaccine". The following deploy command will package,
//   install, approve, and commit the javascript chaincode, all the actions it takes
//   to deploy a chaincode to a channel.
//         ===> from directory /fabric-samples/test-network
//         ./network.sh deployCC -ccn vaccine -ccp ../techtalk/chaincode-javascript/ -ccl javascript
// - Be sure that node.js is installed
//         ===> from directory /fabric-samples/techtalk/application-javascript
//         node -v
// - npm installed code dependencies
//         ===> from directory /fabric-samples/techtalk/application-javascript
//         npm install
// - to run this test application
//         ===> from directory /fabric-samples/techtalk/application-javascript
//         node app.js

// NOTE: If you see  kind an error like these:
/*
		2020-08-07T20:23:17.590Z - error: [DiscoveryService]: send[canal-tt] - Channel:canal-tt received discovery error:access denied
		******** FAILED to run the application: Error: DiscoveryService: canal-tt error: access denied

	 OR

	 Failed to register user : Error: fabric-ca request register failed with errors [[ { code: 20, message: 'Authentication failure' } ]]
	 ******** FAILED to run the application: Error: Identity not found in wallet: appUser
*/
// Delete the /fabric-samples/techtalk/application-javascript/wallet directory
// and retry this application.
//
// The certificate authority must have been restarted and the saved certificates for the
// admin and application user are not valid. Deleting the wallet store will force these to be reset
// with the new certificate authority.
//

/**
 *  A test application to show basic queries operations with any of the vaccine chaincodes
 *   -- How to submit a transaction
 *   -- How to query and check the results
 *
 * To see the SDK workings, try setting the logging to show on the console before running
 *        export HFC_LOGGING='{"debug":"console"}'
 */
async function main() {
	try {
		// build an in memory object with the network configuration (also known as a connection profile)
		const ccp = buildCCPOrg1();

		// build an instance of the fabric ca services client based on
		// the information in the network configuration
		const caClient = buildCAClient(FabricCAServices, ccp, 'ca.org1.example.com');

		// setup the wallet to hold the credentials of the application user
		const wallet = await buildWallet(Wallets, walletPath);

		// in a real application this would be done on an administrative flow, and only once
		await enrollAdmin(caClient, wallet, mspOrg1);

		// in a real application this would be done only when a new user was required to be added
		// and would be part of an administrative flow
		await registerAndEnrollUser(caClient, wallet, mspOrg1, org1UserId, 'org1.department1');

		// Create a new gateway instance for interacting with the fabric network.
		// In a real application this would be done as the backend server session is setup for
		// a user that has been verified.
		const gateway = new Gateway();

		try {
			// setup the gateway instance
			// The user will now be able to create connections to the fabric network and be able to
			// submit transactions and query. All transactions submitted by this gateway will be
			// signed by this user using the credentials stored in the wallet.
			await gateway.connect(ccp, {
				wallet,
				identity: org1UserId,
				discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
			});

			// Build a network instance based on the channel where the smart contract is deployed
			const network = await gateway.getNetwork(channelName);

			// Get the contract from the network.
			const vaccineSmartContract = network.getContract(chaincodeVaccine);
			const userSmartContract = network.getContract(chaincodeUser);

			console.log('\n--> Submit Transaction: InitLedger for vaccine');
			await vaccineSmartContract.submitTransaction('InitLedger');
			console.log('*** Result: committed');

			console.log('\n--> Submit Transaction: InitLedger for user');
			await userSmartContract.submitTransaction('InitLedger');
			console.log('*** Result: committed');

			console.log('\n--> Evaluate Transaction: FindAll, function returns all the current assets on the ledger');
			let result = await vaccineSmartContract.evaluateTransaction('FindAll');
			console.log(`*** Result: ${prettyJSONString(result.toString())}`);


			console.log('\n--> Submit Transaction: Create, creates new asset with key, id, name, dose, scheme');
			result = await vaccineSmartContract.submitTransaction('Create', '20-1-sch1', '20', 'soberana-2', '1', 'sch1');
			console.log('*** Result: committed');
			if (`${result}` !== '') {
				console.log(`*** Result: ${prettyJSONString(result.toString())}`);
			}

			console.log('\n--> Evaluate Transaction: FindOne, function returns an asset with a given assetID');
			result = await vaccineSmartContract.evaluateTransaction('FindOne', '20-1-sch1');
			console.log(`*** Result: ${prettyJSONString(result.toString())}`);


			console.log('\n--> Submit Transaction: Update 20-1-sch1, change the name');
			await vaccineSmartContract.submitTransaction('Update', '20-1-sch1', '20', 'phizer', '1', 'sch1');
			console.log('*** Result: committed');

			console.log('\n--> Evaluate Transaction: FindOne, function returns "20-1-sch1" attributes');
			result = await vaccineSmartContract.evaluateTransaction('FindOne', '20-1-sch1');
			console.log(`*** Result: ${prettyJSONString(result.toString())}`);

			console.log('\n--> Evaluate Transaction: FindAll, function returns all the current assets on the ledger for last time :)');
			result = await vaccineSmartContract.evaluateTransaction('FindAll');
			console.log(`*** Result: ${prettyJSONString(result.toString())}`);

		} finally {
			// Disconnect from the gateway when the application is closing
			// This will close all connections to the network
			gateway.disconnect();
		}
	} catch (error) {
		console.error(`******** FAILED to run the application: ${error}`);
	}
}

main();
