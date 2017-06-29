# badkv

A [b.a.d.](https://github.com/bruceadowns) Multi-Tenant Distributed Redundant Key-Value Store written using  [golang](https://golang.org) and its standard libraries.

## Programming Assignment

#### Overview

The assignment is to design and build a proof of concept of a multi-tenant distributed redundant key-value store. Clients interact with the K-V store through a to be REST interaction model, to be designed as part of the assignment. The objective is to have a single service role, which embodies functionality needed to deliver the required capabilities: multi-tenant, distributed, redundant K-V storage and interaction model.

#### Objectives

* Clearly articulate the overall problem statement, use cases and requirements.
* Explain through the technical design, architecture, technical decisions and choices made.
* Reason through the technical tradeoffs made and the impact on the usage scenarios.
* Demonstrate a working proof of concept.
* Answer and defend questions regarding the design and proof of concept by a forum of peers.

## Minimum Viable POC

* Cluster size of 1 or 3
* Synchronized in-memory data structures
* One known leader
* Followers forward write traffic to leader
* Leader replicates to followers
* Client may query any node
* Write succeeds upon propagation

#### .Next Thoughts

* persistence - wal, backup/restore, anti-corrupt, fsync'ed, stripped, indexed
* leadership - raft/paxos/membership/swim/coördination-free
* tls endpoints via pki, self-signed ca, signed certs
* authenticated join via otp or user/password
* member discovery via arguments, discovery service, dns srv records, udp broadcast, etc
* snappy compress replication data based on empirical threshold
* load balance via proxy, F5, AWS ELB
* node and cluster uniqueness, reject duplicates
* multi-key operations, compare-and-set
* consider read/write latency and throughput tradeoffs
* refine log/health/metric strategy
* lexigraphical ordering considerations

## REST API

#### Get Key/Value

/api/v1/keys/{tenant}/{name} - HTTP GET

* return latest [tenant-name] value.data from map[string]value or error

#### Put Key/Value

/api/v1/keys/{tenant}/{name}/{timestamp} - HTTP POST

* request body is the value
* timestamp - optional, default to now
* save [tenant-name]value to map[string]value
* return accepted or error

#### Remove Key

/api/v1/keys/{tenant}/{name}/{timestamp} - HTTP DELETE

* timestamp - optional, default to now
* tombstone [tenant-name] in map[string]value
* return accepted or error

#### Take

/api/v1/take/{timestamp} - HTTP GET

* Return map[string]value records from timestamp forward.
* Run periodically (i.e. every 10s) for eventually consistency.

#### Admin

```
/admin/stats
/admin/metrics
/admin/members
```

## Execution

#### Build

```
$ make all
```

#### Standalone

```
$ ./badkv -name badkv1 -listen localhost:12345
```

#### Clustered (via goreman/foreman)

```
$ goreman start
```

## Dependencies

* https://golang.org - base language and standard library
* http://www.gorillatoolkit.org/pkg/mux - routing support for http requests
