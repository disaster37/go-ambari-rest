# go-amabri-rest
Rest client for Ambari API in Golang
It provide cli and library


All development is base on the following documentation:
- https://github.com/apache/ambari/blob/trunk/ambari-server/docs/api/v1/index.md
- https://github.com/apache/ambari/blob/trunk/ambari-server/docs/api/v1/update-hostcomponent.md


## Contribute

You PR are always welcome. Please use develop branch to do PR (git flow pattern)
Don't forget to add test if you add some functionalities.

To build, you can use the following command line:
```sh
make build
```

To lauch golang test, you can use the folowing command line:
```sh
make test-api
```

To lauch cli test, you can use the following command line:
```sh
make test-cli
```




## CLI

### Global options

The following parameters are available for all commands line :
- **--ambari-url**: The Ambari URL. For exemple https://srv1:8443
- **--ambari-login**: The Ambari login to connect on Ambari API
- **--ambari-password**: The Ambari password to connect on Ambari API
- **--debug**: Enable the debug mode
- **--help**: Display help for the current command

### Create or update repository

This command line permit to create or update the repository to get HDP RPM files.
it has the following parameters:
- **--repository-file**: The Json file that describe the repository to add or update
- **--use-spacewalk**: Use it if your repository is managed by Spacewalk


Sample of `repository.json`:
```json
{
	"stackName": "HDP",
	"stackVersion": "2.6",
	"name": "HDP-2.6.4.0",
	"version": "2.6.4.0-91",
	"operatingSystems": [{
		"osName": "redhat7",
		"repositories": [{
			"repositoryName": "HDP",
			"repositoryId": "HDP-2.6.4.0",
			"repositoryBaseUrl": "http://<my_private_url>/hdp/HDP-2.6.4.0"
		}, {
			"repositoryName": "HDP-UTILS",
			"repositoryId": "HDP-UTILS-1.1.0.22",
			"repositoryBaseUrl": "http://<my_private_url>/hdp-utils/HDP-UTILS-1.1.0.22"
		}, {
			"repositoryName": "HDP-GPL",
			"repositoryId": "HDP-GPL-2.6.4.0",
			"repositoryBaseUrl": "http://<my_private_url>/hdp-gpl/HDP-GPL-2.6.4.0"
		}, {
			"repositoryName": "HDP-SOLR",
			"repositoryId": "HDP-SOLR-2.6-100",
			"repositoryBaseUrl": "http://<my_private_url>/hdp-search/HDP-SOLR-2.6-100"
		}, {
			"repositoryName": "HDP-UTILS-GPL",
			"repositoryId": "HDP-UTILS-GPL-1.1.0.22",
			"repositoryBaseUrl": "http://<my_private_url>/hdp-utils-gpl/HDP-UTILS-GPL-1.1.0.22"
		}]
	}]
}
```

Sample of how to use this command line
```sh
./ambari-cli_linux_amd64 --ambari-url https://ambari-server:8443/api/v1 --ambari-login admin --ambari-password admin create-or-update-repository --repository-file repository.json
```

### Create cluster if not exist

This command line permit to create new cluster with Ambari blueprint API and so deploy a complete topology of HDP
First, you need to add some repository before to deploy new cluster.
If cluster already exist, it do nothing and not return error. It just skip the creation.

It has the following parameters:
- **--cluster-name**: The name of the cluster you should to create
- **--blueprint-file**: The json file that describe the HDP topologie
- **--hosts-template-file**: The Json file that describe the role of each HDP server

