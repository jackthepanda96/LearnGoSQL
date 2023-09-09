package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type DBConfig struct {
	Username string
	Password string
	Host     string
	Table    string
	Port     int
}

func loadData() *DBConfig {
	var res = new(DBConfig)
	if val, found := os.LookupEnv("DB_PORT"); found {
		cnv, _ := strconv.Atoi(val)
		res.Port = cnv
	}

	if val, found := os.LookupEnv("DB_USER"); found {
		res.Username = val
	}

	if val, found := os.LookupEnv("DB_PASSWORD"); found {
		res.Password = val
	}

	if val, found := os.LookupEnv("DB_TABLE"); found {
		res.Table = val
	}

	if val, found := os.LookupEnv("DB_HOST"); found {
		res.Host = val
	}

	return res
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Cannot read .env file. ", err.Error())
	}

	var configs = loadData()
	fmt.Println(configs)

	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", configs.Username, configs.Password, configs.Host, configs.Port, configs.Table)

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		fmt.Println("Some error occured-1", err.Error())
	}
	defer db.Close()

	var qry1 = "SELECT id, nama FROM user"
	prpQry, err := db.Prepare(qry1)

	if err != nil {
		fmt.Println("Prepare query error", err.Error())
	}

	defer prpQry.Close()

	rows, err := prpQry.Query()

	if err != nil {
		fmt.Println("Retrieve Query result error", err.Error())
	}

	var tmp = []any{}
	for rows.Next() {
		type tmpScan struct {
			Id   int
			Nama string
		}
		var scanTmp tmpScan
		err = rows.Scan(&scanTmp.Id, &scanTmp.Nama)
		if err != nil {
			fmt.Println(err)
		}
		tmp = append(tmp, scanTmp)
	}

	fmt.Println(tmp)

	// res, err := db.Exec("INSERT INTO user (nama, hp, password) values (?, ?, ?)", "nurul", "123", "nurul123")
	// if err != nil {
	// 	fmt.Println("Insert execute error ", err.Error())
	// }

	res, err := db.Exec("UPDATE user SET nama=? WHERE id= ?", 1, "joko")

	if err != nil {
		fmt.Println("Insert execute error ", err.Error())
	}

	if count, _ := res.RowsAffected(); count == 0 {
		fmt.Println("Query failed")
	}

	// db.Exec() -> execute (INSERT UDPATE DELETE | CREATE TABLE | DROP TABLE | SELECT)
	// db.Query() -> SELECT yang MENGAHSIL BANYAK JAWABAN
	// db.QueryRow() -> SELECT yang MENGHASILKAN 1 BARIS JAWABAN
}
