package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Active   bool   `json:"active"`
}

type Product struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Inventory int32   `json:"inventory"`
}

type Order struct {
	ID        int64   `json:"id"`
	UserID    int64   `json:"user_id"`
	ProductID int64   `json:"product_id"`
	Quantity  int32   `json:"quantity"`
	Total     float64 `json:"total"`
}

type UserService interface {
	GetUser(ctx context.Context, userID int64) (*User, error)
	ValidateUser(ctx context.Context, userID int64) (bool, error)
}

type ProductService interface {
	GetProduct(ctx context.Context, productID int64) (*Product, error)
	CheckInventory(ctx context.Context, productID int64, quantity int32) (bool, error)
}

type UserServiceServer struct {
	users map[int64]*User
}

type ProductServiceServer struct {
	products map[int64]*Product
}

func NewUserServiceServer() *UserServiceServer {
	users := map[int64]*User{
		1: {ID: 1, Username: "alice", Email: "alice@example.com", Active: true},
		2: {ID: 2, Username: "bob", Email: "bob@example.com", Active: true},
		3: {ID: 3, Username: "charlie", Email: "charlie@example.com", Active: false},
	}

	return &UserServiceServer{users: users}
}

func (s *UserServiceServer) GetUser(ctx context.Context, userID int64) (*User, error) {
	user, exists := s.users[userID]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	return user, nil
}

func (s *UserServiceServer) ValidateUser(ctx context.Context, userID int64) (bool, error) {
	user, exists := s.users[userID]
	if !exists {
		return false, status.Errorf(codes.NotFound, "user not found")
	}
	return user.Active, nil
}

func NewProductServiceServer() *ProductServiceServer {
	products := map[int64]*Product{
		1: {ID: 1, Name: "Laptop", Price: 999.99, Inventory: 10},
		2: {ID: 2, Name: "Phone", Price: 499.99, Inventory: 20},
		3: {ID: 3, Name: "Headphones", Price: 99.99, Inventory: 0},
	}
	return &ProductServiceServer{products: products}
}

func (p *ProductServiceServer) GetProduct(ctx context.Context, productID int64) (*Product, error) {
	product, exists := p.products[productID]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "product not found")
	}

	return product, nil
}

func (p *ProductServiceServer) CheckInventory(ctx context.Context, productID int64, quantity int32) (bool, error) {
	product, exists := p.products[productID]
	if !exists {
		return false, status.Errorf(codes.NotFound, "product not found")
	}

	if product.Inventory < quantity {
		return false, status.Errorf(codes.ResourceExhausted, "quantity is not enough")
	}

	return true, nil
}

type GetUserRequest struct {
	UserId int64 `json:"user_id"`
}

type GetUserResponse struct {
	User *User `json:"user"`
}

type ValidateUserRequest struct {
	UserID int64 `json:"user_id"`
}

type ValidateUserResponse struct {
	Valid bool `json:"valid"`
}

func (s *UserServiceServer) GetUserRPC(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error) {
	user, err := s.GetUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &GetUserResponse{User: user}, nil
}

func (s *UserServiceServer) ValidateUserRPC(ctx context.Context, req *ValidateUserRequest) (*ValidateUserResponse, error) {
	valid, err := s.ValidateUser(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	return &ValidateUserResponse{Valid: valid}, nil
}

type GetProductRequest struct {
	ProductID int64 `json:"product_id"`
}

type GetProductResponse struct {
	Product *Product `json:"product"`
}

type CheckInventoryRequest struct {
	ProductID int64 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}

type CheckInventoryResponse struct {
	Available bool `json:"available"`
}

func (p *ProductServiceServer) GetProductRPC(ctx context.Context, req *GetProductRequest) (*GetProductResponse, error) {
	product, err := p.GetProduct(ctx, req.ProductID)
	if err != nil {
		return nil, err
	}

	return &GetProductResponse{Product: product}, nil
}

func (p *ProductServiceServer) CheckInventoryRPC(ctx context.Context, req *CheckInventoryRequest) (*CheckInventoryResponse, error) {
	available, err := p.CheckInventory(ctx, req.ProductID, req.Quantity)
	if err != nil {
		return nil, err
	}

	return &CheckInventoryResponse{Available: available}, nil
}

type OrderService struct {
	userClient    UserService
	productClient ProductService
	orders        map[int64]*Order
	nextOrderID   int64
}

func NewOrderService(userClient UserService, productClient ProductService) *OrderService {
	return &OrderService{
		userClient:    userClient,
		productClient: productClient,
		orders:        make(map[int64]*Order),
		nextOrderID:   1,
	}
}

func (o *OrderService) CreateOrder(ctx context.Context, userId, productID int64, quantity int32) (*Order, error) {
	exists, err := o.userClient.ValidateUser(ctx, productID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, status.Errorf(codes.PermissionDenied, "user not active")
	}

	product, err := o.productClient.GetProduct(ctx, productID)
	if err != nil {
		return nil, err
	}

	ok, err := o.productClient.CheckInventory(ctx, productID, quantity)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "input invalid")
	}

	orderID := int64(o.nextOrderID + 1)

	order := &Order{
		ID:        orderID,
		UserID:    userId,
		ProductID: userId,
		Quantity:  quantity,
		Total:     float64(quantity) * product.Price,
	}

	o.orders[orderID] = order

	return order, nil
}

