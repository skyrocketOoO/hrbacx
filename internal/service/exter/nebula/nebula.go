package nebulaservice

import (
	"fmt"
	"log"
	"time"

	nebula "github.com/vesoft-inc/nebula-go/v3"
)

const (
	SPACE    = "my_space"
	Address  = "127.0.0.1"
	Port     = 9669
	Username = "root"
	Password = "nebula"
	UseHTTP2 = false
)

var SessionPool *nebula.SessionPool

const BasicSchema = `
	CREATE TAG IF NOT EXISTS user();
	CREATE TAG IF NOT EXISTS role();
	CREATE TAG IF NOT EXISTS obj();

	CREATE EDGE IF NOT EXISTS has_permission(type STRING NOT NULL);
	CREATE EDGE IF NOT EXISTS belongs_to();
	CREATE EDGE IF NOT EXISTS leader_of();
`

func New() error {
	prepareSpace()
	NewSessionPool()
	// defer SessionPool.Close()

	if _, err := Exec(BasicSchema); err != nil {
		log.Fatal(err)
	}
	time.Sleep(5 * time.Second)

	return nil
}

func CheckResultSet(prefix string, res *nebula.ResultSet) {
	if !res.IsSucceed() {
		log.Fatal(fmt.Sprintf("%s, ErrorCode: %v, ErrorMsg: %s", prefix, res.GetErrorCode(), res.GetErrorMsg()))
	}
}

func NewSessionPool() {
	hostAddress := nebula.HostAddress{Host: Address, Port: Port}

	// Create configs for session pool
	config, err := nebula.NewSessionPoolConf(
		"root",
		"nebula",
		[]nebula.HostAddress{hostAddress},
		SPACE,
		nebula.WithHTTP2(UseHTTP2),
	)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to create session pool config, %s", err.Error()))
	}

	// create session pool
	SessionPool, err = nebula.NewSessionPool(*config, nebula.DefaultLogger{})
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to initialize session pool, %s", err.Error()))
	}
}

// Just a helper function to create a space for this example to run.
func prepareSpace() {
	hostAddress := nebula.HostAddress{Host: Address, Port: Port}
	hostList := []nebula.HostAddress{hostAddress}
	// Create configs for connection pool using default values
	testPoolConfig := nebula.GetDefaultConf()
	testPoolConfig.UseHTTP2 = UseHTTP2

	// Initialize connection pool
	pool, err := nebula.NewConnectionPool(hostList, testPoolConfig, nebula.DefaultLogger{})
	if err != nil {
		log.Fatal(fmt.Sprintf("Fail to initialize the connection pool, host: %s, port: %d, %s",
			Address, Port, err.Error()))
	}
	// Close all connections in the pool
	defer pool.Close()

	// Create session
	session, err := pool.GetSession(Username, Password)
	if err != nil {
		log.Fatal(
			fmt.Sprintf("Fail to create a new session from connection pool, username: %s, password: %s, %s",
				Username, Password, err.Error()))
	}
	// Release session and return connection back to connection pool
	defer session.Release()

	checkResultSet := func(prefix string, res *nebula.ResultSet) {
		if !res.IsSucceed() {
			log.Fatal(
				fmt.Sprintf("%s, ErrorCode: %v, ErrorMsg: %s",
					prefix, res.GetErrorCode(), res.GetErrorMsg()))
		}
	}

	{
		// Prepare the query
		createSchema := fmt.Sprintf(`
			CREATE SPACE IF NOT EXISTS %s (vid_type=FIXED_STRING(20)); 
		`, SPACE)

		// Execute a query
		resultSet, err := session.Execute(createSchema)
		if err != nil {
			fmt.Print(err.Error())
			return
		}
		checkResultSet(createSchema, resultSet)
	}

	log.Println("Space my_space was created")
	time.Sleep(5 * time.Second)
}

func Exec(schema string) (res *nebula.ResultSet, err error) {
	res, err = SessionPool.Execute(schema)
	CheckResultSet(schema, res)
	return
}
