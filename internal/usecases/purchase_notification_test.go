package usecases

import (
	"fmt"
	"testing"
)

func Test_CountryCode(t *testing.T) {
	data := []string{"BY", "RU", "GB", "SE", "UA"}
	for _, c := range data {
		fmt.Println(countryName(c))
	}
}
