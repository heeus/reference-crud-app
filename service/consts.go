/*
 * Copyright (c) 2019-present Heeus authors
 */

package service

const MaxReconnectCount = 100

//DefaultPort s.e.
const DefaultPort = 80

//DefaultKeyspaceName s.e.
const DefaultKeyspaceName = "heeustst"

//DefaultClass s.e.
const DefaultClass = "SimpleStrategy"

//DefaultReplicationFactor s.e.
const DefaultReplicationFactor = 1

//DefaultHost s.e.
const DefaultHost = "127.0.0.1"

//DefaultPathPattern s.e.
const DefaultPathPattern = "/api/{region}/{zone}/{user}/{app}/{service}/{wsid}/{module}/{consistency}/{function}"

//ReadDefaultFunc s.e.
const ReadDefaultFunc = "YcsbView"

//InsertDefaultFunc s.e.
const InsertDefaultFunc = "YcsbAdd"

//UpdateDefaultFunc s.e.
const UpdateDefaultFunc = "YcsbMod"

//ScanDefaultFunc s.e.
const ScanDefaultFunc = "YcsbScan"

//DeleteDefaultFunc s.e.
const DeleteDefaultFunc = "YcsbDel"

//PathPatternEnvironmentProperty s.e
const ServiceDriverEnvironmentProperty = "SERVICE_DRIVER"

//PathPatternEnvironmentProperty s.e
const ServicePortEnvironmentProperty = "SERVICE_PORT"

//PathPatternEnvironmentProperty s.e
const PathPatternEnvironmentProperty = "SERVICE_PATH_PATTERN"

//ServiceInsertFuncEnvironmentProperty s.e
const ServiceInsertFuncEnvironmentProperty = "SERVICE_INSERT_FUNC_NAME"

//ServiceReadFuncEnvironmentProperty s.e
const ServiceReadFuncEnvironmentProperty = "SERVICE_READ_FUNC_NAME"

//ServiceUpdateFuncEnvironmentProperty s.e
const ServiceUpdateFuncEnvironmentProperty = "SERVICE_UPDATE_FUNC_NAME"

//ServiceScanFuncEnvironmentProperty s.e
const ServiceScanFuncEnvironmentProperty = "SERVICE_SCAN_FUNC_NAME"

//ServiceDeleteFuncEnvironmentProperty s.e
const ServiceDeleteFuncEnvironmentProperty = "SERVICE_DELETE_FUNC_NAME"

//LoggerLevelEnvironmentProperty s.e.
const LoggerLevelEnvironmentProperty = "SERVICE_LOGGER_LEVEL"

//KeyspaceEnvironmentProperty s.e.
const KeyspaceEnvironmentProperty = "DB_KEYSPACE"

//HostsEnvironmentProperty s.e.
const HostsEnvironmentProperty = "DB_SERVERS" // hosts
//UserEnvironmentProperty s.e.
const UserEnvironmentProperty = "DB_USER"

//PasswordEnvironmentProperty s.e.
const PasswordEnvironmentProperty = "DB_PASSWORD"

//ClassEnvironmentProperty s.e.
const ClassEnvironmentProperty = "DB_CAS_CLASS"

//ConsistencyEnvironmentProperty s.e.
const ConsistencyEnvironmentProperty = "DB_CAS_CONSISTENCY"

//ReplicationFactorEnvironmentProperty s.e.
const ReplicationFactorEnvironmentProperty = "DB_REP_FACTOR"

//LightWeightTransactionAttribute s.e.
const LightWeightTransactionEnvironmentProperty = "DB_LWT"

const NoopServiceEnvironmentProperty = "SERVICE_NOP"

//ServiceDriverAttribute s.e
const ServiceDriverAttribute = "-d"

//ServicePortAttribute s.e
const ServicePortAttribute = "-p"

//HostsAttribute s.e.
const HostsAttribute = "--hosts" // "cassandra.hosts";
//KeyspaceAttribute s.e.
const KeyspaceAttribute = "--ks" // "cassandra.keyspace";
//UserAttribute s.e.
const UserAttribute = "--user"

//PasswordAttribute s.e.
const PasswordAttribute = "--pass"

//ClassAttribute s.e.
const ClassAttribute = "--cs"

//ReplicationFactorAttribute s.e.
const ReplicationFactorAttribute = "--rf"

//ConsistencyAttribute s.e.
const ConsistencyAttribute = "--c"

//LightWeightTransactionAttribute s.e.
const LightWeightTransactionAttribute = "--lwt"

const PathPatternAttribute = "-pp"

//ServiceInsertFuncAttribute s.e
const ServiceInsertFuncAttribute = "-ifn"

//ServiceReadFuncAttribute s.e
const ServiceReadFuncAttribute = "-rfn"

//ServiceUpdateFuncAttribute s.e
const ServiceUpdateFuncAttribute = "-ufn"

//ServiceScanFuncAttribute s.e
const ServiceScanFuncAttribute = "-sfn"

//ServiceDeleteFuncAttribute s.e
const ServiceDeleteFuncAttribute = "-dfn"

//ServiceDeleteFuncAttribute s.e
const LoggerLevelAttribute = "-ll"

const NoopServiceAttribute = "-nop"

//HTTPMethods s.e.
var HTTPMethods = []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS", "PATCH"}
