package discover

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/optimism-java/utp-go"
	"github.com/prysmaticlabs/go-bitfield"

	"github.com/ethereum/go-ethereum/internal/testlog"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/discover/portalwire"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

type MockStorage struct {
	db map[string][]byte
}

// func (m *MockStorage) ContentId(contentKey []byte) []byte {
// 	return crypto.Keccak256(contentKey)
// }

func (m *MockStorage) Get(contentId []byte) ([]byte, error) {
	if content, ok := m.db[string(contentId)]; ok {
		return content, nil
	}
	return nil, ContentNotFound
}

func (m *MockStorage) Put(contentId []byte, content []byte) error {
	m.db[string(contentId)] = content
	return nil
}

func setupLocalPortalNode(addr string, bootNodes []*enode.Node) (*PortalProtocol, error) {
	conf := DefaultPortalProtocolConfig()
	if addr != "" {
		conf.ListenAddr = addr
	}
	if bootNodes != nil {
		conf.BootstrapNodes = bootNodes
	}

	contentQueue := make(chan *ContentElement, 50)
	portalProtocol, err := NewPortalProtocol(conf, portalwire.HistoryNetwork, newkey(), &MockStorage{db: make(map[string][]byte)}, contentQueue)
	if err != nil {
		return nil, err
	}

	return portalProtocol, nil
}

func TestPortalWireProtocolUdp(t *testing.T) {
	node1, err := setupLocalPortalNode(":8777", nil)
	assert.NoError(t, err)
	node1.log = testlog.Logger(t, log.LvlTrace)
	err = node1.Start()
	assert.NoError(t, err)

	node2, err := setupLocalPortalNode(":8778", []*enode.Node{node1.localNode.Node()})
	assert.NoError(t, err)
	node2.log = testlog.Logger(t, log.LvlTrace)
	err = node2.Start()
	assert.NoError(t, err)

	node3, err := setupLocalPortalNode(":8779", []*enode.Node{node1.localNode.Node()})
	assert.NoError(t, err)
	node3.log = testlog.Logger(t, log.LvlTrace)
	err = node3.Start()
	assert.NoError(t, err)
	time.Sleep(10 * time.Second)

	assert.Equal(t, 2, len(node1.table.Nodes()))
	assert.Equal(t, 2, len(node2.table.Nodes()))
	assert.Equal(t, 2, len(node3.table.Nodes()))

	node1Addr, _ := utp.ResolveUTPAddr("utp", "127.0.0.1:8777")
	node2Addr, _ := utp.ResolveUTPAddr("utp", "127.0.0.1:8778")
	node3Addr, _ := utp.ResolveUTPAddr("utp", "127.0.0.1:8779")

	cid := uint32(12)
	cliSendMsgWithCid := "there are connection id : 12!"
	cliSendMsgWithRandomCid := "there are connection id: random!"
	//
	serverEchoWithCid := "accept connection sends back msg: echo"
	//serverEchoWithRandomCid := "ccept connection with random cid sends msg: echo"

	largeTestContent := make([]byte, 1199)
	_, err = rand.Read(largeTestContent)
	assert.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		var acceptConn *utp.Conn
		defer func() {
			wg.Done()
			_ = acceptConn.Close()
		}()
		acceptConn, err := node3.utp.AcceptUTPWithConnId(cid)
		if err != nil {
			panic(err)
		}
		buf := make([]byte, 100)
		n, err := acceptConn.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		assert.Equal(t, cliSendMsgWithCid, string(buf[:n]))
		_, _ = acceptConn.Write([]byte(serverEchoWithCid))
	}()
	go func() {
		var randomConnIdConn net.Conn
		defer func() {
			wg.Done()
			_ = randomConnIdConn.Close()
		}()
		randomConnIdConn, err := node1.utp.Accept()
		if err != nil {
			panic(err)
		}
		buf := make([]byte, 100)
		n, err := randomConnIdConn.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		assert.Equal(t, cliSendMsgWithRandomCid, string(buf[:n]))

		_, _ = randomConnIdConn.Write(largeTestContent)
	}()

	go func() {
		var connWithConnId net.Conn
		defer func() {
			wg.Done()
			if connWithConnId != nil {
				_ = connWithConnId.Close()
			}
		}()
		connWithConnId, err := utp.DialUTPOptions("utp", node2Addr, node3Addr, utp.WithConnId(cid), utp.WithSocketManager(node2.utpSm))
		if err != nil {
			panic(err)
		}
		_, err = connWithConnId.Write([]byte("there are connection id : 12!"))
		if err != nil && err != io.EOF {
			panic(err)
		}
		buf := make([]byte, 100)
		n, err := connWithConnId.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		assert.Equal(t, serverEchoWithCid, string(buf[:n]))
	}()
	go func() {
		var randomConnIdConn net.Conn
		defer func() {
			wg.Done()
			//_ = randomConnIdConn.Close()
		}()
		randomConnIdConn, err := utp.DialUTPOptions("utp", node2Addr, node1Addr, utp.WithSocketManager(node2.utpSm))
		if err != nil && err != io.EOF {
			panic(err)
		}
		_, err = randomConnIdConn.Write([]byte(cliSendMsgWithRandomCid))
		if err != nil {
			panic(err)
		}

		data := make([]byte, 0)
		buf := make([]byte, 1024)
		for {
			var n int
			n, err = randomConnIdConn.Read(buf)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
			}
			data = append(data, buf[:n]...)
		}
		assert.Equal(t, largeTestContent, data)
	}()
	wg.Wait()
	fmt.Println("done")
}

