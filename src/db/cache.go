package db

import (
	"database/sql"
	"os"

	"github.com/rzhade3/dep-recon/src/manifest"
	_ "modernc.org/sqlite"
)

type DbConfig struct {
	DbFilename string
}

func InitializeCache(cacheFilepath string) (DbConfig, error) {
	// Create a new SQLite database if it doesn't exist
	// Create a new table for READMEs if it doesn't exist
	// Create a new table for dependencies if it doesn't exist
	// Create a new table for dependencies if it doesn't exist
	// Check if the database exists
	if _, err := os.Stat(cacheFilepath); os.IsNotExist(err) {
		// Create a new database
		os.Create(cacheFilepath)
	}

	db, err := sql.Open("sqlite", cacheFilepath)
	if err != nil {
		return DbConfig{}, err
	}

	// Create a new table if it doesn't exist
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS dependencies (id INTEGER PRIMARY KEY, package_name TEXT, ecosystem TEXT, readme BLOB)")
	if err != nil {
		return DbConfig{}, err
	}
	return DbConfig{
		DbFilename: cacheFilepath,
	}, nil
}

func (d *DbConfig) WriteToCache(packageName, ecosystem, readme string) error {
	db, err := sql.Open("sqlite", d.DbFilename)
	if err != nil {
		return err
	}
	defer db.Close()

	// Insert a new row
	_, err = db.Exec("INSERT INTO dependencies (package_name, ecosystem, readme) VALUES (?, ?, ?)", packageName, ecosystem, readme)
	if err != nil {
		return err
	}

	return nil
}

func (d *DbConfig) ReadFromCache(packageName, ecosystem string) (string, error) {
	db, err := sql.Open("sqlite", d.DbFilename)
	if err != nil {
		return "", err
	}
	defer db.Close()

	// Query the database
	rows, err := db.Query("SELECT readme FROM dependencies WHERE package_name = ? AND ecosystem = ?", packageName, ecosystem)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		var readme string
		err = rows.Scan(&readme)
		if err != nil {
			return "", err
		}
		// There will only be one row, so we can return it immediately
		return readme, nil
	}
	return "", nil
}

func (d *DbConfig) DeleteFromCache(packageName, ecosystem string) error {
	db, err := sql.Open("sqlite", d.DbFilename)
	if err != nil {
		return err
	}
	defer db.Close()
	// Delete the row
	_, err = db.Exec("DELETE FROM dependencies WHERE package_name = ? AND ecosystem = ?", packageName, ecosystem)
	if err != nil {
		return err
	}
	return nil
}

// Fetches the README for a dependency from the cache if it exists, or pulls it from the package manager's registry
func (d *DbConfig) FetchDependencyReadme(manifest manifest.Language, dependency, version string) (string, error) {
	readme, err := d.ReadFromCache(dependency, manifest.GetEcosystem())
	if err != nil {
		return "", err
	}
	if readme == "" {
		readme, err = manifest.PullDependencyReadme(dependency, version)
		if err != nil {
			return "", err
		}
		// If the README just doesn't exist, don't try to write anything
		if readme == "" {
			return "", nil
		}
		err = d.WriteToCache(dependency, manifest.GetEcosystem(), readme)
		if err != nil {
			return "", err
		}
	}
	return readme, nil
}
