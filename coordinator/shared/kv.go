package shared

import (
	"github.com/kigland/OpenHPC/coordinator/models/dbmod"
	"github.com/kigland/OpenHPC/lib/store"
)

var GCKVStore store.IKVStore

func initGCKVStore() {
	GCKVStore = dbmod.NewGCStore(DB)
}