func TestPortalWireProtocol(t *testing.T) {
	node1, err := setupLocalPortalNode(":7777", nil)
	assert.NoError(t, err)
	node1.log = testlog.Logger(t, log.LvlTrace)
	err = node1.Start()
	assert.NoError(t, err)
	fmt.Println(node1.localNode.Node().String())

	node2, err := setupLocalPortalNode(":7778", []*enode.Node{node1.localNode.Node()})
	assert.NoError(t, err)
	node2.log = testlog.Logger(t, log.LvlTrace)
	err = node2.Start()
	assert.NoError(t, err)
	fmt.Println(node2.localNode.Node().String())

	node3, err := setupLocalPortalNode(":7779", []*enode.Node{node1.localNode.Node()})
	assert.NoError(t, err)
	node3.log = testlog.Logger(t, log.LvlTrace)
	err = node3.Start()
	assert.NoError(t, err)
	fmt.Println(node3.localNode.Node().String())
	time.Sleep(10 * time.Second)

	assert.Equal(t, 2, len(node1.table.Nodes()))
	assert.Equal(t, 2, len(node2.table.Nodes()))
	assert.Equal(t, 2, len(node3.table.Nodes()))

	slices.ContainsFunc(node1.table.Nodes(), func(n *enode.Node) bool {
		return n.ID() == node2.localNode.Node().ID()
	})
	slices.ContainsFunc(node1.table.Nodes(), func(n *enode.Node) bool {
		return n.ID() == node3.localNode.Node().ID()
	})

	slices.ContainsFunc(node2.table.Nodes(), func(n *enode.Node) bool {
		return n.ID() == node1.localNode.Node().ID()
	})
	slices.ContainsFunc(node2.table.Nodes(), func(n *enode.Node) bool {
		return n.ID() == node3.localNode.Node().ID()
	})

	slices.ContainsFunc(node3.table.Nodes(), func(n *enode.Node) bool {
		return n.ID() == node1.localNode.Node().ID()
	})
	slices.ContainsFunc(node3.table.Nodes(), func(n *enode.Node) bool {
		return n.ID() == node2.localNode.Node().ID()
	})

	err = node1.storage.Put(node1.toContentId([]byte("test_key")), []byte("test_value"))
	assert.NoError(t, err)

	flag, content, err := node2.findContent(node1.localNode.Node(), []byte("test_key"))
	assert.NoError(t, err)
	assert.Equal(t, portalwire.ContentRawSelector, flag)
	assert.Equal(t, []byte("test_value"), content)

	flag, content, err = node2.findContent(node3.localNode.Node(), []byte("test_key"))
	assert.NoError(t, err)
	assert.Equal(t, portalwire.ContentEnrsSelector, flag)
	assert.Equal(t, 1, len(content.([]*enode.Node)))
	assert.Equal(t, node1.localNode.Node().ID(), content.([]*enode.Node)[0].ID())

	// create a byte slice of length 1199 and fill it with random data
	// this will be used as a test content
	largeTestContent := make([]byte, 2000)
	_, err = rand.Read(largeTestContent)
	assert.NoError(t, err)

	err = node1.storage.Put(node1.toContentId([]byte("large_test_key")), largeTestContent)
	assert.NoError(t, err)

	flag, content, err = node2.findContent(node1.localNode.Node(), []byte("large_test_key"))
	assert.NoError(t, err)
	assert.Equal(t, largeTestContent, content)
	assert.Equal(t, portalwire.ContentConnIdSelector, flag)

	testEntry1 := &ContentEntry{
		ContentKey: []byte("test_entry1"),
		Content:    []byte("test_entry1_content"),
	}

	testEntry2 := &ContentEntry{
		ContentKey: []byte("test_entry2"),
		Content:    []byte("test_entry2_content"),
	}

	testTransientOfferRequest := &TransientOfferRequest{
		Contents: []*ContentEntry{testEntry1, testEntry2},
	}

	offerRequest := &OfferRequest{
		Kind:    TransientOfferRequestKind,
		Request: testTransientOfferRequest,
	}

	contentKeys, err := node1.offer(node3.localNode.Node(), offerRequest)
	assert.Equal(t, uint64(2), bitfield.Bitlist(contentKeys).Count())
	assert.NoError(t, err)

	contentElement := <-node3.contentQueue
	assert.Equal(t, node1.localNode.Node().ID(), contentElement.Node)
	assert.Equal(t, testEntry1.ContentKey, contentElement.ContentKeys[0])
	assert.Equal(t, testEntry1.Content, contentElement.Contents[0])
	assert.Equal(t, testEntry2.ContentKey, contentElement.ContentKeys[1])
	assert.Equal(t, testEntry2.Content, contentElement.Contents[1])
}
