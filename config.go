package posmigrator

type Config struct {
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string

	MigrationsPath string
}
