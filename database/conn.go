package database

import (
	"beli-mang/config"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

func NewDatabase(cfg *config.Config) (*sqlx.DB, error) {
	dbd, err := sql.Open("postgres", cfg.DB.ConnectionString())
	if err != nil {
		return nil, err
	}

	/*
				MaxIdleConns: This parameter determines the maximum number of idle (unused) connections that the connection pool maintains.
			Set it based on the expected idle connections during normal operation. A value between 2 to 4 times the number of CPU cores on your EC2 instance is often a good starting point.
		For t3a.medium with 2 vCPUs, consider setting MaxIdleConns to 4 to 8.
	*/
	dbd.SetMaxIdleConns(10) // Set maximum idle connections
	/*
			MaxOpenConns: This parameter defines the maximum number of open connections allowed in the pool simultaneously.
		Set it based on the anticipated concurrency of your API endpoints and the capacity of your RDS instance.
		For a target of 500 RPS, you may need to experiment with different values to find the optimal setting.
		Start with a conservative value and gradually increase it while monitoring the database and application server metrics for signs of overload or contention.
		A value between 50 to 100 could be a reasonable starting point, but you may need to adjust it based on performance testing and monitoring results.
	*/
	dbd.SetMaxOpenConns(100) // Set maximum open connections
	db := sqlx.NewDb(dbd, "postgres")
	return db, err
}
