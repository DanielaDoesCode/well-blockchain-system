"""
Locust load-test definitions for the Well EHR API.

Models a realistic mixed read workload against the REST backend:

  DoctorUser   (weight 3): identity resolve + record-key retrieve + download
  PatientUser  (weight 2): download own record (+ recovery-key read, + optional upload)

Callers are authenticated with the seeded X-User-Address / X-Session-Token pair
taken from the seed CSV (see seed_locust.js). Sessions are pre-provisioned
on-chain (Authenticated events the listener persists into the backend's
in-memory store), so Locust itself never mints sessions — that keeps the load
test focused on the API paths, not on session bootstrapping.

Usage:
  LOCUST_HOST=http://localhost:8081 LOCUST_SEED_CSV=... locust --headless \\
    -u 30 --spawn-rate 5 --run-time 5m --html report.html --csv results

Env knobs:
  LOCUST_HOST          default http://localhost:8081  (point at the API relay
                       http://localhost:18443 for latency-injected runs)
  LOCUST_SEED_CSV      default scripts/benchmarks/results/locust/seed.csv
  LOCUST_MIN_WAIT      think-time lower bound   (default 0.5s)
  LOCUST_MAX_WAIT      think-time upper bound   (default 3.0s)
  LOCUST_WRITE_TASKS   1 adds the /records/upload task (default 0)
  LOCUST_PRE_TASKS     1 adds the proxy re-encryption tasks: a doctor
                       triggering the backend PRE transform
                       (POST /records/{addr}/reencrypt) and a patient
                       re-delegating (POST /records/{addr}/delegate).
                       Default 0, so the read-only control stays comparable
                       with the pre-PRE baseline. The read tasks themselves are
                       PRE-aware in both modes: record.keyRetrieve returns the
                       re-encrypted lockbox C_doctor.
  LOCUST_SHAPE=soak    use the staged soak LoadTestShape instead of --users
"""

import csv
import os
import random

from locust import HttpUser, LoadTestShape, between, task

HOST = os.environ.get("LOCUST_HOST", "http://localhost:8081")
SEED_CSV = os.environ.get(
    "LOCUST_SEED_CSV",
    os.path.join(os.path.dirname(os.path.abspath(__file__)), "..", "results", "locust", "seed.csv"),
)
WRITE_TASKS = os.environ.get("LOCUST_WRITE_TASKS", "0") == "1"
PRE_TASKS = os.environ.get("LOCUST_PRE_TASKS", "0") == "1"

STAFF_ROLES = {"doctor", "first_responder", "hospital"}
PATIENT_ROLES = {"patient"}


def _load_rows():
    if not os.path.exists(SEED_CSV):
        raise SystemExit(
            f"seed CSV missing: {SEED_CSV}\n"
            "Run the seeder first:\n"
            "  bash scripts/benchmarks/locust/run_locust_load.sh\n"
            "  (or manually: source ./.env && node scripts/benchmarks/locust/seed_locust.js)\n"
        )
    with open(SEED_CSV, newline="") as fh:
        return list(csv.DictReader(fh))


_ROWS = _load_rows()


def _records(row):
    return [a for a in (row.get("records") or "").split("|") if a]


def _nid_pool():
    return [r["national_id"] for r in _ROWS if r.get("national_id")]


class _AuthUser(HttpUser):
    abstract = True
    host = HOST
    wait_time = between(
        float(os.environ.get("LOCUST_MIN_WAIT", "0.5")),
        float(os.environ.get("LOCUST_MAX_WAIT", "3.0")),
    )

    def _rows_for_role(self):
        allowed = STAFF_ROLES if self._is_staff else PATIENT_ROLES
        return [r for r in _ROWS if r["role"] in allowed]

    def _pick_row(self):
        rows = self._rows_for_role()
        if not rows:
            raise RuntimeError(f"no seed rows for role {self._is_staff}")
        # cheap stable distribution across users
        return rows[random.randrange(len(rows))]

    def on_start(self):
        self.row = self._pick_row()
        self.records = _records(self.row)
        self.client.headers.update(
            {
                "X-User-Address": self.row["address"],
                "X-Session-Token": self.row["token"],
            }
        )


