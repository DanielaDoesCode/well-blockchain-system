// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Registry.sol";
import "../src/UserRegistry.sol";
import "../src/HistoryFactory.sol";
import "../src/Types.sol";
import "../src/History.sol";
import "../src/Record.sol";
import "../src/RecordPermissions.sol";
import "../src/PolicyVerifier.sol";
import "../src/lib/PoseidonT3.sol";

contract PolicyEnforcementTest is Test {
    WellRegistry internal registry;
    UserRegistry internal userRegistry;
    HistoryFactory internal historyFactory;
    PolicyVerifier internal policyVerifier;

    address internal patient = address(0x5);
    address internal cardiologist = address(0x1); // proof A/C bound to 0x...0001
    address internal otherDoctor = address(0x2);
    address internal attacker = address(0x6);

    bytes32 internal patientCommitment = keccak256("patient-commitment");
    bytes32 internal patientInternalID = keccak256("patient-internal-id");

    // Proof A (full policy: dept 0x1111...11, inst 0x2222...22), doctor 0x...0001
    uint256[2] internal proofA_pA = [
        uint256(17749563003958815978953565070961340461836647835011941111424515940428376907718),
        uint256(4343310052084520144226717958615460601535835760159473856528055549674236178475)
    ];
    uint256[2][2] internal proofA_pB = [
        [uint256(12450916388111255360561991966519019842536551587698933560671727733323273889458),
         uint256(18928262754817229482333214559706493925645187840605793615228364354707553787885)],
        [uint256(4524717860817071861674610336557751705437115145070303733957411309972002126946),
         uint256(20458358698452852315033311969357989944415481030582626379549270293295068073066)]
    ];
    uint256[2] internal proofA_pC = [
        uint256(15200803100421121808902431370918802603255786112118688551529861049755358388782),
        uint256(18090051901104465583333513076733811991829267041731573868887924952520343706137)
    ];

    bytes32 internal deptCommitmentCardio =
        bytes32(uint256(0x18f447fe0780f480be2efdec415f8eac78a73bbaa10085b99a9852354330219b));
    bytes32 internal instCommitmentGeneral =
        bytes32(uint256(0x15cd25d14834cb5b5acee14d67fb0660ae6f080ccc3ff1cbcc7e1ac4cb6dd62d));
    bytes32 internal doctorCommitmentCardio =
        bytes32(uint256(0x2eaa00ff94aa542b726083d39a960fa6c3f79dded2d9bb24511ac009971ee7ed));

    // Proof C (wildcard policy: requiredDept 0, requiredInst 0), doctor 0x...0001
    uint256[2] internal proofC_pA = [
        uint256(20986976941541981685726103504807053431535806275284912517733887727741358330682),
        uint256(20417869664483805163525707514859395779591749919775029984745837668039598769198)
    ];
    uint256[2][2] internal proofC_pB = [
        [uint256(11970993134954215716961056959372164482734160509940813133947504561748551655180),
         uint256(6724208088431634695251177652088215898446575979568478915322412244159596301727)],
        [uint256(8940473049938370432298891352133229952434896298151635036580573409657064393942),
         uint256(2070599948814936597018194581118594440617075271270644531942401135022201207395)]
    ];
    uint256[2] internal proofC_pC = [
        uint256(21593111544282038810514267833173034556003087472432926796912939240679730134100),
        uint256(13540084426283367957829309982316975708986350905403842008022672720178363341004)
    ];

    // Proof D (dept 0x3333...33 only), doctor 0x...0001
    bytes32 internal deptCommitmentPeds =
        bytes32(uint256(0x2c62d22fbbf65258c48af1e41eb39339e7a27711ca8f19de483abc82fd71ab97));
    uint256[2] internal proofD_pA = [
        uint256(7730955266887831200950736679334143071476485855399079395903497989159601236703),
        uint256(9948278528525427548122746797439900866073888357197004528829077187395687702866)
    ];
    uint256[2][2] internal proofD_pB = [
        [uint256(21109619852526186723089926964747711449179122576016623217774971028966951026282),
         uint256(6156122135177029853732106968309639384834012327232292957640946010397147098665)],
        [uint256(489586211021651886587755900942100523467277609089685676654440807633619098525),
         uint256(9272348004440275829081991018744378084404748596692995242235792014361830006545)]
    ];
    uint256[2] internal proofD_pC = [
        uint256(4878430349959876478764422855201135931082386435261099527228671632106891025050),
        uint256(1444114688938617572666625552207473469344668851710644758190818213335678028287)
    ];

    uint48 internal constant EVALUATED_AT = 1750000000;

    event PolicyUpdated(bytes32 requiredDeptCommitment, bytes32 requiredInstCommitment);

    Record internal record;
    History internal history;

    function setUp() public {
        policyVerifier = new PolicyVerifier();

        userRegistry = new UserRegistry();
        historyFactory = new HistoryFactory(address(userRegistry));
        historyFactory.setPolicyVerifier(address(policyVerifier));
        registry = new WellRegistry(address(userRegistry), address(historyFactory));
        historyFactory.setWellRegistry(address(registry));

        vm.prank(patient);
        address historyAddress = registry.register(
            patientCommitment,
            patientInternalID,
            Types.Roles.Patient
        );
        history = History(historyAddress);

        Types.RecordMetadata memory metadata = Types.RecordMetadata({
            recordType: Types.RecordType.LabResult,
            createdAt: 0,
            updatedAt: 0
        });

        vm.prank(patient);
        record = registry.createRecord(
            cardiologist,
            metadata,
            keccak256("content-hash"),
            "ipfs://cid"
        );
    }

    function testPoseidonCommitmentsMatchCircuit() public {
        assertEq(
            PoseidonT3.hash([uint256(0x1111111111111111111111111111111111111111111111111111111111111111), 1]),
            uint256(deptCommitmentCardio),
            "circuit deptCommitment != on-chain PoseidonT3"
        );
        assertEq(
            PoseidonT3.hash([uint256(0x2222222222222222222222222222222222222222222222222222222222222222), 2]),
            uint256(instCommitmentGeneral),
            "circuit instCommitment != on-chain PoseidonT3"
        );
        assertEq(
            PoseidonT3.hash([uint256(uint160(cardiologist)), 3]),
            uint256(doctorCommitmentCardio),
            "circuit doctorCommitment != on-chain PoseidonT3"
        );
    }

    function testSelfGrantWithValidProof() public {
        vm.prank(patient);
        record.setPolicy(deptCommitmentCardio, instCommitmentGeneral);

        vm.warp(EVALUATED_AT);
        vm.prank(cardiologist);
        record.grantAccessWithProof(proofA_pA, proofA_pB, proofA_pC, deptCommitmentCardio, instCommitmentGeneral, doctorCommitmentCardio, EVALUATED_AT);

        assertTrue(record.permissions().hasAccess(cardiologist));
        assertFalse(record.permissions().hasAccess(otherDoctor));
    }

    function testSelfGrantWithWildcardPolicy() public {
        vm.prank(patient);
        record.setPolicy(bytes32(0), bytes32(0));

        // Proof C was generated with requiredDept=requiredInst=0, so its public
        // signals are the commitments of those values (not bytes32(0)).
        bytes32 deptC = bytes32(uint256(0x1bd20834f5de9830c643778a2e88a3a1363c8b9ac083d36d75bf87c49953e65e));
        bytes32 instC = bytes32(uint256(0x065e2c6cc08a36c4a943286bc91c216054a1981eb4f7570f67394ef8937a21b8));

        vm.warp(EVALUATED_AT);
        vm.prank(cardiologist);
        record.grantAccessWithProof(proofC_pA, proofC_pB, proofC_pC, deptC, instC, doctorCommitmentCardio, EVALUATED_AT);

        assertTrue(record.permissions().hasAccess(cardiologist));
    }

    function testReplayByDifferentDoctorReverts() public {
        vm.prank(patient);
        record.setPolicy(deptCommitmentCardio, instCommitmentGeneral);

        vm.warp(EVALUATED_AT);
        vm.prank(otherDoctor);
        vm.expectRevert(Record.PolicyMismatch.selector);
        record.grantAccessWithProof(proofA_pA, proofA_pB, proofA_pC, deptCommitmentCardio, instCommitmentGeneral, doctorCommitmentCardio, EVALUATED_AT);
    }

    function testWrongDepartmentCommitmentReverts() public {
        vm.prank(patient);
        record.setPolicy(deptCommitmentPeds, instCommitmentGeneral);

        vm.warp(EVALUATED_AT);
        vm.prank(cardiologist);
        vm.expectRevert(Record.PolicyMismatch.selector);
        record.grantAccessWithProof(proofA_pA, proofA_pB, proofA_pC, deptCommitmentCardio, instCommitmentGeneral, doctorCommitmentCardio, EVALUATED_AT);
    }

    function testStaleProofReverts() public {
        vm.prank(patient);
        record.setPolicy(deptCommitmentCardio, instCommitmentGeneral);

        vm.warp(EVALUATED_AT + record.POLICY_PROOF_TTL() + 1);
        vm.prank(cardiologist);
        vm.expectRevert(Record.StaleProof.selector);
        record.grantAccessWithProof(proofA_pA, proofA_pB, proofA_pC, deptCommitmentCardio, instCommitmentGeneral, doctorCommitmentCardio, EVALUATED_AT);
    }

    function testProofForFutureEvaluatedAtReverts() public {
        vm.prank(patient);
        record.setPolicy(deptCommitmentCardio, instCommitmentGeneral);

        vm.warp(EVALUATED_AT - 1);
        vm.prank(cardiologist);
        vm.expectRevert(Record.StaleProof.selector);
        record.grantAccessWithProof(proofA_pA, proofA_pB, proofA_pC, deptCommitmentCardio, instCommitmentGeneral, doctorCommitmentCardio, EVALUATED_AT);
    }

    function testTamperedProofReverts() public {
        vm.prank(patient);
        record.setPolicy(deptCommitmentCardio, instCommitmentGeneral);

        vm.warp(EVALUATED_AT);
        uint256[2] memory tamperedA = proofA_pA;
        tamperedA[0] = tamperedA[0] ^ 1;
        vm.prank(cardiologist);
        vm.expectRevert(Record.InvalidProof.selector);
        record.grantAccessWithProof(tamperedA, proofA_pB, proofA_pC, deptCommitmentCardio, instCommitmentGeneral, doctorCommitmentCardio, EVALUATED_AT);
    }

    function testWildcardPolicyWithForcedCommitmentMismatchReverts() public {
        // Wildcard policy (0,0) accepts any proof, but a proof bound to a
        // different doctor (proof D is doctor 0x...0001, replay by 0x...0002)
        // must still revert on the doctor-binding check.
        vm.prank(patient);
        record.setPolicy(bytes32(0), bytes32(0));

        vm.warp(EVALUATED_AT);
        vm.prank(otherDoctor);
        vm.expectRevert(Record.PolicyMismatch.selector);
        record.grantAccessWithProof(proofD_pA, proofD_pB, proofD_pC, bytes32(0), bytes32(0), doctorCommitmentCardio, EVALUATED_AT);
    }

    function testOwnerGrantStillWorksAlongsidePolicy() public {
        vm.prank(patient);
        record.setPolicy(deptCommitmentCardio, instCommitmentGeneral);

        vm.prank(patient);
        record.grantAccess(otherDoctor);
        assertTrue(record.permissions().hasAccess(otherDoctor));
    }

    function testPolicyUpdatedEmitsEvent() public {
        vm.expectEmit();
        emit PolicyUpdated(deptCommitmentCardio, instCommitmentGeneral);
        vm.prank(patient);
        record.setPolicy(deptCommitmentCardio, instCommitmentGeneral);

        (bytes32 dept, bytes32 inst) = record.getPolicy();
        assertEq(dept, deptCommitmentCardio);
        assertEq(inst, instCommitmentGeneral);
    }

    function testGrantRevokeCycleWithProof() public {
        vm.prank(patient);
        record.setPolicy(deptCommitmentCardio, instCommitmentGeneral);

        vm.warp(EVALUATED_AT);
        vm.prank(cardiologist);
        record.grantAccessWithProof(proofA_pA, proofA_pB, proofA_pC, deptCommitmentCardio, instCommitmentGeneral, doctorCommitmentCardio, EVALUATED_AT);
        assertTrue(record.permissions().hasAccess(cardiologist));

        vm.prank(patient);
        history.revokeRecordAccess(record, cardiologist);
        assertFalse(record.permissions().hasAccess(cardiologist));
    }
}
