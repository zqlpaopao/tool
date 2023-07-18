package etcd

// NewEtcd creates a Etcd cli
// This can mask the differences between different versions of Etcd
func NewEtcd[C, R any](conf C, f func(C) (R, error)) (R, error) {
	return f(conf)
}
