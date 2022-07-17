package database

func (db *Database) Update() error {
	var version int
	err := db.Instance.QueryRowx("SELECT database_version FROM servers WHERE id = 1").Scan(&version)
	if err != nil {
		return err
	}

	var request string

	//if version < 2 {
	//	request += "...;"
	//}

	request += "UPDATE servers SET database_version = $1 WHERE id = 1;"

	_, err = db.Instance.Exec(request, DatabaseVersion)
	return err
}
