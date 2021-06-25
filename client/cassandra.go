package client

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gocql/gocql"
)

var session *gocql.Session

func init() {
	// env.Load()

	cluster := gocql.NewCluster(os.Getenv("CASSANDRA_URL"))
	cluster.Keyspace = os.Getenv("CASSANDRA_KEYSPACE")
	cluster.Consistency = gocql.Quorum
	cluster.ConnectTimeout = time.Second * 20
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: os.Getenv("CASSANDRA_USERNAME"),
		Password: os.Getenv("CASSANDRA_PASSWORD"),
	}

	s, err := cluster.CreateSession()
	if err != nil {
		log.Fatalln("err cluster:", err)
	}
	// defer s.Close()
	fmt.Println("Cassandra cluster is ready!")
	session = s
}

func CassandraSession() *gocql.Session {
	return session
}
