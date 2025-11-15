package exchange

// MapperMexcToSpot accepts base and quote tokens and returns a specific mexc spot pair string.
//
// BTCUSDT
func MapperMexcSpot(base, quote string) string {
	return base + quote
}

func MapperMexcFutures(base, quote string) string {
	return base + "_" + quote
}
