package discovery

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/clientv3"
)

type Register struct {
	EtcdAddrs   []string
	DialTimeout int
	clostCh     chan struct{}

	leasesID    clientv3.LeaseID
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse

	serInfo Server
	svrTTL  int64
	cli     *clientv3.Client
}

func NewRegister(etcdAddrs []string) *Register {
	return &Register{
		EtcdAddrs:   etcdAddrs,
		DialTimeout: 3,
	}
}

func (r *Register) Register(srvInfo Server, ttl int64) (chan<- struct{}, error) {
	var err error

	if strings.Split(srvInfo.Addr, ":")[0] == "" {
		return nil, errors.New("invalid ip address")
	}

	if r.cli, err = clientv3.New(clientv3.Config{
		Endpoints:   r.EtcdAddrs,
		DialTimeout: time.Duration(r.DialTimeout) * time.Second,
	}); err != nil {
		return nil, err
	}
	r.serInfo = srvInfo
	r.svrTTL = ttl

	if err = r.register(); err != nil {
		return nil, err
	}
	r.clostCh = make(chan struct{})

	go r.KeepAlive()

	return r.clostCh, nil
}

func (r *Register) register() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.DialTimeout)*time.Second)
	defer cancel()

	leaseRresp, err := r.cli.Grant(ctx, r.svrTTL)
	if err != nil {
		return err
	}
	r.leasesID = leaseRresp.ID
	if r.keepAliveCh, err = r.cli.KeepAlive(context.Background(), r.leasesID); err != nil {
		return err
	}

	data, err := json.Marshal(r.serInfo)
	if err != nil {
		return err
	}
	_, err = r.cli.Put(context.Background(), BuildRegisterPath(r.serInfo), string(data), r.cli.WithLease(r.leasesID))
	return err
}
func (r *Register) Stop() {
	r.clostCh <- struct{}{}
}

func (r *Register) UnRegister() error {
	_, err := r.cli.Delete(context.Background(), BuildRegisterPath(r.serInfo))
	return err
}

func (r *Register) KeepAlive() {
	ticker := time.NewTicker(time.Duration(r.svrTTL) * time.Second)
	for {
		select {
		case <-r.clostCh:
			if err := r.UnRegister(); err != nil {
				fmt.Print("err")
			}
			if _, err := r.cli.Remoke(context.Background(), r.leasesID); err != nil {
				fmt.Print("err")
			}
		case res := <-r.keepAliveCh:
			if res == nil {
				if err := r.register(); err != nil {
					fmt.Print("err")
				}
			}
		case <-ticker.C:
			if r.keepAliveCh != nil {
				if err := r.register(); err != nil {
					fmt.Print("err")
				}
			}
		}
	}
}

func (r *Register) UpdateHandler() http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		wightstr := req.URL.Query().Get("weight")
		weight, err := strconv.Atoi(wightstr)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		var update = func() error {
			r.serInfo.Weight = int64(weight)
			data, err := json.Marshal(r.serInfo)
			if err != nil {
				return err
			}
			_, err = r.cli.Put(context.Background(), BuildRegisterPath(r.serInfo), string(data), clientV3.WithLease(r.leasesID))
			return err
		}
		if err := update(); err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		_, _ = writer.Write([]byte("update succsess"))
	})
}

func (r *Register) GetServerInfo() (Server, error) {
	resp, err := r.cli.Get(context.Background(), BuildRegisterPath(r.serInfo))
	if err != nil {
		return r.serInfo, err
	}
	server := Server{}
	if resp.Count >= 1 {
		if err := json.Unmarshal(resp.Kvs[0].Value, &server); err != nil {
			return server, nil
		}
	}
	return server, err
}
