# Welcome to Gator!

Gator is a command-line rss feed aggregator written in Go.

## Prerequisites

You'll need Postgres and Go v 1.24.2 to run this program.

## Installation

```sh
go install github.com/michaeldebetaz/gator
```

## Configuration

You'll need to create a `.gatorconfig.json` for the program to connect to your configured Postgres database. On Linux, you can easly initiate a postgres database by using the default services.

First install all the necessary packages:

```sh
sudo apt update
sudo apt install postgresql postgresql-contrib
psql --version
```

Start the postgres service:

```sh
sudo passwd postgres
sudo service postgresql start
```

Run `psql`:

```sh
sudo -u postgres psql
```

Create the gator database and update the `postgres` user password:

```sh
postgres=# CREATE DATABASE gator;
postgres=# \c gator;
postgres=# ALTER USER postgres PASSWORD 'postgres';
```

The `.gatorconfig.json` file should be located in your `$HOME` (`~`) directory.

```json
{
  "db_url": "postgres://postgres:postgres@localhost:5433/gator?sslmode=disable",
} 
```

## Usage

```sh
gator register <username>               # register a new user
gator login <username>                  # login as an existing user
gator reset                             # delete all users
gator users                             # list all users
gator agg <time between req>            # aggregate new rss feed posts
gator addfeed <feed name> <feed url>    # add a new feed to the user
gator feeds                             # list all feeds
gator follow <url>                      # add a feed follow to the feed url for the user
gator following                         # list all followed feeds for the user
gator browse <limit>                    # list all the posts of the followed feeds of the user
```



