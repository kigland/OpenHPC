package utils

func ShortId(cid string) string {
	if len(cid) > 12 {
		return cid[:12]
	}
	return cid
}
