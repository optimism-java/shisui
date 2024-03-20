package beacon

import (
	"github.com/ethereum/go-ethereum/p2p/discover"
)

type API struct {
	*discover.PortalProtocolAPI
}

func (p *API) BeaconRoutingTableInfo() *discover.RoutingTableInfo {
	return p.RoutingTableInfo()
}

func (p *API) BeaconAddEnr(enr string) (bool, error) {
	return p.AddEnr(enr)
}

func (p *API) BeaconGetEnr(nodeId string) (string, error) {
	return p.GetEnr(nodeId)
}

func (p *API) BeaconDeleteEnr(nodeId string) (bool, error) {
	return p.DeleteEnr(nodeId)
}

func (p *API) BeaconLookupEnr(nodeId string) (string, error) {
	return p.LookupEnr(nodeId)
}

func (p *API) BeaconPing(enr string) (*discover.PortalPongResp, error) {
	return p.Ping(enr)
}

func (p *API) BeaconFindNodes(enr string, distances []uint) ([]string, error) {
	return p.FindNodes(enr, distances)
}

func (p *API) BeaconFindContent(enr string, contentKey string) (interface{}, error) {
	return p.FindContent(enr, contentKey)
}

func (p *API) BeaconOffer(enr string, contentKey string, contentValue string) (string, error) {
	return p.Offer(enr, contentKey, contentValue)
}

func (p *API) BeaconRecursiveFindNodes(nodeId string) ([]string, error) {
	return p.RecursiveFindNodes(nodeId)
}

func (p *API) BeaconRecursiveFindContent(contentKeyHex string) (*discover.ContentInfo, error) {
	return p.RecursiveFindContent(contentKeyHex)
}

func (p *API) BeaconLocalContent(contentKeyHex string) (string, error) {
	return p.LocalContent(contentKeyHex)
}

func (p *API) BeaconStore(contentKeyHex string, contextHex string) (bool, error) {
	return p.Store(contentKeyHex, contextHex)
}

func (p *API) BeaconGossip(contentKeyHex, contentHex string) (int, error) {
	return p.Gossip(contentKeyHex, contentHex)
}

func (p *API) BeaconTraceRecursiveFindContent(contentKeyHex string) {
	p.TraceRecursiveFindContent(contentKeyHex)
}

func NewBeaconNetworkAPI(BeaconAPI *discover.PortalProtocolAPI) *API {
	return &API{
		BeaconAPI,
	}
}
