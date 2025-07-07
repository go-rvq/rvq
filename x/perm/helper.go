package perm

func OkOrDanied(ok bool) error {
	if !ok {
		return PermissionDenied
	}
	return nil
}
