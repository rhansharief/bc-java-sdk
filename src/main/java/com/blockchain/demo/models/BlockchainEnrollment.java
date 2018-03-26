package com.blockchain.demo.models;

import java.io.Serializable;
import java.security.PrivateKey;

import org.hyperledger.fabric.sdk.Enrollment;

public class BlockchainEnrollment implements Enrollment, Serializable {

    private static final long serialVersionUID = -2784835212445309006L;
    private final PrivateKey privateKey;
    private final String certificate;

    public BlockchainEnrollment(PrivateKey privateKey, String certificate) {

        this.certificate = certificate;

        this.privateKey = privateKey;
    }

    @Override
    public PrivateKey getKey() {

        return privateKey;
    }

    @Override
    public String getCert() {
        return certificate;
    }
}
