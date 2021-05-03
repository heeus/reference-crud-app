/*
 * Copyright (c) 2019-present Heeus authors
 */

package service

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gocql/gocql"
)

//CasandraPartitionedDriver s.e.
type CasandraPartitionedDriver struct {
	cluster *gocql.ClusterConfig
	session *gocql.Session

	args map[string]string

	hosts             []string
	keyspace          string
	class             string
	consistency       gocql.Consistency
	replicationFactor int64
	lightWeight       int64

	logger *Logger
}

//Name s.e.
func (d *CasandraPartitionedDriver) Name() string {
	return "Cassandra sharded driver"
}

//Info s.e.
func (d *CasandraPartitionedDriver) Info() string {
	str := "Casandra sharded driver info: \n\n"

	str += fmt.Sprintf("Hosts: %v\n", d.hosts)
	str += fmt.Sprintf("Keyspace name: %v\n", d.keyspace)
	str += fmt.Sprintf("Class: %v\n", d.class)
	str += fmt.Sprintf("Consistency: %v\n", d.consistency)
	str += fmt.Sprintf("Replication factor: %v\n", d.replicationFactor)
	str += fmt.Sprintf("lightWeight mode: %v\n", d.lightWeight)

	str += "\n\n --- end --- \n\n"

	return str
}

//Init s.e.
func (d *CasandraPartitionedDriver) Init(args map[string]string) error {
	d.args = args

	err := d.initParams()

	if err != nil {
		return err
	}

	cluster := gocql.NewCluster(d.hosts...)
	cluster.Keyspace = "system"

	session, err := cluster.CreateSession()
	defer session.Close()

	if err != nil {
		d.logger.Error(err.Error())
		return err
	}

	q := fmt.Sprintf("CREATE KEYSPACE IF NOT EXISTS %v WITH replication = {'class': '%s', 'replication_factor' : %v}", d.keyspace, d.class, d.replicationFactor)
	reconnectCount := 0

	for {
		if err = session.Query(q).Exec(); err != nil {
			time.Sleep(500 * time.Millisecond)
		} else {
			break
		}

		reconnectCount++

		if reconnectCount >= MaxReconnectCount {
			break
		}
	}

	if err != nil {
		d.logger.Error(err.Error())
		return err
	}

	d.cluster = gocql.NewCluster(d.hosts...)
	d.cluster.Consistency = d.consistency
	d.cluster.Keyspace = d.keyspace

	d.session, err = d.cluster.CreateSession()

	if err != nil {
		d.logger.Error(err.Error())
		return err
	}

	q = "CREATE TABLE IF NOT EXISTS records_p ( key text, partition int, version int, values blob, type text, weight int, PRIMARY KEY (partition, key))"
	reconnectCount = 0

	for {
		if err := d.session.Query(q).Exec(); err != nil {
			time.Sleep(500 * time.Millisecond)
		} else {
			break
		}

		reconnectCount++

		if reconnectCount >= MaxReconnectCount {
			break
		}
	}

	if err != nil {
		d.logger.Error(err.Error())
		return err
	}

	d.logger.Log("Casandra sharded driver initialized")

	d.logger.Debug("Cassandra hosts: %v", d.hosts)
	d.logger.Debug("Cassandra replication factor: %v", d.replicationFactor)
	d.logger.Debug("Cassandra keyspace: %v", d.keyspace)
	d.logger.Debug("Cassandra class: %v", d.class)
	d.logger.Debug("Light weight transaction mode: %v", d.lightWeight)

	return nil
}

//Free s.e.
func (d *CasandraPartitionedDriver) Free() error {
	d.logger.Log("Casandra sharded driver freed")

	d.session.Close()

	return nil
}

//Clean s.e.
func (d *CasandraPartitionedDriver) Clean(r *DBRequest) *DBResponse {
	if err := d.session.Query(`TRUNCATE records_p;`).Exec(); err != nil {
		return &DBResponse{Error: err.Error()}
	}

	return &DBResponse{Status: 200}
}

//Read s.e.
func (d *CasandraPartitionedDriver) Read(r *DBRequest) *DBResponse {
	var records []*Record

	if r == nil {
		return &DBResponse{Status: 400, Error: "wrong request data"}
	}

	if len(r.ViewViews) > 0 {
		records = make([]*Record, len(r.ViewViews))

		for i, v := range r.ViewViews {
			rec, err := d.read(r.Partition, &v)

			if err != nil {
				return &DBResponse{Status: 400, Error: err.Error()}
			}

			if rec != nil {
				records[i] = rec
			}
		}
	}

	return &DBResponse{Status: 200, Records: records}
}

func (d *CasandraPartitionedDriver) read(partition int64, view *ViewView) (*Record, error) {
	var err error

	if view.ViewType == "" {
		return nil, fmt.Errorf("record ViewType malformed")
	}

	key, err := buildKey(view.PartitionKey, view.ClusterKey)

	if err != nil {
		return nil, err
	}

	r, err := d.get(key, partition, view.ViewType)

	if err != nil {
		return nil, err
	}

	return r, nil
}

