package node

import (
	"docker.io/go-docker/api/types"
	"github.com/appcelerator/amp/pkg/docker"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server is used to implement log.LogServer
type Server struct {
	Docker *docker.Docker
}

// GetNodes implements Node.GetNodes
func (s *Server) NodeGet(ctx context.Context, in *GetRequest) (*GetReply, error) {
	list, err := s.Docker.NodeList(ctx, types.NodeListOptions{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	nodeList := &GetReply{}
	for _, item := range list {
		node := &NodeEntry{
			Id:           item.ID,
			Name:         item.Spec.Name,
			Hostname:     item.Description.Hostname,
			Role:         string(item.Spec.Role),
			Architecture: item.Description.Platform.Architecture,
			Os:           item.Description.Platform.OS,
			Engine:       item.Description.Engine.EngineVersion,
			Addr:         item.Status.Addr,
			Status:       string(item.Status.State),
			Availability: string(item.Spec.Availability),
			Labels:       item.Spec.Annotations.Labels,
		}
		if item.ManagerStatus != nil {
			node.Leader = item.ManagerStatus.Leader
			node.Reachability = string(item.ManagerStatus.Reachability)
		}
		nodeList.Entries = append(nodeList.Entries, node)
	}
	return nodeList, nil
}
