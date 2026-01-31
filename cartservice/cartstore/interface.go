package cartstore

import (
	pb "cartservice/proto"
	"context"
	"sync"
)

// 数据保存在内存中结构体，使用嵌套map保存
type memoryCartStore struct {
	// 读写锁
	sync.RWMutex
	carts map[string]map[string]int32
}

type CartStore interface {
	AddItem(ctx context.Context, userID, productID string, quantity int32, out *pb.Empty) (r *pb.Empty, err error)
	EmptyCart(ctx context.Context, userID string) (r *pb.Empty, err error)
	GetCart(ctx context.Context, userID string) (r *pb.Cart, err error)
}

func NewMemoryCartStore() CartStore {
	return &memoryCartStore{
		carts: make(map[string]map[string]int32),
	}
}
