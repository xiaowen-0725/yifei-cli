package db

import (
	"database/sql"
	"fmt"

	_ "github.com/microsoft/go-mssqldb"
)

type DB struct{ conn *sql.DB }

func Open(dsn string) (*DB, error) {
	conn, err := sql.Open("sqlserver", dsn)
	if err != nil {
		return nil, fmt.Errorf("打开数据库失败: %w", err)
	}
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("数据库连接失败 (检查 config): %w", err)
	}
	return &DB{conn: conn}, nil
}

func (d *DB) Close() error { return d.conn.Close() }

func (d *DB) Query(query string, limit int) ([]string, [][]any, error) {
	rows, err := d.conn.Query(query)
	if err != nil {
		return nil, nil, fmt.Errorf("SQL 执行失败: %w", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}
	var out [][]any
	for rows.Next() {
		if limit > 0 && len(out) >= limit {
			break
		}
		scan := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range scan {
			ptrs[i] = &scan[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, nil, err
		}
		out = append(out, scan)
	}
	return cols, out, rows.Err()
}
