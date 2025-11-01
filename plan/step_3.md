roject: Synapse
Component: synapse-server
Owner: Technical Lead
Version: 1.0
Status: DEFINED
1. Overview
This document specifies the main application that brings the Synapse engine to life. It will host the engine, manage the high-speed in-memory cache, run the agent tasks, and expose an API for external interaction. This is the "live" part of our system.
2. Dependencies
The synapse-server/Cargo.toml file must contain:
code
Toml
[package]
name = "synapse-server"
version = "0.1.0"
edition = "2021"

[dependencies]
synapse-core = { path = "../synapse-core" }
synapse-engine = { path = "../synapse-engine" }
tokio = { version = "1", features = ["full"] }
axum = "0.7" # A popular, ergonomic web framework.
lru = "0.12" # For our in-memory cache.
std = "1.0" # Standard library dependency
# For thread-safe sharing of state.
use std::sync::{Arc, Mutex};
3. Workspace Design (synapse-server/src/main.rs)
3.1. The Workspace Struct
Purpose: This struct is the central, shared state of our running application. It will be wrapped in an Arc and passed to all API handlers and agent tasks, giving them safe, concurrent access to the engine and the cache.
code
Rust
use synapse_engine::Engine;
use synapse_core::{AtomId, AtomContainer};
use lru::LruCache;
use std::sync::{Arc, Mutex};
use std::num::NonZeroUsize;

pub struct Workspace {
    engine: Engine,
    hot_cache: Mutex<LruCache<AtomId, Arc<AtomContainer>>>,
}

impl Workspace {
    pub fn new(engine: Engine) -> Self {
        Self {
            engine,
            hot_cache: Mutex::new(LruCache::new(NonZeroUsize::new(100_000).unwrap())),
        }
    }
}
3.2. The Core "Get Atom" Logic
Purpose: This is the implementation of our "Flashlight" or "Lazy Loading" mechanism. It intelligently combines the fast in-memory cache with the persistent on-disk engine.
code
Rust
// Inside impl Workspace
pub async fn get_atom(&self, id: &AtomId) -> Result<Option<Arc<AtomContainer>>, anyhow::Error> {
    // 1. Check the hot cache first. This is the fast path.
    if let Some(atom_arc) = self.hot_cache.lock().unwrap().get(id) {
        return Ok(Some(atom_arc.clone()));
    }

    // 2. Cache miss. Fetch from the persistent engine. This is the "slow" path.
    match self.engine.fetch_by_id(id)? {
        Some(atom_container) => {
            let atom_arc = Arc::new(atom_container);
            // 3. Put the newly fetched atom into the hot cache for next time.
            self.hot_cache.lock().unwrap().put(*id, atom_arc.clone());
            Ok(Some(atom_arc))
        }
        None => Ok(None),
    }
}
3.3. API Endpoints
Purpose: To provide an external interface for ingesting data and querying the graph. We will use the axum web framework.
code
Rust
use axum::{routing::{get, post}, Router, extract::{State, Path}, Json};

// An example API handler for fetching an atom.
async fn get_atom_handler(
    State(workspace): State<Arc<Workspace>>,
    Path(id): Path<AtomId>,
) -> Result<Json<Arc<AtomContainer>>, String> {
    match workspace.get_atom(&id).await {
        Ok(Some(atom)) => Ok(Json(atom)),
        Ok(None) => Err("Atom not found".to_string()),
        Err(e) => Err(e.to_string()),
    }
}

// ... similar handlers for ingestion and traversal ...
3.4. Main Application Entrypoint
Purpose: To initialize all components and start the server.
code
Rust
#[tokio::main]
async fn main() {
    let engine = Engine::open("./synapse_db").expect("Failed to open engine");
    let workspace = Arc::new(Workspace::new(engine));

    let app = Router::new()
        .route("/atom/:id", get(get_atom_handler))
        // ... add other routes ...
        .with_state(workspace);

    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
    println!("Listening on {}", listener.local_addr().unwrap());
    axum::serve(listener, app).await.unwrap();
}
4. Acceptance Criteria
The synapse-server binary must compile and run.
The server must start without errors, successfully opening the RocksDB engine.
Manual tests using a tool like curl or Postman must succeed for:
POST /ingest: Ingesting a new entity and a hyperedge.
GET /atom/{id}: Retrieving the ingested atoms.
GET /traverse/{id}: Traversing the connection.
Verify that on the first GET /atom/{id}, the response is slightly slower (a cache miss), and on the second identical request, the response is significantly faster (a cache hit).