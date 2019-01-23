package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-pg/pg"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"time"
)

func InitPostgresContainer(ctx context.Context) testcontainers.Container {
	log.Printf("Starting Postgres container...")
	user := "postgres_user"
	password := "postgres_password"
	dbName := "postgres_db"
	envVariables := map[string]string{
		"POSTGRES_USER":     user,
		"POSTGRES_PASSWORD": password,
		"POSTGRES_DB":       dbName,
	}
	req := testcontainers.ContainerRequest{
		Image:        "postgres:10.6",
		ExposedPorts: []string{"5432/tcp"},
		Env:          envVariables,
		WaitingFor:   &PostgresWaitStrategy{User: user, Database: dbName, Password: password},
	}
	container, host := StartContainerAndGetHost(ctx, req, "5432")

	_ = os.Setenv("DB_HOST", host)
	_ = os.Setenv("DB_USER", user)
	_ = os.Setenv("DB_USER_PASSWORD", password)
	_ = os.Setenv("DB_NAME", dbName)

	log.Printf("Test container with Postgres has been started. Host[%s] user[%s] password[%s] dbName[%s]",
		host, user, password, dbName)
	return container
}

type PostgresWaitStrategy struct {
	User     string
	Password string
	Database string
}

func (strategy *PostgresWaitStrategy) WaitUntilReady(ctx context.Context, target wait.StrategyTarget) (err error) {
	ip, _ := target.Host(ctx)
	port, _ := target.MappedPort(ctx, "5432")
	host := fmt.Sprintf("%s:%s", ip, port.Port())
	db := pg.Connect(&pg.Options{
		User:     strategy.User,
		Password: strategy.Password,
		Database: strategy.Database,
		Addr:     host,
	})

	attempt := 0
	for attempt < 60 {
		attempt++
		log.Printf("Checking Postgres Up&Running on host [%s]... Attempt #[%d]", host, attempt)
		err := CheckPostgresConnection(db)
		if err == nil {
			log.Println("Successfully connected to Postgres DB.")
			return nil
		}
		log.Println(err)
		time.Sleep(time.Second)
	}
	return errors.New(fmt.Sprintf("failed to connect to Postrgres container using host[%s] ", host))
}

func CheckPostgresConnection(db *pg.DB) error {
	conn := db.Conn()
	defer conn.Close()
	_, err := conn.Exec("SELECT 1")
	return err
}
