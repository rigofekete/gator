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

<img width="640" height="81" alt="Image" src="https://github.com/user-attachments/assets/c5468da7-d0bf-4c61-be3e-18529c5a0509" />
<br>

List users:

```bash
gator users
```

<img width="418" height="112" alt="Image" src="https://github.com/user-attachments/assets/1cb5b9e8-4819-43f8-90f8-22995236855d" />
<br>

Login user:

```bash
gator login <registered username>
```

<img width="570" height="78" alt="Image" src="https://github.com/user-attachments/assets/467555db-2aa8-425d-9d0c-50c265088708" />
<br>

Add a feed and follow it automatically:

```
gator addfeed <name> <url>
```

<img width="1357" height="448" alt="Image" src="https://github.com/user-attachments/assets/603ab050-56a4-4691-9eb0-0ceda6b6d476" />
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

<img width="712" height="622" alt="Image" src="https://github.com/user-attachments/assets/8d28858e-32cf-4080-834b-eec440917542" />
<br>

List followed feeds by currently logged in user:

```bash
gator following
```

<img width="439" height="94" alt="Image" src="https://github.com/user-attachments/assets/6c4a6c35-7cf0-4cba-9c68-1a5c82c5fdf4" />
<br>

Fetch and aggregate posts at a fixed interval to avoid server DOS:

```bash
gator agg <duration>
# examples:
# gator agg 1m
# gator agg 10s
# gator agg 5m30s
```

<img width="1156" height="237" alt="Image" src="https://github.com/user-attachments/assets/efde1aff-456a-472f-848b-185163d91d03" />

<br>

Browse aggregated posts and post them by limit (if no limit is provided it defaults to 2):

```bash
gator browse <limit>
# examples:
# gator browse  - limit defaults to 2
# gator browse 3
```

<img width="1911" height="604" alt="Image" src="https://github.com/user-attachments/assets/73fc8da9-e280-42ea-9175-7d42b160ec59" />

<br>

<img width="1920" height="388" alt="Image" src="https://github.com/user-attachments/assets/0bcb3b43-b1f3-462b-aa44-1ceb801da77f" />
<br>

Reset database and apply clean schema:

```bash
gator reset
```
