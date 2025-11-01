ocument 4 of 4: The Agentic Layer (Future Work Specification)
Project: Synapse
Component: Agentic Layer & Collaboration
Owner: Technical Lead
Version: 0.1 (Conceptual)
Status: PROPOSED
1. Overview
This document outlines the design for the agentic layer that will operate on top of the Synapse workspace. This phase focuses on the "live intelligence" of the system.
2. Core Concepts
Agents as tokio Tasks: Each agent (Supervisor, Specialist, etc.) will be a long-running, asynchronous task spawned by the synapse-server.
Shared State: All agents will receive an Arc<Workspace> handle, giving them shared access to the cognitive environment.
In-Memory Communication (Cognitive Momentum):
The AtomContainer struct in synapse-core will be modified to include a volatile, in-memory field:
code
Rust
#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct AtomContainer {
    // ... existing fields ...
    #[serde(skip)] // This field is NOT persisted to RocksDB.
    pub volatile_tags: Arc<Mutex<HashMap<String, serde_json::Value>>>,
}
Agents will collaborate by acquiring a lock on this Mutex and adding/reading temporary tags (e.g., {"status": "under_review_by_agent_A"}). Because all agents working on a cached atom share the same Arc pointer, this communication is instantaneous.
Event-Driven Collaboration (Stigmergy):
The Workspace will contain a tokio::sync::broadcast::Sender.
The commit method in the Engine will be modified to send a CortexEvent (e.g., AtomCreated(AtomContainer)) to this channel after a successful write to RocksDB.
Agents can subscribe to this broadcast channel to be awakened by relevant, permanent state changes in the graph.
3. Implementation Plan
Define the agent traits and structs.
Implement the Supervisor agent to manage the overall query lifecycle.
Implement a simple Specialist agent that can perform a traversal.
Integrate the SLM-based Intent Engine to translate natural language into agentic tasks.
Build the in-memory tagging and event bus mechanisms.
4. Acceptance Criteria
A complex query can be submitted to the server.
A Supervisor agent can be observed (via logging) breaking the query down.
A Specialist agent can be observed traversing the graph via the Workspace API.
Agents can successfully communicate by applying and reading volatile tags on shared, in-memory Atoms.
