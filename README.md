## Overview

A test reference app 

## Arguments

`-p` - port number for listening
`-d` - driver selection; available three options: 
- `mem` - memory driver;
- `cas` - default; enables cassandra driver
- `light` - light driver that just sends `Ok` status for all operations

`mem` and `light` drivers are not supports any options or arguments.

`-pp` (env.v. `SERVICE_PATH_PATTERN`)- string; handler path pattern; default is `/api/{region}/{zone}/{user}/{app}/{service}/{wsid}/{module}/{consistency}/{function}/`
`-ifn` (env.v. `SERVICE_INSERT_FUNC_NAME`) - string; insert function name; default is `YcsbAdd`
`-rfn` (env.v. `SERVICE_READ_FUNC_NAME`) - string; read function name; default is `YcsbView`
`-ufn` (env.v. `SERVICE_UPDATE_FUNC_NAME`) - string; update function name; default is `YcsbUpd` (not implemented yet)
`-sfn` (env.v. `SERVICE_SCAN_FUNC_NAME`) - string; scan function name; default is `YcsbScan` (not implemented yet)
`-dfn` (env.v. `SERVICE_DELETE_FUNC_NAME`) - string; deelte function name; default is `YcsbDel`

## Cassandra driver setup

### arguments

`--hosts` - hosts IPs separated with comma
`--ks` - keyspace name; default is `heeus`
`--user` - user login; not implemented;
`--pass` - user password; not implemented;
`--cs` - strategy class; available values: `SimpleStrategy`(default), `NetworkTopologyStrategy`
`--rf` - replication factor; default is 3
`--c` - consistency level; available values: `any`, `one`, `two`, `three`, `quorum`, `all`, `lquorum`, `equorum`, `lone`; `all` is default
`--lwt` - if 1 is given the light weight transaction mode will be enabled

### enviroment variables

`DB_KEYSPACE` - keyspace name; default is `heeus`
`DB_SERVERS` - hosts IPs separated with comma
`DB_USER` - user login; not implemented;
`DB_PASSWORD` - user password; not implemented;
`DB_CAS_CLASS` - strategy class; available values: `SimpleStrategy`(default), `NetworkTopologyStrategy`
`DB_CAS_CONSISTENCY` - consistency level; available values: `any`, `one`, `two`, `three`, `quorum`, `all`, `lquorum`, `equorum`, `lone`; `all` is default
`DB_REP_FACTOR` - replication factor; default is 3

You can use both arguments and e.variables setup method, but e.variables will be used in priority 