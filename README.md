# SDK FABRIC NETWORK

Fabric network for SDK v1.0.0

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. There will be a separate deployment for notes on how to deploy the project on a live system.

### Prerequisites

Be sure to install Hyperledger Fabric [Prerequisites](http://hyperledger-fabric.readthedocs.io/en/release/prereqs.html)

Execute following command in current directory

```
curl -sSL https://goo.gl/kFFqh5 | bash -s 1.0.6
```

### Installing

Generate certifcates

```
./byfn.sh -m generate
```

Start Blockchain network

```
./byfn.sh -m up
```

End with an example of getting some data out of the system or using it for a little demo

## Other operations

Stopping blockchain network

```
./byfn.sh -m down
```

Help

```
./byfn.sh -h
```

## Blockchain Explorer
Before going further, please take note to start fabric with a cli timeout set to some considerable time.
```
./byfn.sh -m up -t 600000
```

Explorer is already configured and can be access via - `http://localhost:8080`

In order to trigger an update on Explorer, we need to simulate an invoke.

1. On a seperate terminal enter the following:
```
docker exec -it cli bash
```
After a successful execution, you should be redirected to
```
root@c451fee40fe5:/opt/gopath/src/github.com/hyperledger/fabric/peer#
```
2. Use the following command to invoke the chaincode. Replace `args` with the function you want to call. In this example, it is `createAsset`.

```
peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C mychannel -n mycc -c '{"Args": ["createAsset", "{\"Id\":2, \"Name\":\"Asset2\"}"]}'
```
After a successful invoke, an new transaction record should be added on Explorer.