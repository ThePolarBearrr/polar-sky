package node

type Node struct {
	IP              string
	Name            string
	MemoryCapacity  int
	MemoryAllocated int
	DiskCapacity    int
	DiskAllocated   int
	TaskCount       int
}
