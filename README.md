# Event Trigger Platform

A backend application for managing scheduled and API-based event triggers. Built with Go, PostgreSQL, and Docker.

## Table of Contents
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Local Setup](#local-setup)
- [Database Schema](#database-schema)
- [Deployment](#deployment)


## Features
- Create/manage scheduled triggers (cron jobs)
- Create/manage API triggers with payload validation
- Test triggers manually
- Event logging with 48-hour retention
- Archived events view

## Prerequisites
- Go 1.22+
- Docker 20.10+
- PostgreSQL 13+
- Google Cloud account (for deployment)

## Local Setup

1. Clone the repository:
```bash
git clone https://github.com/your-username/event-trigger.git
cd event-trigger
```
2. Start services using Docker Compose:

```
docker-compose up --build
```
3. Access endpoints:

- API: [http://localhost:8080](http://localhost:8080)
- Frontend: [http://localhost:8080/web](http://localhost:8080)
## API Documentation
Example Requests  
Create Scheduled Trigger:
```
POST /api/triggers HTTP/1.1
Content-Type: application/json

{
  "type": "scheduled",
  "cron_expression": "*/5 * * * *",
  "is_recurring": true
}
```
Test API Trigger:

```
POST /api/triggers/{id}/test HTTP/1.1
Content-Type: application/json

{
  "payload": {
    "key1": "value1",
    "key2": "value2"
  }
}
```
## Database Schema

### Triggers Table
| Column          | Type        | Description                    |
|-----------------|-------------|--------------------------------|
| id              | UUID        | Primary key                    |
| type            | VARCHAR(20) | Scheduled or API               |
| cron_expression | TEXT        | Cron schedule                  |
| next_run        | TIMESTAMPTZ | Next execution time            |
| is_recurring    | BOOLEAN     | Recurring flag                 |
| api_payload     | JSONB       | Payload schema for API triggers|
| is_active       | BOOLEAN     | Active status                  |

### Event Logs Table
| Column      | Type        | Description                 |
|-------------|-------------|-----------------------------|
| id          | UUID        | Primary key                 |
| trigger_id  | UUID        | Foreign key to triggers     |
| triggered_at| TIMESTAMPTZ | Event execution timestamp   |
| payload     | JSONB       | API payload data            |
| is_test     | BOOLEAN     | Test execution flag         |
| status      | VARCHAR(10) | Active or archived          |


## Deployment
### Production Deployment
- API: [https://backend-498987061135.us-central1.run.app](https://backend-498987061135.us-central1.run.app)
- FrontEnd: [https://backend-498987061135.us-central1.run.app/web](https://backend-498987061135.us-central1.run.app/web)