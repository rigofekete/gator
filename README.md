# gator – RSS Feed Aggregator (CLI)

gator is a Go CLI that aggregates RSS feeds into a PostgreSQL database. It supports managing users, subscribing to feeds, fetching posts, and resetting state for clean testing.

## Tools and Dependencies
- Go
- PostgreSQL
- goose (migrations)
- sqlc (type-safe queries)

## Requirements

- Go 1.21+ 
- PostgreSQL 14+

## Install

```bash
git clone https://github.com/rigofekete/gator
```

```bash
go install .
```

## Database Setup

You need a running PostgreSQL instance and an empty database. The app manages tables/migrations, not database creations.

### 1 - Create the gator database

Using psql
```bash
# connect as your Postgres user
psql -U postgres -h localhost
```

-- inside psql:
```
CREATE DATABASE gator;
```

Create a gator folder in the .config of your home path with a .gatorconfig.json file inside. 

```bash
~/.config/gator/.gatorconfig.json
```

Add the database url (replace fields with your own Postgres user and password):

```json 
{
    "db_url": "postgres://USER:PASSWORD@localhost:5432/gator?sslmode=disable"
}
```

### 2 - Apply database schema

from the terminal, run:

```bash
gator reset
```

## Usage

Create a user:
```bash
gator user add <username>
```

List users:
```bash
gator users
```

Follow a feed manually:
```bash
gator follow <url>
````

Add a feed and follow it automatically:
```
gator addfeed <name> <url>
```

Unfollow a feed:
```bash
gator unfollow <url>
```

List followed feeds:
```bash
gator following
```

Fetch and aggregate posts at a fixed interval to avoid server DOS:
```bash
gator agg <duration>
# examples:
# gator agg 1m
# gator agg 10s
# gator agg 5m30s
```

Browse aggregated posts and post them by limit (if no limit is provided it defaults to 2):
```bash
gator browse <limit>
# examples:
# gator browse  - limit defaults to 2
# gator browse 3
```

Reset database and apply clean schema:
```bash
gator reset
```








