package ports

func SetIfNotNil[T any](updates map[string]any, key string, v *T, transform func(T) (any, error)) error {
	if v == nil {
		return nil
	}
	val, err := transform(*v)
	if err != nil {
		return err
	}
	updates[key] = val
	return nil
}
