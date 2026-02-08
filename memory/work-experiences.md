# Work Experiences & Problem Solving

This document tracks significant issues, bugs, and problems encountered at work that demonstrate problem-solving skills and technical expertise. Use these for future interviews and reference.

---

## Browser Extension Conflict Bug

**Problem:** Single user experiencing issues while all others worked fine.

**Investigation:** After thorough investigation, discovered the issue was caused by a browser extension interfering with our styling package.

**Resolution:** Identified the specific extension and its interaction with the styling library.

**Skills:** Debugging, user-specific issue isolation, understanding browser extension behavior

---

## False DDoS Alert - Frontend Performance Issue

**Problem:** Company suspected a DDoS attack due to unusual traffic patterns. Another developer had been investigating for days without finding the cause.

**Investigation:** Took over the investigation a few days after the initial developer.

**Root Cause:** Discovered a missing debounce on a frontend slider component, causing excessive API requests.

**Resolution:** Implemented proper debouncing on the slider.

**Skills:** Performance debugging, frontend optimization, API traffic analysis, persistence in problem-solving

---

## AWS Account Compromise - Leaked Access Keys

**Problem:** Old AWS IAM user access keys were leaked and exploited by a bot to send mass SES (Simple Email Service) messages. This resulted in account-level restrictions on Lambda usage.

**Investigation:**

- Used AWS CloudTrail to trace unauthorized activity
- Identified the compromised user and attack vector

**Resolution:**

- Immediately deactivated the compromised IAM user
- Cleaned up other unused IAM users and reviewed account security
- Activated MFA (Multi-Factor Authentication) across the account
- Stayed up late coordinating with AWS Support to lift account restrictions
- Successfully restored full account access

**Impact:** Prevented further unauthorized usage, restored critical Lambda functionality, learned AWS security best practices.

**Skills:** AWS security, CloudTrail log analysis, IAM management, incident response, vendor communication, working under pressure

---

## V1 to V2 Migration - Production Data Incident

**Problem:** During a major system migration from V1 to V2, accidentally deleted more production data than intended while cleaning the database. The migration involved continuous data migration and cleaning scripts over 2 weeks.

**Context:**

- Migration required ongoing data scripts to move and clean data in the new database
- High pressure to complete tasks quickly
- Direct manipulation of production database became the norm

**Investigation:** Immediately recognized the scope of the deletion and assessed what data was affected.

**Resolution:**

- Adapted and repurposed existing migration scripts to recover the deleted data
- Restored all affected data within 2 hours
- Successfully recovered from the incident

**Lessons Learned:**

- **Never treat production databases casually** - even under time pressure
- Importance of proper migration planning and staging environments
- Always have rollback plans and backups before production changes
- Test deletion scripts thoroughly in non-production environments first
- Time pressure should never compromise safety protocols
- Keep detailed logs of all database operations

**Impact:** Learned critical lessons about database safety, disaster recovery, and the importance of proper planning over speed. Developed better practices for handling production data.

**Skills:** Database management, data recovery, script adaptation, working under pressure, incident recovery, learning from mistakes, PostgreSQL/MySQL operations

## Full Entity Update Race Condition

**Problem:**  
Intermittent data inconsistencies occurred when the same entity was updated multiple times in quick succession, causing fields to be unexpectedly overwritten with stale data.

**Investigation:**

- Analyzed application logs and database timestamps to identify overlapping update operations
- Reproduced the issue by triggering concurrent update requests from different flows
- Inspected request payloads and noticed full entity data being sent with outdated values

**Root Cause:**  
Race condition caused by **full entity updates** where concurrent requests were based on stale reads and the last write overwrote newer data.

**Resolution:**

- Refactored update logic to use partial updates instead of full entity replacements
- COULD HAVE Added optimistic locking/version checks to prevent stale writes

**Impact:**  
Eliminated silent data loss, improved data consistency under concurrency, and increased system reliability in high-traffic scenarios.

**Skills:**  
Concurrency debugging, race condition handling, API design, data consistency, backend architecture, defensive programming

## Node.js Temporary File Memory Leak (Unflushed Writes)

**Problem:**  
An EC2 instance began running out of memory and becoming unstable over time. The issue escalated under load, eventually requiring instance restarts.

**Investigation:**

- Monitored EC2 memory usage and noticed a steady, non-recovering increase
- Analyzed Node.js process memory and system-level disk usage
- Reviewed recent changes involving file generation and temporary file handling
- Identified heavy use of in-memory buffers and temporary files written during request processing

**Root Cause:**  
Node.js was **writing files in memory but not properly flushing or closing streams**, causing temporary files to accumulate. This led to memory pressure and unbounded growth over time, especially under concurrent workloads.

**Resolution:**

- Refactored file handling to use proper stream-based writes instead of large in-memory buffers
- Ensured all write streams were explicitly closed and flushed
- Implemented cleanup logic for temporary files after use
- Added monitoring and alerts for memory usage and disk growth on the EC2 instance

**Impact:**  
Stopped the memory leak, stabilized the EC2 instance, and significantly improved system reliability under load. Gained deeper understanding of Node.js memory management, stream lifecycle, and the importance of resource cleanup in long-running services.

**Skills:**  
Node.js internals, memory leak debugging, stream handling, temporary file management, EC2 monitoring, performance optimization, production incident diagnosis

---

## Template for Future Entries

**Problem:** [Brief description of what went wrong]

**Investigation:** [How you approached the problem]

**Root Cause:** [What you discovered]

**Resolution:** [How you fixed it]

**Impact:** [What was the result/learning]

**Skills:** [Key technical and soft skills demonstrated]

---
