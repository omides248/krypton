package cmd

import (
	"errors"
	"fmt"
	"krypton/wallet/config"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)

type RunServerFunc func() error

func NewRootCmd(runServer RunServerFunc) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "wallet-service",
		Short: "Main entry point for the wallet service",
		Run: func(cmd *cobra.Command, args []string) {
			if err := runServer(); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error running server: %v\n", err)
				os.Exit(1)
			}
		},
	}

	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Database migration tool",
	}

	migrateUpCmd := &cobra.Command{
		Use:   "up",
		Short: "Apply all available migrations",
		Run: func(cmd *cobra.Command, args []string) {
			runMigrations("up")
		},
	}

	migrateDownCmd := &cobra.Command{
		Use:   "down",
		Short: "Revert the last migration",
		Run: func(cmd *cobra.Command, args []string) {
			runMigrations("down")
		},
	}

	migrateCmd.AddCommand(migrateUpCmd, migrateDownCmd)
	rootCmd.AddCommand(migrateCmd)

	return rootCmd
}

func runMigrations(direction string) {
	cfg, err := config.LoadConfig("./config")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	databaseURL := cfg.Database.Postgresql.URI
	migrationsPath := cfg.Database.Migration.Path

	m, err := migrate.New(migrationsPath, databaseURL)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	if direction == "up" {
		fmt.Println("Applying migrations...")
		err = m.Up()
	} else {
		fmt.Println("Reverting last migration...")
		err = m.Steps(-1)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		fmt.Println("No new migrations to apply.")
		return
	}
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("Migration finished successfully.")
}
