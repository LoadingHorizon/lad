package syncer

import (
	"encoding/json"
	"fmt"
	"github.com/siddontang/go-mysql/canal"
	"github.com/siddontang/go-mysql/mysql"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type MysqlSyncer struct {
	canal *canal.Canal
	canal.DummyEventHandler
	done      chan bool
	syncpoint string
}

func (s *MysqlSyncer) Init(serverId uint32, flavor string, host string, port uint16,
	user string, password string, database string, tables []string, syncpoint string) (err error) {
	s.syncpoint = syncpoint
	s.canal, err = canal.NewCanal(&canal.Config{
		ServerID: serverId,
		Flavor:   flavor,
		Addr:     fmt.Sprintf("%v:%v", host, port),
		User:     user,
		Password: password,
		Dump: canal.DumpConfig{
			TableDB: database,
			Tables:  tables,
		},
	})
	s.canal.SetEventHandler(s)
	return err
}

func (s *MysqlSyncer) Dump() error {
	syncpoint, err := s.LoadSyncpoint()
	if err != nil {
		return err
	}
	return s.canal.StartFrom(*syncpoint)
}

func (s *MysqlSyncer) SaveSyncpoint() error {
	syncpointBk := fmt.Sprintf("%v.bk", s.syncpoint)
	if err := os.Rename(s.syncpoint, syncpointBk); err != nil && !os.IsNotExist(err) {
		return err
	}

	fp, err := os.Create(s.syncpoint)
	if err != nil {
		return err
	}
	defer fp.Close()

	_, err = fp.WriteString(fmt.Sprintf("%v\t%v", s.canal.SyncedPosition().Name, s.canal.SyncedPosition().Pos))
	return err
}

func (s *MysqlSyncer) LoadSyncpoint() (*mysql.Position, error) {
	buf, err := ioutil.ReadFile(s.syncpoint)
	if os.IsNotExist(err) {
		return &mysql.Position{Name: "", Pos: 0}, nil
	}
	if err != nil {
		return nil, err
	}
	txt := string(buf)
	values := strings.Split(txt, "\t")
	if len(values) != 2 {
		return nil, woklog.Errorf("len(values) should be 2. txt[]", values)
	}

	pos, err := strconv.ParseUint(values[1], 10, 32)
	if err != nil {
		return nil, err
	}

	return &mysql.Position{Name: values[0], Pos: uint32(pos)}, nil
}

// @DummyEventHandler.OnRow
func (s *MysqlSyncer) OnRow(event *canal.RowsEvent) error {
	eventObj := map[string]interface{}{}
	eventObj["action"] = event.Action
	eventObj["table"] = event.Table.Name

	if event.Action == "insert" || event.Action == "delete" {
		values := map[string]interface{}{}
		for i := range event.Table.Columns {
			values[event.Table.Columns[i].Name] = event.Rows[0][i]
		}
		eventObj["values"] = values
	} else if event.Action == "update" {
		{
			values := map[string]interface{}{}
			for i := range event.Table.Columns {
				values[event.Table.Columns[i].Name] = event.Rows[0][i]
			}
			eventObj["oldValues"] = values
		}
		{
			values := map[string]interface{}{}
			for i := range event.Table.Columns {
				values[event.Table.Columns[i].Name] = event.Rows[1][i]
			}
			eventObj["values"] = values
		}
	}

	{
		primaryKey := []string{}
		for _, i := range event.Table.PKColumns {
			primaryKey = append(primaryKey, event.Table.Columns[i].Name)
		}
		eventObj["primaryKey"] = primaryKey
	}

	{
		syncpoint := map[string]interface{}{}
		syncpoint["filename"] = s.canal.SyncedPosition().Name
		syncpoint["offset"] = s.canal.SyncedPosition().Pos
		eventObj["syncpoint"] = syncpoint
	}

	jsonByte, err := json.Marshal(eventObj)
	if err != nil {
		fmt.Println(err)
	}

	acclog.Info(string(jsonByte))
	return nil
}

// @DummyEventHandler.OnXID
func (s *MysqlSyncer) OnXID(mysql.Position) error {
	return s.SaveSyncpoint()
}

// @DummyEventHandler.String
func (s *MysqlSyncer) String() string {
	return "MysqlSyncer"
}
