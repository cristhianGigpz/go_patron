package grpc_server

import (
	"context"
	"log"
	"product-service/internal/dto"
	"product-service/internal/usecase"
	productpb "product-service/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server implementa los métodos definidos en ProductService.
// Depende del caso de uso y no directamente de la base de datos.
type Server struct {
	productpb.UnimplementedProductServiceServer
	usecase usecase.ProductUseCaseInterface
}

func NewProductGRPCServer(productUseCase usecase.ProductUseCaseInterface) *Server {
	return &Server{usecase: productUseCase}
}

// GetProduct obtiene un producto por su ID.
func (s *Server) GetProduct(ctx context.Context, req *productpb.ProductRequest) (*productpb.ProductResponse, error) {
	if req == nil || req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "el ID debe ser un número positivo")
	}

	product, err := s.usecase.FindByID(uint(req.GetId()))
	if err != nil {
		return nil, status.Error(codes.NotFound, "producto no encontrado")
	}

	return toProtoProduct(product), nil
}

// GetProducts envía todos los productos usando un stream.
func (s *Server) GetProducts(_ *productpb.Empty, stream productpb.ProductService_GetProductsServer) error {
	products, err := s.usecase.FindAll()
	if err != nil {
		return status.Error(codes.Internal, "no se pudieron obtener los productos")
	}

	for _, product := range products {
		if err := stream.Context().Err(); err != nil {
			return status.Error(codes.Canceled, "conexión cancelada por el cliente")
		}

		if err := stream.Send(toProtoProduct(product)); err != nil {
			log.Printf("Error al enviar producto ID %d: %v", product.ID, err)
			return status.Error(codes.Internal, "no se pudo transmitir el producto")
		}
	}

	return nil
}

func toProtoProduct(product dto.ProductResponse) *productpb.ProductResponse {
	return &productpb.ProductResponse{
		Id:          int32(product.ID),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       int32(product.Stock),
	}
}
