package shared

import (
	"github.com/kigland/HPC-Scheduler/coodinator/models/dbmod"
	"github.com/kigland/HPC-Scheduler/lib/store"
)

var GCKVStore store.IKVStore

func initGCKVStore() {
	GCKVStore = dbmod.NewGCStore(DB)
}
