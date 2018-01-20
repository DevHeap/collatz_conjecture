package cache

import (
	"os"
	"database/sql"
	"log"
	_ "github.com/lib/pq"
	"math/big"
	"encoding/json"
	"strings"
)

type Database struct {
	db *sql.DB
	loadStmt   *sql.Stmt
	insertStmt *sql.Stmt
	deleteStmt *sql.Stmt
}

type DatabaseCacheEntry struct {
	number *big.Int
	path   []*big.Int
}

func NewDatabase() (*Database, error) {
	log.Print("Opening connection to DB")
	connUri := strings.Trim(os.Getenv("POSTGRES_URI"), " ")
	db, err := sql.Open("postgres", connUri)

	if err != nil {
		log.Fatal("Couldn't open the connection to database: ", err)
		return nil, err
	}

	log.Print("Establishing connection to DB")
	err = db.Ping()

	if err != nil {
		log.Fatal("Couldn't establish the connection to database: ", err)
		return nil, err
	}
	
	log.Print("Creating schema")
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS \"cache\" (\"number\" varchar(1024), \"path\" text)")

	if err != nil {
		log.Fatal("Failed to create table: ", err)
		return nil, err
	}

	log.Print("Preparing statements")
	loadStmt, err := db.Prepare("SELECT * FROM cache")

	if err != nil {
		log.Fatal("Failed to prepare load statement: ", err)
		return nil, err
	}

	insertStmt, err := db.Prepare("INSERT INTO \"cache\" VALUES($1, $2)")

	if err != nil {
		log.Fatal("Failed to prepare insert statement: ", err)
		return nil, err
	}

	deleteStmt, err := db.Prepare("DELETE FROM \"cache\" WHERE \"number\" = $1")

	if err != nil {
		log.Fatal("Failed to prepare delete statement: ", err)
		return nil, err
	}

	log.Print("Database initialized successfuly")
	return &Database {
		db,
		loadStmt,
		insertStmt,
		deleteStmt,
	}, nil
}

func (d *Database) Load() (map[*big.Int][]*big.Int, error) {
	rows, err := d.loadStmt.Query()
	if err != nil {
		log.Fatal("Failed to load cache from the database: ", err)
		return nil, err
	}

	cache := make(map[*big.Int][]*big.Int)

	for rows.Next() {
		var number string
		var path []byte
		err = rows.Scan(&number, &path)
		if err != nil {
			log.Fatal("Failed to scan from DB rows, possibly faulty layout!")
			return nil, err
		}
		num, success := new(big.Int).SetString(number, 10)
		if !success {
			log.Fatal("Failed to parse ", number, " into a big.Int")
			continue
		}
		p, err := deserializeBigints(path)
		if err != nil {
			log.Fatal("Failed to parse path from database")
			continue
		}
		cache[num] = p
	}

	return cache, nil
}

func (d *Database) InsertAsync(number *big.Int, path []*big.Int) {
	go d.Insert(number, path)
}

func (d *Database) Insert(number *big.Int, path []*big.Int) error {
	key := number.String()
	data, err := serializeBigints(path)

	if err != nil {
		return err
	}

	_, err = d.insertStmt.Exec(key, data)

	if err != nil {
		log.Fatal("Failed to insert path for ", key, ": ", err)
		return err
	}

	log.Print("Inserted path for ", key, " into the database")

	return nil
}


func (d *Database) DeleteByKeyAsync(key string) {
	go d.DeleteByKey(key)
}

func (d *Database) DeleteByKey(key string) error {
	_, err := d.deleteStmt.Exec(key)

	if err != nil {
		log.Fatal("Failed to delete path for ", key, ": ", err)
		return err
	}

	return nil
}

func (d *Database) DeleteAsync(number *big.Int) {
	go d.Delete(number)
}

func (d *Database) Delete(number *big.Int) error {
	key := number.String()

	return d.DeleteByKey(key)
}

func serializeBigints(bigintArray []*big.Int) ([]byte, error) {
	var path []string

	for i := 0; i < len(bigintArray); i++ {
		path = append(path, bigintArray[i].String())
	}

	jsonPath, err := json.Marshal(path)

	if err != nil {
		log.Fatal("Failed to marshal path")
		return nil, err
	}

	return jsonPath, nil
}

func deserializeBigints(serialized []byte) ([]*big.Int, error) {
	var strPath []string

	err := json.Unmarshal(serialized, &strPath)
	if err != nil {
		log.Fatal("Failed to unmarshal path: ", err)
		return nil, err
	}

	var path []*big.Int

	for i := 0; i < len(strPath); i++ {
		bigInt, err := new(big.Int).SetString(strPath[i], 10)
		if !err {
			log.Fatal("Couldn't parse big.Int received from database")
			// TODO custom error type
			return nil, nil
		}
		path = append(path, bigInt)
	}

	return path, nil
}
