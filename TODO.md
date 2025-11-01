We need AI for simple operations -> Ditch seperate backend, Add Go Code to talk to AI

## AI Operations
1. Decide when to do what -> Decide what information to retreive
2. Summarization of results -> Decide which information is most relevant from results

### Best retrieval? 
- How? 
- Interactions of the user with the website
    - How long user was on the post?
    - How long/did he read comments?
    - Upvote or Downvote
    - Click Analysis
    - ????

Based on all this, change retrieval mechanism.


Immediate next steps
1. Connect with AuraDB
    1. Create new repo for GraphDB 
    2. Link with service layer
    3. CRUD operations 
2. Create Query Page on Website
3. Figure out how to convert Natural Language to GraphQuery
    - Neo4j MCP Server??
    - Tools
    - Code to do fixed operations
4. How do recommendation systems work? 
    - How to change weights in the graph? 
5. How to measure interactions with website? 
6. How to decide how they affect the RAG?
7. How to make this actually affect the RAG?


- We are generating tags from AI? 
- tags that we get are different from the tags we have
- solution 
    - send LLm existing tags
        - Update the tags constantly- > issue
    - 

text -> Model -> numbers
numbers ???
- Change Model, Find better model? 


Embeddings? 
Algorithm 
Additive Manufactuting + 3D printing




# What are the next steps?
- How to extract entities? 
    - Nouns
    - LLM extraction
- How to merge Tags
    - Graph Techniques
        - Similiar neighbours
    - Fixed lists
    - Embeddings
- Retrieval business logic
    - When to retrieve what and how? 
    - tools 
    - Database + Graph

Store + Retrieve 
