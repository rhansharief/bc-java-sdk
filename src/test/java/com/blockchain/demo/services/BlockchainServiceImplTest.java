package com.blockchain.demo.services;

import java.io.File;
import java.io.FileInputStream;
import java.io.IOException;
import java.io.Reader;
import java.io.StringReader;
import java.io.StringWriter;
import java.net.MalformedURLException;
import java.security.NoSuchAlgorithmException;
import java.security.NoSuchProviderException;
import java.security.PrivateKey;
import java.security.spec.InvalidKeySpecException;
import java.util.Collection;
import java.util.HashMap;
import java.util.HashSet;
import java.util.LinkedList;
import java.util.Map;
import java.util.Properties;
import java.util.Set;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.CompletionException;

import org.hyperledger.fabric.sdk.BlockEvent;
import org.hyperledger.fabric.sdk.Peer.PeerRole;

import org.apache.commons.io.IOUtils;
import org.bouncycastle.asn1.pkcs.PrivateKeyInfo;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.bouncycastle.openssl.PEMParser;
import org.bouncycastle.openssl.PEMWriter;
import org.bouncycastle.openssl.jcajce.JcaPEMKeyConverter;
import org.hyperledger.fabric.sdk.ChaincodeID;
import org.hyperledger.fabric.sdk.ChaincodeResponse;
import org.hyperledger.fabric.sdk.Channel;
import org.hyperledger.fabric.sdk.Enrollment;
import org.hyperledger.fabric.sdk.HFClient;
import org.hyperledger.fabric.sdk.NetworkConfig;
import org.hyperledger.fabric.sdk.Peer;
import org.hyperledger.fabric.sdk.ProposalResponse;
import org.hyperledger.fabric.sdk.TransactionProposalRequest;
import org.hyperledger.fabric.sdk.exception.ProposalException;
import org.hyperledger.fabric.sdk.security.CryptoSuite;
import org.hyperledger.fabric_ca.sdk.EnrollmentRequest;
import org.hyperledger.fabric_ca.sdk.HFCAClient;
import org.hyperledger.fabric_ca.sdk.RegistrationRequest;
import org.junit.Assert;
import org.junit.Test;

import com.blockchain.demo.models.BlockchainEnrollment;
import com.blockchain.demo.models.BlockchainUser;
import com.blockchain.demo.models.Organization;


import static java.lang.String.format;
import static java.nio.charset.StandardCharsets.UTF_8;

public class BlockchainServiceImplTest {
    Map<String, Properties> clientTLSProperties = new HashMap<>();

