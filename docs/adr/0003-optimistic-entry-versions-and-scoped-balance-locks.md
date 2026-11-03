# Optimistic entry versions and scoped balance locks

Use optimistic version checks for expense and repayment mutations and database row locks for affected balances scoped to workspace and member. Keep entry changes, balance effects, revisions, and durable request-key results in one PostgreSQL transaction. The user chose this direction on 2026-11-03 instead of the proposed exclusive workspace lock for every mutation.

This permits independent writes within a workspace while retaining conflict detection and balance integrity. It requires consistent lock ordering and coordination with membership, authorization, and currency changes; those details must be settled before financial implementation. Balance locks alone do not solve those races. A complex automatic retry framework is not required initially.
