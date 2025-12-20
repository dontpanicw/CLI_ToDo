package app

import (
	"CLI_ToDo/backend/config"
	"CLI_ToDo/backend/internal/adapter/repo/postgres"
	"CLI_ToDo/backend/internal/cases"
	migrations "CLI_ToDo/backend/pkg/migrations/postgres"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
)

// Start инициализирует зависимости и обрабатывает CLI-команды
func Start(cfg config.Config, logger *zap.Logger) error {
	db, err := sql.Open("pgx", cfg.PgConnStr)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	if err := migrations.Migrate(db); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(10 * time.Minute)

	logger.Info("database ready")

	taskRepo := postgres.NewPostgresRepository(db)
	taskUC := cases.NewTaskUseCase(taskRepo)

	if len(os.Args) < 2 {
		printUsage()
		return nil
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "create":
		fs := flag.NewFlagSet("create", flag.ExitOnError)
		name := fs.String("name", "", "название задачи")
		description := fs.String("desc", "", "описание задачи")
		if err := fs.Parse(args); err != nil {
			return fmt.Errorf("failed to parse flags: %w", err)
		}

		id, err := taskUC.Create(*name, *description)
		if err != nil {
			return err
		}
		fmt.Printf("Задача создана с id=%d\n", id)
	case "read":
		fs := flag.NewFlagSet("read", flag.ExitOnError)
		id := fs.Int("id", 0, "id задачи")
		if err := fs.Parse(args); err != nil {
			return fmt.Errorf("failed to parse flags: %w", err)
		}

		task, err := taskUC.Read(*id)
		if err != nil {
			return err
		}
		if task == nil {
			fmt.Println("Задача не найдена")
			return nil
		}
		fmt.Printf("ID: %d\nName: %s\nDescription: %s\n", task.Id, task.Name, task.Description)
	case "update":
		fs := flag.NewFlagSet("update", flag.ExitOnError)
		id := fs.Int("id", 0, "id задачи")
		name := fs.String("name", "", "название задачи")
		description := fs.String("desc", "", "описание задачи")
		if err := fs.Parse(args); err != nil {
			return fmt.Errorf("failed to parse flags: %w", err)
		}

		updatedID, err := taskUC.Update(*id, *name, *description)
		if err != nil {
			return err
		}
		fmt.Printf("Задача обновлена id=%d\n", updatedID)
	case "delete":
		fs := flag.NewFlagSet("delete", flag.ExitOnError)
		id := fs.Int("id", 0, "id задачи")
		if err := fs.Parse(args); err != nil {
			return fmt.Errorf("failed to parse flags: %w", err)
		}

		if err := taskUC.Delete(*id); err != nil {
			return err
		}
		fmt.Printf("Задача удалена id=%d\n", *id)
	case "done":
		fs := flag.NewFlagSet("done", flag.ExitOnError)
		id := fs.Int("id", 0, "id задачи")
		if err := fs.Parse(args); err != nil {
			return fmt.Errorf("failed to parse flags: %w", err)
		}

		if err := taskUC.MarkDone(*id); err != nil {
			return err
		}
		fmt.Printf("Задача отмечена как выполненная id=%d\n", *id)
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Println("Неизвестная команда:", cmd)
		printUsage()
	}

	return nil
}

func printUsage() {
	fmt.Println("CLI ToDo. Команды:")
	fmt.Println("  create -name \"Название\" -desc \"Описание\"")
	fmt.Println("  read   -id <id>")
	fmt.Println("  update -id <id> -name \"Название\" -desc \"Описание\"")
	fmt.Println("  delete -id <id>")
	fmt.Println("  done   -id <id>")
	fmt.Println("  help")
}