Sample of blueprint.json
```json
{
	"configurations": [{
		"ams-grafana-env": {
			"properties": {
				"metrics_grafana_password": "P@sswOrd$007"
			}
		}
	}, {
		"hive-site": {
			"properties": {
				"javax.jdo.option.ConnectionUserName": "hive",
				"javax.jdo.option.ConnectionPassword": "hive",
				"javax.jdo.option.ConnectionDriverName": "org.postgresql.Driver",
				"javax.jdo.option.ConnectionURL": "jdbc:postgresql://postgres-server:5432/hive",
				"hive.server2.enable.doAs": "true",
				"hive.default.fileformat.managed": "ORC",
				"hive.plan.serialization.format": "kryo"
			}
		}
	}, {
		"hive-env": {
			"properties": {
				"hive_database": "Existing PostgreSQL Database",
				"hive_database_name": "hive",
				"hive_database_type": "postgres"
			}
		}
	}, {
		"oozie-env": {
			"properties": {
				"oozie_database": "Existing PostgreSQL Database"
			}
		}
	}, {
		"oozie-site": {
			"properties": {
				"oozie.service.JPAService.jdbc.driver": "org.postgresql.Driver",
				"oozie.service.JPAService.jdbc.password": "oozie",
				"oozie.service.JPAService.jdbc.url": "jdbc:postgresql://postgres-server:5432/oozie",
				"oozie.service.JPAService.jdbc.username": "oozie",
				"oozie.zookeeper.connection.string": "%HOSTGROUP::master0%:2181,%HOSTGROUP::master1%:2181,%HOSTGROUP::master2%:2181",
				"oozie.services.ext": "org.apache.oozie.service.ZKLocksService,org.apache.oozie.service.ZKXLogStreamingService,org.apache.oozie.service.ZKJobsConcurrencyService",
				"oozie.base.url": "http://%HOSTGROUP::master1%:11000/oozie"
			}
		}
	}, {
		"admin-properties": {
			"properties": {
				"SQL_CONNECTOR_JAR": "/usr/share/java/postgresql-jdbc.jar",
				"DB_FLAVOR": "POSTGRES",
				"db_host": "postgres-server",
				"db_name": "ranger",
				"db_user": "ranger",
				"db_password": "ranger",
				"db_root_password": ""
			}
		}
	}, {
		"ranger-admin-site": {
			"properties": {
				"ranger.jpa.jdbc.driver": "org.postgresql.Driver",
				"ranger.jpa.jdbc.url": "jdbc:postgresql://postgres-server:5432/ranger"
			}
		}
	}, {
		"ranger-env": {
			"properties": {
				"ranger-hdfs-plugin-enabled": "Yes",
				"ranger-hbase-plugin-enabled": "Yes",
				"ranger-hive-plugin-enabled": "Yes",
				"ranger-knox-plugin-enabled": "Yes",
				"ranger-yarn-plugin-enabled": "Yes",
				"ranger-atlas-plugin-enabled": "Yes",
				"ranger-storm-plugin-enabled": "No",
				"ranger-kafka-plugin-enabled": "No",
				"ranger_admin_password": "P@sswOrd$007",
				"xasecure.audit.destination.db": "false",
				"xasecure.audit.destination.solr": "true",
				"xasecure.audit.destination.hdfs": "false",
				"is_solrCloud_enabled": "true",
				"create_db_dbuser": "false"
			}
		}
	}, {
		"knox-env": {
			"properties": {
				"knox_master_secret": "P@sswOrd$007"
			}
		}
	}, {
		"core-site": {
			"properties": {
				"fs.defaultFS": "hdfs://cluster",
				"ha.zookeeper.quorum": "%HOSTGROUP::master0%:2181,%HOSTGROUP::master1%:2181,%HOSTGROUP::master2%:2181",
				"fs.trash.interval": "4320",
				"fs.protected.directories": "/,/mapred,/mapred/system,/tmp,/user,/user/apps,/user/apps/hbase,/user/apps/hbase/data,/user/apps/hive"
			}
		}
	}, {
		"hdfs-site": {
			"properties": {
				"dfs.client.failover.proxy.provider.cluster": "org.apache.hadoop.hdfs.server.namenode.ha.ConfiguredFailoverProxyProvider",
				"dfs.ha.automatic-failover.enabled": "true",
				"dfs.ha.fencing.methods": "shell(/bin/true)",
				"dfs.ha.namenodes.cluster": "nn1,nn2",
				"dfs.namenode.http-address": "%HOSTGROUP::master0%:50070",
				"dfs.namenode.http-address.cluster.nn1": "%HOSTGROUP::master0%:50070",
				"dfs.namenode.http-address.cluster.nn2": "%HOSTGROUP::master2%:50070",
				"dfs.namenode.https-address": "%HOSTGROUP::master0%:50470",
				"dfs.namenode.https-address.cluster.nn1": "%HOSTGROUP::master0%:50470",
				"dfs.namenode.https-address.cluster.nn2": "%HOSTGROUP::master2%:50470",
				"dfs.namenode.rpc-address.cluster.nn1": "%HOSTGROUP::master0%:8020",
				"dfs.namenode.rpc-address.cluster.nn2": "%HOSTGROUP::master2%:8020",
				"dfs.namenode.shared.edits.dir": "qjournal://%HOSTGROUP::master0%:8485;%HOSTGROUP::master1%:8485;%HOSTGROUP::master2%:8485/cluster",
				"dfs.nameservices": "cluster",
				"dfs.namenode.name.dir": "/data/1/hadoop/hdfs/namenode",
				"dfs.datanode.data.dir": "/data/1/hadoop/hdfs/data,/data/2/hadoop/hdfs/data,/data/3/hadoop/hdfs/data,/data/4/hadoop/hdfs/data,/data/5/hadoop/hdfs/data,/data/6/hadoop/hdfs/data,/data/7/hadoop/hdfs/data,/data/8/hadoop/hdfs/data",
				"dfs.datanode.failed.volumes.tolerated": "1",
				"dfs.namenode.safemode.threshold-pct": "0.99f",
				"dfs.namenode.checkpoint.period": "3600",
				"dfs.namenode.checkpoint.txns": "10000000",
				"fs.permissions.umask-mode": "077",
				"dfs.namenode.acls.enabled": "false"
			}
		}
	}, {
		"yarn-site": {
			"properties": {
				"hadoop.registry.rm.enabled": "true",
				"hadoop.registry.zk.quorum": "%HOSTGROUP::master0%:2181,%HOSTGROUP::master1%:2181,%HOSTGROUP::master2%:2181",
				"yarn.log.server.url": "http://%HOSTGROUP::master1%:19888/jobhistory/logs",
				"yarn.resourcemanager.address": "%HOSTGROUP::master2%:8050",
				"yarn.resourcemanager.admin.address": "%HOSTGROUP::master2%:8141",
				"yarn.resourcemanager.cluster-id": "yarn-cluster",
				"yarn.resourcemanager.ha.automatic-failover.zk-base-path": "/yarn-leader-election",
				"yarn.resourcemanager.ha.enabled": "true",
				"yarn.resourcemanager.ha.rm-ids": "rm1,rm2",
				"yarn.resourcemanager.hostname": "%HOSTGROUP::master2%",
				"yarn.resourcemanager.recovery.enabled": "true",
				"yarn.resourcemanager.resource-tracker.address": "%HOSTGROUP::master2%:8025",
				"yarn.resourcemanager.scheduler.address": "%HOSTGROUP::master2%:8030",
				"yarn.resourcemanager.store.class": "org.apache.hadoop.yarn.server.resourcemanager.recovery.ZKRMStateStore",
				"yarn.resourcemanager.webapp.address": "%HOSTGROUP::master2%:8088",
				"yarn.resourcemanager.webapp.https.address": "%HOSTGROUP::master2%:8090",
				"yarn.timeline-service.address": "%HOSTGROUP::master1%:10200",
				"yarn.timeline-service.webapp.address": "%HOSTGROUP::master1%:8188",
				"yarn.timeline-service.webapp.https.address": "%HOSTGROUP::master1%:8190",
				"yarn.resourcemanager.hostname.rm1": "%HOSTGROUP::master0%",
				"yarn.resourcemanager.hostname.rm2": "%HOSTGROUP::master2%",
				"yarn.resourcemanager.zk-address": "%HOSTGROUP::master0%:2181,%HOSTGROUP::master1%:2181,%HOSTGROUP::master2%:2181",
				"yarn.timeline-service.leveldb-state-store.path": "/var/hadoop/yarn/timeline",
				"yarn.timeline-service.leveldb-timeline-store.path": "/var/hadoop/yarn/timeline",
				"yarn.timeline-service.generic-application-history.save-non-am-container-meta-info": "false",
				"yarn.scheduler.minimum-allocation-mb": "12800",
				"yarn.scheduler.maximum-allocation-mb": "192000",
				"yarn.nodemanager.resource.memory-mb": "192000"
			}
		}
	}, {
		"hbase-site": {
			"properties": {
				"hbase.rootdir": "hdfs://cluster/hbase",
				"hbase.ipc.client.tcpnodelay": "true",
				"hbase.ipc.server.tcpnodelay": "true"
			}
		}
	}, {
		"application-properties": {
			"properties": {
				"atlas.rest.address": "http://%HOSTGROUP::master1%:21000,http://%HOSTGROUP::master2%:21000",
				"atlas.server.ha.enabled": "true",
				"atlas.server.ids": "id1,id2",
				"atlas.server.address.id1": "%HOSTGROUP::master1%:21000",
				"atlas.server.address.id2": "%HOSTGROUP::master2%:21000",
				"atlas.server.bind.address": "%HOSTGROUP::master1%"
			}
		}
	}, {
		"tez-site": {
			"properties": {
				"tez.task.am.heartbeat.interval-ms.max": "200"
			}
		}
	}, {
		"hst-server-conf": {
			"properties": {
				"customer.account.name": "None",
				"customer.smartsense.id": "None",
				"customer.notification.email": "None"
			}
		}
	}, {
		"hadoop-env": {
			"properties": {
				"namenode_heapsize": "16332m",
				"namenode_opt_newsize": "2048m",
				"namenode_opt_maxnewsize": "2048m",
				"dtnode_heapsize": "4096m"
			}
		}
	}, {
		"mapred-site": {
			"properties": {
				"mapreduce.map.memory.mb": "12800",
				"mapreduce.map.java.opts": "-Xmx10240m",
				"mapreduce.reduce.memory.mb": "12800",
				"mapreduce.reduce.java.opts": "-Xmx10240m",
				"yarn.app.mapreduce.am.resource.mb": "12800",
				"yarn.app.mapreduce.am.command-opts": "-Xmx10240m -Dhdp.version=${hdp.version}",
				"mapreduce.task.io.sort.mb": "5120"
			}
		}
	}, {
		"activity-zeppelin-shiro": {
			"properties": {
				"users.admin": "P@sswOrd$007"
			}
		}
	}],
	"host_groups": [{
		"name": "master0",
		"configurations": [],
		"components": [{
			"name": "ZOOKEEPER_CLIENT"
		}, {
			"name": "ZOOKEEPER_SERVER"
		}, {
			"name": "INFRA_SOLR_CLIENT"
		}, {
			"name": "METRICS_MONITOR"
		}, {
			"name": "NAMENODE"
		}, {
			"name": "ZKFC"
		}, {
			"name": "JOURNALNODE"
		}, {
			"name": "HDFS_CLIENT"
		}, {
			"name": "HBASE_MASTER"
		}, {
			"name": "PHOENIX_QUERY_SERVER"
		}, {
			"name": "HBASE_CLIENT"
		}, {
			"name": "YARN_CLIENT"
		}, {
			"name": "RESOURCEMANAGER"
		}, {
			"name": "MAPREDUCE2_CLIENT"
		}, {
			"name": "SPARK2_THRIFTSERVER"
		}, {
			"name": "LIVY2_SERVER"
		}, {
			"name": "SPARK2_CLIENT"
		}, {
			"name": "TEZ_CLIENT"
		}, {
			"name": "HIVE_SERVER"
		}, {
			"name": "HIVE_METASTORE"
		}, {
			"name": "HIVE_CLIENT"
		}, {
			"name": "WEBHCAT_SERVER"
		}, {
			"name": "HCAT"
		}, {
			"name": "ATLAS_CLIENT"
		}, {
			"name": "OOZIE_CLIENT"
		}, {
			"name": "OOZIE_SERVER"
		}, {
			"name": "PIG"
		}, {
			"name": "SLIDER"
		}, {
			"name": "HST_AGENT"
		}],
		"cardinality": "1"
	}, {
		"name": "master1",
		"configurations": [],
		"components": [{
			"name": "ZOOKEEPER_CLIENT"
		}, {
			"name": "ZOOKEEPER_SERVER"
		}, {
			"name": "INFRA_SOLR_CLIENT"
		}, {
			"name": "METRICS_MONITOR"
		}, {
			"name": "METRICS_GRAFANA"
		}, {
			"name": "METRICS_COLLECTOR"
		}, {
			"name": "HDFS_CLIENT"
		}, {
			"name": "JOURNALNODE"
		}, {
			"name": "HBASE_CLIENT"
		}, {
			"name": "PHOENIX_QUERY_SERVER"
		}, {
			"name": "YARN_CLIENT"
		}, {
			"name": "APP_TIMELINE_SERVER"
		}, {
			"name": "MAPREDUCE2_CLIENT"
		}, {
			"name": "HISTORYSERVER"
		}, {
			"name": "SPARK2_THRIFTSERVER"
		}, {
			"name": "SPARK2_CLIENT"
		}, {
			"name": "SPARK2_JOBHISTORYSERVER"
		}, {
			"name": "LIVY2_SERVER"
		}, {
			"name": "TEZ_CLIENT"
		}, {
			"name": "HIVE_SERVER"
		}, {
			"name": "HIVE_CLIENT"
		}, {
			"name": "WEBHCAT_SERVER"
		}, {
			"name": "HCAT"
		}, {
			"name": "ATLAS_CLIENT"
		}, {
			"name": "ATLAS_SERVER"
		}, {
			"name": "OOZIE_CLIENT"
		}, {
			"name": "OOZIE_SERVER"
		}, {
			"name": "RANGER_ADMIN"
		}, {
			"name": "RANGER_TAGSYNC"
		}, {
			"name": "RANGER_USERSYNC"
		}, {
			"name": "PIG"
		}, {
			"name": "SLIDER"
		}, {
			"name": "ACTIVITY_ANALYZER"
		}, {
			"name": "ACTIVITY_EXPLORER"
		}, {
			"name": "HST_SERVER"
		}, {
			"name": "HST_AGENT"
		}],
		"cardinality": "1"
	}, {
		"name": "master2",
		"configurations": [],
		"components": [{
			"name": "ZOOKEEPER_CLIENT"
		}, {
			"name": "ZOOKEEPER_SERVER"
		}, {
			"name": "INFRA_SOLR_CLIENT"
		}, {
			"name": "METRICS_MONITOR"
		}, {
			"name": "NAMENODE"
		}, {
			"name": "ZKFC"
		}, {
			"name": "JOURNALNODE"
		}, {
			"name": "HDFS_CLIENT"
		}, {
			"name": "HBASE_CLIENT"
		}, {
			"name": "HBASE_MASTER"
		}, {
			"name": "PHOENIX_QUERY_SERVER"
		}, {
			"name": "YARN_CLIENT"
		}, {
			"name": "RESOURCEMANAGER"
		}, {
			"name": "MAPREDUCE2_CLIENT"
		}, {
			"name": "SPARK2_THRIFTSERVER"
		}, {
			"name": "SPARK2_CLIENT"
		}, {
			"name": "LIVY2_SERVER"
		}, {
			"name": "TEZ_CLIENT"
		}, {
			"name": "HIVE_SERVER"
		}, {
			"name": "HIVE_METASTORE"
		}, {
			"name": "HIVE_CLIENT"
		}, {
			"name": "WEBHCAT_SERVER"
		}, {
			"name": "HCAT"
		}, {
			"name": "ATLAS_CLIENT"
		}, {
			"name": "ATLAS_SERVER"
		}, {
			"name": "OOZIE_CLIENT"
		}, {
			"name": "PIG"
		}, {
			"name": "SLIDER"
		}, {
			"name": "HST_AGENT"
		}],
		"cardinality": "1"
	}, {
		"name": "services",
		"configurations": [],
		"components": [{
			"name": "ZOOKEEPER_CLIENT"
		}, {
			"name": "METRICS_MONITOR"
		}, {
			"name": "HDFS_CLIENT"
		}, {
			"name": "HBASE_CLIENT"
		}, {
			"name": "YARN_CLIENT"
		}, {
			"name": "MAPREDUCE2_CLIENT"
		}, {
			"name": "SPARK2_CLIENT"
		}, {
			"name": "TEZ_CLIENT"
		}, {
			"name": "HIVE_CLIENT"
		}, {
			"name": "HCAT"
		}, {
			"name": "KAFKA_BROKER"
		}, {
			"name": "KNOX_GATEWAY"
		}, {
			"name": "PIG"
		}, {
			"name": "SLIDER"
		}, {
			"name": "HST_AGENT"
		}],
		"cardinality": "2"
	}, {
		"name": "ingests",
		"configurations": [],
		"components": [{
			"name": "ZOOKEEPER_CLIENT"
		}, {
			"name": "METRICS_MONITOR"
		}, {
			"name": "HDFS_CLIENT"
		}, {
			"name": "HBASE_CLIENT"
		}, {
			"name": "YARN_CLIENT"
		}, {
			"name": "MAPREDUCE2_CLIENT"
		}, {
			"name": "SPARK2_CLIENT"
		}, {
			"name": "TEZ_CLIENT"
		}, {
			"name": "HIVE_CLIENT"
		}, {
			"name": "HCAT"
		}, {
			"name": "KAFKA_BROKER"
		}, {
			"name": "PIG"
		}, {
			"name": "SLIDER"
		}, {
			"name": "HST_AGENT"
		}],
		"cardinality": "2"
	}, {
		"name": "workers",
		"configurations": [],
		"components": [{
			"name": "ZOOKEEPER_CLIENT"
		}, {
			"name": "INFRA_SOLR"
		}, {
			"name": "METRICS_MONITOR"
		}, {
			"name": "DATANODE"
		}, {
			"name": "HDFS_CLIENT"
		}, {
			"name": "HBASE_CLIENT"
		}, {
			"name": "HBASE_REGIONSERVER"
		}, {
			"name": "YARN_CLIENT"
		}, {
			"name": "NODEMANAGER"
		}, {
			"name": "MAPREDUCE2_CLIENT"
		}, {
			"name": "SPARK2_CLIENT"
		}, {
			"name": "TEZ_CLIENT"
		}, {
			"name": "HIVE_CLIENT"
		}, {
			"name": "HCAT"
		}, {
			"name": "PIG"
		}, {
			"name": "SLIDER"
		}, {
			"name": "SOLR_SERVER"
		}, {
			"name": "HST_AGENT"
		}],
		"cardinality": "4"
	}],
	"Blueprints": {
		"stack_name": "HDP",
		"stack_version": "2.6"
	}
}
```

