package server

import (
	"github.com/danielealbano/svdb/engine-worker/collection"
	shared_proto_build_collection "github.com/danielealbano/svdb/shared/proto/build/collection"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"unsafe"
)

type collectionGrpcServerImplementation struct {
	shared_proto_build_collection.UnimplementedCollectionServer
	collection     *collection.Collection
	collectionPath string
}

func vectorToPB(v []float32) *shared_proto_build_collection.Vector {
	return &shared_proto_build_collection.Vector{Values: v}
}

func RegisterCollectionGrpcServerImplementation(
	server *GrpcServer,
	path string) {
	shared_proto_build_collection.RegisterCollectionServer(server.grpcServer, &collectionGrpcServerImplementation{
		collection:     server.collection,
		collectionPath: path,
	})
}

func (s *collectionGrpcServerImplementation) Search(
	_ context.Context,
	req *shared_proto_build_collection.SearchRequest) (*shared_proto_build_collection.SearchResponse, error) {
	if req == nil || req.Query == nil {
		return &shared_proto_build_collection.SearchResponse{},
			status.Errorf(codes.InvalidArgument, "request empty or missing arguments")
	}

	if len(req.Query.Values) != int(s.collection.Config.Dimensions) {
		return &shared_proto_build_collection.SearchResponse{},
			status.Errorf(
				codes.InvalidArgument,
				"expected %d dimensions, got %d",
				len(req.Query.Values),
				s.collection.Config.Dimensions)
	}

	keys, distances, err := s.collection.Search(req.Query.Values, req.Limit)
	if err != nil {
		return nil, err
	}

	return &shared_proto_build_collection.SearchResponse{
		Keys:      *(*[]uint64)(unsafe.Pointer(&keys)),
		Distances: distances,
	}, nil
}

func (s *collectionGrpcServerImplementation) Add(
	_ context.Context,
	req *shared_proto_build_collection.AddRequest) (*shared_proto_build_collection.AddResponse, error) {
	if req == nil || req.Vector == nil {
		return &shared_proto_build_collection.AddResponse{},
			status.Errorf(codes.InvalidArgument, "request empty or missing arguments")
	}

	if len(req.Vector.Values) != int(s.collection.Config.Dimensions) {
		return &shared_proto_build_collection.AddResponse{},
			status.Errorf(
				codes.InvalidArgument,
				"expected %d dimensions, got %d",
				len(req.Vector.Values),
				s.collection.Config.Dimensions)
	}

	_, isFull, err := s.collection.Add(collection.Key(req.Key), req.Vector.Values)
	return &shared_proto_build_collection.AddResponse{
		ShardFull: isFull,
	}, err
}

func (s *collectionGrpcServerImplementation) AddMulti(
	_ context.Context,
	req *shared_proto_build_collection.AddMultiRequest) (*shared_proto_build_collection.AddMultiResponse, error) {
	if req == nil || req.Vectors == nil || req.Keys == nil {
		return &shared_proto_build_collection.AddMultiResponse{},
			status.Errorf(codes.InvalidArgument, "request empty or missing arguments")
	}

	if len(req.Keys) != len(req.Vectors) {
		return &shared_proto_build_collection.AddMultiResponse{},
			status.Errorf(codes.InvalidArgument, "keys and vectors must have the same length")
	}

	if len(req.Keys) == 0 {
		return &shared_proto_build_collection.AddMultiResponse{},
			status.Errorf(codes.InvalidArgument, "no data provided")
	}

	vectors := make([][]float32, len(req.Vectors))
	for i, v := range req.Vectors {
		if len(req.Vectors[0].Values) != int(s.collection.Config.Dimensions) {
			return &shared_proto_build_collection.AddMultiResponse{},
				status.Errorf(
					codes.InvalidArgument,
					"vector %d, expected %d dimensions, got %d",
					i,
					len(req.Vectors[0].Values),
					s.collection.Config.Dimensions)
		}

		vectors[i] = v.Values
	}

	inserted, isFull, err := s.collection.AddMulti(
		*(*[]collection.Key)(unsafe.Pointer(&req.Keys)),
		*(*[]collection.Vector)(unsafe.Pointer(&vectors)))

	if err != nil {
		var err2 error
		serr := status.Newf(codes.Internal, "failed to add vectors: %v", err)
		serr, err2 = serr.WithDetails(&shared_proto_build_collection.AddMultiResponse{
			Inserted:  inserted,
			ShardFull: isFull,
		})
		if err2 != nil {
			return &shared_proto_build_collection.AddMultiResponse{},
				status.Errorf(
					codes.Internal,
					"unable to build response with details when failed to add vectors: %v", err)
		}

		return &shared_proto_build_collection.AddMultiResponse{}, serr.Err()
	}

	return &shared_proto_build_collection.AddMultiResponse{
		Inserted:  inserted,
		ShardFull: isFull,
	}, err
}

