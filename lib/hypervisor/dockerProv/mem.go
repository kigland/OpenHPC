package dockerProv

func (ops StartContainerOptions) WithMaxMemory(MB int64) StartContainerOptions {
	ops.Resources.Memory = MB * 1024 * 1024
	return ops
}

func (ops StartContainerOptions) WithMaxMemoryByte(B int64) StartContainerOptions {
	ops.Resources.Memory = B
	return ops
}
