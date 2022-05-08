package arrayHelper

type ArrayType interface {
	~string | ~int | ~float32
}

func ArrayReverse[T ArrayType, V []T](arr V) {
	i := 0
	j := len(arr) - 1

	for i < j {
		arr[i], arr[j] = arr[j], arr[i]
		i++
		j--
	}
}
