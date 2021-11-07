#!/bin/bash

function _exit(){
    printf "Exiting:%s\n" "$1"
    exit -1
}

# Exit on first error, print all commands.
set -ev
set -o pipefail

# Where am I?
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export FABRIC_CFG_PATH="${PWD}/../config"

cd "${DIR}/../test-network/"

./network.sh down
./network.sh up createChannel -ca -s couchdb -c canal-tt

# Behind the scenes, this script uses the chaincode lifecycle to package, install, query installed chaincode, approve chaincode for both Org1 and Org2, and finally commit the chaincode
./network.sh deployCC -ccn vaccine -ccp ../techtalk/chaincode-vaccine -ccl go -c canal-tt
./network.sh deployCC -ccn user -ccp ../techtalk/chaincode-user -ccl go -c canal-tt
