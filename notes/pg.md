# PostgreSQL
It is open source, very reliable, and has some odern features:
- support for array & JSON data types
- full-text search
- geospatial queries

## setup pg
To install PostgreSQL on your computer:
```sh
# on Ubuntu
sudo apt install postgresql

# on windows
choco install postgresql
```
When PG is installed, a `psql` binary gets created on your system. This contains a *terminal_based front-end* for working with PG.
```sh
psql --version
# psql (PostgreSQL) 16.14 (Ubuntu 16.14-0ubuntu0.24.04.1)
```

During installation, *an operating system user* named `postgres` should also have been created:
```sh
cat /etc/passwd | grep 'postgres'
# postgres:x:105:109:PostgreSQL administrator,,,:/var/lib/postgresql:/bin/bash
```
By default, PG uses an authentication scheme called **peer authentication** for any connection from the local machine.
Hence, if we switch to the operating user called `postgres`, we chould be able to connect to PG using `psql` without needing any further authentication.

```sh
sudo -u postgres psql
# psql (16.14 (Ubuntu 16.14-0ubuntu0.24.04.1))
# Type "help" for help.

# postgres=# SELECT current_user;
#  current_user
# --------------
#  postgres
# (1 row)
# \q
```

NOTE: `\` is a meta command.
`\l` list all databases, `\dt` list tables, `\du` list users, `\?`

### create database, users, and extensions
```sql
postgres=# CREATE DATABASE pixeldb;
-- CREATE DATABASE

postgres=# \c pixeldb;
-- You are now connected to database "pixeldb" as user "postgres".
-- pixeldb=#
```

Create `pixel_user`, without superuser permissions. We connect to the db from our Go application via `pixel_user`.
We want to set up this new user to use **password-based authentication**, instead of *peer authentication*.

```sql
pixeldb=# CREATE ROLE pixel_user WITH LOGIN PASSWORD 'ChangeM3';
-- CREATE ROLE

pixeldb=# CREATE EXTENSION IF NOT EXISTS citext;
-- CREATE EXTENSION

pixeldb=# \q
```

Connecting as the new user:
```sh
psql --host=localhost --dbname=pixeldb --username=pixel_user
# Password for user pixel_user:
# psql (16.14 (Ubuntu 16.14-0ubuntu0.24.04.1))
# SSL connection (protocol: TLSv1.3, cipher: TLS_AES_256_GCM_SHA384, compression: off)
# Type "help" for help.

# pixeldb=> SELECT current_user;
#  current_user
# --------------
#  pixel_user
# (1 row)

```

Other approaches:
```sh
PGPASSWORD='ChangeM3' psql -h localhost -d pixeldb -U pixel_user

# -----
export PGPASSWORD='ChangeM3'

psql -h localhost -d pixeldb -U pixel_user
```

## connecting to PostgreSQL

To work with a SQL database, we need to use a **database driver** to act as the *middleman* between Go and the database itself.

```sh
go get github.com/lib/pq@v1
```

## SQL migrations
To manage SQL migrations in this project, we're going to use the `migrate` command-line tool.

```sh
# https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

curl -L https://github.com/golang-migrate/migrate/releases/download/$version/migrate.$os-$arch.tar.gz | tar xvz

uname -m
dpkg --print-architecture

# uname -m → x86_64 ✅ (the kernel reports a 64-bit x86 CPU)
# dpkg --print-architecture → amd64 ✅ (Debian/Ubuntu package architecture)

cd /tmp
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
mv migrate ~/go/bin/

# verify:
migrate -version
# 4.16.2
```

### working with SQL migrations
Generate a pair of migration files:
```sh
migrate create -seq -ext=.sql -dir=./migrations create_movies_table

ls migrations
# 000001_create_movies_table.down.sql  000001_create_movies_table.up.sql
```

Adding CHECK constraints to enforce some of our business rules at the database-level:
```sh
migrate create -seq -ext=.sql -dir=./migrations add_movies_check_constraints
# 000002_add_movies_check_constraints.down.sql  000002_add_movies_check_constraints.up.sql
```

If the CHECK fails, the database driver will throw an error:
```sh
pq: new row for relation "movies" violates check constraint "genres_length_check"
```

### executing the migrations
```sh
migrate -path=./migrations -database=$PIXEL_DB_DSN up
# error: pq: permission denied for schema public in line 0: CREATE TABLE IF NOT EXISTS "public"."schema_migrations"
```
Fix 👇
```sh
sudo -u postgres psql

# docker exec -it <postgres-container> psql -U postgres

```
Once you're at the `postgres=#` prompt, run:
```sql
ALTER DATABASE pixeldb OWNER TO pixel_user;

# verify
\l

\q
```
Try the migration again:
```sh
migrate -path=./migrations -database="$PIXEL_DB_DSN" up
# 1/u create_movies_table (27.465904ms)
# 2/u add_movies_check_constraints (40.311514ms)
```

Verify
```sh
psql $PIXEL_DB_DSN
# psql (16.14 (Ubuntu 16.14-0ubuntu0.24.04.1))
# SSL connection (protocol: TLSv1.3, cipher: TLS_AES_256_GCM_SHA384, compression: off)
# Type "help" for help.

# pixeldb=> \dt 🎉
#                 List of relations
#  Schema |       Name        | Type  |   Owner
# --------+-------------------+-------+------------
#  public | movies            | table | pixel_user
#  public | schema_migrations | table | pixel_user
# (2 rows)
```
The `schema_migrations` table is automatically generated by the `migrate` tool and used to track which migrations have been applied:
```sql
pixeldb=> SELECT * FROM schema_migrations;
 version | dirty
---------+-------
       2 | f
(1 row)


pixeldb=> \d movies
/*
                                     Table "public.movies"
   Column   |            Type             | Collation | Nullable |           Default
------------+-----------------------------+-----------+----------+------------------------------
 id         | bigint                      |           | not null | generated always as identity
 created_at | timestamp(0) with time zone |           | not null | now()
 title      | text                        |           | not null |
 year       | integer                     |           | not null |
 runtime    | integer                     |           | not null |
 genres     | text[]                      |           | not null |
 version    | integer                     |           | not null | 1
Indexes:
    "movies_pkey" PRIMARY KEY, btree (id)
Check constraints:
    "genres_length_check" CHECK (array_length(genres, 1) >= 1 AND array_length(genres, 1) <= 5)
    "movies_runtime_check" CHECK (runtime >= 0)
    "movies_year_check" CHECK (year >= 1888 AND year::double precision <= date_part('year'::text, now()))
*/
```
## MiSK

### optimizing PostgreSQL settings
The default settings that PG ships with are quite conservative. We can improve the performance of our database by tweaking the values in the `postgresql.conf` file.

```sh
sudo -u postgres psql -c 'SHOW config_file;'
#                config_file
# -----------------------------------------
#  /etc/postgresql/16/main/postgresql.conf
# (1 row)
```
