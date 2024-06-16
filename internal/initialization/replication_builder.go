package initialization

import (
	"errors"
	"time"

	"github.com/rs/zerolog"

	"picodb/internal/config"
	"picodb/internal/database/storage/replication"
	"picodb/internal/database/storage/wal"
	"picodb/internal/network"
)

const defaultReplicationType = "master"
const defaultReplicationMasterAddress = "localhost:3232"
const defaultReplicationSyncInterval = time.Second

func CreateReplica(
	replicationCfg *config.ReplicationConfig,
	walCfg *config.WALConfig,
	logger *zerolog.Logger,
	stream chan<- []wal.LogData,
) (interface{}, error) {
	replicaType := defaultReplicationType
	masterAddress := defaultReplicationMasterAddress
	syncInterval := defaultReplicationSyncInterval
	walDirectory := defaultWALDataDirectory

	if replicationCfg != nil {
		if replicationCfg.ReplicaType != "" {
			if replicationCfg.ReplicaType != "master" && replicationCfg.ReplicaType != "slave" {
				return nil, errors.New("replica type is incorrect")
			}

			replicaType = replicationCfg.ReplicaType
		}

		if replicationCfg.MasterAddress != "" {
			masterAddress = replicationCfg.MasterAddress
		}

		if replicationCfg.SyncInterval != 0 {
			syncInterval = replicationCfg.SyncInterval
		}
	}

	if walCfg != nil && walCfg.DataDirectory != "" {
		walDirectory = walCfg.DataDirectory
	}

	const maxReplicasNumber = 5
	const maxMessageSize = 16 << 20
	idleTimeout := syncInterval * 3

	if replicaType == "master" {
		server, err := network.NewTCPServer(masterAddress, maxReplicasNumber, maxMessageSize, idleTimeout, logger)
		if err != nil {
			return nil, err
		}

		return replication.NewMaster(server, walDirectory, logger)
	}

	client, err := network.NewTCPClient(masterAddress, maxMessageSize, idleTimeout)
	if err != nil {
		return nil, err
	}

	return replication.NewSlave(client, stream, walDirectory, syncInterval, logger)
}
