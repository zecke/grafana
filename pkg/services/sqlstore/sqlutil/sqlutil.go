package sqlutil

type TestDB struct {
	DriverName string
	ConnStr    string
}

var TestDB_Sqlite3 = TestDB{DriverName: "sqlite3", ConnStr: ":memory:"}
var TestDB_Mssql = TestDB{DriverName: "mssql", ConnStr: "server=localhost;port=1433;database=grafanatest;user id=grafana;password=Password!"}