Sample of host-template.json
```json
{
	"blueprint": "my_cluster",
	"config_recommendation_strategy": "ALWAYS_APPLY_DONT_OVERRIDE_CUSTOM_VALUES",
	"host_groups": [{
		"name": "master0",
		"hosts": [{
			"fqdn": "master0.domain.com"
		}]
	}, {
		"name": "master1",
		"hosts": [{
			"fqdn": "master1.domain.com"
		}]
	}, {
		"name": "master2",
		"hosts": [{
			"fqdn": "master2.domain.com"
		}]
	}, {
		"name": "services",
		"hosts": [{
			"fqdn": "service0.domain.com"
		}, {
			"fqdn": "service1..domain.com"
		}]
	}, {
		"name": "ingests",
		"hosts": [{
			"fqdn": "ingest0.domain.com"
		}, {
			"fqdn": "ingest1..domain.com"
		}]
	}, {
		"name": "workers",
		"hosts": [{
			"fqdn": "worker0.domain.com"
		}, {
			"fqdn": "worker1.domain.com"
		}, {
			"fqdn": "worker2.domain.com"
		}, {
			"fqdn": "worker3.domain.com"
		}]
	}]
}
```

Sample of how to use this command line
```sh
./ambari-cli_linux_amd64 --ambari-url https://ambari-server:8443/api/v1 --ambari-login admin --ambari-password admin create-cluster-if-not-exist --cluster-name my_cluster --blueprint-file blueprint.json --hosts-template-file host-template.json
```

