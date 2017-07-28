package syncer

import (
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
	"context"
	"time"
	"os"
)

type MysqlSyncer struct {
	streamer *replication.BinlogStreamer
}

func (s *MysqlSyncer) Init(serverId uint32, flavor string, host string, port uint16, user string, password string) (err error) {
	rsyncer := replication.NewBinlogSyncer(&replication.BinlogSyncerConfig{
		ServerID: serverId,
		Flavor:   flavor,
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
	})

	s.streamer, err = rsyncer.StartSync(mysql.Position{})
	return err
}

func (s *MysqlSyncer) Dump() {
	for {
		ev, _ := s.streamer.GetEvent(context.Background())
		ev.Dump(os.Stdout)
	}

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		ev, err := s.streamer.GetEvent(ctx)
		cancel()

		if err == context.DeadlineExceeded {
			continue
		}

		ev.Dump(os.Stdout)
	}
}
