# SMS-panel

Creating a New Migration
To create a new migration, follow these steps:

Open a terminal and navigate to the project's database/migrations directory.

Run the following command to create a new migration file:

migrate create -ext sql -dir migrations -seq <migration_name>
Replace <migration_name> with a descriptive name for your migration.

Two files will be created: <timestamp>_<migration_name>.down.sql and <timestamp>_<migration_name>.up.sql. The up file contains the SQL statements for applying the migration, while the down file contains the SQL statements for rolling back the migration.

Open the up and down SQL files and write the necessary SQL statements to define your migration. Update the tables, columns, constraints, or any other changes you want to make to the database schema.

Applying Migrations
To apply the migrations and update the database schema, run the following command:

migrate -path database/migrations -database "<database_connection_string>" up
Replace <database_connection_string> with the actual database connection string.
like this migrate -path database/migrations -database "postgres://user:password@host:port/database?sslmode=disable" up


Rolling Back Migrations
To roll back the last applied migration and revert the changes made to the database schema, run the following command:


migrate -path database/migrations -database "<database_connection_string>" down
Replace <database_connection_string> with the actual database connection string.
like this migrate -path database/migrations -database "<database_connection_string>" down

That's it! With these instructions, you can create new migrations and apply them to your database, as well as roll back migrations if needed.

## Swagger
start using swagger [echo-swagger man page](https://github.com/swaggo/echo-swagger)

