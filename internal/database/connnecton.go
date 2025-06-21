package database

import (
    "database/sql"
    "fmt"
    "io/ioutil"
    "path/filepath"
    "sort"
    
    _ "github.com/lib/pq"
)

func Connect(databaseURL string) (*sql.DB, error) {
    db, err := sql.Open("postgres", databaseURL)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }
    
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
    
    return db, nil
}

func RunMigrations(db *sql.DB) error {
    migrationDir := "database/migrations"
    
    files, err := ioutil.ReadDir(migrationDir)
    if err != nil {
        return fmt.Errorf("failed to read migration directory: %w", err)
    }
    
    // Sort files to ensure they run in order
    sort.Slice(files, func(i, j int) bool {
        return files[i].Name() < files[j].Name()
    })
    
    for _, file := range files {
        if filepath.Ext(file.Name()) != ".sql" {
            continue
        }
        
        migrationPath := filepath.Join(migrationDir, file.Name())
        content, err := ioutil.ReadFile(migrationPath)
        if err != nil {
            return fmt.Errorf("failed to read migration file %s: %w", file.Name(), err)
        }
        
        if _, err := db.Exec(string(content)); err != nil {
            return fmt.Errorf("failed to execute migration %s: %w", file.Name(), err)
        }
        
        fmt.Printf("Applied migration: %s\n", file.Name())
    }
    
    return nil
}
