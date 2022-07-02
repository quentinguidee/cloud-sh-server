package database

func (db *Database) CreateBucketsTable() {
	_, _ = db.Instance.Exec(`
		CREATE TABLE IF NOT EXISTS buckets (
			id        INTEGER UNIQUE PRIMARY KEY AUTOINCREMENT,
			name      VARCHAR(255),
			root_node VARCHAR(63),
			type      VARCHAR(63),
			FOREIGN KEY(root_node) REFERENCES buckets_nodes(uuid)
		)
	`)

	_, _ = db.Instance.Exec(`
		CREATE TABLE IF NOT EXISTS buckets_access (
		    id        INTEGER UNIQUE PRIMARY KEY AUTOINCREMENT,
		    bucket_id INTEGER,
		    user_id   INTEGER,
		    access_type VARCHAR(63),
		    FOREIGN KEY(bucket_id) REFERENCES buckets(id),
		    FOREIGN KEY(user_id)   REFERENCES users(id)
		)
	`)

	_, _ = db.Instance.Exec(`
		CREATE TABLE IF NOT EXISTS buckets_nodes (
		    uuid      VARCHAR(63) UNIQUE PRIMARY KEY,
		    name      VARCHAR(255),
		    type      VARCHAR(63),
		    bucket_id INTEGER,
		    FOREIGN KEY(bucket_id) REFERENCES buckets(id)
		)
	`)

	_, _ = db.Instance.Exec(`
		CREATE TABLE IF NOT EXISTS buckets_nodes_associations (
		    id        INTEGER UNIQUE PRIMARY KEY AUTOINCREMENT,
		    from_node VARCHAR(63),
		    to_node   VARCHAR(63),
		    FOREIGN KEY(from_node) REFERENCES buckets_nodes(uuid),
		    FOREIGN KEY(to_node)   REFERENCES buckets_nodes(uuid)
		)
	`)
}
