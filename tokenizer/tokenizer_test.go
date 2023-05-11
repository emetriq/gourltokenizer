package tokenizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	DefaultStopWordFunc = IsGermanStopWord
}
func Test_tokenizeCorrectPath(t *testing.T) {
	path := "/some-thing/very/interesting?queryparam2=1&queryparam2=3"
	result := tokenize(path)
	assert.ElementsMatch(t, []string{"some", "thing", "very", "interesting"}, result)
}

func Test_tokenizePathWithDashes(t *testing.T) {
	path := "/some-thing/very/interesting"
	result := tokenize(path)
	assert.ElementsMatch(t, []string{"some", "thing", "very", "interesting"}, result)
}

func Test_tokenizePathWithDashes2(t *testing.T) {
	path := "/hsv-fussball"
	result := tokenize(path)
	assert.ElementsMatch(t, []string{"hsv", "fussball"}, result)
}

func Test_tokenizeEmptyString(t *testing.T) {
	path := ""
	result := tokenize(path)
	assert.ElementsMatch(t, []string{}, result)
}

func Test_filterStopWorlds(t *testing.T) {
	result := filterStopWords([]string{"hallo", "cms", "titel", "welt"}, func(val string) bool {
		if val == "cms" {
			return true
		}
		if val == "titel" {
			return true
		}
		return false
	})
	assert.ElementsMatch(t, []string{"hallo", "welt"}, result)
}

func Test_URLTokenizer(t *testing.T) {
	result := Tokenize("http://example.com/path/sport/hsv-fussball?bla=1")
	assert.ElementsMatch(t, []string{"path", "sport", "hsv", "fussball", "example.com"}, result)
}

func Test_URLTokenizerOneWord(t *testing.T) {
	result := Tokenize("http://example.com/sport")
	assert.ElementsMatch(t, []string{"example.com", "sport"}, result)
}

func Test_URLTokenizerOneWordMinSize(t *testing.T) {

	result := Tokenize("http://www.test-page.de/aaa/bbb/bc/ccc")
	assert.ElementsMatch(t, []string{"www.test-page.de", "aaa", "bbb", "ccc"}, result)
}

func Test_URLTokenizerWithScapedChars(t *testing.T) {
	result := Tokenize("http://example.com/%3ahttps%3A%2F%2Fwww.emetriq.com%2F", IsGermanStopWord)
	assert.ElementsMatch(t, []string{"emetriq", "com", "example.com"}, result)
}

func Test_URLTokenizerWithWrongEscapedChars(t *testing.T) {
	result := Tokenize("http://example.com/%%ssomething/usefull")
	assert.Equal(t, []string{"ssomething", "usefull", "example.com"}, result)
}
func Test_URLTokenizerWithWrongEscapedChars2(t *testing.T) {
	DefaultStopWordFunc = IsGermanStopWord
	result := Tokenize("https://www.morgenpost.de/vermischtes/article233484549/marisa-burger-rosenheim-cops-schaupielerin.html?service=amp#aoh=16333619698076&csi=0&referrer=https://www.google.com&amp_tf=Von %1$s")
	assert.Equal(t, []string{
		"vermischtes",
		"marisa",
		"burger",
		"rosenheim",
		"cops",
		"schaupielerin",
		"www.morgenpost.de",
	}, result)
}
func Test_URLTokenizerWithWrongHostEscapedChars(t *testing.T) {
	result := Tokenize("http://..example.com/something")
	assert.Equal(t, []string{"something", "example.com"}, result)
}

func Test_URLTokenizerWithCapitalChars(t *testing.T) {
	DefaultStopWordFunc = IsGermanStopWord
	result := Tokenize("mailto://www.Subdomain.example.com/HSV-fussbal%3asome/a")
	assert.ElementsMatch(t, []string{"subdomain", "hsv", "fussbal", "some", "www.subdomain.example.com"}, result)
}

func Test_URLWithoutHTTP(t *testing.T) {
	result := Tokenize("www.Subdomain.example.com")
	assert.ElementsMatch(t, []string{"subdomain", "www.subdomain.example.com"}, result)
}

func Test_URLWithoutHTTPAndWithoutSubdomain(t *testing.T) {
	result := Tokenize("www.example.com")
	assert.ElementsMatch(t, []string{"www.example.com"}, result)
}

func Test_URLWithoutHTTPAndSubdomain(t *testing.T) {
	result := Tokenize("sport.fussball.example.com")
	assert.ElementsMatch(t, []string{"sport", "fussball", "sport.fussball.example.com"}, result)
}

func Test_URLWithoutHTTPButWithPath(t *testing.T) {
	result := Tokenize("www.ironsrc.com/sports")
	assert.ElementsMatch(t, []string{"sports", "www.ironsrc.com"}, result)
}

func Test_SkipWordsWithNumbers(t *testing.T) {
	result := Tokenize("https://www.autoscout24.at/angebote/seat-altea-xl-reference-1-4-tfsi-motorschaden-benzin-grau-b82ebced-cb95-4f49-8038-5eb1c098e652")
	// no 'ebced'
	assert.ElementsMatch(t, []string{"angebote", "seat", "altea", "reference", "tfsi", "motorschaden", "benzin", "grau", "www.autoscout24.at"}, result)
	assert.NotContains(t, result, "ebced")

	result = Tokenize("https://www.coches.net/123nissan-interstar-25dci-120-pro-l2h2-3500-diesel-2009-en-barcelona-52386149-fuvivo.aspx1")
	// no '123nissan', 'dci' and 'aspx1'
	assert.ElementsMatch(t, []string{"interstar", "pro", "diesel", "barcelona", "fuvivo", "www.coches.net"}, result)
	assert.NotContains(t, result, "dci")
}

func BenchmarkEscapedURLTokenizer(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Tokenize("http://example.com/path/sport/hsv-fussball?bla=1&escaped=%2C%2C%3A%3A%3B%3B")
	}
}

func BenchmarkURLTokenizer(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Tokenize("http://example.com/path/sport/hsv-fussball?bla=1")
	}
}

func BenchmarkURLTokenizerFast(b *testing.B) {
	for n := 0; n < b.N; n++ {
		TokenizeFast("http://example.com/path/sport/hsv-fussball?bla=1")
	}
}

func BenchmarkTokenizerV3(b *testing.B) {
	for n := 0; n < b.N; n++ {
		tokenize("http://example.com/path/sport/hsv-fussball?bla=1")
	}
}