### Create or update privileges

This command line permit to create or update the privileges in HDP cluster.
it has the following parameters:
- **--privileges-file**: The Json file that describe the privileges to add or update


Sample of `privileges.json`:
```json
{
	"clusterName": "nemodatahubdev2",
	"privileges": [{
		"permission": "CLUSTER.ADMINISTRATOR",
		"type": "GROUP",
		"name": "hdp_admin"
	}, {
		"permission": "CLUSTER.USER",
		"type": "GROUP",
		"name": "hdp_user"
	}, {
		"permission": "CLUSTER.OPERATOR",
		"type": "GROUP",
		"name": "hdp_operator"
	}]
}
```

Sample of how to use this command line
```sh
./ambari-cli_linux_amd64 --ambari-url https://ambari-server:8443/api/v1 --ambari-login admin --ambari-password admin create-or-update-privileges --privileges-file privileges.json
```

### Add new node in existing cluster deployed with Blueprint API

This command line permit to add new node in existing HDP cluster deployed with Blueprint API.
You need to map the new node with one of the blueprint API.
it has the following parameters:
- **--cluster-name**: The HDP cluster name where you should to add the new node
- **--blueprint-name**: The Blueprint name use to deploy the HDP cluster
- **--hostname**: The name of the new node. Per default, it's the node's FQDN.
- **--role**: The Blueprint `host_group` to affect the new node.
- **--rack** (optionnal): The node rack for rack awerness.


