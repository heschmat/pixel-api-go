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
