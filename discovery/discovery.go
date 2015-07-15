package discovery

import (
	"github.com/coreos/go-etcd/etcd"
	"path"
	"strings"
	"time"
)

func parse(rawurl string) (scheme string, addrs []string, prefix string) {
	parts := strings.SplitN(rawurl, "://", 2)
	// nodes:port,node2:port => nodes://node1:port,node2:port
	if len(parts) == 1 {
		scheme = "node"
		return
	}
	scheme = parts[0]
	parts = strings.SplitN(parts[1], "/", 2)
	addrs = strings.Split(parts[0], ",")
	if len(parts) == 2 {
		prefix = parts[1]
	}
	return
}

type Store struct {
	root string
	kv   *etcd.Client
}

type Event struct {
	Action string
	Node   *etcd.Node
}

func makeEndPoints(addrs []string, scheme string) (entries []string) {
	for _, addr := range addrs {
		entries = append(entries, "http"+"://"+addr)
	}
	return
}

func New(rawurl string, root string) (s *Store, err error) {
	_, addrs, prefix := parse(rawurl)
	s = &Store{path.Join(prefix, root), etcd.NewClient(makeEndPoints(addrs, "http"))}
	return
}

func (s *Store) Set(key string, value string, ttl time.Duration) error {
	_, err := s.kv.Set(path.Join(s.root, key), value, uint64(ttl.Seconds()))
	return err
}

func (s *Store) Delete(key string) error {
	_, err := s.kv.Delete(path.Join(s.root, key), false)
	return err
}

func (s *Store) WatchTree(key string, index uint64) (<-chan *Event, chan<- bool) {
	receiver := make(chan *etcd.Response)
	stopper := make(chan bool)
	go s.kv.Watch(path.Join(s.root, key), index, true, receiver, stopper)

	event := make(chan *Event)
	go func() {
		defer close(event)
		for resp := range receiver {
			event <- &Event{Action: resp.Action, Node: resp.Node}
		}
	}()

	return event, stopper
}