class DoctorUser(_AuthUser):
    """Staff reading the identity index, granted record keys, and record blobs."""

    _is_staff = True
    weight = 3

    @task(3)
    def resolve_identity(self):
        nid = random.choice(_nid_pool())
        with self.client.get(
            f"/api/v1/identities?national_id={nid}", name="identity.resolve", catch_response=True
        ) as resp:
            if resp.status_code < 200 or resp.status_code >= 300:
                resp.failure(f"resolve {nid}: HTTP {resp.status_code}")

    @task(3)
    def read_record_key(self):
        if not self.records:
            return
        addr = random.choice(self.records)
        with self.client.get(
            f"/api/v1/records/{addr}/key?doctor_address={self.row['key_doctor']}",
            name="record.keyRetrieve",
            catch_response=True,
        ) as resp:
            if resp.status_code < 200 or resp.status_code >= 300:
                resp.failure(f"key {addr}: HTTP {resp.status_code}")

    @task(1)
    def download_record(self):
        if not self.records:
            return
        addr = random.choice(self.records)
        with self.client.get(
            f"/api/v1/records/{addr}/download", name="record.download", catch_response=True
        ) as resp:
            if resp.status_code < 200 or resp.status_code >= 300:
                resp.failure(f"download {addr}: HTTP {resp.status_code}")

    @task(1)
    def reencrypt_key(self):
        """Server-side proxy re-encryption: the backend runs the AFGH
        transform (lockbox + rk -> C_doctor) and upserts the result. Only
        counted when LOCUST_PRE_TASKS=1, so the control condition keeps its
        read-only profile."""
        if not PRE_TASKS or not self.records:
            return
        addr = random.choice(self.records)
        with self.client.post(
            f"/api/v1/records/{addr}/reencrypt",
            json={"doctor_address": self.row["key_doctor"] or self.row["address"]},
            name="record.reencrypt",
            catch_response=True,
        ) as resp:
            if resp.status_code < 200 or resp.status_code >= 300:
                resp.failure(f"reencrypt {addr}: HTTP {resp.status_code}")


class PatientUser(_AuthUser):
    """Patients pulling their own records (and, optionally, the escrowed keys)."""

    _is_staff = False
    weight = 2

    @task(3)
    def download_record(self):
        if not self.records:
            return
        addr = random.choice(self.records)
        with self.client.get(
            f"/api/v1/records/{addr}/download", name="record.download", catch_response=True
        ) as resp:
            if resp.status_code < 200 or resp.status_code >= 300:
                resp.failure(f"download {addr}: HTTP {resp.status_code}")

    @task(1)
    def read_recovery_key(self):
        if self.row.get("has_recovery") != "1" or not self.row.get("summary"):
            return
        with self.client.get(
            f"/api/v1/emergency/recovery-key?emergency_summary={self.row['summary']}",
            name="emergency.readRecoveryKey",
            catch_response=True,
        ) as resp:
            if resp.status_code < 200 or resp.status_code >= 300:
                resp.failure(f"recovery-key: HTTP {resp.status_code}")

    @task(1)
    def upload_record(self):
        if not WRITE_TASKS:
            return
        payload = b"LOCUST-EHR-" + random.randbytes(256)
        with self.client.post(
            "/api/v1/records/upload",
            files={"file": ("ehr.bin", payload, "application/octet-stream")},
            name="record.upload",
            catch_response=True,
        ) as resp:
            if resp.status_code < 200 or resp.status_code >= 300:
                resp.failure(f"upload: HTTP {resp.status_code}")

    @task(1)
    def delegate_key(self):
        """Re-delegate a patient's re-encryption key to their doctor.
        The seed CSV carries one pre-derived rk per record (`re_encryption_keys`,
        parallel to `records`), so the request always pairs a record with the rk
        derived for *its* lockbox — the rk is only meaningful for that pairing.
        This measures the endpoint (G2 parse + ownership check + upsert), not
        client-side crypto; client PRE cost is measured by bench_pre.js."""
        if not PRE_TASKS or not self.records:
            return
        rks = (self.row.get("re_encryption_keys") or "").split("|")
        if not self.row.get("key_doctor") or not rks or not rks[0]:
            return
        i = random.randrange(len(self.records))
        if i >= len(rks):
            i = len(rks) - 1
        with self.client.post(
            f"/api/v1/records/{self.records[i]}/delegate",
            json={
                "patient_address": self.row["address"],
                "doctor_address": self.row["key_doctor"],
                "re_encryption_key": rks[i],
            },
            name="record.delegate",
            catch_response=True,
        ) as resp:
            if resp.status_code < 200 or resp.status_code >= 300:
                resp.failure(f"delegate: HTTP {resp.status_code}")


if os.environ.get("LOCUST_SHAPE") == "soak":

    class StagedSoak(LoadTestShape):
        """Optional staged soak used when LOCUST_SHAPE=soak (ignores -u/--spawn-rate).

        Ramps to target in 60s, holds 4 minutes, then ramps back down.
        Set LOCUST_HOST/LOCUST_SEED_CSV and run with --run-time 6m --headless.
        """

        # pylint: disable=no-self-use
        def tick(self):
            run_time = self.get_run_time()
            target = int(os.environ.get("LOCUST_USERS", "30"))
            if run_time < 60:
                return target * (run_time / 60), int(target * (run_time / 60))
            if run_time < 300:
                return target, target // 5
            if run_time < 325:
                return max(1, int(target * (325 - run_time) / 25)), target // 5
            return None