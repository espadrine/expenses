package expense

import (
	"database/sql"
	"embed"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
}

func NewDB() (*Store, error) {
	db, isInit, err := initConfig()
	if err != nil {
		return nil, err
	}
	store := Store{db}
	if !isInit {
		err = store.initDB()
		if err != nil {
			return nil, err
		}
	}
	return &store, nil
}

func initConfig() (db *sql.DB, alreadyInit bool, err error) {
	alreadyInit = true

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}
	configPath := filepath.Join(homeDir, ".config")
	expensePath := filepath.Join(configPath, "expense")
	defaultDBPath := filepath.Join(expensePath, "db.sqlite")
	expenseDBEnv := os.Getenv("EXPENSE_DB")

	dbPath := defaultDBPath
	if expenseDBEnv != "" {
		dbPath = expenseDBEnv
	}

	// Ensure ~/.config is there.
	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		alreadyInit = false
		err := os.Mkdir(configPath, 0755)
		if err != nil {
			return nil, alreadyInit, err
		}
	} else if err != nil {
		return
	}

	// Ensure ~/.config/expense is there.
	_, err = os.Stat(expensePath)
	if os.IsNotExist(err) {
		alreadyInit = false
		err := os.Mkdir(expensePath, 0700)
		if err != nil {
			return nil, alreadyInit, err
		}
	} else if err != nil {
		return
	}

	// Check whether the database file exists.
	_, err = os.Stat(dbPath)
	if os.IsNotExist(err) {
		alreadyInit = false
	}

	db, err = sql.Open("sqlite3", dbPath)
	return
}

func (s *Store) initDB() (err error) {
	err = migrateDB(s.db)
	if err != nil {
		return
	}
	// If the user is not defined, let’s add it.
	username := os.Getenv("USER")
	_, err = s.getUser(username)
	if err == sql.ErrNoRows {
		err = nil
		_, err = s.createUser(username)
		return
	}
	return
}

//go:embed sql/*.sql
var sqlMigrations embed.FS

func migrateDB(db *sql.DB) error {
	migration, err := sqlMigrations.ReadFile("sql/v0.sql")
	if err != nil {
		return err
	}
	// Execute the migration one statement at a time.
	sqlStatements := strings.Split(string(migration), ";")
	for _, stmt := range sqlStatements {
		_, err = db.Exec(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) Close() {
	s.db.Close()
}
