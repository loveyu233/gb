package etcd

import (
	"fmt"
	"github.com/loveyu233/gb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/net/context"
	"testing"
)

func init() {
	gb.InitEtcd(gb.WithEtcdEndpointsOpt([]string{"127.0.0.1:2379"}), gb.WithEtcdUsernameOpt("root"), gb.WithEtcdPasswordOpt("etcd123"))
}

func TestA1(t *testing.T) {
	watch := gb.InsEtcd.Watch(context.Background(), "v/", clientv3.WithPrefix())
	for {
		select {
		case resp := <-watch:
			for _, event := range resp.Events {
				fmt.Println(event.Kv.ModRevision)
			}
		}
	}
}
