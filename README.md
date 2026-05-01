# realtime_1v1

Real-time 1v1 coding/trivia game built with a Go backend and a React + Vite frontend.

Current state:
- Frontend flow is implemented but largely unstyled (UI polish is still WIP)
- Backend core gameplay flow is close to complete

## Stack

- **Backend:** Go, Gorilla Mux, Gorilla WebSocket, pgx (PostgreSQL), godotenv
- **Frontend:** React 19, React Router, Vite

## Features (Current)

- Account creation and login
- Create and join game rooms
- Pre-game lobby readiness flow via WebSocket
- Start game by selecting a niche topic
- Timed general-question rounds
- Niche/final answer submission
- End-game result broadcast and score display

## Project Structure

- `backend/` Go API server, game/room logic, websocket logic, DB layer
- `backend/data/` question/problem JSON datasets
- `frontend/` React client app and routing

## Environment Variables

### Backend (`backend/.env`)

Required:
- `POSTGRES_CONN_STRING`
- `PROBLEMS_GENERAL`
- `PROBLEMS_NICHE`

Optional:
- `OLLAMA_URL`

Example:

```env
POSTGRES_CONN_STRING=postgres://user:password@localhost:5432/realtime_1v1?sslmode=disable
PROBLEMS_GENERAL=./data/questions.json
PROBLEMS_NICHE=./data/problems.json
OLLAMA_URL=http://localhost:11434
```

### Frontend (`frontend/.env`)

Required:
- `VITE_BACKEND_URL`
- `VITE_BACKEND_WS`

Example:

```env
VITE_BACKEND_URL=http://localhost:3002
VITE_BACKEND_WS=ws://localhost:3002
```

## Database Setup

The backend currently expects a `users` table with at least:
- `username` (unique)
- `password`

Starter SQL:

```sql
CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  username TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL
);
```

## Run Locally

### 1) Start backend

```bash
cd backend
go run ./
```

Backend listens on `:3002`.

### 2) Start frontend

```bash
cd frontend
npm install
npm run dev
```

## API Routes (Current)

- `POST /newroom`
- `POST /createuser`
- `POST /login`
- `POST /addplayer?roomID=...`
- `POST /tokenforws`
- `GET /websocketconn?token=...&roomid=...`
- `POST /game/{id}?username=...`
- `POST /startgame`
- `POST /answergeneralquestion`
- `POST /answernichequestion`

## Known Gaps / Next Work

- Frontend visual design and responsive polish
- Better API validation and error consistency
- Password/token security hardening
- Automated tests (backend + frontend)
- Deployment setup

## License

No license file is currently included.
