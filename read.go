package posmigrator

import (
	"fmt"
	"os"
	"strconv"
)

type migrationFile struct {
	id       int64
	name     string
	fileName string
}

func (m migrationFile) sql(config Config) ([]byte, error) {
	return os.ReadFile(config.MigrationsPath + "/" + m.fileName)
}

func readMigrationFiles(config Config) ([]migrationFile, error) {
	dirs, err := os.ReadDir(config.MigrationsPath)
	if err != nil {
		return nil, err
	}

	migrations := make([]migrationFile, 0, len(dirs))
	for _, dir := range dirs {
		fileName := dir.Name()
		if dir.IsDir() {
			if fileName != "views" {
				yellow.Printf("Found directory \"%s\" in migrations directory\n", fileName)
			}
			continue
		}
		if !migrationRegex.MatchString(fileName) {
			red.Printf(invalidMigrationFileName, fileName)
			os.Exit(1)
		}
		matches := migrationRegex.FindStringSubmatch(fileName)
		if len(matches) != 3 {
			return nil, fmt.Errorf(invalidMigrationFileName, fileName)
		}
		name := matches[2]
		id, err := strconv.ParseInt(matches[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf(invalidMigrationFileName, fileName)
		}
		migrations = append(migrations, migrationFile{
			id:       id,
			name:     name,
			fileName: fileName,
		})
	}

	return migrations, nil
}
