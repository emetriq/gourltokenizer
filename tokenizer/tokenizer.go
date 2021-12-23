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

//IsStopWord returns true if word is stop word
func IsGermanStopWord(word string) bool {
	switch word {
	case "www":
		return true
	case "jenen":
		return true
	case "manchem":
		return true
	case "euren":
		return true
	case "ihres":
		return true
	case "war":
		return true
	case "meinem":
		return true
	case "jeden":
		return true
	case "thirdparty":
		return true
	case "bei":
		return true
	case "manchen":
		return true
	case "ander":
		return true
	case "solchem":
		return true
	case "habe":
		return true
	case "koennen":
		return true
	case "den":
		return true
	case "anders":
		return true
	case "muss":
		return true
	case "haben":
		return true
	case "demselben":
		return true
	case "aus":
		return true
	case "in":
		return true
	case "allen":
		return true
	case "keinem":
		return true
	case "während":
		return true
	case "eurer":
		return true
	case "derer":
		return true
	case "anderem":
		return true
	case "nichts":
		return true
	case "instantarticles":
		return true
	case "jeder":
		return true
	case "ein":
		return true
	case "eine":
		return true
	case "solches":
		return true
	case "von":
		return true
	case "denselben":
		return true
	case "andere":
		return true
	case "indem":
		return true
	case "eurem":
		return true
	case "selbst":
		return true
	case "zum":
		return true
	case "poid":
		return true
	case "getrenderedemetriqcontent":
		return true
	case "auch":
		return true
	case "keinen":
		return true
	case "alle":
		return true
	case "cms":
		return true
	case "htm":
		return true
	case "welchen":
		return true
	case "deines":
		return true
	case "anderr":
		return true
	case "derselben":
		return true
	case "sollte":
		return true
	case "könnte":
		return true
	case "wirst":
		return true
	case "eures":
		return true
	case "fuer":
		return true
	case "meinen":
		return true
	case "wo":
		return true
	case "ihrer":
		return true
	case "man":
		return true
	case "dazu":
		return true
	case "der":
		return true
	case "euer":
		return true
	case "will":
		return true
	case "sehr":
		return true
	case "ob":
		return true
	case "dem":
		return true
	case "ins":
		return true
	case "aber":
		return true
	case "einen":
		return true
	case "sonst":
		return true
	case "was":
		return true
	case "manche":
		return true
	case "static":
		return true
	case "im":
		return true
	case "weiter":
		return true
	case "eines":
		return true
	case "mancher":
		return true
	case "wir":
		return true
	case "würden":
		return true
	case "derselbe":
		return true
	case "deinem":
		return true
	case "wie":
		return true
	case "wieder":
		return true
	case "seine":
		return true
	case "mich":
		return true
	case "hatten":
		return true
	case "hatte":
		return true
	case "jener":
		return true
	case "daß":
		return true
	case "embedded":
		return true
	case "und":
		return true
	case "seinen":
		return true
	case "uns":
		return true
	case "forum":
		return true
	case "thread":
		return true
	case "jedem":
		return true
	case "meiner":
		return true
	case "über":
		return true
	case "jetzt":
		return true
	case "diese":
		return true
	case "ich":
		return true
	case "keines":
		return true
	case "aller":
		return true
	case "durch":
		return true
	case "meine":
		return true
	case "damit":
		return true
	case "weg":
		return true
	case "sondern":
		return true
	case "unseren":
		return true
	case "wollte":
		return true
	case "widget":
		return true
	case "anderm":
		return true
	case "wuerden":
		return true
	case "als":
		return true
	case "unserem":
		return true
	case "da":
		return true
	case "ueber":
		return true
	case "waehrend":
		return true
	case "koennte":
		return true
	case "zwar":
		return true
	case "hab":
		return true
	case "wuerde":
		return true
	case "anderes":
		return true
	case "so":
		return true
	case "lightbox":
		return true
	case "dieselben":
		return true
	case "dein":
		return true
	case "diesen":
		return true
	case "alles":
		return true
	case "wollen":
		return true
	case "zur":
		return true
	case "kein":
		return true
	case "etwas":
		return true
	case "mit":
		return true
	case "an":
		return true
	case "jedes":
		return true
	case "deine":
		return true
	case "oder":
		return true
	case "dort":
		return true
	case "bis":
		return true
	case "einiges":
		return true
	case "kann":
		return true
	case "waren":
		return true
	case "hin":
		return true
	case "das":
		return true
	case "wenn":
		return true
	case "php":
		return true
	case "dies":
		return true
	case "ihre":
		return true
	case "euch":
		return true
	case "unter":
		return true
	case "anderer":
		return true
	case "solchen":
		return true
	case "für":
		return true
	case "jenem":
		return true
	case "hinter":
		return true
	case "welcher":
		return true
	case "dieses":
		return true
	case "wird":
		return true
	case "pid":
		return true
	case "doch":
		return true
	case "dieselbe":
		return true
	case "werde":
		return true
	case "noch":
		return true
	case "ihren":
		return true
	case "machen":
		return true
	case "jenes":
		return true
	case "einige":
		return true
	case "einigen":
		return true
	case "welchem":
		return true
	case "ist":
		return true
	case "jene":
		return true
	case "um":
		return true
	case "ihnen":
		return true
	case "html":
		return true
	case "jede":
		return true
	case "du":
		return true
	case "es":
		return true
	case "zwischen":
		return true
	case "einer":
		return true
	case "nach":
		return true
	case "anderen":
		return true
	case "dass":
		return true
	case "jsp":
		return true
	case "seinem":
		return true
	case "manches":
		return true
	case "unsere":
		return true
	case "gegen":
		return true
	case "iframe":
		return true
	case "https":
		return true
	case "ihrem":
		return true
	case "weil":
		return true
	case "ihn":
		return true
	case "werden":
		return true
	case "andern":
		return true
	case "keine":
		return true
	case "desselben":
		return true
	case "viel":
		return true
	case "downloads":
		return true
	case "bin":
		return true
	case "deinen":
		return true
	case "hat":
		return true
	case "gewesen":
		return true
	case "nicht":
		return true
	case "diesem":
		return true
	case "ohne":
		return true
	case "welches":
		return true
	case "einigem":
		return true
	case "dann":
		return true
	case "einig":
		return true
	case "tid":
		return true
	case "zu":
		return true
	case "einmal":
		return true
	case "seines":
		return true
	case "er":
		return true
	case "mir":
		return true
	case "auf":
		return true
	case "dessen":
		return true
	case "sid":
		return true
	case "mein":
		return true
	case "seiner":
		return true
	case "musste":
		return true
	case "nur":
		return true
	case "einiger":
		return true
	case "nun":
		return true
	case "dich":
		return true
	case "stats":
		return true
	case "deiner":
		return true
	case "welche":
		return true
	case "unseres":
		return true
	case "am":
		return true
	case "warst":
		return true
	case "bist":
		return true
	case "würde":
		return true
	case "solche":
		return true
	case "einem":
		return true
	case "denn":
		return true
	case "diff":
		return true
	case "also":
		return true
	case "sie":
		return true
	case "hier":
		return true
	case "ihr":
		return true
	case "vor":
		return true
	case "des":
		return true
	case "allem":
		return true
	case "keiner":
		return true
	case "unser":
		return true
	case "titel":
		return true
	case "sein":
		return true
	case "vom":
		return true
	case "widgets":
		return true
	case "dieser":
		return true
	case "sind":
		return true
	case "meines":
		return true
	case "dir":
		return true
	case "eure":
		return true
	case "archiv":
		return true
	case "ihm":
		return true
	case "solcher":
		return true
	case "die":
		return true
	case "dasselbe":
		return true
	case "können":
		return true
	case "sich":
		return true
	case "http":
		return true
	case "soll":
		return true
	default:
		return false
	}
}

