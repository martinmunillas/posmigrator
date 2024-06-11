
# posmigrator

`posmigrator` is a command-line tool designed to manage database migrations for PostgreSQL. It provides functionalities to run and ensure the validity of database migrations.

## Installation

To install `posmigrator`, you need to have Go installed on your system. Clone the repository and build the project using the following commands:

```sh
git clone https://github.com/martinmunillas/posmigrator.git
cd posmigrator
go build -o posmigrator
```

## Usage

The tool provides two main commands: `migrate` and `ensure`. You can run these commands with the appropriate flags to manage your database migrations.

### Commands

#### `migrate`

The `migrate` command runs the migrations on the specified PostgreSQL database.

```sh
posmigrator migrate --dbhost <host> --dbport <port> --dbuser <user> --dbpassword <password> --dbname <database> --migrationspath <path>
```

#### `ensure`

The `ensure` command verifies that all migrations have been run and are valid.

```sh
posmigrator ensure --dbhost <host> --dbport <port> --dbuser <user> --dbpassword <password> --dbname <database> --migrationspath <path>
```

### Flags

- `--dbhost`: The host of the PostgreSQL database.
- `--dbport`: The port of the PostgreSQL database.
- `--dbuser`: The user to connect to the PostgreSQL database.
- `--dbpassword`: The password to connect to the PostgreSQL database.
- `--dbname`: The name of the PostgreSQL database.
- `--migrationspath`: The path to the directory containing the migration files.

### Examples

To run migrations:

```sh
posmigrator migrate --dbhost localhost --dbport 5432 --dbuser user --dbpassword password --dbname mydb --migrationspath ./migrations
```

To ensure all migrations ran and are valid:

```sh
posmigrator ensure --dbhost localhost --dbport 5432 --dbuser user --dbpassword password --dbname mydb --migrationspath ./migrations
```

