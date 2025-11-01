Document 1 of 4: The Core Data Model (synapse-core)
Project: Synapse
Component: synapse-core
Owner: Technical Lead
Version: 1.0
Status: DEFINED
1. Overview
This document specifies the universal, canonical data structures for the Synapse engine. These structures are the "DNA" of our system. They must be stable, serializable, and expressive enough to model a complex hypergraph of knowledge. This crate will have no logic; it is purely for data definitions.
2. Dependencies
The synapse-core/Cargo.toml file must contain the following:
code
Toml
[package]
name = "synapse-core"
version = "0.1.0"
edition = "2021"

[dependencies]
# For converting structs to and from binary format for storage.
bincode = "1.3"
# For flexible, schemaless property bags.
serde_json = "1.0"
# The foundation for serialization.
serde = { version = "1.0", features = ["derive"] }
# For universally unique identifiers.
uuid = { version = "1.8", features = ["v4", "serde"] }
3. Data Structure Definitions (synapse-core/src/lib.rs)
3.1. Type Aliases
AtomId: A type alias for uuid::Uuid. This will be the unique identifier for every single node and edge in our graph. Using a type alias makes the code more readable.
code
Rust
pub type AtomId = uuid::Uuid;
3.2. The AtomContainer Struct
Purpose: This is the top-level wrapper for every object that will be stored in our database. It contains the object's unique ID and its data payload. Every key in our primary atoms database table will be an AtomId, and every value will be a serialized AtomContainer.
code
Rust
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)]
pub struct AtomContainer {
    pub id: AtomId,
    pub payload: Atom,
}
3.3. The Atom Enum
Purpose: This enum differentiates between the two fundamental types of objects in our hypergraph: "things" (Entities) and the "events that connect them" (Hyperedges). This strong typing at the core level is critical for the engine's logic.
code
Rust
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)]
pub enum Atom {
    Entity(EntityData),
    Hyperedge(HyperedgeData),
}
3.4. EntityData Struct
Purpose: Represents a "noun" or a node in the graph. This could be a person, a project, a technology, or any other distinct entity.
code
Rust
use std::collections::HashMap;

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)]
pub struct EntityData {
    // A machine-readable category for the entity (e.g., "User", "Project").
    pub entity_type: String,
    // The primary, human-readable name for this entity.
    pub primary_label: String,
    // A flexible, schemaless key-value store for any other metadata.
    // This allows us to add custom attributes without changing the core struct.
    pub properties: HashMap<String, serde_json::Value>,
}
3.5. HyperedgeData Struct
Purpose: Represents a "verb" or a complex relationship that connects multiple entities. This is the core of the hypergraph model.
code
Rust
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)]
pub struct HyperedgeData {
    // A machine-readable category for the event (e.g., "Contribution", "Dependency").
    pub edge_type: String,
    // A list of all entities participating in this event, and the roles they play.
    pub connections: Vec<Connection>,
}
3.6. Connection Struct
Purpose: A helper struct used within HyperedgeData. It explicitly defines the link between a hyperedge and a single participating entity, assigning a specific role to that entity within the context of the event.
code
Rust
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)]
pub struct Connection {
    // The role this entity plays (e.g., "contributor", "project_worked_on").
    pub role: String,
    // The unique ID of the entity that is connected.
    pub entity_id: AtomId,
}
4. Acceptance Criteria
The synapse-core crate must compile successfully.
All defined structs and enums must derive Serialize, Deserialize, Clone, Debug, and PartialEq.
A simple test case should be written that constructs an Entity, a Hyperedge, and an AtomContainer, then successfully serializes it to bytes with bincode::serialize and deserializes it back, asserting that the original and the result are equal.