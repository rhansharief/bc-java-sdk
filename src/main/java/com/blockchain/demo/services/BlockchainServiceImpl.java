package com.blockchain.demo.services;

import java.io.IOException;
import java.io.StringWriter;
import java.net.MalformedURLException;
import java.nio.file.Paths;
import java.security.PrivateKey;
import java.util.HashMap;
import java.util.HashSet;
import java.util.Map;
import java.util.Properties;
import java.util.Set;

import org.bouncycastle.openssl.PEMWriter;
import org.hyperledger.fabric.sdk.Enrollment;
import org.hyperledger.fabric.sdk.security.CryptoSuite;
import org.hyperledger.fabric_ca.sdk.EnrollmentRequest;
import org.hyperledger.fabric_ca.sdk.HFCAClient;
import org.hyperledger.fabric_ca.sdk.HFCAInfo;
import org.hyperledger.fabric_ca.sdk.RegistrationRequest;
import org.hyperledger.fabric_ca.sdk.exception.InvalidArgumentException;
import org.springframework.stereotype.Service;

import com.blockchain.demo.models.BlockchainUser;
import com.blockchain.demo.models.Organization;
import com.blockchain.demo.models.SampleStore;
import com.blockchain.demo.models.SampleUser;

import static java.lang.String.format;
import static java.nio.charset.StandardCharsets.UTF_8;

@Service
public class BlockchainServiceImpl implements BlockchainService {

    Map<String, Properties> clientTLSProperties = new HashMap<>();

    @Override
    public void invoke() {
        try {

            Organization organization = new Organization("Org1", "Org1MSP");

            String caName = organization.getCAName(); //Try one of each name and no name.
            if (caName != null && !caName.isEmpty()) {
                    organization.setCAClient(HFCAClient.createNewInstance(caName, organization.getCALocation(), organization.getCAProperties()));

            } else {
                organization.setCAClient(HFCAClient.createNewInstance(organization.getCALocation(), organization
                    .getCAProperties()));
            }


            enrollUsersSetup(organization);


        } catch (MalformedURLException e) {
            e.printStackTrace();
        } catch (InvalidArgumentException e) {
            e.printStackTrace();
        } catch (Exception e) {
            e.printStackTrace();
        }


    }
    public void enrollUsersSetup(Organization organization) throws Exception {

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

        if (!user.isRegistered()) {  // users need to be registered AND enrolled
            RegistrationRequest rr = new RegistrationRequest(user.getName(), "org1.department1");
            user.setEnrollmentSecret(ca.register(rr, admin));
        }
        if (!user.isEnrolled()) {
            user.setEnrollment(ca.enroll(user.getName(), user.getEnrollmentSecret()));
        }

    }

    static String getPEMStringFromPrivateKey(PrivateKey privateKey) throws IOException {
        StringWriter pemStrWriter = new StringWriter();
        PEMWriter pemWriter = new PEMWriter(pemStrWriter);

        pemWriter.writeObject(privateKey);

        pemWriter.close();

        return pemStrWriter.toString();
    }

}
