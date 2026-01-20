## Notification System (Python / FastAPI)

This repo contains a **Python 3.10+** notification system backend using **FastAPI + SQLAlchemy + PostgreSQL**.

The older Go version is still present, but the Python version is the target implementation.

### Architecture

Layered architecture:
- **handlers** → HTTP layer
- **service** → business logic
- **repository** → DB access only
- **notifier** → Factory + Strategy for channels

### Crash resistance / idempotency

- Uses **status flag**: `PENDING` → `NOTIFIED`
- Uses a **DB transaction** per subscription so we never mark `NOTIFIED` before the notification rows are stored.
- If restock is called twice, only `PENDING` subscriptions are processed (no duplicates).
- If **all channels fail**, subscription stays `PENDING`.
- `notifications.json` is written as **best-effort**; PostgreSQL is the source of truth.

### Setup

- Python 3.10+
- PostgreSQL

Install:

```bash
python -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
```

Create DB + apply schema:

```bash
createdb notifications
psql postgresql://postgres:postgres@localhost/notifications -f notification_system/schema.sql
```

Optional env var:

```bash
export DATABASE_URL="postgresql+psycopg2://postgres:postgres@localhost/notifications"
```

Run:

```bash
uvicorn main:app --reload --port 8080
```

### APIs

- **Subscribe**: `POST /subscribe`

```json
{"user_id": 1, "item_id": 100, "channels": ["email", "push"]}
```

- **Restock**: `POST /inventory/restock`

```json
{"item_id": 100, "new_stock": 20}
```

- **Health**: `GET /health`

### Notes

- Inventory is mocked (same as the Go version) — restock triggers notification processing.
- To receive notifications for a later restock, the user must re-subscribe after being notified.
