package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"time"
)

func QueryWithRetry(db *sql.DB, ctx context.Context, query string, args ...any) (sql.Result, error) {
	retries := 0

RETRY:
	time.Sleep(time.Duration(retries) * time.Second)

	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		if regexp.MustCompile("SQL logic error: no such table").MatchString(err.Error()) {
			if retries < 3 {
				retries++

				log.Printf("models.QueryWithRetry.retry: %d\n", retries)

				goto RETRY
			}
		}
		return nil, fmt.Errorf("models.QueryWithRetry.1: %s", err.Error())
	}

	return result, nil
}
