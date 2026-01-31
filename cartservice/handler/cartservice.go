package handler

import (
	"cartservice/cartstore"
	pb "cartservice/proto"
	"context"
)

type CartService struct {
	Store cartstore.CartStore
	pb.UnimplementedCartServiceServer
}

func (s *CartService) AddItem(ctx context.Context, in *pb.AddItemRequest) (out *pb.Empty, err error) {
	out = new(pb.Empty)
	return s.Store.AddItem(ctx, in.UserId, in.Item.ProductId, in.Item.Quantity, out)
}

func (s *CartService) GetCart(ctx context.Context, in *pb.GetCartRequest) (out *pb.Cart, err error) {
	cart, err := s.Store.GetCart(ctx, in.UserId)
	out = new(pb.Cart)
	if err != nil {
		return out, err
	}
	out.UserId = cart.UserId
	out.Items = cart.Items
	return out, nil
}

func (s *CartService) EmptyCart(ctx context.Context, in *pb.EmptyCartRequest) (out *pb.Empty, err error) {
	return s.Store.EmptyCart(ctx, in.UserId)
}
