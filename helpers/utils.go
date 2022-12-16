package helpers

import "strconv"

func EtherToWei(amount float64) string {
	return strconv.Itoa(int(amount * 1000000000000000000))
}
