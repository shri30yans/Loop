# Loop
I'd love to learn about all the various projects people are building around me, yet I have no way of discovering them. Loop is a platform where builders share their projects and the journey of building them, the rationale behind their choices, the challenges they faced, and how their ideas evolved.
Users regularly share project updates to keep everyone in the "loop". Feedback on these updates would hold them accountable, encourage "learning by doing" and create a community where people can connect with like-minded individuals.

### Techstack
- Next.js + NextUI
- Go using Mux 
- Postgres DB

# Project Setup

## Getting Started

Follow these steps to set up the development environment for the project.

### 1. Clone the Repository
Clone the repository to your local machine:

```bash
git clone <repository-url>
cd <repository-directory>
```
2. Installations
- Download [Go](https://go.dev/dl/) and (Node.js)[https://nodejs.org/en]
- Download and install [PostgreSQL](https://www.postgresql.org/download/) on your local machine.

After installation, create a new PostgreSQL database:
```
CREATE DATABASE loop;
```
3. Configure Environment Variables
Create a .env file in the root directory of the project and add the following database connection details:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=Welcome1
DB_NAME=loop
```
4. . Run the Backend
Install Go modules:
```
go mod tidy
```
Run the backend server:
```
go run server.go
```
Run the frontend
```
npm run dev
```

