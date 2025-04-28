package server

import (
	"github.com/danielealbano/svdb/shared/collection"
	"github.com/danielealbano/svdb/shared/grpc_server"
	shared_proto_build_frontend "github.com/danielealbano/svdb/shared/proto/build/frontend"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type frontendGrpcServerImplementation struct {
	shared_proto_build_frontend.UnimplementedFrontendServer
	collectionConfig *shared_collection.CollectionConfig
}

func vectorToPB(v []float32) *shared_proto_build_frontend.Vector {
	return &shared_proto_build_frontend.Vector{Values: v}
}

func RegisterFrontendGrpcServerImplementation(
	server *shared_grpc_server.GrpcServer,
	collectionConfig *shared_collection.CollectionConfig) {
	shared_proto_build_frontend.RegisterFrontendServer(server.GrpcServer, &frontendGrpcServerImplementation{
		collectionConfig: collectionConfig,
	})
}

func (s *frontendGrpcServerImplementation) Search(
	_ context.Context,
	req *shared_proto_build_frontend.SearchRequest) (*shared_proto_build_frontend.SearchResponse, error) {
	if req == nil || req.Query == nil {
		return &shared_proto_build_frontend.SearchResponse{},
			status.Errorf(codes.InvalidArgument, "request empty or missing arguments")
	}

	if len(req.Query.Values) != int(s.collectionConfig.Dimensions) {
		return &shared_proto_build_frontend.SearchResponse{},
			status.Errorf(
				codes.InvalidArgument,
				"expected %d dimensions, got %d",
				len(req.Query.Values),
				s.collectionConfig.Dimensions)
	}

	//keys, distances, err := s.collection.Search(req.Query.Values, req.Limit)
	//if err != nil {
	//	return nil, err
	//}
	//
	//return &shared_proto_build_frontend.SearchResponse{
	//	Keys:      *(*[]uint64)(unsafe.Pointer(&keys)),
	//	Distances: distances,
	//}, nil
}

func (s *frontendGrpcServerImplementation) Add(
	_ context.Context,
	req *shared_proto_build_frontend.AddRequest) (*shared_proto_build_frontend.Empty, error) {
	if req == nil || req.Vector == nil {
		return &shared_proto_build_frontend.Empty{},
			status.Errorf(codes.InvalidArgument, "request empty or missing arguments")
	}

	if len(req.Vector.Values) != int(s.collectionConfig.Dimensions) {
		return &shared_proto_build_frontend.Empty{},
			status.Errorf(
				codes.InvalidArgument,
				"expected %d dimensions, got %d",
				len(req.Vector.Values),
				s.collectionConfig.Dimensions)
	}

	//_, isFull, err := s.collection.Add(shared_collection.Key(req.Key), req.Vector.Values)
	//return &shared_proto_build_frontend.AddResponse{
	//	ShardFull: isFull,
	//}, err
}

func (s *frontendGrpcServerImplementation) AddMulti(
	_ context.Context,
	req *shared_proto_build_frontend.AddMultiRequest) (*shared_proto_build_frontend.AddMultiResponse, error) {
	if req == nil || req.Vectors == nil || req.Keys == nil {
		return &shared_proto_build_frontend.AddMultiResponse{},
			status.Errorf(codes.InvalidArgument, "request empty or missing arguments")
	}

	if len(req.Keys) != len(req.Vectors) {
		return &shared_proto_build_frontend.AddMultiResponse{},
			status.Errorf(codes.InvalidArgument, "keys and vectors must have the same length")
	}

	if len(req.Keys) == 0 {
		return &shared_proto_build_frontend.AddMultiResponse{},
			status.Errorf(codes.InvalidArgument, "no data provided")
	}

	//vectors := make([][]float32, len(req.Vectors))
	//for i, v := range req.Vectors {
	//	if len(req.Vectors[0].Values) != int(s.collectionConfig.Dimensions) {
	//		return &shared_proto_build_frontend.AddMultiResponse{},
	//			status.Errorf(
	//				codes.InvalidArgument,
	//				"vector %d, expected %d dimensions, got %d",
	//				i,
	//				len(req.Vectors[0].Values),
	//				s.collectionConfig.Dimensions)
	//	}
	//
	//	vectors[i] = v.Values
	//}
	//
	//inserted, isFull, err := s.collection.AddMulti(
	//	*(*[]shared_collection.Key)(unsafe.Pointer(&req.Keys)),
	//	*(*[]shared_collection.Vector)(unsafe.Pointer(&vectors)))
	//
	//if err != nil {
	//	var err2 error
	//	serr := status.Newf(codes.Internal, "failed to add vectors: %v", err)
	//	serr, err2 = serr.WithDetails(&shared_proto_build_frontend.AddMultiResponse{
	//		Inserted:  inserted,
	//		ShardFull: isFull,
	//	})
	//	if err2 != nil {
	//		return &shared_proto_build_frontend.AddMultiResponse{},
	//			status.Errorf(
	//				codes.Internal,
	//				"unable to build response with details when failed to add vectors: %v", err)
	//	}
	//
	//	return &shared_proto_build_frontend.AddMultiResponse{}, serr.Err()
	//}
	//
	//return &shared_proto_build_frontend.AddMultiResponse{
	//	Inserted:  inserted,
	//	ShardFull: isFull,
	//}, err
}

func (s *frontendGrpcServerImplementation) Get(
	_ context.Context,
	req *shared_proto_build_frontend.GetRequest) (*shared_proto_build_frontend.GetResponse, error) {
	if req == nil {
		return &shared_proto_build_frontend.GetResponse{},
			status.Errorf(codes.InvalidArgument, "request empty or missing arguments")
	}

	if req.Count <= 0 {
		return &shared_proto_build_frontend.GetResponse{},
			status.Errorf(codes.InvalidArgument, "count must be greater than 0")
	}

	//vec, err := s.collection.Get(shared_collection.Key(req.Key), uint(req.Count))
	//if err != nil {
	//	return nil, err
	//}
	//
	//return &shared_proto_build_frontend.GetResponse{Vector: vectorToPB(vec)}, nil
}

func (s *frontendGrpcServerImplementation) Has(
	_ context.Context,
	req *shared_proto_build_frontend.HasRequest) (*shared_proto_build_frontend.HasResponse, error) {
	if req == nil {
		return &shared_proto_build_frontend.HasResponse{},
			status.Errorf(codes.InvalidArgument, "request empty or missing arguments")
	}

	//return &shared_proto_build_frontend.HasResponse{
	//	Ok: s.collection.Has(shared_collection.Key(req.Key)),
	//}, nil
}

func (s *frontendGrpcServerImplementation) Delete(
	_ context.Context,
	req *shared_proto_build_frontend.DeleteRequest) (*shared_proto_build_frontend.DeleteResponse, error) {
	if req == nil {
		return &shared_proto_build_frontend.DeleteResponse{},
			status.Errorf(codes.InvalidArgument, "request empty or missing arguments")
	}

	//err := s.collection.Delete(shared_collection.Key(req.Key))
	//
	//return &shared_proto_build_frontend.DeleteResponse{
	//	Ok: err == nil,
	//}, err
}

func (s *frontendGrpcServerImplementation) Save(
	_ context.Context,
	_ *shared_proto_build_frontend.Empty) (*shared_proto_build_frontend.Empty, error) {
	//err := s.collection.Save(s.collectionPath)
	//return &shared_proto_build_frontend.Empty{}, err
}

func (s *frontendGrpcServerImplementation) Length(
	_ context.Context,
	_ *shared_proto_build_frontend.Empty) (*shared_proto_build_frontend.LengthResponse, error) {
	//length, err := s.collection.Length()
	//if err != nil {
	//	return nil, err
	//}
	//return &shared_proto_build_frontend.LengthResponse{Length: uint64(length)}, nil
}

func (s *frontendGrpcServerImplementation) Size(
	_ context.Context,
	_ *shared_proto_build_frontend.Empty) (*shared_proto_build_frontend.SizeResponse, error) {
	//size, err := s.collection.Size()
	//if err != nil {
	//	return nil, err
	//}
	//return &shared_proto_build_frontend.SizeResponse{Size: uint64(size)}, nil
}