func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("Request received: %s", info.FullMethod)
	start := time.Now()
	resp, err := handler(ctx, req)
	log.Printf("Request completed: %s in %v", info.FullMethod, time.Since(start))
	return resp, err
}

func AuthInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// Add auth token to metadata
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer token123")
	return invoker(ctx, method, req, reply, cc, opts...)
}

func StartUserService(port string) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return nil, fmt.Errorf("failed to listen %v", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(LoggingInterceptor))

	userServer := NewUserServiceServer()

	mux := http.NewServeMux()

	mux.HandleFunc("/user/get", func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.URL.Query().Get("id")
		userID, _ := strconv.ParseInt(userIDStr, 10, 64)

		user, err := userServer.GetUser(r.Context(), userID)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	})

	mux.HandleFunc("/user/validate", func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.URL.Query().Get("id")
		userID, _ := strconv.ParseInt(userIDStr, 10, 64)

		valid, err := userServer.ValidateUser(r.Context(), userID)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"valid": valid})
	})

	go func() {
		log.Printf("user service HTTP server listening on %s", port)
		if err := http.Serve(lis, mux); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	return s, nil
}

func StartProductService(port string) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Errorf("failed to listen %v", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(LoggingInterceptor))

	productServer := NewProductServiceServer()
	mux := http.NewServeMux()
	mux.HandleFunc("/product/get", func(w http.ResponseWriter, r *http.Request) {
		productIDStr := r.URL.Query().Get("id")
		productID, _ := strconv.ParseInt(productIDStr, 10, 64)

		product, err := productServer.GetProduct(r.Context(), productID)

		if err != nil {
			if status.Code(err) == codes.NotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "Application/json")
		json.NewEncoder(w).Encode(product)
	})

	mux.HandleFunc("/product/check", func(w http.ResponseWriter, r *http.Request) {
		productIDStr := r.URL.Query().Get("id")
		productQtyStr := r.URL.Query().Get("qty")
		productID, _ := strconv.ParseInt(productIDStr, 10, 64)
		quantity, _ := strconv.ParseInt(productQtyStr, 10, 64)
		valid, err := productServer.CheckInventory(r.Context(), productID, int32(quantity))
		if err != nil {
			if status.Code(err) == codes.NotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"valid": valid})
	})

	go func() {
		log.Printf("user service HTTP server listening on %s", port)
		if err := http.Serve(lis, mux); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	return s, nil
}

func ConnectToServices(userServiceAddr, productServiceAddr string) (*OrderService, error) {
	userCon, err := grpc.Dial(userServiceAddr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(AuthInterceptor))
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "failed to connect to user service: %v", err)
	}
	productCon, err := grpc.Dial(productServiceAddr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(AuthInterceptor))
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "failed to connect to product service: %v", err)
	}
	userClient := NewUserServiceClient(userCon)
	productClient := NewProductServiceClient(productCon)
	return NewOrderService(userClient, productClient), nil
}

type UserServiceClient struct {
	baseURL string
}

// GetUser implements UserService.
func (u *UserServiceClient) GetUser(ctx context.Context, userID int64) (*User, error) {
	resp, err := http.Get(fmt.Sprintf("%s/user/get?id=%d", u.baseURL, userID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

// ValidateUser implements UserService.
func (u *UserServiceClient) ValidateUser(ctx context.Context, userID int64) (bool, error) {
	resp, err := http.Get(fmt.Sprintf("%s/user/validate?id=%d", u.baseURL, userID))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return false, status.Errorf(codes.NotFound, "user not found")
	}
	var result map[string]bool
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}
	return result["valid"], nil
}

type ProductServiceClient struct {
	conn *grpc.ClientConn
}

// CheckInventory implements ProductService.
func (p *ProductServiceClient) CheckInventory(ctx context.Context, productID int64, quantity int32) (bool, error) {
	resp, err := http.Get(fmt.Sprintf("%s/product/check?id=%d&quantity=%d", p.conn.Target(), productID, quantity))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return false, status.Errorf(codes.NotFound, "product not found")
	}

	var result map[string]bool

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}
	return result["valid"], nil
}

// GetProduct implements ProductService.
func (p *ProductServiceClient) GetProduct(ctx context.Context, productID int64) (*Product, error) {
	resp, err := http.Get(fmt.Sprintf("%s/product/get?id=%d", p.conn.Target(), productID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, status.Errorf(codes.NotFound, "product not found")
	}
	var product Product
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, err
	}
	return &product, nil
}

func NewProductServiceClient(conn *grpc.ClientConn) ProductService {
	return &ProductServiceClient{conn: conn}
}

func NewUserServiceClient(conn *grpc.ClientConn) UserService {
	return &UserServiceClient{baseURL: "http://localhost:50051"}
}

func RegisterUserServiceServer(s *grpc.Server, srv *UserServiceServer) {
	// In a real implementation, this would be generated code
	// For this challenge, we'll manually handle the registration
}

func RegisterProductServiceServer(s *grpc.Server, srv *ProductServiceServer) {
	// In a real implementation, this would be generated code
	// For this challenge, we'll manually handle the registration
}

func main() {

}
