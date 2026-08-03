---
document_id: decision-guides-caching-and-database
title: Caching and Database Strategy Guide
ecosystem: cross-cutting
audience: [human, agent]
last_reviewed: 2026-08-03
---

# Caching and Database Strategy Guide

This guide defines the options matrix for selecting database optimization patterns, cache eviction strategies, and asynchronous data processing systems.

---

## 1. Caching & Storage Patterns

| Target Parameter | Selected Strategy | Tool / Technology | Eviction Policy |
| :--- | :--- | :--- | :--- |
| Static metadata, site settings | **In-memory cache** | Application runtime | Evict on change event |
| Heavy database read queries, user sessions | **Distributed Cache** | Redis | Least Recently Used (LRU) / TTL |
| Financial transactions, logs, audits | **Relational ledger** | PostgreSQL / MySQL | None (immutable, append-only) |
| Asynchronous jobs, notifications | **Message Queue** | Redis queues / RabbitMQ | Evict on job acknowledgment |

---

## 2. Decision Guidelines

### Database Optimization Hierarchy
Before adding caching layers to solve latency bounds:
1. **Explain the Plan**: Run database `EXPLAIN` query plans to find missing indexes or sequential scan blockages.
2. **Optimize Schema**: Ensure correct indexing keys are applied to target query filters (refer to [Database Optimization and Indexing](../performance/02-database-optimization.md)).
3. **Refactor Queries**: Split large multi-table joins or use keyset pagination (refer to [Performance Engineering Standard](../core/10-performance-engineering-standard.md)).
4. **Cache as Fallback**: Introduce Redis caching only if database scaling boundaries are reached.

---

## 3. References
- Database Design Standards: [Database Engineering Standard](../core/06-database-engineering-standard.md)
- Queue and Cache Rules: [Caching and Background Queue Workers](../performance/03-caching-and-queues.md)
