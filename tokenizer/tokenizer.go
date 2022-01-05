package tokenizer

import (
	"fmt"
	"net/url"
	"strings"
)

var MinWordSize = 3

var DefaultStopWordFunc = IsEnglishStopWord

func isRuneAllowed(r rune, isDotCountMode bool) bool {
	if r >= 'a' && r <= 'z' {
		return true
	}
	return isDotCountMode && r != '.' && r != '/'
}

//TokenizeV2 splits URL to host and path parts and tokenize path and host part
//all terms are returned in lower case
func TokenizeV2(encodedURL string, stopwordfunc ...func(string) bool) []string {
	encodedURLLower := strings.ToLower(encodedURL)
	decodedURL, err := url.QueryUnescape(encodedURLLower)
	if err != nil {
		escapedEncodedURL := url.QueryEscape(encodedURL)
		decodedURL, err = url.QueryUnescape(escapedEncodedURL)
	}

	if err != nil {
		return []string{}
	}

	result := filterStopWords(tokenizeV2(decodedURL), stopwordfunc...)

	return result
}

//TokenizeFastV2 splits URL to host and path parts and tokenize path and host part without url decoding
//all terms are returned in lower case
func TokenizeFastV2(encodedURL string, stopwordfunc ...func(string) bool) []string {
	urlLower := strings.ToLower(encodedURL)
	result := tokenizeV2(urlLower)
	if len(stopwordfunc) > 0 {
		result = filterStopWords(result, stopwordfunc[0])
	}
	return result
}

//TokenizeURL splits URL to host and path parts and tokenize path part
//all terms are returned in lower case
func TokenizeV1(url string, stopwordfunc ...func(string) bool) []string {
	urlToParse := url
	if !strings.HasPrefix(url, "http") && !strings.HasPrefix(url, "mailto") {
		urlToParse = fmt.Sprintf("http://%s", urlToParse)
	}
	urlLower := strings.ToLower(urlToParse)
	host, path, err := parseURL(urlLower)
	if err != nil {
		return []string{}
	}

	result := filterStopWords(tokenizeV1(path), stopwordfunc...)

	return append(result, host)
}

func tokenizeV2(str string) []string {
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
	domainNameEndIndex := -1
	domainNameStartIndex := startIndex
	for idx, r := range str {
		if idx < startIndex {
			continue
		}

		if isRuneAllowed(r, isDotCountMode) {
			if start == -1 {
				start = idx
			}
			if idx == lastIndex && ((lastIndex-start+1) >= MinWordSize || isDotCountMode) {
				result = append(result, str[start:strLen])
			}
		} else if ((idx-start) >= MinWordSize || isDotCountMode) && start > -1 {
			result = append(result, str[start:idx])
			start = -1
		} else {
			start = -1
		}
		if r == '/' && isDotCountMode {
			isDotCountMode = false
			domainNameEndIndex = idx
			dotCounter = len(result) - 1
		}

		if r == '?' { // skip query params
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

func tokenizeV1(str string) []string {
	strLen := len(str)
	lastIndex := strLen - 1
	result := make([]string, 0, strLen/MinWordSize)
	start := -1

	for idx, r := range str {
		if isRuneAllowed(r, false) {
			if start == -1 {
				start = idx
			}
			if idx == lastIndex && (lastIndex-start+1) >= MinWordSize {
				result = append(result, str[start:strLen])
			}
		} else if (idx-start) >= MinWordSize && start > -1 {
			result = append(result, str[start:idx])
			start = -1
		} else {
			start = -1
		}
	}
	return result
}

func parseURL(str string) (host string, path string, err error) {
	u, err := url.Parse(str)
	if err != nil {
		return "", "", err
	}
	return u.Host, u.Path, nil
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
