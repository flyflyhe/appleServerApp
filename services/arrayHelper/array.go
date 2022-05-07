package arrayHelper

func ArrayReverse[V []string](arr V) {
	i := 0
	j := len(arr) - 1

	for i < j {
		arr[i], arr[j] = arr[j], arr[i]
		i++
		j--
	}
}