//Insert s.e.
func (d *CasandraPartitionedDriver) Insert(r *DBRequest) *DBResponse {
	if r == nil {
		return &DBResponse{Status: 400, Error: "wrong request data"}
	}

	if len(r.ViewMods) > 0 {
		for _, v := range r.ViewMods {
			err := d.insert(r.Partition, &v)

			if err != nil {
				return &DBResponse{Status: 400, Error: fmt.Sprintf("Insert error: %v", err.Error())}
			}
		}
	}

	return &DBResponse{Status: 200}
}

func (d *CasandraPartitionedDriver) insert(partition int64, view *ViewMod) error {
	d.logger.Debug("insert request: %v", view)

	if view.ViewType == "" {
		return fmt.Errorf("record ViewType name malformed")
	}

	key, err := buildKey(view.PartitionKey, view.ClusterKey)

	if err != nil {
		return err
	}

	return d.set(key, partition, view.ViewType, view.Values)
}

//Update s.e.
func (d *CasandraPartitionedDriver) Update(r *DBRequest) *DBResponse {
	var err error
	var key string

	if r == nil {
		return &DBResponse{Status: 400, Error: "wrong request data"}
	}

	if len(r.ViewMods) > 0 {
		for _, v := range r.ViewMods {

			key, err = buildKey(v.PartitionKey, v.ClusterKey)

			if err != nil {
				return &DBResponse{Status: 400, Error: err.Error()}
			}

			switch d.lightWeight {
			case 2:
				_, err = d.updLwL(key, r.Partition, v.ViewType, v.Values)
			case 1:
				_, err = d.updLw(key, r.Partition, v.ViewType, v.Values)
			default:
				_, err = d.upd(key, r.Partition, v.ViewType, v.Values)
			}

			if err != nil {
				return &DBResponse{Status: 400, Error: err.Error()}
			}
		}
	}

	return &DBResponse{Status: 200}
}

//Scan s.e.
func (d *CasandraPartitionedDriver) Scan(r *DBRequest) *DBResponse {
	// not implemented yet
	return &DBResponse{Status: 200}
}

//Delete s.e.
func (d *CasandraPartitionedDriver) Delete(r *DBRequest) *DBResponse {
	if r == nil {
		return &DBResponse{Status: 400, Error: "wrong request data"}
	}

	if len(r.ViewViews) > 0 {

		for _, v := range r.ViewViews {
			if key, e := buildKey(v.PartitionKey, v.ClusterKey); e == nil {
				err := d.delete(
					key,
					r.Partition,
					v.ViewType)

				if err != nil {
					return &DBResponse{Status: 400, Error: err.Error()}
				}
			} else {
				return &DBResponse{Status: 400, Error: e.Error()}
			}
		}
	}

	return &DBResponse{Status: 200}
}

func (d *CasandraPartitionedDriver) delete(key string, partition int64, vtype string) error {

	/*

		 session, err := d.cluster.CreateSession()

		 if err != nil {
			 return err
		 }

		 defer session.Close()

	*/

	if err := d.session.Query(`DELETE FROM records_p WHERE key = ? and partition = ? and type = ?`, key, partition, vtype).Exec(); err != nil {
		return err
	}

	return nil
}

func (d *CasandraPartitionedDriver) get(key string, partition int64, vtype string) (*Record, error) {
	var values []byte
	var version int

	/*

		 session, err := d.cluster.CreateSession()

		 if err != nil {
			 return nil, err
		 }

		 defer session.Close()

	*/

	if err := d.session.Query(`SELECT values, version FROM records_p WHERE key = ? and partition = ?`, key, partition).Scan(&values, &version); err != nil {
		return nil, err
	}

	r := Record{Key: key, Version: version}

	if err := json.Unmarshal(values, &r.Values); err != nil {
		return nil, err
	}

	return &r, nil
}

func (d *CasandraPartitionedDriver) set(key string, partition int64, vtype string, values map[string]interface{}) error {
	b, e := json.Marshal(values)

	if e != nil {
		return e
	}

	if err := d.session.Query(`INSERT INTO records_p (key, partition, version, type, values, weight) VALUES (?, ?, ?, ?, ?, ?)`, key, partition, 0, vtype, b, 0).Exec(); err != nil {
		d.logger.Error("Set error %v", err.Error())
		return err
	}

	return nil
}

func (d *CasandraPartitionedDriver) upd(key string, partition int64, vtype string, values map[string]interface{}) (bool, error) {
	b, e := json.Marshal(values)

	if e != nil {
		return false, e
	}

	var q = d.session.Query(`UPDATE records_p SET version=?, values=? WHERE key = ? and partition = ?`, 0, b, key, partition)

	if err := q.Exec(); err != nil {
		return false, err
	}

	return true, nil
}

