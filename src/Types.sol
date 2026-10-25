// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.8.0 <0.9.0;

library Types {
    enum Roles {
        Patient,
        Doctor,
        FirstResponder,
        Admin,
        Hospital
    }

    struct ProviderAttributes {
        bytes32 department;
        bytes32 institution;
        address issuedBy;
        uint48 validUntil;
    }

    struct RecordRef {
        address recordContract;
        bytes32 ipfsCid;
        uint48 createdAt;
        bool active;
    }

    struct Incident {
        bytes32 id;
        address hospital;
        bytes32 locationHash;
        uint48 createdAt;
        uint48 expiresAt;
        bool active;
    }

    enum RecordType {
        LabResult,
        Prescript,
        EmergencySummary,
        Imaging,
        Consultation,
        Vaccination,
        Surgery,
        Other
    }

    struct RecordMetadata {
        RecordType recordType;
        uint256 createdAt;
        uint256 updatedAt;
    }
}
