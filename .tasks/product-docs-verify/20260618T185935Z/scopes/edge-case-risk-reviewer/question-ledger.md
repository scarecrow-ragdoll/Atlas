# Question Ledger — edge-case-risk-reviewer

| ID | Scope | Severity | Question | Why It Matters | Source Or Report | Status |
| --- | --- | --- | --- | --- | --- | --- |
| Q-EDGE-01 | edge-case-risk | non-blocking | What happens when PIN hash is missing but PIN is enabled? | Edge case in settings validation | worker-attempt-1 | open |
| Q-EDGE-02 | edge-case-risk | non-blocking | What is the PIN brute-force protection policy? | Security risk | worker-attempt-1 | open |
| Q-EDGE-03 | edge-case-risk | non-blocking | What is the session TTL for PIN-protected sessions? | Session management boundary | worker-attempt-1 | open |
| Q-EDGE-04 | edge-case-risk | non-blocking | What happens when workout save is interrupted (network/DB)? | Data loss risk | worker-attempt-1 | open |
| Q-EDGE-05 | edge-case-risk | non-blocking | How are duplicate exercise names handled? | Data integrity | worker-attempt-1 | open |
| Q-EDGE-06 | edge-case-risk | non-blocking | How is concurrent edit handled (two tabs)? | Data consistency risk | worker-attempt-1 | open |
| Q-EDGE-07 | edge-case-risk | non-blocking | What is the import behavior when data already exists? | Silent data loss risk | worker-attempt-1 | open |
| Q-EDGE-08 | edge-case-risk | non-blocking | Are EXIF metadata stripped from exported photos? | Privacy risk | worker-attempt-1 | open |
| Q-EDGE-09 | edge-case-risk | non-blocking | What is the max file size for media uploads? | Resource boundary | worker-attempt-1 | open |
| Q-EDGE-10 | edge-case-risk | non-blocking | What is the backup export timeout for large datasets? | Export reliability | worker-attempt-1 | open |
| Q-EDGE-11 | edge-case-risk | non-blocking | How is media cleanup handled on exercise deletion? | Orphaned data risk | worker-attempt-1 | open |
| Q-EDGE-12 | edge-case-risk | non-blocking | What is the migration strategy for schema version changes? | Long-term data integrity | worker-attempt-1 | open |