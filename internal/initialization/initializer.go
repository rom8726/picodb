package initialization

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	"picodb/internal/config"
	"picodb/internal/database"
	"picodb/internal/database/compute"
	"picodb/internal/database/storage"
	"picodb/internal/database/storage/replication"
	walPkg "picodb/internal/database/storage/wal"
	"picodb/internal/network"
)

type Initializer struct {
	wal    storage.WAL
	engine storage.Engine
	server *network.TCPServer
	slave  *replication.Slave
	master *replication.Master
	logger *zerolog.Logger
	stream chan []walPkg.LogData
}

func NewInitializer(cfg config.Config) (*Initializer, error) {
	stream := make(chan []walPkg.LogData, 1)

	logger, err := CreateLogger(cfg.Logging)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	wal, err := CreateWAL(cfg.WAL, logger, stream)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize wal: %w", err)
	}

	dbEngine, err := CreateEngine(cfg.Engine, logger, stream)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize engine: %w", err)
	}

	tcpServer, err := CreateNetwork(cfg.Network, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize network: %w", err)
	}

	replica, err := CreateReplica(cfg.Replication, cfg.WAL, logger, stream)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize replication: %w", err)
	}

	initializer := &Initializer{
		wal:    wal,
		engine: dbEngine,
		server: tcpServer,
		logger: logger,
		stream: stream,
	}

	initializer.initializeReplication(replica)

	return initializer, nil
}

func (i *Initializer) StartDatabase(ctx context.Context) error {
	defer close(i.stream)

	computeLayer := i.createComputeLayer()

	strg, err := i.createStorageLayer(ctx)
	if err != nil {
		return err
	}

	db := database.NewDatabase(computeLayer, strg, i.logger)

	group, groupCtx := errgroup.WithContext(ctx)
	if i.master != nil {
		group.Go(func() error {
			return i.master.Start(groupCtx)
		})
	}

	group.Go(func() error {
		strg.Start(groupCtx)

		return nil
	})

	group.Go(func() error {
		return i.server.HandleQueries(groupCtx, func(ctx context.Context, query []byte) []byte {
			response := db.HandleQuery(ctx, string(query))

			return []byte(response)
		})
	})

	return group.Wait()
}

func (i *Initializer) createComputeLayer() *compute.Compute {
	queryParser := compute.NewParser(i.logger)
	queryAnalyzer := compute.NewAnalyzer(i.logger)

	return compute.NewCompute(queryParser, queryAnalyzer, i.logger)
}

func (i *Initializer) createStorageLayer(context.Context) (*storage.Storage, error) {
	// if i.slave != nil {
	// i.slave.StartSynchronization(ctx) // TODO:
	// }

	strg, err := storage.NewStorage(i.engine, i.wal, i.storageReplicaSlave(), i.logger)
	if err != nil {
		i.logger.Error().Err(err).Msg("failed to initialize storage layer")

		return nil, err
	}

	return strg, nil
}

func (i *Initializer) initializeReplication(replica interface{}) {
	if replica == nil {
		return
	}

	if i.wal == nil {
		i.logger.Error().Msg("wal is required for replication")

		return
	}

	switch v := replica.(type) {
	case *replication.Slave:
		i.slave = v
	case *replication.Master:
		i.master = v
	default:
		i.logger.Error().Msg("incorrect replication type")
	}
}

func (i *Initializer) storageReplicaSlave() storage.Replica {
	if i.slave == nil {
		return nil
	}

	return i.slave
}
