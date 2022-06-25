package database

import (
	"self-hosted-cloud/server/models/storage"
	"strings"
)

func (db *Database) CreateBucketsTable() {
	_, _ = db.instance.Exec(`
		CREATE TABLE IF NOT EXISTS buckets (
			id        INTEGER UNIQUE PRIMARY KEY AUTOINCREMENT,
			name      VARCHAR(255),
			root_node INTEGER,
			FOREIGN KEY(root_node) REFERENCES buckets_nodes(id)
		)
	`)

	_, _ = db.instance.Exec(`
		CREATE TABLE IF NOT EXISTS buckets_access (
		    id          INTEGER UNIQUE PRIMARY KEY AUTOINCREMENT,
		    bucket_id   INTEGER,
		    user_id     INTEGER,
		    access_type VARCHAR(63),
		    FOREIGN KEY(bucket_id) REFERENCES buckets(id),
		    FOREIGN KEY(user_id)   REFERENCES users(id)
		)
	`)

	_, _ = db.instance.Exec(`
		CREATE TABLE IF NOT EXISTS buckets_nodes (
		    id                INTEGER UNIQUE PRIMARY KEY AUTOINCREMENT,
		    filename          VARCHAR(255),
		    filetype          VARCHAR(63),
		    internal_filename VARCHAR(255),
		    bucket_id         INTEGER,
		    FOREIGN KEY(bucket_id) REFERENCES buckets(id)
		)
	`)

	_, _ = db.instance.Exec(`
		CREATE TABLE IF NOT EXISTS buckets_nodes_associations (
		    id        INTEGER UNIQUE PRIMARY KEY AUTOINCREMENT,
		    from_node INTEGER,
		    to_node   INTEGER(255),
		    FOREIGN KEY(from_node) REFERENCES buckets_nodes(id),
		    FOREIGN KEY(to_node)   REFERENCES buckets_nodes(id)
		)
	`)
}

func (db *Database) GetBucket(bucketId string) (storage.Bucket, error) {
	request := "SELECT id, name, root_node FROM buckets WHERE buckets.id = ?"

	var bucket storage.Bucket
	err := db.instance.QueryRow(request, bucketId).Scan(&bucket.Id, &bucket.Name, &bucket.RootNode)
	if err != nil {
		return storage.Bucket{}, err
	}
	return bucket, nil
}

func (db *Database) GetNodes(fromNode int) ([]storage.Node, error) {
	request := `
		SELECT nodes.id, nodes.filename, nodes.filetype, nodes.internal_filename
		FROM buckets_nodes nodes, buckets_nodes_associations associations
		WHERE associations.from_node = ?
          AND associations.to_node = nodes.id
	`

	rows, err := db.instance.Query(request, fromNode)
	if err != nil {
		return nil, err
	}

	var nodes []storage.Node
	for rows.Next() {
		var node storage.Node
		err := rows.Scan(
			&node.Id,
			&node.Filename,
			&node.Filetype,
			&node.InternalFilename)

		if err != nil {
			return nil, err
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}

func (db *Database) GetNode(fromNode int, filename string) (storage.Node, error) {
	request := `
		SELECT nodes.id, nodes.filename, nodes.filetype, nodes.internal_filename
		FROM buckets_nodes nodes, buckets_nodes_associations associations
		WHERE associations.from_node = ?
		  AND associations.to_node = nodes.id
		  AND nodes.filename = ?
	`

	var node storage.Node
	err := db.instance.QueryRow(request, fromNode, filename).Scan(
		&node.Id,
		&node.Filename,
		&node.Filetype,
		&node.InternalFilename)

	if err != nil {
		return storage.Node{}, err
	}

	return node, nil
}

func (db *Database) GetFiles(bucketId string, path string) ([]storage.Node, error) {
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	bucket, err := db.GetBucket(bucketId)
	if err != nil {
		return nil, err
	}

	node := bucket.RootNode
	for _, filename := range strings.Split(path, "/") {
		node, err := db.GetNode(node, filename)
		if err != nil {
			return nil, err
		}
		if node.Filetype != "directory" {
			return []storage.Node{node}, nil
		}
	}

	nodes, err := db.GetNodes(node)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}
