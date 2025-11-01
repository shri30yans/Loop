Document 2 of 4: The Storage Engine (synapse-engine)
Project: Synapse
Component: synapse-engine
Owner: Technical Lead
Version: 1.0
Status: DEFINED
1. Overview
This document specifies the design and implementation of the Synapse storage engine. This engine is responsible for the durable, persistent storage of the hypergraph on disk. It will be built on top of the RocksDB embedded database and provide atomic, transactional APIs for writing and reading graph data. It will manage all the custom indexes required for high-performance traversal.
2. Dependencies
The synapse-engine/Cargo.toml file must contain:
code
Toml
[package]
name = "synapse-engine"
version = "0.1.0"
edition = "2021"

[dependencies]
# The core data structures we will be storing.
synapse-core = { path = "../synapse-core" }
# The high-performance embedded database.
rocksdb = { version = "0.22", features = ["multi-threaded-cf"] }
# For serializing/deserializing data structures.
bincode = "1.3"
# For unique IDs.
uuid = "1.8"
# For error handling.
thiserror = "1.0"
3. Engine Design & Implementation (synapse-engine/src/lib.rs)
3.1. Constants and Error Handling
Purpose: Define constants for our Column Family names to avoid magic strings. Define a custom error type for the engine for clean error handling.
code
Rust
use thiserror::Error;

// Define the names of our "tables" within RocksDB.
const CF_ATOMS: &str = "atoms";
const CF_TOPOLOGY: &str = "topology";
const CF_LABEL_TO_ID: &str = "label_to_id";

#[derive(Error, Debug)]
pub enum EngineError {
    #[error("RocksDB error: {0}")]
    DbError(#[from] rocksdb::Error),
    #[error("Serialization error: {0}")]
    SerializationError(#[from] Box<bincode::ErrorKind>),
    #[error("Column Family '{0}' not found")]
    CfNotFound(String),
    // ... other error types
}
3.2. The Engine Struct
Purpose: The main struct that holds the database connection and provides the public API for all storage operations.
code
Rust
use rocksdb::{DB, Options, ColumnFamilyDescriptor, WriteBatch};
use synapse_core::{AtomId, AtomContainer, Atom, HyperedgeData};
use std::sync::Arc;

pub struct Engine {
    db: Arc<DB>,
}
3.3. Initialization (Engine::open)
Purpose: To open an existing RocksDB database or create a new one if it doesn't exist. This function is responsible for defining and opening all the necessary Column Families.
code
Rust
impl Engine {
    pub fn open(path: &str) -> Result<Self, EngineError> {
        let mut db_opts = Options::default();
        db_opts.create_if_missing(true);
        db_opts.create_missing_column_families(true);

        let cf_names = vec![CF_ATOMS, CF_TOPOLOGY, CF_LABEL_TO_ID];
        let cf_descriptors: Vec<ColumnFamilyDescriptor> = cf_names
            .iter()
            .map(|name| ColumnFamilyDescriptor::new(*name, Options::default()))
            .collect();

        let db = DB::open_cf_descriptors(&db_opts, path, cf_descriptors)?;
        Ok(Engine { db: Arc::new(db) })
    }
}
3.4. The commit Method
Purpose: This is the primary write method. It must be atomic. It takes an AtomContainer, saves it to the main atoms table, and updates all relevant indexes (topology and label_to_id) in a single, all-or-nothing transaction using a WriteBatch.
code
Rust
// Inside impl Engine
pub fn commit(&self, atom_container: &AtomContainer) -> Result<(), EngineError> {
    let mut batch = WriteBatch::default();

    // Get handles to all required column families.
    let cf_atoms = self.db.cf_handle(CF_ATOMS).ok_or_else(|| EngineError::CfNotFound(CF_ATOMS.to_string()))?;
    let cf_topology = self.db.cf_handle(CF_TOPOLOGY).ok_or_else(|| EngineError::CfNotFound(CF_TOPOLOGY.to_string()))?;
    let cf_label_to_id = self.db.cf_handle(CF_LABEL_TO_ID).ok_or_else(|| EngineError::CfNotFound(CF_LABEL_TO_ID.to_string()))?;

    // 1. Add the main atom data to the batch.
    let atom_bytes = bincode::serialize(atom_container)?;
    batch.put_cf(&cf_atoms, atom_container.id.as_bytes(), &atom_bytes);

    // 2. Update indexes based on the atom type.
    match &atom_container.payload {
        Atom::Entity(entity_data) => {
            // Update the label-to-ID index.
            batch.put_cf(&cf_label_to_id, entity_data.primary_label.as_bytes(), atom_container.id.as_bytes());
        }
        Atom::Hyperedge(hyperedge_data) => {
            // Update the topology index for every connected entity.
            for connection in &hyperedge_data.connections {
                let entity_id = &connection.entity_id;
                let current_connections_bytes = self.db.get_cf(&cf_topology, entity_id.as_bytes())?;
                let mut connections_vec: Vec<AtomId> = match current_connections_bytes {
                    Some(bytes) => bincode::deserialize(&bytes)?,
                    None => Vec::new(),
                };
                connections_vec.push(atom_container.id);
                let updated_connections_bytes = bincode::serialize(&connections_vec)?;
                batch.put_cf(&cf_topology, entity_id.as_bytes(), &updated_connections_bytes);
            }
        }
    }

    // 3. Execute the batch write. This is an atomic operation.
    self.db.write(batch)?;
    Ok(())
}
3.5. Read Methods (fetch_by_id, fetch_by_label, get_connections_for)
Purpose: Provide the necessary APIs to retrieve data and traverse the graph.
code
Rust
// Inside impl Engine
pub fn fetch_by_id(&self, id: &AtomId) -> Result<Option<AtomContainer>, EngineError> {
    // ... implementation to get from cf_atoms and deserialize ...
}

pub fn fetch_by_label(&self, label: &str) -> Result<Option<AtomContainer>, EngineError> {
    // ... implementation to get ID from cf_label_to_id, then call fetch_by_id ...
}

pub fn get_connections_for(&self, id: &AtomId) -> Result<Vec<AtomId>, EngineError> {
    // ... implementation to get and deserialize the connection vector from cf_topology ...
}
4. Acceptance Criteria
The synapse-engine crate must compile.
Unit tests must cover:
Opening and creating a database.
Committing a single Entity and fetching it back by both ID and label.
Committing a Hyperedge that connects two Entities.
Calling get_connections_for on one of the entities and verifying that the ID of the hyperedge is returned.
Error handling for non-existent atoms.
