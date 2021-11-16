package sharding

type Partition struct {
	Key   string
	Shard string
}

func PARTITION_ETH() *Partition { return &Partition{Key: "eth", Shard: "ethereum"} }

func PARTITION_FAB() *Partition { return &Partition{Key: "fab", Shard: "fabric"} }

func PARTITION_DEFAULT() *Partition { return &Partition{Key: "fab", Shard: "fabric"} }