Sample of how to use this command line
```sh
./ambari-cli_linux_amd64 --ambari-url https://ambari-server:8443/api/v1 --ambari-login admin --ambari-password admin add-host-in-cluster --cluster-name test --blueprint-name test --hostname "node10.domain.com" --role worker --rack "dc1/r1"
```

### Stop one service

This command line permit to stop one service in HDP cluster.
it has the following parameters:
- **--cluster-name**: The HDP cluster name
- **--service-name**: The service name you should to stop
- **--enable-maintenance** (optionnal): Put service in maintenance state.
- **--force** (optionnal): Remove maintenance state in service before to stop it.


Sample of how to use this command line
```sh
./ambari-cli_linux_amd64 --ambari-url https://ambari-server:8443/api/v1 --ambari-login admin --ambari-password admin stop-service --cluster-name test --service-name HBASE --enable-maintenance --force
```

### Start one service

This command line permit to start one service in HDP cluster.
it has the following parameters:
- **--cluster-name**: The HDP cluster name
- **--service-name**: The service name you should to start
- **--disable-maintenance** (optionnal): Remove maintenance state in service before to start it.


Sample of how to use this command line
```sh
./ambari-cli_linux_amd64 --ambari-url https://ambari-server:8443/api/v1 --ambari-login admin --ambari-password admin start-service --cluster-name test --service-name HBASE --disable-maintenance
```

