package com.blockchain.demo.models;

import java.util.Set;

import org.hyperledger.fabric.sdk.Enrollment;
import org.hyperledger.fabric.sdk.User;

public class BlockchainUser implements User {
    private String name;
    private Set<String> roles;
    private String account;
    private String affiliation;
    private Enrollment enrollment;
    private String mspId;
    private String enrollmentSecret;
    private boolean enrolled;
    private boolean registered;

    public void setName(final String name) {
        this.name = name;
    }

    public void setRoles(final Set<String> roles) {
        this.roles = roles;
    }

    public void setAccount(final String account) {
        this.account = account;
    }

    public void setAffiliation(final String affiliation) {
        this.affiliation = affiliation;
    }

    public void setEnrollment(final Enrollment enrollment) {
        this.enrollment = enrollment;
    }

    public void setMspId(final String mspId) {
        this.mspId = mspId;
    }

    @Override
    public String getName() {
        return name;
    }

    @Override
    public Set<String> getRoles() {
        return roles;
    }

    @Override
    public String getAccount() {
        return account;
    }

    @Override
    public String getAffiliation() {
        return affiliation;
    }

    @Override
    public Enrollment getEnrollment() {
        return enrollment;
    }

    @Override
    public String getMspId() {
        return mspId;
    }

    public String getEnrollmentSecret() {
        return enrollmentSecret;
    }

    public void setEnrollmentSecret(final String enrollmentSecret) {
        this.enrollmentSecret = enrollmentSecret;
    }

    public boolean isEnrolled() {
        return enrolled;
    }

    public void setEnrolled(final boolean enrolled) {
        this.enrolled = enrolled;
    }

    public boolean isRegistered() {
        return registered;
    }

    public void setRegistered(final boolean registered) {
        this.registered = registered;
    }
}
