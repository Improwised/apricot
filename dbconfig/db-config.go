/*
Package dbconfig reads database settings from a database.yml file (following the rails database.yml convention)
and generates a connection string for the github.com/lib/pq and github.com/go-sql-driver/mysql drivers.
*/
package dbconfig

import (
	"strings"
)

/*
Settings returns the database settings from the database.yml file from a given application
and the corresponding application enviroment. The location to the database.yml file and
the enviroment is configured in the settings json configuration file.
If environment is NOT configured, you can set the environment variable APPLICATION_ENV (on os level).
If this is also not defined "development" is the default.
*/
func Settings(path string) map[string]string {

	jsonConf := LoadJSONConfig(path)
	dbConfig := LoadYamlConfig(jsonConf.Database_file)
	environment := jsonConf.Environment

	return dbConfig[environment]
}

/*
PostgresConnectionString returns the connection string to open a sql session
used by the github.com/lib/pq package like for example:
"host=dbserver.org password=password user=dbuser dbname=blog_production sslmode=disable"
The first parameter is the path to the database settings configuration (json) file
and the second paramater defines the sslmode.
*/
func PostgresConnectionString(path string, sslmode string) string {
	settings := Settings(path)

	connection := []string{
		"host=", settings["host"], " ",
		"password=", settings["password"], " ",
		"user=", settings["username"], " ",
		"dbname=", settings["database"], " ",
		"sslmode=", sslmode}

	return strings.Join(connection, "")
}
