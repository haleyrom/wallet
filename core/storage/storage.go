package storage

// Storage 存储相关接口
type Storage interface {
	Init(addr, prefix string) error
}