### Stop all services in cluster

This command line permit to stop all services in HDP cluster.
it has the following parameters:
- **--cluster-name**: The HDP cluster name
- **--enable-maintenance** (optionnal): Put all services in maintenance state.
- **--force** (optionnal): Remove maintenance state in all services before to stop them.


Sample of how to use this command line
```sh
./ambari-cli_linux_amd64 --ambari-url https://ambari-server:8443/api/v1 --ambari-login admin --ambari-password admin stop-all-services --cluster-name test --enable-maintenance --force
```

### Start all services in cluster

This command line permit to start all services in HDP cluster.
it has the following parameters:
- **--cluster-name**: The HDP cluster name
- **--disable-maintenance** (optionnal): Remove maintenance state for all services before to start them.


Sample of how to use this command line
```sh
./ambari-cli_linux_amd64 --ambari-url https://ambari-server:8443/api/v1 --ambari-login admin --ambari-password admin start-all-services --cluster-name test --disable-maintenance
```


### Stop all components on node

This command line permit to stop all components on particular node.
it has the following parameters:
- **--cluster-name**: The HDP cluster name
- **--hostname**: The node where you should to stop components. Per default, it's the FQDN.
- **--enable-maintenance** (optionnal): Put node in maintenance state.
- **--force** (optionnal): Remove maintenance state on node before to stop components.


