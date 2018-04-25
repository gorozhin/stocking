package auth

type Container map[string]string

func (c Container)Valid(l, p string) bool {
	for k, v := range c {
		if k == l && v == p {
			return true
		}
	}
	return false
}