func (s *collectionGrpcServerImplementation) Get(
	_ context.Context,
	req *shared_proto_build_collection.GetRequest) (*shared_proto_build_collection.GetResponse, error) {
	if req == nil {
		return &shared_proto_build_collection.GetResponse{},
			status.Errorf(codes.InvalidArgument, "request empty or missing arguments")
	}

	if req.Count <= 0 {
		return &shared_proto_build_collection.GetResponse{},
			status.Errorf(codes.InvalidArgument, "count must be greater than 0")
	}

	vec, err := s.collection.Get(collection.Key(req.Key), uint(req.Count))
	if err != nil {
		return nil, err
	}

	return &shared_proto_build_collection.GetResponse{Vector: vectorToPB(vec)}, nil
}

func (s *collectionGrpcServerImplementation) Has(
	_ context.Context,
	req *shared_proto_build_collection.HasRequest) (*shared_proto_build_collection.HasResponse, error) {
	if req == nil {
		return &shared_proto_build_collection.HasResponse{},
			status.Errorf(codes.InvalidArgument, "request empty or missing arguments")
	}

	return &shared_proto_build_collection.HasResponse{
		Ok: s.collection.Has(collection.Key(req.Key)),
	}, nil
}

func (s *collectionGrpcServerImplementation) Delete(
	_ context.Context,
	req *shared_proto_build_collection.DeleteRequest) (*shared_proto_build_collection.DeleteResponse, error) {
	if req == nil {
		return &shared_proto_build_collection.DeleteResponse{},
			status.Errorf(codes.InvalidArgument, "request empty or missing arguments")
	}

	err := s.collection.Delete(collection.Key(req.Key))

	return &shared_proto_build_collection.DeleteResponse{
		Ok: err == nil,
	}, err
}

func (s *collectionGrpcServerImplementation) Save(
	_ context.Context,
	_ *shared_proto_build_collection.Empty) (*shared_proto_build_collection.Empty, error) {
	err := s.collection.Save(s.collectionPath)
	return &shared_proto_build_collection.Empty{}, err
}

func (s *collectionGrpcServerImplementation) Length(
	_ context.Context,
	_ *shared_proto_build_collection.Empty) (*shared_proto_build_collection.LengthResponse, error) {
	length, err := s.collection.Length()
	if err != nil {
		return nil, err
	}
	return &shared_proto_build_collection.LengthResponse{Length: uint64(length)}, nil
}

func (s *collectionGrpcServerImplementation) Capacity(
	_ context.Context,
	_ *shared_proto_build_collection.Empty) (*shared_proto_build_collection.CapacityResponse, error) {
	capacity, err := s.collection.Capacity()
	if err != nil {
		return nil, err
	}
	return &shared_proto_build_collection.CapacityResponse{Capacity: uint64(capacity)}, nil
}

func (s *collectionGrpcServerImplementation) Size(
	_ context.Context,
	_ *shared_proto_build_collection.Empty) (*shared_proto_build_collection.SizeResponse, error) {
	size, err := s.collection.Size()
	if err != nil {
		return nil, err
	}
	return &shared_proto_build_collection.SizeResponse{Size: uint64(size)}, nil
}
