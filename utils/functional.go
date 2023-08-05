package utils

func First[T any](fst T, _ any) T {
	return fst
}

func Second[T any](_ any, snd T) T {
	return snd
}
