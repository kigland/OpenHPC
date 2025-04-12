package container

func (u *UserSpecificConf) PrepareContainerRequest(vCore, vGPU, vMem int) error {
	u.RequestedCPU = vCore
	u.RequestedGPU = vGPU
	u.RequestedMem = vMem

	u.Normalize()
	err := u.Validate()

	if err != nil {
		return err
	}

	return nil
}
