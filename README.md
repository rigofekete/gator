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
cd gator
go install .
```

> **Note:** With `go install` the binary (`gator`) becomes globally available, you can run it from anywhere in the terminal by name.

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

Change directory to sql/schema and migrate the database:

```
goose postgres "postgres://USER:PASSWORD@localhost:5432/gator" up
```

Then, reset the database to apply the clean schema:

```bash
gator reset
```

## Usage

Create a user:

```bash
gator register <username>
```

<img width="616" height="91" alt="Image" src="https://github.com/user-attachments/assets/9ed57271-354b-4338-acc1-257e9cbd0360" />
<br>

List users:

```bash
gator users
```

<img width="349" height="106" alt="Image" src="https://github.com/user-attachments/assets/28447c30-fe52-4fb5-9c24-8a915db4a56f" />
<br>

Login user:

```bash
gator login <registered username>
```

<img width="492" height="81" alt="Image" src="https://github.com/user-attachments/assets/49dcba78-c711-4fc3-b0fc-4521c999f7db" />
<br>

Add a feed and follow it automatically:

```
gator addfeed <name> <url>
```

<img width="1351" height="454" alt="Image" src="https://github.com/user-attachments/assets/9a772119-a7dc-4eb1-9fe0-52ac437efc15" />
<br>

Follow a feed manually:

```bash
gator follow <url>
````

Unfollow a feed:

```bash
gator unfollow <url>
```

List available feeds:

```bash
gator feeds
```

<img width="694" height="618" alt="Image" src="https://github.com/user-attachments/assets/ae3f3f52-8ea4-4ca6-8784-5c9fed4a7107" />
<br>

List followed feeds by currently logged in user:

```bash
gator following
```

<img width="394" height="99" alt="Image" src="https://github.com/user-attachments/assets/fe459834-1a7e-494f-984c-e96615015303" />
<br>

Fetch and aggregate posts at a fixed interval to avoid server DOS:

```bash
gator agg <duration>
# examples:
# gator agg 1m
# gator agg 10s
# gator agg 5m30s
```

<img width="1143" height="246" alt="Image" src="https://github.com/user-attachments/assets/a28a20b3-2112-4540-843f-b83653dab0e4" />

<br>

Browse aggregated posts and post them by limit (if no limit is provided it defaults to 2):

```bash
gator browse <limit>
# examples:
# gator browse  - limit defaults to 2
# gator browse 3
```

<img width="1912" height="652" alt="Image" src="https://github.com/user-attachments/assets/f22dbd57-1292-4609-9d1f-39be94f0b424" />

<br>

<img width="1554" height="361" alt="Image" src="https://github.com/user-attachments/assets/9124275c-26af-491a-bc83-f976c76fe4ee" />

<br>

Reset database and apply clean schema:

```bash
gator reset
```
