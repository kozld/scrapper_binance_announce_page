package database

const CreateTableQuery = `CREATE TABLE IF NOT EXISTS binance_page (
	hash bytea UNIQUE NOT NULL,
	text text NOT NULL,
	time timestamp with time zone default current_timestamp,
	PRIMARY KEY (hash)
);`

const InsertQuery = `INSERT INTO binance_page (hash, text) VALUES ($1, $2);`

const SelectQuery = `SELECT hash FROM binance_page WHERE hash = $1;`
