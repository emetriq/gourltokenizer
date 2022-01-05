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
	result := tokenizeV2(path)
	assert.ElementsMatch(t, []string{"some", "thing", "very", "interesting"}, result)
}

func Test_tokenizePathWithDashes(t *testing.T) {
	path := "/some-thing/very/interesting"
	result := tokenizeV2(path)
	assert.ElementsMatch(t, []string{"some", "thing", "very", "interesting"}, result)
}

func Test_tokenizePathWithDashes2(t *testing.T) {
	path := "/hsv-fussball"
	result := tokenizeV2(path)
	assert.ElementsMatch(t, []string{"hsv", "fussball"}, result)
}

func Test_tokenizeEmptyString(t *testing.T) {
	path := ""
	result := tokenizeV2(path)
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
	result := TokenizeV2("http://example.com/path/sport/hsv-fussball?bla=1")
	assert.ElementsMatch(t, []string{"path", "sport", "hsv", "fussball", "example.com"}, result)
}

func Test_URLTokenizerOneWord(t *testing.T) {
	result := TokenizeV2("http://example.com/sport")
	assert.ElementsMatch(t, []string{"example.com", "sport"}, result)
}

func Test_URLTokenizerOneWordMinSize(t *testing.T) {

	result := TokenizeV2("http://www.test-page.de/aaa/bbb/bc/ccc")
	assert.ElementsMatch(t, []string{"www.test-page.de", "aaa", "bbb", "ccc"}, result)
}

func Test_URLTokenizerWithScapedChars(t *testing.T) {
	result := TokenizeV2("http://example.com/%3ahttps%3A%2F%2Fwww.emetriq.com%2F", IsGermanStopWord)
	assert.ElementsMatch(t, []string{"emetriq", "com", "example.com"}, result)
}

func Test_URLTokenizerWithWrongEscapedChars(t *testing.T) {
	result := TokenizeV2("http://example.com/%%ssomething/usefull")
	assert.Equal(t, []string{"ssomething", "usefull", "example.com"}, result)
}
func Test_URLTokenizerWithWrongEscapedChars2(t *testing.T) {
	DefaultStopWordFunc = IsGermanStopWord
	result := TokenizeV2("https://www.morgenpost.de/vermischtes/article233484549/marisa-burger-rosenheim-cops-schaupielerin.html?service=amp#aoh=16333619698076&csi=0&referrer=https://www.google.com&amp_tf=Von %1$s")
	assert.Equal(t, []string{
		"vermischtes",
		"article",
		"marisa",
		"burger",
		"rosenheim",
		"cops",
		"schaupielerin",
		"www.morgenpost.de",
	}, result)
}
func Test_URLTokenizerWithWrongHostEscapedChars(t *testing.T) {
	result := TokenizeV2("http://..example.com/something")
	assert.Equal(t, []string{"something", "example.com"}, result)
}

func Test_URLTokenizerWithCapitalChars(t *testing.T) {
	DefaultStopWordFunc = IsGermanStopWord
	result := TokenizeV2("mailto://www.Subdomain.example.com/HSV-fussbal%3asome/a")
	assert.ElementsMatch(t, []string{"subdomain", "hsv", "fussbal", "some", "www.subdomain.example.com"}, result)
}

func Test_URLWithoutHTTP(t *testing.T) {
	result := TokenizeV2("www.Subdomain.example.com")
	assert.ElementsMatch(t, []string{"subdomain", "www.subdomain.example.com"}, result)
}

func Test_URLWithoutHTTPAndWithoutSubdomain(t *testing.T) {
	result := TokenizeV2("www.example.com")
	assert.ElementsMatch(t, []string{"www.example.com"}, result)
}

func Test_URLWithoutHTTPAndSubdomain(t *testing.T) {
	result := TokenizeV2("sport.fussball.example.com")
	assert.ElementsMatch(t, []string{"sport", "fussball", "sport.fussball.example.com"}, result)
}

func Test_URLWithoutHTTPButWithPath(t *testing.T) {
	result := TokenizeV2("www.ironsrc.com/sports")
	assert.ElementsMatch(t, []string{"sports", "www.ironsrc.com"}, result)
}

func BenchmarkURLTokenizerV2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		TokenizeV2("http://example.com/path/sport/hsv-fussball?bla=1")
	}
}

func BenchmarkURLTokenizerV2Fast(b *testing.B) {
	for n := 0; n < b.N; n++ {
		TokenizeFastV2("http://example.com/path/sport/hsv-fussball?bla=1")
	}
}

func BenchmarkURLTokenizerV1(b *testing.B) {
	for n := 0; n < b.N; n++ {
		TokenizeV1("http://example.com/path/sport/hsv-fussball?bla=1")
	}
}

func BenchmarkTokenizerV1(b *testing.B) {
	for n := 0; n < b.N; n++ {
		tokenizeV1("http://example.com/path/sport/hsv-fussball?bla=1")
	}
}
func BenchmarkTokenizerV2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		tokenizeV2("http://example.com/path/sport/hsv-fussball?bla=1")
	}
}
