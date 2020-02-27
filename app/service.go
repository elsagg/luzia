package app

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"time"

	"github.com/elsagg/luzia/datacell"
	"github.com/elsagg/luzia/domain"
	"google.golang.org/grpc"
)

// Server is a http server
type Server struct {
	RpcServer *grpc.Server
}

// GetLatestDataCell will find the latest version of a data cell and return it
func (s *Server) GetLatestDataCell(ctx context.Context, request *datacell.GetLatestDataCellRequest) (*datacell.GetDataCellResponse, error) {
	log.Printf("Received: %v", request)

	d := domain.NewDomain(request.GetDataSource())

	dc, err := d.GetCellLatest(request.GetRowKey(), request.GetColumnKey())

	if err != nil {
		return nil, err
	}

	var body map[string]interface{}

	err = dc.UnmarshalBody(&body)

	if err != nil {
		return nil, err
	}

	js, err := json.Marshal(body)

	if err != nil {
		return nil, err
	}

	return &datacell.GetDataCellResponse{
		AddedID:   dc.AddedID,
		RowKey:    dc.RowKey,
		ColumnKey: dc.ColumnKey,
		RefKey:    dc.RefKey,
		Body:      string(js),
		CreatedAt: dc.CreatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// GetDataCell will find a data cell and return it
func (s *Server) GetDataCell(ctx context.Context, request *datacell.GetDataCellRequest) (*datacell.GetDataCellResponse, error) {
	log.Printf("Received: %v", request)

	d := domain.NewDomain(request.GetDataSource())

	dc, err := d.GetCell(request.GetRowKey(), request.GetColumnKey(), request.GetRefKey())

	if err != nil {
		return nil, err
	}

	var body map[string]interface{}

	err = dc.UnmarshalBody(&body)

	if err != nil {
		return nil, err
	}

	js, err := json.Marshal(body)

	if err != nil {
		return nil, err
	}

	return &datacell.GetDataCellResponse{
		AddedID:   dc.AddedID,
		RowKey:    dc.RowKey,
		ColumnKey: dc.ColumnKey,
		RefKey:    dc.RefKey,
		Body:      string(js),
		CreatedAt: dc.CreatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// PutDataCell will put a new data cell in the specified data source
func (s *Server) PutDataCell(ctx context.Context, request *datacell.PutDataCellRequest) (*datacell.GetDataCellResponse, error) {
	log.Printf("Received: %v", request)

	d := domain.NewDomain(request.GetDataSource())

	var body map[string]interface{}

	err := json.Unmarshal([]byte(request.GetBody()), &body)

	if err != nil {
		return nil, err
	}

	dc, err := d.PutCell(request.GetRowKey(), request.GetColumnKey(), body)

	if err != nil {
		return nil, err
	}

	var rtBody map[string]interface{}

	err = dc.UnmarshalBody(&rtBody)

	if err != nil {
		return nil, err
	}

	js, err := json.Marshal(rtBody)

	if err != nil {
		return nil, err
	}

	return &datacell.GetDataCellResponse{
		AddedID:   dc.AddedID,
		RowKey:    dc.RowKey,
		ColumnKey: dc.ColumnKey,
		RefKey:    dc.RefKey,
		Body:      string(js),
		CreatedAt: dc.CreatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// Start will start the server on the designated port
func (s *Server) Start(port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s.RpcServer = grpc.NewServer()
	datacell.RegisterDataCellServiceServer(s.RpcServer, s)
	if err := s.RpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// NewServer creates a new server
func NewServer() *Server {
	return &Server{}
}
