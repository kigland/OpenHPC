package frpc

import "os/exec"

type Frpc struct {
	BinPath    string `json:"bin_path"`
	ConfigPath string `json:"config_path"`
}

// doc: https://gofrp.org/zh-cn/docs/features/common/client/#%E5%8A%A8%E6%80%81%E9%85%8D%E7%BD%AE%E6%9B%B4%E6%96%B0
func (f Frpc) Refresh() error {
	return exec.Command(f.BinPath, "reload", "-c", f.ConfigPath).Run()
}
