#!/bin/bash

# Se puede usar peer CLI para interactuar con la red blockchain
cd "${PWD}/../test-network"

# para acceder a los binarios (como peer)
export PATH=${PWD}/../bin:$PATH

export FABRIC_CFG_PATH=$PWD/../config/

# Environment variables for Org1
# export CORE_PEER_TLS_ENABLED=true
# export CORE_PEER_LOCALMSPID="Org1MSP"
# export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
# export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
# export CORE_PEER_ADDRESS=localhost:7051


# Environment variables for Org2
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
export CORE_PEER_ADDRESS=localhost:9051

function invoke(){
  # $1 is the chaincode name
  peer chaincode invoke -o localhost:7050   \
    -C canal-tt -n $1   \
    --ordererTLSHostnameOverride orderer.example.com  \
    --tls   \
    --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem   \
    --peerAddresses localhost:7051  \
    --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt  \
    --peerAddresses localhost:9051  \
    --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt  \
    -c '{"function":"InitLedger","Args":[]}'  
}
# Init ledgers
invoke vaccine
invoke user

# waiting to invoke finished :)
sleep 2

echo "Vaccine:"
peer chaincode query -C canal-tt -n vaccine -c '{"Args":["FindAll"]}'
echo "User:"
peer chaincode query -C canal-tt -n user -c '{"Args":["FindAll"]}'



