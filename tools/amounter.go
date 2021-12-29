package tools

func GetAmount(len, limit int) int {
	amount := len / limit
	if len%limit != 0 {
		amount++
	}
	return amount
}
