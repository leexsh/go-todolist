package discovery

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/grpc/resolver"
)

type Server struct {
	name    string `json:"name"`
	Addr    string `json:"addr"`
	Version string `json:"version"`
	Weight  int64  `json:"weight"`
}

func BuildPrefix(server Server) string {
	if server.Version == "" {
		return fmt.Sprintf("/%s/", server.name)
	}
	return fmt.Sprintf("%s%s", server.name, server.Version)
}

func BuildRegisterPath(s Server) string {
	return fmt.Sprintf("%s%s", BuildPrefix(s), s.Addr)
}

func ParseValue(val []byte) (Server, error) {
	s := Server{}
	if err := json.Unmarshal(val, &s); err != nil {
		return s, err
	}
	return s, nil
}

func SpiltPath(path string) (Server, error) {
	s := Server{}
	strs := strings.Split(path, "/")
	if len(strs) == 0 {
		return s, errors.New("invalid path")
	}
	s.Addr = strs[len(strs)-1]
	return s, nil
}

func Exist(l []resolver.Address, addr resolver.Address) bool {
	for i := range l {
		if l[i].Addr == addr.Addr {
			return true
		}
	}
	return false
}

func Remove(l []resolver.Address, addr resolver.Address) ([]resolver.Address, bool) {
	for i := range l {
		if l[i].Addr == addr.Addr {
			l[i] = l[len(l)-1]
			return l[:len(l)-1], true
		}
	}
	return nil, false
}
func BuildResolverUrl(app string) string {
	return "etcd" + ":///" + app
}
