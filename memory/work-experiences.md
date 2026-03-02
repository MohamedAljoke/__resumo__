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

## A1 Certificate PDF Signing and Validation Flow

**Problem:**  
Needed to support secure PDF signing with A1 digital certificates while meeting Brazilian ITI/ICP-Brasil validation requirements and protecting sensitive certificate material.

**Investigation:**

- Studied the A1 certificate signing flow using `.pfx` files
- Reviewed how to validate certificate chains and document authenticity through AIA-based verification
- Evaluated secure storage requirements for private certificate files at rest
- Verified how RSA encryption and certificate handling fit into the PDF signing process
- Assessed memory risk from concurrent PDF generation and considered scaling alternatives such as asynchronous processing or offloading generation to Lambda

**Resolution:**

- Implemented encryption at rest for `.pfx` certificate files
- Added document and certificate validation using AIA verification
- Used RSA-based cryptography in the signing flow
- Built PDF signing using A1 certificates in a valid Brazilian ITI/ICP-Brasil-compliant pattern
- Added semaphore-based concurrency control around PDF generation, allowing 4 active jobs and 16 queued jobs before returning `503 Service Unavailable`
- Chose the semaphore approach as a simple, effective safeguard for the existing traffic level while keeping future options open for asynchronous processing or dedicated Lambda-based generation

**Impact:**  
Enabled secure, standards-compliant PDF signing while improving protection of sensitive certificate assets, strengthening validation of signed documents, and preventing PDF-generation spikes from exhausting instance memory.

**Skills:**  
Digital signatures, A1 certificates, PDF signing, RSA cryptography, certificate-chain validation, AIA verification, encryption at rest, `.pfx` handling, ICP-Brasil compliance, concurrency control, backpressure design, memory protection, capacity planning

---

## Go Internal Proxy Panic from Reused `fasthttp` Headers

**Problem:**  
An internal proxy server written in Go began panicking under high concurrency while forwarding requests through `fasthttp`.

**Investigation:**

- Investigated failures that only appeared under heavier concurrent traffic
- Traced the issue to request-header handling inside the proxy path
- Identified that `fasthttp` optimizes allocations by reusing internal buffers rather than duplicating header data automatically
- Built a focused test case to reliably reproduce the panic and confirm the failure mode

**Root Cause:**  
The proxy flow retained or reused header values without making an explicit copy. Under concurrency, the underlying `fasthttp` buffers could be reused and mutated, creating unsafe shared-memory behavior and causing panics.

**Resolution:**

- Fixed the header handling by reallocating/copying the data into independent memory when it needed to outlive the current request context
- Adjusted the implementation to respect `fasthttp`'s zero-allocation design and buffer-lifetime expectations

**Impact:**  
Eliminated concurrency-related panics in the proxy path, added regression coverage for the failure mode, and improved understanding of memory ownership, request lifetimes, and zero-copy tradeoffs in high-performance Go HTTP libraries.

**Skills:**  
Go, `fasthttp`, concurrency debugging, memory ownership, zero-copy APIs, proxy servers, high-throughput systems, panic diagnosis, regression testing

---

## Bulk Compliance Reprocessing with Observability and Backpressure Monitoring

**Problem:**  
Some accounts had not gone through required compliance checks, creating a need to reprocess them reliably at scale without overwhelming the existing system.

**Investigation:**

- Reviewed the existing compliance flow to reuse proven application behavior instead of creating a separate one-off process
- Identified how to trigger many compliance checks in bulk while preserving normal service contracts
- Monitored Grafana dashboards during execution to track queue backpressure, microservice memory usage, and overall system stability
- Used the higher-load execution path to expose a latent `fasthttp` header-memory issue in an internal proxy service

**Resolution:**

- Built a script to re-run compliance checks for the affected accounts using the already-existing business flow
- Designed the script so many accounts could be processed in one run while still observing system limits
- Actively monitored queue depth and microservice memory consumption in Grafana during execution
- Investigated and fixed the `fasthttp` concurrency bug uncovered during the bulk processing work

**Impact:**  
Completed compliance reprocessing for affected accounts without creating a parallel workflow, while validating the system under load and uncovering a production-risk concurrency issue before it caused broader failures.

**Skills:**  
Operational automation, bulk processing, compliance workflows, observability, Grafana, queue backpressure, microservice monitoring, load-aware execution, incident discovery, system reliability

---

## Template for Future Entries

**Problem:** [Brief description of what went wrong]

**Investigation:** [How you approached the problem]

**Root Cause:** [What you discovered]

**Resolution:** [How you fixed it]

**Impact:** [What was the result/learning]

**Skills:** [Key technical and soft skills demonstrated]

---
