package database

func (db *Database) CreateBucketsTable() {
	_, _ = db.Instance.Exec(`
		CREATE TABLE IF NOT EXISTS buckets (
			id        INTEGER UNIQUE PRIMARY KEY AUTOINCREMENT,
			name      VARCHAR(255),
			root_node INTEGER,
			type      VARCHAR(63),
			FOREIGN KEY(root_node) REFERENCES buckets_nodes(id)
		)
	`)

	_, _ = db.Instance.Exec(`
		CREATE TABLE IF NOT EXISTS buckets_access (
		    id          INTEGER UNIQUE PRIMARY KEY AUTOINCREMENT,
		    bucket_id   INTEGER,
		    user_id     INTEGER,
		    access_type VARCHAR(63),
		    FOREIGN KEY(bucket_id) REFERENCES buckets(id),
		    FOREIGN KEY(user_id)   REFERENCES users(id)
		)
	`)

	_, _ = db.Instance.Exec(`
		CREATE TABLE IF NOT EXISTS buckets_nodes (
		    id                INTEGER UNIQUE PRIMARY KEY AUTOINCREMENT,
		    filename          VARCHAR(255),
		    filetype          VARCHAR(63),
		    bucket_id         INTEGER,
		    FOREIGN KEY(bucket_id) REFERENCES buckets(id)
		)
	`)

	_, _ = db.Instance.Exec(`
		CREATE TABLE IF NOT EXISTS buckets_nodes_associations (
		    id        INTEGER UNIQUE PRIMARY KEY AUTOINCREMENT,
		    from_node INTEGER,
		    to_node   INTEGER(255),
		    FOREIGN KEY(from_node) REFERENCES buckets_nodes(id),
		    FOREIGN KEY(to_node)   REFERENCES buckets_nodes(id)
		)
	`)
}
