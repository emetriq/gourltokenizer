package tokenizer

import (
	"net/url"
	"strings"
)

var MinWordSize = 3

var DefaultStopWordFunc = IsEnglishStopWord

func isByteAllowed(b byte, isDotCountMode bool) bool {
	if b >= 'a' && b <= 'z' {
		return true
	}
	return isDotCountMode && b != '.' && b != '/'
}

// faster solution than using strings.Contains(), because we are only looking
// for a single char and can leave the loop after
func stringContainsByteChar(s string, r byte) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == r {
			return true
		}
	}
	return false
}

// Tokenize splits URL to host and path parts and tokenize path and host part
// all terms are returned in lower case. If numbers are within a word, the complete
// word is filtered out.
func Tokenize(encodedURL string, stopwordfunc ...func(string) bool) []string {
	encodedURLLower := strings.ToLower(encodedURL)
	var result []string

	// check if url needs unescaping
	if stringContainsByteChar(encodedURLLower, '%') {
		decodedURL, err := url.QueryUnescape(encodedURLLower)
		if err != nil {
			escapedEncodedURL := url.QueryEscape(encodedURL)
			decodedURL, err = url.QueryUnescape(escapedEncodedURL)
		}

		if err != nil {
			return []string{}
		}

		result = filterStopWords(tokenize(decodedURL), stopwordfunc...)
	} else {
		result = filterStopWords(tokenize(encodedURLLower), stopwordfunc...)
	}

	return result
}

// TokenizeFastV3 splits URL to host and path parts and tokenize path and host part
// all terms are returned in lower case
func TokenizeFast(encodedURL string, stopwordfunc ...func(string) bool) []string {
	urlLower := strings.ToLower(encodedURL)
	result := tokenize(urlLower)
	if len(stopwordfunc) > 0 {
		result = filterStopWords(result, stopwordfunc[0])
	}
	return result
}

func tokenize(str string) []string {
	// remove protocol
	startIndex := strings.Index(str, "://")
	if startIndex < 7 && startIndex > 0 && len(str) > startIndex+3 {
		startIndex = startIndex + 3
	} else {
		startIndex = 0
	}

	strLen := len(str)
	lastIndex := strLen - 1
	result := make([]string, 0, strLen/MinWordSize)
	start := -1
	dotCounter := 0
	isDotCountMode := true
	isContainingNumber := false
	domainNameEndIndex := -1
	domainNameStartIndex := startIndex
	var b byte
	for idx := 0; idx < len(str); idx++ {
		b = str[idx]
		if idx < startIndex {
			continue
		}

		if isByteAllowed(b, isDotCountMode) {
			if start == -1 {
				start = idx
			}
			if idx == lastIndex && ((lastIndex-start+1) >= MinWordSize || isDotCountMode) {
				if !isContainingNumber {
					result = append(result, str[start:strLen])
				}
				isContainingNumber = false
			}
		} else if b >= '0' && b <= '9' && !isDotCountMode {
			isContainingNumber = true
		} else if ((idx-start) >= MinWordSize || isDotCountMode) && start > -1 {
			if !isContainingNumber {
				result = append(result, str[start:idx])
			}

			isContainingNumber = false
			start = -1
		} else {
			isContainingNumber = false
			start = -1
		}
		if b == '/' && isDotCountMode {
			isDotCountMode = false
			domainNameEndIndex = idx
			dotCounter = len(result) - 1
		}

		if b == '?' { // skip query params
			break
		}
	}

	if isDotCountMode {
		dotCounter = len(result) - 1
		domainNameEndIndex = len(str)
	}

	if dotCounter > 0 && len(result) > 1 {
		result = append(result[:(dotCounter-1)], result[dotCounter+1:]...)
		if domainNameEndIndex-domainNameStartIndex > 3 { // if domain name is longer than 3 chars
			for len(str) > domainNameStartIndex && str[domainNameStartIndex] == '.' {
				domainNameStartIndex++
			}
			result = append(result, str[domainNameStartIndex:domainNameEndIndex])
		}
	}
	return result
}

func filterStopWords(terms []string, stopwordfunc ...func(string) bool) []string {
	filter := DefaultStopWordFunc
	if len(stopwordfunc) > 0 {
		filter = stopwordfunc[0]
	} else if filter == nil {
		return terms
	}

	for i := 0; len(terms) > i; i++ {
		if filter(terms[i]) || filter(terms[i][1:]) {
			terms = append(terms[:i], terms[i+1:]...)
			i--
		}
	}
	return terms
}
