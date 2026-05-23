# tinykv Runbook

## Leader Election Storms
**Symptom:** rapid term increments in metrics; client errors spike.
**Cause:** network partition or overloaded leader.
**Action:** check `raft_term`, `raft_leader_changes_total`. Drain noisy node. Verify clock skew <100ms across nodes.

## Disk Full
**Symptom:** writes fail with `no space left on device`.
**Cause:** log not compacting, large snapshots.
**Action:** confirm snapshot threshold; trigger manual snapshot; rotate disk.

## Network Partition
**Symptom:** minority partition reports `ErrNoLeader`.
**Cause:** expected — minority can't elect.
**Action:** restore network; verify majority continued serving; check linearizability test results.

## Cert Rotation
**Symptom:** TLS handshake failures after rotation.
**Cause:** missed reload.
**Action:** SIGHUP triggers reload; verify `tls_cert_expiry_seconds` metric.
