package beacon

import (
	"encoding/json"
	"net"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/discover/portalwire"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/portalnetwork/storage"
	"github.com/stretchr/testify/assert"
)

func setupBeaconNetwork(addr string, bootNodes []*enode.Node) (*BeaconNetwork, error) {
	glogger := log.NewGlogHandler(log.NewTerminalHandler(os.Stderr, true))
	slogVerbosity := log.FromLegacyLevel(5)
	glogger.Verbosity(slogVerbosity)
	log.SetDefault(log.NewLogger(glogger))
	conf := discover.DefaultPortalProtocolConfig()
	if addr != "" {
		conf.ListenAddr = addr
	}
	if bootNodes != nil {
		conf.BootstrapNodes = bootNodes
	}

	addr1, err := net.ResolveUDPAddr("udp", conf.ListenAddr)
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", addr1)
	if err != nil {
		return nil, err
	}

	privKey, err := crypto.GenerateKey()
	if err != nil {
		panic("couldn't generate key: " + err.Error())
	}

	discCfg := discover.Config{
		PrivateKey:  privKey,
		NetRestrict: conf.NetRestrict,
		Bootnodes:   conf.BootstrapNodes,
	}

	nodeDB, err := enode.OpenDB(conf.NodeDBPath)
	if err != nil {
		return nil, err
	}

	localNode := enode.NewLocalNode(nodeDB, privKey)
	localNode.SetFallbackIP(net.IP{127, 0, 0, 1})
	localNode.Set(discover.Tag)

	var addrs []net.Addr
	if conf.NodeIP != nil {
		localNode.SetStaticIP(conf.NodeIP)
	} else {
		addrs, err = net.InterfaceAddrs()

		if err != nil {
			return nil, err
		}

		for _, address := range addrs {
			// check ip addr is loopback addr
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					localNode.SetStaticIP(ipnet.IP)
					break
				}
			}
		}
	}

	discV5, err := discover.ListenV5(conn, localNode, discCfg)
	if err != nil {
		return nil, err
	}

	contentQueue := make(chan *discover.ContentElement, 50)

	portalProtocol, err := discover.NewPortalProtocol(conf, string(portalwire.BeaconLightClientNetwork), privKey, conn, localNode, discV5, &storage.MockStorage{Db: make(map[string][]byte)}, contentQueue)
	if err != nil {
		return nil, err
	}

	return NewBeaconNetwork(portalProtocol), nil
}

func TestBeaconNetworkContent(t *testing.T) {
	//logger := testlog.Logger(t, log.LvlTrace)
	node1, err := setupBeaconNetwork(":6677", nil)
	assert.NoError(t, err)
	//node1.log = logger
	err = node1.Start()
	assert.NoError(t, err)

	time.Sleep(10 * time.Second)

	node2, err := setupBeaconNetwork(":6776", []*enode.Node{node1.portalProtocol.Self()})
	assert.NoError(t, err)
	//node2.log = logger
	err = node2.Start()
	assert.NoError(t, err)

	time.Sleep(10 * time.Second)

	filePath := "testdata/light_client_updates_by_range.json"
	jsonStr, _ := os.ReadFile(filePath)

	var result map[string]interface{}
	err = json.Unmarshal(jsonStr, &result)
	assert.NoError(t, err)

	var kBytes []byte
	var vBytes []byte

	id := node1.portalProtocol.Self().ID()
	for _, v := range result {
		kBytes, _ = hexutil.Decode(v.(map[string]interface{})["content_key"].(string))
		vBytes, _ = hexutil.Decode(v.(map[string]interface{})["content_value"].(string))
	}

	contentKey := storage.NewContentKey(LightClientUpdate, kBytes).Encode()
	gossip, err := node1.portalProtocol.NeighborhoodGossip(&id, [][]byte{contentKey}, [][]byte{vBytes})
	assert.NoError(t, err)

	assert.Equal(t, 1, gossip)

	time.Sleep(1 * time.Second)

	contentId := node2.portalProtocol.ToContentId(contentKey)
	get, err := node2.portalProtocol.Get(contentKey, contentId)
	assert.NoError(t, err)
	assert.Equal(t, vBytes, get)

	//filePath1 := "testdata/light_client_finality_update.json"
	//jsonStr1, _ := os.ReadFile(filePath1)
	//
	//var result1 map[string]interface{}
	//err = json.Unmarshal(jsonStr1, &result1)
	//assert.NoError(t, err)
	//
	//for _, v := range result1 {
	//	kBytes1, _ := hexutil.Decode(v.(map[string]interface{})["content_key"].(string))
	//	vBytes1, _ := hexutil.Decode(v.(map[string]interface{})["content_value"].(string))
	//
	//	contentKey1 := storage.NewContentKey(LightClientFinalityUpdate, kBytes1).Encode()
	//
	//	id := node1.portalProtocol.Self().ID()
	//	gossip, err := node1.portalProtocol.RandomGossip(&id, [][]byte{contentKey1}, [][]byte{vBytes1})
	//	assert.NoError(t, err)
	//
	//	assert.Equal(t, 1, gossip)
	//
	//	time.Sleep(1 * time.Second)
	//
	//	contentId1 := node2.portalProtocol.ToContentId(contentKey1)
	//	get1, err := node2.portalProtocol.Get(contentKey1, contentId1)
	//	assert.NoError(t, err)
	//	assert.Equal(t, vBytes1, get1)
	//}
}
