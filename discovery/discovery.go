package discovery

import (
	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"path"
	"strings"
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

type Options struct {
	Root   string
	Config *store.Config
}

type Store struct {
	root string
	kv   store.Store
}

func New(rawurl string, opts *Options) (s *Store, err error) {
	scheme, addrs, prefix := parse(rawurl)

	root := path.Join(prefix, opts.Root)

	kv, err := libkv.NewStore(
		store.Backend(scheme),
		addrs,
		opts.Config,
	)
	s = &Store{root, kv}

	return
}

func (s *Store) Put(key string, value []byte, options *store.WriteOptions) error {
	return s.kv.Put(path.Join(s.root, key), value, options)
}