Sample of how to use this command line
```sh
./ambari-cli_linux_amd64 --ambari-url https://ambari-server:8443/api/v1 --ambari-login admin --ambari-password admin stop-all-components-in-host --cluster-name test --hostname worker01.domain.com --enable-maintenance --force
```

### Start all components on node

This command line permit to start all components on particular node.
it has the following parameters:
- **--cluster-name**: The HDP cluster name
- **--hostname**: The node where you should to stop components. Per default, it's the FQDN.
- **--disable-maintenance** (optionnal): Remove maintenance state on node before start components.


Sample of how to use this command line
```sh
./ambari-cli_linux_amd64 --ambari-url https://ambari-server:8443 --ambari-login admin --ambari-password admin start-all-components-in-host --cluster-name test --hostname worker01.domain.com --disable-maintenance
```

### Stop one component on node

This command line permit to stop one component on particular node.
it has the following parameters:
- **--cluster-name**: The HDP cluster name
- **--hostname**: The node where you should to stop component. Per default, it's the FQDN.
- **--component-name**: The component you should to stop.


Sample of how to use this command line
```sh
./ambari-cli_linux_amd64 --ambari-url https://ambari-server:8443/api/v1 --ambari-login admin --ambari-password admin stop-component-in-host --cluster-name test --hostname worker01.domain.com --component-name ZOOKEEPER_SERVER
```

### Start one component on node

This command line permit to start one component on particular node.
it has the following parameters:
- **--cluster-name**: The HDP cluster name
- **--hostname**: The node where you should to start component. Per default, it's the FQDN.
- **--component-name**: The component you should to start.


Sample of how to use this command line
```sh
./ambari-cli_linux_amd64 --ambari-url https://ambari-server:8443/api/v1 --ambari-login admin --ambari-password admin start-component-in-host --cluster-name test --hostname worker01.domain.com --component-name ZOOKEEPER_SERVER
```