package database

import (
	"fmt"
	"self-hosted-cloud/server/models/storage"
	"strings"
)

func (db *Database) CreateBucketsTable() {
	_, _ = db.instance.Exec(`
		CREATE TABLE IF NOT EXISTS buckets (
			id        INTEGER UNIQUE PRIMARY KEY AUTOINCREMENT,
			name      VARCHAR(255),
			root_node INTEGER,
			type      VARCHAR(63),
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

func (db *Database) GetUserBucket(userId int) (storage.Bucket, error) {
	request := `
		SELECT buckets.id, buckets.name, buckets.root_node, buckets.type
		FROM buckets, buckets_access access
		WHERE buckets.id = access.bucket_id
		  AND buckets.type = 'user_bucket'
		  AND access.user_id = ?
	`

	var bucket storage.Bucket
	err := db.instance.QueryRow(request, userId).Scan(
		&bucket.Id,
		&bucket.Name,
		&bucket.RootNode,
		&bucket.Type)

	if err != nil {
		return storage.Bucket{}, err
	}

	return bucket, nil
}

func (db *Database) GetBucket(bucketId int) (storage.Bucket, error) {
	request := "SELECT id, name, root_node, type FROM buckets WHERE buckets.id = ?"

	var bucket storage.Bucket
	err := db.instance.QueryRow(request, bucketId).Scan(&bucket.Id, &bucket.Name, &bucket.RootNode, &bucket.Type)
	if err != nil {
		return storage.Bucket{}, err
	}
	return bucket, nil
}

func (db *Database) GetNodesFromNode(fromNode int) ([]storage.Node, error) {
	request := `
		SELECT nodes.id, nodes.filename, nodes.filetype
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
			&node.Filetype)

		if err != nil {
			return nil, err
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}

func (db *Database) GetNodeFromNode(fromNode int, filename string) (storage.Node, error) {
	request := `
		SELECT nodes.id, nodes.filename, nodes.filetype
		FROM buckets_nodes nodes, buckets_nodes_associations associations
		WHERE associations.from_node = ?
		  AND associations.to_node = nodes.id
		  AND nodes.filename = ?
	`

	var node storage.Node
	err := db.instance.QueryRow(request, fromNode, filename).Scan(
		&node.Id,
		&node.Filename,
		&node.Filetype)

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
	if err != nil {
		return nil, err
	}

	nodes, err := db.GetNodesFromNode(node.Id)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func (db *Database) CreateBucket(userId int) (storage.Bucket, error) {
	request := "INSERT INTO buckets(name, type) VALUES (?, 'user_bucket') RETURNING id"
	bucket := storage.Bucket{
		Name: fmt.Sprintf("Main bucket"),
	}
	err := db.instance.QueryRow(request, bucket.Name).Scan(&bucket.Id)
	if err != nil {
		return storage.Bucket{}, err
	}

	node := storage.Node{
		Filename: "/",
		Filetype: "directory",
	}
	request = `
		INSERT INTO buckets_nodes(filename, filetype, bucket_id)
		VALUES ('/', 'directory', ?)
		RETURNING id
	`

	err = db.instance.QueryRow(request, bucket.Id).Scan(&node.Id)
	if err != nil {
		return storage.Bucket{}, err
	}

	request = `
		UPDATE buckets
		SET root_node = ?
		WHERE id = ?
	`

	_, err = db.instance.Exec(request, node.Id, bucket.Id)
	if err != nil {
		return storage.Bucket{}, err
	}

	request = `
		INSERT INTO buckets_access(bucket_id, user_id, access_type)
		VALUES (?, ?, 'admin')
	`

	_, err = db.instance.Exec(request, bucket.Id, userId)
	if err != nil {
		return storage.Bucket{}, err
	}

	return bucket, nil
}

func (db *Database) CreateNode(directoryId int, node storage.Node) error {
	request := `
		INSERT INTO buckets_nodes(filename, filetype, bucket_id)
		VALUES (?, ?, ?)
		RETURNING id
	`

	err := db.instance.QueryRow(request,
		node.Filename,
		node.Filetype,
		node.BucketId,
	).Scan(&node.Id)

	if err != nil {
		return err
	}

	request = `
		INSERT INTO buckets_nodes_associations(from_node, to_node)
		VALUES (?, ?)
	`

	_, err = db.instance.Exec(request, directoryId, node.Id)
	if err != nil {
		return err
	}

	return nil
}
