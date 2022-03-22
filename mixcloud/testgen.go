package mixcloud

import (
	"math/rand"
	"time"
)

var charSet = "abcdefghijkklmnopqrstuvwxyz"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func GeneratePage(length int, paging bool) Response {
	items := []Mix{}
	response := Response{}
	for i := 0; i < length; i++ {
		name := StringWithCharset(10, charSet)
		m := Mix{
			name + "key",
			name + "url",
			name,
		}
		items = append(items, m)
	}
	response.Data = items
	if paging {
		response.Paging = Paging{Next: "yes"}
	}

	return response

}

// function to generate an array of responses for iterating over in tests
func GeneratePages(count int, length int) []Response {
	res := []Response{}
	for i := 0; i < count-1; i++ {
		res = append(res, GeneratePage(length, true))
	}
	res = append(res, GeneratePage(length, false))

	return res

}
