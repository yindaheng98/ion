package main

import (
	"net/http"
	_ "net/http/pprof"

	log "github.com/pion/ion-log"
	conf "github.com/pion/ion/pkg/conf/avp"
	"github.com/pion/ion/pkg/discovery"
	"github.com/pion/ion/pkg/node/avp"
)

func init() {
	fixByFile := []string{"asm_amd64.s", "proc.go", "icegatherer.go"}
	fixByFunc := []string{}
	log.Init(conf.Avp.Log.Level, fixByFile, fixByFunc)

	avp.InitAVP(conf.Avp)
}

func main() {
	log.Infof("--- Starting AVP Node ---")

	if conf.Global.Pprof != "" {
		go func() {
			log.Infof("Start pprof on %s", conf.Global.Pprof)
			err := http.ListenAndServe(conf.Global.Pprof, nil)
			if err != nil {
				log.Errorf("http.ListenAndServe err=%v", err)
			}
		}()
	}

	serviceNode := discovery.NewServiceNode(conf.Etcd.Addrs, conf.Global.Dc)
	serviceNode.RegisterNode("avp", "node-avp", "avp-channel-id")
	avp.Init(conf.Global.Dc, serviceNode.NodeInfo().Info["id"], conf.Nats.URL)

	select {}
}
