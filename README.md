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

Swagger Documentation

Updating Swagger Documentation
To update the Swagger documentation for the API, follow these steps:

Make the necessary changes to your API endpoints and their annotations in your code.

Open a terminal and navigate to the package where your handlers are located.

Run the swag init command to regenerate the Swagger JSON file:

swag init
This command scans your Go files and updates the Swagger JSON file based on the annotations in your code. Make sure you have the swag tool installed and available in your system's PATH.

After running the swag init command, the Swagger JSON file will be updated with your latest code changes.

Accessing Swagger Documentation
To access the Swagger documentation and interact with the API, follow these steps:

Build and run your application.

Open a web browser and navigate to the Swagger UI URL. The default URL is typically http://localhost:8080/swagger/index.html.

The Swagger UI provides a user-friendly interface to explore the API endpoints, view request/response details, and even test the API by sending requests directly from the Swagger UI.