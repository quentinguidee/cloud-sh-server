package database

import (
	"database/sql"
	"self-hosted-cloud/server/models/storage"
	"strings"
)

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

func (db *Database) GetNodesFromNode(fromNode int) ([]storage.Node, error) {
	request := `
		SELECT nodes.id, nodes.filename, nodes.filetype, nodes.bucket_id
		FROM buckets_nodes nodes, buckets_nodes_associations associations
		WHERE associations.from_node = ?
          AND associations.to_node = nodes.id
	`

	rows, err := db.Instance.Query(request, fromNode)
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
			&node.BucketId)

		if err != nil {
			return nil, err
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}

func (db *Database) GetNodeFromNode(fromNode int, filename string) (storage.Node, error) {
	request := `
		SELECT nodes.id, nodes.filename, nodes.filetype, nodes.bucket_id
		FROM buckets_nodes nodes, buckets_nodes_associations associations
		WHERE associations.from_node = ?
		  AND associations.to_node = nodes.id
		  AND nodes.filename = ?
	`

	var node storage.Node
	err := db.Instance.QueryRow(request, fromNode, filename).Scan(
		&node.Id,
		&node.Filename,
		&node.Filetype,
		&node.BucketId)

	if err != nil {
		return storage.Node{}, err
	}

	return node, nil
}

func (db *Database) GetNode(bucket storage.Bucket, path string) (storage.Node, error) {
	node := storage.Node{
		Id: bucket.RootNode,
	}

	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	if len(path) == 0 {
		return node, nil
	}

	for _, filename := range strings.Split(path, "/") {
		var err error
		node, err = db.GetNodeFromNode(node.Id, filename)
		if err != nil {
			return storage.Node{}, err
		}
		if node.Filetype != "directory" {
			return node, nil
		}
	}

	return node, nil
}

func (db *Database) GetFiles(bucket storage.Bucket, path string) ([]storage.Node, error) {
	node, err := db.GetNode(bucket, path)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	nodes, err := db.GetNodesFromNode(node.Id)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}
