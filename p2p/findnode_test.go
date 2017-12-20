package p2p

import (
	"bytes"
	"github.com/UnrulyOS/go-unruly/assert"
	"github.com/UnrulyOS/go-unruly/log"
	"github.com/UnrulyOS/go-unruly/p2p/node"
	"github.com/UnrulyOS/go-unruly/p2p/pb"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestFindNodeProtocolCore(t *testing.T) {

	// let there be 3 nodes - node1, node2 and node 3
	node1Local, node1Remote := GenerateTestNode(t)
	node2Local, node2Remote := GenerateTestNode(t)
	_, node3Remote := GenerateTestNode(t)

	// node 1 know about node 2 and node 3
	node1Local.GetSwarm().RegisterNode(node.NewRemoteNodeData(node2Remote.String(), node2Remote.TcpAddress()))
	node1Local.GetSwarm().RegisterNode(node.NewRemoteNodeData(node3Remote.String(), node3Remote.TcpAddress()))

	// node 2 knows node 1
	node2Local.GetSwarm().RegisterNode(node.NewRemoteNodeData(node1Remote.String(), node1Remote.TcpAddress()))

	// node 2 doesn't know about node 3 and asks node 1 to find it
	reqId := []byte(uuid.New().String())
	callback := make(chan *pb.FindNodeResp)
	node2Local.GetSwarm().getFindNodeProtocol().Register(callback)
	node2Local.GetSwarm().getFindNodeProtocol().FindNode(reqId, node1Remote.String(), node3Remote.String())

Loop:
	for {
		select {
		case c := <-callback:
			if bytes.Equal(c.GetMetadata().ReqId, reqId) {

				log.Info("Got find response: `%v`", c)

				assert.NotNil(t, "expected response slice")
				assert.True(t, len(c.NodeInfos) >= 1, "expected at least 1 node")

				for _, d := range c.NodeInfos {
					if bytes.Equal(d.NodeId,node3Remote.Id()) {
						log.Info("Found node 3 :-)")
						break Loop
					}
				}
			}

		case <-time.After(time.Second * 10):
			t.Fatalf("Test timed out")
		}

	}
}