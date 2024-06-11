package main

import (
	"fmt"
	"os"

	"github.com/martinmunillas/posmigrator"
	"github.com/spf13/cobra"
)

var cfg posmigrator.Config

func validateFlags(cfg posmigrator.Config) error {
	if cfg.DbHost == "" {
		return fmt.Errorf("dbhost is required")
	}
	if cfg.DbPort == "" {
		return fmt.Errorf("dbport is required")
	}
	if cfg.DbUser == "" {
		return fmt.Errorf("dbuser is required")
	}
	if cfg.DbPassword == "" {
		return fmt.Errorf("dbpassword is required")
	}
	if cfg.DbName == "" {
		return fmt.Errorf("dbname is required")
	}
	if cfg.MigrationsPath == "" {
		return fmt.Errorf("migrationspath is required")
	}
	return nil
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run migrations",
	Run: func(cmd *cobra.Command, args []string) {
		err := validateFlags(cfg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		conn, err := posmigrator.ConnectPostgres(cfg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = posmigrator.RunMigrations(conn, cfg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var ensureCmd = &cobra.Command{
	Use:   "ensure",
	Short: "Ensure all migrations ran and are valid",
	Run: func(cmd *cobra.Command, args []string) {
		err := validateFlags(cfg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		conn, err := posmigrator.ConnectPostgres(cfg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = posmigrator.EnsureAllMigrationsRanAndAreValid(conn, cfg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "posmigrator",
		Short: "A tool to manage database migrations",
	}

	rootCmd.PersistentFlags().StringVar(&cfg.DbHost, "dbhost", "", "Database host")
	rootCmd.PersistentFlags().StringVar(&cfg.DbPort, "dbport", "", "Database port")
	rootCmd.PersistentFlags().StringVar(&cfg.DbUser, "dbuser", "", "Database user")
	rootCmd.PersistentFlags().StringVar(&cfg.DbPassword, "dbpassword", "", "Database password")
	rootCmd.PersistentFlags().StringVar(&cfg.DbName, "dbname", "", "Database name")
	rootCmd.PersistentFlags().StringVar(&cfg.MigrationsPath, "migrationspath", "", "Path to migrations")

	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(ensureCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
