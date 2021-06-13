package tuntap

import (
	"encoding/json"

	"github.com/yggdrasil-network/yggdrasil-go/src/admin"
)

type GetTUNRequest struct{}
type GetTUNResponse map[string]TUNEntry

type TUNEntry struct {
	MTU uint64 `json:"mtu"`
}

func (t *TunAdapter) getTUNHandler(req *GetTUNRequest, res *GetTUNResponse) error {
	*res = GetTUNResponse{
		t.Name(): TUNEntry{
			MTU: t.MTU(),
		},
	}
	return nil
}

func (t *TunAdapter) SetupAdminHandlers(a *admin.AdminSocket) {
	_ = a.AddHandler("getTunTap", []string{}, func(in json.RawMessage) (interface{}, error) {
		req := &GetTUNRequest{}
		res := &GetTUNResponse{}
		if err := json.Unmarshal(in, &req); err != nil {
			return nil, err
		}
		if err := t.getTUNHandler(req, res); err != nil {
			return nil, err
		}
		return res, nil
	})
	//_ = a.AddHandler("getNodeInfo", []string{"key"}, t.proto.nodeinfo.nodeInfoAdminHandler)
	//_ = a.AddHandler("debug_remoteGetSelf", []string{"key"}, t.proto.getSelfHandler)
	//_ = a.AddHandler("debug_remoteGetPeers", []string{"key"}, t.proto.getPeersHandler)
	//_ = a.AddHandler("debug_remoteGetDHT", []string{"key"}, t.proto.getDHTHandler)
}
