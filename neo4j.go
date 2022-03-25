package neo4j

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"go.k6.io/k6/js/modules"
)

// Register the extension on module initialization, available to
// import from JS as "k6/x/neo4j".
func init() {
	modules.Register("k6/x/neo4j", new(Neo4j))
}

// Neo4j is the k6 extension for a Neo4j driver.
type Neo4j struct{}

// Driver is the neo4j driver wrapper.
type Driver struct {
	driver neo4j.Driver
}

type RunCypherResult struct {
	Result neo4j.Result `js:"result"`
	Error  error        `js:"error"`
}

type Config struct {
	Address                 string        `js:"address"`
	User                    string        `js:"user"`
	HTTPPort                int           `js:"httpPort"`
	HTTPSPort               int           `js:"httpsPort"`
	Password                string        `js:"password"`
	MaxTransactionRetryTime time.Duration `js:"maxTransactionRetryTime"`
}

// XDriver represents the Driver constructor (i.e. `neo4.j.NewDriver()`) and
func (*Neo4j) XDriver(ctxPtr *context.Context, cfg Config) interface{} {
	dbURI := buildDbURI(cfg.Address, cfg.HTTPPort, cfg.HTTPSPort)

	driver, err := neo4j.NewDriver(dbURI, neo4j.BasicAuth(cfg.User, cfg.Password, ""), func(config *neo4j.Config) {
		config.MaxTransactionRetryTime = cfg.MaxTransactionRetryTime
	})

	if err != nil {
		panic(fmt.Errorf("failed to create driver for Neo4j: %w", err))
	}

	return &Driver{driver: driver}
}

func (d *Driver) RunCypherInSession(cypher string) RunCypherResult {
	runResult := RunCypherResult{}
	session := d.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func() {
		if err := session.Close(); err != nil {
			log.Fatal(fmt.Errorf("could not close session: %w", err))
		}
	}()

	runResult.Result, runResult.Error = session.Run(cypher, map[string]interface{}{})
	if runResult.Error != nil {
		log.Fatal(fmt.Errorf("could not close session: %w", runResult.Error))
	}
	return runResult
}

// buildDBUri
func buildDbURI(addr string, httpPort, httpsPort int) string {
	if httpsPort != 0 {
		return fmt.Sprintf("neo4j+s://%s:%s/", addr, strconv.Itoa(httpsPort))
	}

	return fmt.Sprintf("neo4j://%s:%s/", addr, strconv.Itoa(httpPort))
}