func (d *CasandraPartitionedDriver) updLw(key string, partition int64, vtype string, values map[string]interface{}) (bool, error) {
	repeatCount := 0
	version := 0

	for {
		record, err := d.get(key, partition, vtype)

		if err != nil {
			return false, err
		}

		version = record.Version

		for k, v := range values {
			record.Values[k] = v
		}

		b, e := json.Marshal(record.Values)

		if e != nil {
			return false, e
		}

		var q = d.session.Query(`UPDATE records_p SET version=?, values=? WHERE key = ? and partition = ? if version = ?`, version+1, b, key, partition, version)

		err = q.Exec()

		if err != nil {
			repeatCount++
		} else {
			break
		}

		if repeatCount >= LWRepeatCount {
			return false, err
		}
	}

	return true, nil
}

func (d *CasandraPartitionedDriver) updLwL(key string, partition int64, vtype string, values map[string]interface{}) (bool, error) {
	b, e := json.Marshal(values)

	if e != nil {
		return false, e
	}

	var q = d.session.Query(`UPDATE records_p SET values=? WHERE key = ? and partition = ? if weight = ?`, b, key, partition, 0)

	if err := q.Exec(); err != nil {
		return false, err
	}

	return true, nil
}

func (d *CasandraPartitionedDriver) initParams() error {
	if err := d.initHosts(); err != nil {
		return err
	}

	if err := d.initKeyspace(); err != nil {
		return err
	}

	if err := d.initConsistensy(); err != nil {
		return err
	}

	if err := d.initClass(); err != nil {
		return err
	}

	if err := d.initReplicationFactor(); err != nil {
		return err
	}

	if err := d.initLWTMode(); err != nil {
		return err
	}

	return nil
}

func (d *CasandraPartitionedDriver) initHosts() error {
	if h, exists := os.LookupEnv(HostsEnvironmentProperty); exists {
		hosts := strings.Split(strings.TrimSpace(h), ",")

		if len(hosts) >= 1 {
			d.hosts = hosts
			return nil
		}

		return fmt.Errorf("environment variable %v malformed: string of IPs expected (comma as delimiter)", HostsEnvironmentProperty)
	}

	if h, exists := d.args[HostsAttribute]; exists {
		hosts := strings.Split(strings.TrimSpace(h), ",")

		if len(hosts) >= 1 {
			d.hosts = hosts
			return nil
		}

		return fmt.Errorf("attribute %v malformed: string of IPs expected (comma as delimiter)", HostsAttribute)
	}

	d.hosts = []string{DefaultHost}

	return nil
}

func (d *CasandraPartitionedDriver) initKeyspace() error {
	if keyspace, exists := os.LookupEnv(KeyspaceEnvironmentProperty); exists {
		if len(keyspace) > 0 {
			d.keyspace = strings.TrimSpace(keyspace)
			return nil
		}

		return fmt.Errorf("environment variable %v malformed: not empty string expected", KeyspaceEnvironmentProperty)
	}

	if keyspace, exists := d.args[KeyspaceAttribute]; exists {
		if len(keyspace) > 0 {
			d.keyspace = strings.TrimSpace(keyspace)
			return nil
		}

		return fmt.Errorf("argument %v malformed: not empty string expected", KeyspaceEnvironmentProperty)
	}

	d.keyspace = DefaultKeyspaceName

	return nil
}

func (d *CasandraPartitionedDriver) initConsistensy() error {
	if con, exists := os.LookupEnv(ConsistencyEnvironmentProperty); exists {
		d.consistency = d.getConsistency(con)
		return nil
	}

	if con, exists := d.args[ConsistencyAttribute]; exists {
		d.consistency = d.getConsistency(con)
		return nil
	}

	d.consistency = gocql.All

	return nil
}

func (d *CasandraPartitionedDriver) initClass() error {
	d.class = initStringParam(d.args, ClassEnvironmentProperty, ClassAttribute, DefaultClass)

	return nil
}

func (d *CasandraPartitionedDriver) initReplicationFactor() error {
	rf := initIntParam(d.args, ReplicationFactorEnvironmentProperty, ReplicationFactorAttribute, DefaultReplicationFactor)

	if rf <= 0 {
		rf = DefaultReplicationFactor
	}

	d.replicationFactor = rf

	return nil
}
func (d *CasandraPartitionedDriver) initLWTMode() error {
	lwt := initIntParam(d.args, LightWeightTransactionEnvironmentProperty, LightWeightTransactionAttribute, 0)

	if lwt >= 0 && lwt <= 2 {
		d.lightWeight = lwt
	}

	return nil
}

func (d *CasandraPartitionedDriver) getConsistency(s string) gocql.Consistency {
	switch s {
	case "any":
		return gocql.Any
	case "one":
		return gocql.One
	case "two":
		return gocql.Two
	case "three":
		return gocql.Three
	case "quorum":
		return gocql.Quorum
	case "all":
		return gocql.All
	case "lquorum":
		return gocql.LocalQuorum
	case "equorum":
		return gocql.EachQuorum
	case "lone":
		return gocql.LocalOne
	default:
		return gocql.All
	}
}