//IsEnglishStopWord returns true if word is stop word
func IsEnglishStopWord(word string) bool {
	switch word {
	case "www":
		return true
	case "myself":
		return true
	case "our":
		return true
	case "ours":
		return true
	case "ourselves":
		return true
	case "you":
		return true
	case "your":
		return true
	case "yours":
		return true
	case "yourself":
		return true
	case "yourselves":
		return true
	case "him":
		return true
	case "his":
		return true
	case "himself":
		return true
	case "she":
		return true
	case "her":
		return true
	case "hers":
		return true
	case "herself":
		return true
	case "its":
		return true
	case "itself":
		return true
	case "they":
		return true
	case "them":
		return true
	case "their":
		return true
	case "theirs":
		return true
	case "themselves":
		return true
	case "what":
		return true
	case "which":
		return true
	case "who":
		return true
	case "whom":
		return true
	case "this":
		return true
	case "that":
		return true
	case "these":
		return true
	case "those":
		return true
	case "are":
		return true
	case "was":
		return true
	case "were":
		return true
	case "been":
		return true
	case "being":
		return true
	case "have":
		return true
	case "has":
		return true
	case "had":
		return true
	case "having":
		return true
	case "does":
		return true
	case "did":
		return true
	case "doing":
		return true
	case "the":
		return true
	case "and":
		return true
	case "but":
		return true
	case "because":
		return true
	case "until":
		return true
	case "while":
		return true
	case "for":
		return true
	case "with":
		return true
	case "about":
		return true
	case "against":
		return true
	case "between":
		return true
	case "into":
		return true
	case "through":
		return true
	case "during":
		return true
	case "before":
		return true
	case "after":
		return true
	case "above":
		return true
	case "below":
		return true
	case "from":
		return true
	case "down":
		return true
	case "out":
		return true
	case "off":
		return true
	case "over":
		return true
	case "under":
		return true
	case "again":
		return true
	case "further":
		return true
	case "then":
		return true
	case "once":
		return true
	case "here":
		return true
	case "there":
		return true
	case "when":
		return true
	case "where":
		return true
	case "why":
		return true
	case "how":
		return true
	case "all":
		return true
	case "any":
		return true
	case "both":
		return true
	case "each":
		return true
	case "few":
		return true
	case "more":
		return true
	case "most":
		return true
	case "other":
		return true
	case "some":
		return true
	case "such":
		return true
	case "nor":
		return true
	case "not":
		return true
	case "only":
		return true
	case "own":
		return true
	case "same":
		return true
	case "than":
		return true
	case "too":
		return true
	case "very":
		return true
	case "can":
		return true
	case "will":
		return true
	case "just":
		return true
	case "don":
		return true
	case "should":
		return true
	case "now":
		return true
	default:
		return false
	}
}
