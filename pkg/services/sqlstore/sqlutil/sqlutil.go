package sqlutil

type TestDB struct {
	DriverName string
	ConnStr    string
}

var TestDB_Sqlite3 = TestDB{DriverName: "sqlite3", ConnStr: ":memory:"}
var TestDB_Postgres = TestDB{DriverName: "postgres", ConnStr: "user=grafanatest password=grafanatest host=localhost port=5432 dbname=grafanatest sslmode=disable"}
var TestDB_Mssql = TestDB{DriverName: "mssql", ConnStr: "server=localhost;port=1433;database=grafanatest;user id=grafana;password=Password!"}