    @Test
    public void testInvoke() {
        try {
            boolean runningFabricCATLS = true;
            File f = new File("src/main/resources/network-config.json");
            NetworkConfig config = NetworkConfig.fromJsonFile(f);
            Assert.assertNotNull(config);

            Organization organization = new Organization("Org1", "Org1MSP");
            organization.setCAName("ca.org1.example.com");
            organization.setCALocation("https://localhost:7054");

            if (runningFabricCATLS) {
                String cert = "src/main/resources/ca.org1.example.com-cert.pem";
                File cf = new File(cert);
                if (!cf.exists() || !cf.isFile()) {
                    throw new RuntimeException("TEST is missing cert file " + cf.getAbsolutePath());
                }
                Properties properties = new Properties();
                properties.setProperty("pemFile", cf.getAbsolutePath());

                properties.setProperty("allowAllHostNames", "true"); //testing environment only NOT FOR PRODUCTION!

                organization.setCAProperties(properties);
            }

            String caName = organization.getCAName(); //Try one of each name and no name.
            if (caName != null && !caName.isEmpty()) {
                organization.setCAClient(HFCAClient.createNewInstance(caName, organization.getCALocation(), organization.getCAProperties()));

            } else {
                organization.setCAClient(HFCAClient.createNewInstance(organization.getCALocation(), organization
                    .getCAProperties()));
            }


            //CONFIGS
            boolean isRunningFabricTls = true;
            ////////////////////////////
            //Set up USERS

            //SampleUser can be any implementation that implements org.hyperledger.fabric.sdk.User Interface

            ////////////////////////////
            // get users for all orgs

            HFCAClient ca = organization.getCAClient();
            ca.setCryptoSuite(CryptoSuite.Factory.getCryptoSuite());

            if (isRunningFabricTls) {
                //This shows how to get a client TLS certificate from Fabric CA
                // we will use one client TLS certificate for orderer peers etc.
                final EnrollmentRequest enrollmentRequestTLS = new EnrollmentRequest();
                enrollmentRequestTLS.addHost("localhost");
                enrollmentRequestTLS.setProfile("tls");
                final Enrollment enroll = ca.enroll("admin", "adminpw", enrollmentRequestTLS);
                final String tlsCertPEM = enroll.getCert();
                final String tlsKeyPEM = getPEMStringFromPrivateKey(enroll.getKey());

                final Properties tlsProperties = new Properties();

                tlsProperties.put("clientKeyBytes", tlsKeyPEM.getBytes(UTF_8));
                tlsProperties.put("clientCertBytes", tlsCertPEM.getBytes(UTF_8));
                clientTLSProperties.put(organization.getName(), tlsProperties);
                //Save in samplestore for follow on tests.
//            sampleStore.storeClientPEMTLCertificate(organization, tlsCertPEM);
//            sampleStore.storeClientPEMTLSKey(organization, tlsKeyPEM);
            }

            //No need to register admin. Only need to enroll
            BlockchainUser admin = new BlockchainUser();
            admin.setName("admin");
            admin.setMspId("Org1MSP");

            if (!admin.isEnrolled()) {  //Preregistered admin only needs to be enrolled with Fabric caClient.
                admin.setEnrollmentSecret("adminpw");
                admin.setEnrollment(ca.enroll(admin.getName(), "adminpw"));
                admin.setMspId("Org1MSP");
            }

            BlockchainUser user = new BlockchainUser();
            user.setName("user");
            user.setAffiliation("org1.department1");
            user.setAccount("user");
            user.setMspId("Org1MSP");

            Set<String> roles = new HashSet<>();
            roles.add("peer");
            roles.add("app");
            roles.add("user");
            user.setRoles(roles);
//            if (!user.isRegistered()) {  // users need to be registered AND enrolled
//                RegistrationRequest rr = new RegistrationRequest(user.getName(), "org1.department1");
//                user.setEnrollmentSecret(ca.register(rr, admin));
//                System.out.println("ENROLLMENT SECRET");
//                System.out.println(user.getEnrollmentSecret());
//            }
            if (!user.isEnrolled()) {
//                final EnrollmentRequest enrollmentRequestTLS = new EnrollmentRequest();
//                enrollmentRequestTLS.addHost("localhost");
//                enrollmentRequestTLS.setProfile("tls");
                user.setEnrollment(ca.enroll(user.getName(), "qRpDupryRaPo"));
//                user.setEnrollment(ca.enroll(user.getName(), user.getEnrollmentSecret()));
            }


            System.out.println(user);
            Assert.assertNotNull(user);
            Assert.assertNotNull(user.getEnrollment());



            BlockchainUser peerAdmin = new BlockchainUser();
            peerAdmin.setName("Org1Admin");
            peerAdmin.setAffiliation("org1.department1");
            peerAdmin.setAccount("Org1Admin");
            peerAdmin.setMspId("Org1MSP");
            peerAdmin.setRoles(roles);

            File skFile = new File("src/main/resources/614f4a0daaa9af87d0c3e0642b9f7265dd0a9d29f030f22529e19bda2a5b708c_sk");
            File pemFile = new File("src/main/resources/Admin@org1.example.com-cert.pem");

            String certificate = new String(IOUtils.toByteArray(new FileInputStream(pemFile)), "UTF-8");

            PrivateKey privateKey = getPrivateKeyFromBytes(IOUtils.toByteArray(new FileInputStream(skFile)));

            BlockchainEnrollment blockchainEnrollment = new BlockchainEnrollment(privateKey, certificate);
            peerAdmin.setEnrollment(blockchainEnrollment);

            HFClient client = HFClient.createNewInstance();

            client.setCryptoSuite(CryptoSuite.Factory.getCryptoSuite());

            client.setUserContext(peerAdmin);

            Channel channel = client.loadChannelFromConfig("nimble-trx-channel", config);
            Assert.assertNotNull(channel);
            System.out.println(channel);
            channel.getPeers().stream().forEach(peer -> {
                channel.getPeersOptions(peer).setPeerRoles(PeerRole.NO_EVENT_SOURCE);
            });

            channel.initialize();

            String CHAIN_CODE_NAME = "mycc";
            String CHAIN_CODE_PATH = "github.com/hyperledger/fabric/examples/chaincode/go";
            String CHAIN_CODE_VERSION_11 = "1.0";

            ChaincodeID chaincodeID = ChaincodeID.newBuilder().setName(CHAIN_CODE_NAME)
                .setPath(CHAIN_CODE_PATH)
                .setVersion(CHAIN_CODE_VERSION_11)
                .build();

            //invoke
            try {
                Collection<ProposalResponse> successful = new LinkedList<>();
                Collection<ProposalResponse> failed = new LinkedList<>();

                ///////////////
                /// Send transaction proposal to all peers
                TransactionProposalRequest transactionProposalRequest = client.newTransactionProposalRequest();
                transactionProposalRequest.setChaincodeID(chaincodeID);
                transactionProposalRequest.setFcn("createAsset");
                transactionProposalRequest.setArgs(new String[] {"{\"Id\":2, \"Name\":\"Asset2\"}"});
//                transactionProposalRequest.setProposalWaitTime(testConfig.getProposalWaitTime());
                if (user != null) { // specific user use that
                    transactionProposalRequest.setUserContext(user);
                }
//                out("sending transaction proposal to all peers with arguments: move(a,b,%s)", moveAmount);

                Collection<ProposalResponse> invokePropResp = channel.sendTransactionProposal(transactionProposalRequest);
                for (ProposalResponse response : invokePropResp) {
                    if (response.getStatus() == ChaincodeResponse.Status.SUCCESS) {
                        System.out.printf("Successful transaction proposal response Txid: %s from peer %s", response.getTransactionID(), response.getPeer().getName());
                        successful.add(response);
                    } else {
                        failed.add(response);
                    }
                }

                System.out.printf("Received %d transaction proposal responses. Successful+verified: %d . Failed: %d",
                    invokePropResp.size(), successful.size(), failed.size());
                if (failed.size() > 0) {
                    ProposalResponse firstTransactionProposalResponse = failed.iterator().next();

                    throw new ProposalException(format("Not enough endorsers for invoke(move a,b,%s):%d endorser error:%s. Was verified:%b",
                        300, firstTransactionProposalResponse.getStatus().getStatus(), firstTransactionProposalResponse.getMessage(), firstTransactionProposalResponse.isVerified()));

                }
                System.out.println("Successfully received transaction proposal responses.");

                ////////////////////////////
                // Send transaction to orderer
                System.out.printf("Sending chaincode transaction(move a,b,%s) to orderer.", 300);
                if (user != null) {
                    CompletableFuture<BlockEvent.TransactionEvent> test = channel.sendTransaction(successful, user);
                    Assert.assertNotNull(test);
                } else {
                    Assert.assertNotNull(channel.sendTransaction(successful));
                }

            } catch (Exception e) {

                throw new CompletionException(e);

            }



        } catch (MalformedURLException e) {
            e.printStackTrace();
//        } catch (InvalidArgumentException e) {
//            e.printStackTrace();
        } catch (Exception e) {
            e.printStackTrace();
        }

    }

    private String getPEMStringFromPrivateKey(PrivateKey privateKey) throws IOException {
        StringWriter pemStrWriter = new StringWriter();
        PEMWriter pemWriter = new PEMWriter(pemStrWriter);

        pemWriter.writeObject(privateKey);

        pemWriter.close();

        return pemStrWriter.toString();
    }

    static PrivateKey getPrivateKeyFromBytes(byte[] data) throws IOException, NoSuchProviderException, NoSuchAlgorithmException, InvalidKeySpecException {
        final Reader pemReader = new StringReader(new String(data));

        final PrivateKeyInfo pemPair;
        try (PEMParser pemParser = new PEMParser(pemReader)) {
            pemPair = (PrivateKeyInfo) pemParser.readObject();
        }

        PrivateKey privateKey = new JcaPEMKeyConverter().setProvider(BouncyCastleProvider.PROVIDER_NAME).getPrivateKey(pemPair);

        return privateKey;
    }
}
