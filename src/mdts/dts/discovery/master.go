package discovery

import (
	"context"
	"encoding/json"
	"log"
	dproto "mdts/protocols/discovery"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
)

// Master watch etcd to implement service discovery
// etcd path: brokers/<type>_<hash: hash(time)>
type Master struct {
	Path   string
	Nodes  map[string]*BrokerMap
	Client *clientv3.Client
}

// BrokerMap : Master.Nodes { "type": BrokerMap.Brokers{ "hash": BrokerInfo}}
type BrokerMap struct {
	Type    string
	Brokers *SortedMap
}

func NewBrokerMap(key string) *BrokerMap {
	return &BrokerMap{
		Type:    key,
		Brokers: NewSortedMap(),
	}
}

func (bm *BrokerMap) AddBroker(hash string, info *dproto.BrokerInfo) {
	bm.Brokers.AddKV(hash, info)
}

func (bm *BrokerMap) DeleteBroker(hash string) (info *dproto.BrokerInfo) {
	v := bm.Brokers.DeleteKV(hash)
	info, _ = v.(*dproto.BrokerInfo)
	return info
}

func (bm *BrokerMap) Len() int {
	return bm.Brokers.Len()
}

func NewMaster(endpoints []string, watchPath string) (*Master, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 2 * time.Second,
	})

	if err != nil {
		return nil, err
	}

	master := &Master{
		Path:   watchPath,
		Nodes:  make(map[string]*BrokerMap),
		Client: cli,
	}

	return master, err
}

func (m *Master) AddBroker(info *dproto.BrokerInfo) {
	bm, ok := m.Nodes[info.Type]
	if !ok {
		bm = NewBrokerMap(info.Type)
		m.Nodes[info.Type] = bm
	}

	bm.AddBroker(info.Hash, info)
}

func (m *Master) DeleteBroker(Type string, hash string) {
	bm, ok := m.Nodes[Type]
	if !ok {
		return
	}

	bm.DeleteBroker(hash)
	if bm.Len() == 0 {
		delete(m.Nodes, Type)
	}
}

func GetBrokerInfo(ev *clientv3.Event) *dproto.BrokerInfo {
	info := &dproto.BrokerInfo{}
	err := json.Unmarshal([]byte(ev.Kv.Value), info)
	if err != nil {
		log.Println(err)
	}
	return info
}

func GetBrokerInfoFromKey(key string) (Type string, hash string) {
	slashPostion := strings.IndexByte(key, '/')
	key = key[slashPostion+1:]
	items := strings.Split(key, "_")

	return items[0], items[1]
}

func (m *Master) WatchNodes() {
	rch := m.Client.Watch(context.Background(), m.Path, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case clientv3.EventTypePut:
				log.Printf("[%s] %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				info := GetBrokerInfo(ev)
				m.AddBroker(info)
			case clientv3.EventTypeDelete:
				log.Printf("[%s] %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				Type, hash := GetBrokerInfoFromKey(string(ev.Kv.Key))
				m.DeleteBroker(Type, hash)
			}
		}
	}
}
