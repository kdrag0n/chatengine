package util

import (
	"bytes"
	"regexp"
	"strings"
	"time"
	"unicode"
)

//go:generate go run gen/corrections.go

const (
	packedContractions = `ain't	is not
aren't	are not
can't	cannot,can not
cain't	cannot,can notU
could've	could have
couldn't	could not
could of	could haveU
didn't	did not
doesn't	does not
don't	do not
gonna	going toU
gotta	got toU
hadn't	had not
hasn't	has not
haven't	have not
he'd	he would,he had
he'll	he will,he shall
he's	he is,he has
how'd	how did,how would
how'll	how will
how's	how is,how has,how does
I'd	I would,I hadO
I'll	I will,I shall
I'm	I am
I've	I have
isn't	is not
it'd	it would
it'll	it will,it shall
it's	it is,it has
mayn't	may notU
may've	may haveU
may of	may haveU
mightn't	might notU
might've	might have
might of	might haveU
mustn't	must notU
must've	must have
must of	must haveU
needn't	need not
o'clock	of the clock
ol'	oldU
oughtn't	ought notU
shan't	shall notU
she'd	she would,she hadO
she'll	she will,she shall
she's	she is,she has
should've	should have
shouldn't	should not
should of	should have
somebody's	somebody is,somebody has
someone's	someone is,someone has
something's	something is,something has
that'll	that will,that shall
that're	that are
that's	that is,that has
that'd	that would,that had
there'd	there would,there had
there're	there are
there's	there is,there has
these're	these are
they'd	they would,they had
they'll	they will,they shall
they're	they are
they've	they have
this's	this is,this hasU
those're	those areU
'tis	it isU
'twas	it wasU
wasn't	was not
we'd	we would,we hadO
we'll	we willO
we're	we areO
we've	we have
weren't	were not
what'd	what did
what'll	what will,what shall
what're	what are
what's	what is,what has,what does
what've	what haveU
what of	what haveU
when's	when is,when hasU
where'd	where did
where're	where are
where's	where is,where has,where does
where've	where have
where of	where have
which's	which is,which hasU
who'd	who would,who had,who did
who'll	who will,who shall
who're	who are
who's	who is,who has,who does
who've	who have
why'd	why did
why're	why are
why's	why is,why has,why does
willn't	will notU
won't	will not
would've	would have
would of	would have
wouldn't	would not
wanna	want toU
y'all	you all
you'd	you would,you had
you'll	you will,you shall
you're	you are
you've	you have`
)

var (
	iCapsRegexp                = regexp.MustCompile(`(\b)i(\b)`)
	singleSpaceRegexp          = regexp.MustCompile(`\s+`)
	leadingTrailingSpaceRegexp = regexp.MustCompile(`(?:^ | $)`)
	newLineRegexp              = regexp.MustCompile(`\r?\n`)
	firstWordCapRegexp         = regexp.MustCompile(`(?:[〞〟".!?·•>] ?|^)([a-z])`)
	fixPunctSpaceRegexp        = regexp.MustCompile(`(?: [,.]| [:()\[\]"] | !|$ )`)
	multiPunctRegexp           = regexp.MustCompile(`\pP{2,}`)
	urlRegexp                  = regexp.MustCompile(`(?i)https?://(?:[a-z0-9\-]\.(?:[a-z]{2,18}|xn--[a-z0-9]{4,20})|(?:[0-9]{1,3}\.){3}[0-9]{1,3}|\[[a-f0-9:]+\])(?::[0-9]{1,5})?(?:/[^/]+(?:\?(?:[a-z0-9]+=[^=?&/]*&?)+)?)*/?`)
	formatCache                = sfNewSize(64)

	contractionData = func() map[string]Contraction {
		rootSplits := strings.FieldsFunc(packedContractions, func(r rune) bool { return r == '\n' })
		data := make(map[string]Contraction, len(rootSplits))

		for _, packedLine := range rootSplits {
			tabIndex := strings.IndexByte(packedLine, '\t')
			contraction := packedLine[:tabIndex]
			expandedsPacked := packedLine[tabIndex+1:]

			exp := strings.FieldsFunc(expandedsPacked, func(r rune) bool { return r == ',' })
			last := exp[len(exp) - 1]
			lastByte := last[len(last) - 1]
			endsWithU := lastByte == 'U'
			endsWithO := lastByte == 'O'
			if endsWithU || endsWithO {
				exp[len(exp) - 1] = last[:len(last) - 1]
			}

			data[contraction] = Contraction{
				Contracted: contraction,
				Expanded: exp,
				ShouldContract: !endsWithU,
				OnlyIfFirstWord: endsWithO,
			}
		}

		return data
	}()

	contractionAposMap    map[string]Contraction
	contractionAposRegexp = func() *regexp.Regexp {
		aposMap := make(map[string]Contraction, len(contractionData))

		for aposContraction, cont := range contractionData {
			aposMap[strings.ToLower(strings.Replace(aposContraction, "'", "", -1))] = cont
		}
		contractionAposMap = aposMap

		regex := bytes.NewBufferString(`(?i)(\b)(?:`)
		idx, lastIdx := 0, len(aposMap) - 1
		for noApos := range aposMap {
			regex.WriteString(noApos)

			if idx != lastIdx {
				regex.WriteByte('|')
			}
			idx++
		}
		regex.WriteString(`)(\b)`)

		return regexp.MustCompile(regex.String())
	}()

	contractionContractMap    map[string]Contraction
	contractionContractRegexp = func() *regexp.Regexp {
		contractMap := make(map[string]Contraction, len(contractionData))

		for _, expandeds := range contractionData {
			if !expandeds.ShouldContract {
				continue
			}
			
			for _, expanded := range expandeds.Expanded {
				contractMap[strings.ToLower(expanded)] = expandeds
			}
		}
		contractionContractMap = contractMap

		regex := bytes.NewBufferString(`(?i)(\b)(?:`)
		idx, lastIdx := 0, len(contractMap)-1
		for expanded := range contractMap {
			regex.WriteString(expanded)

			if idx != lastIdx {
				regex.WriteByte('|')
			}
			idx++
		}
		regex.WriteString(`)(\b)`)

		return regexp.MustCompile(regex.String())
	}()
)

// FormatResults stores the results from a successful string format.
type FormatResults struct {
	result       string
	creationTime time.Time
}

type Contraction struct {
	Contracted string
	Expanded []string
	ShouldContract bool
	OnlyIfFirstWord bool
}

func safeSliceString(src string, sfrom, sto int) string {
	sLen := len(src)

	if sfrom >= sLen {
		sfrom = sLen - 1
	}
	if sto >= sLen {
		sto = sLen - 1
	}

	return src[sfrom:sto]
}

func contractionAposCallback(src string) func(string) string {
	return func(match string) string {
		replacement := contractionAposMap[strings.TrimSpace(strings.ToLower(match))]
		matches := contractionAposRegexp.FindStringSubmatchIndex(match)

		if len(matches) < 6 {
			return match
		} else if replacement.OnlyIfFirstWord {
			spaceIdx := strings.IndexByte(src, ' ')
			if spaceIdx != -1 && spaceIdx < matches[3] {
				return match
			}
		}

		return safeSliceString(src, matches[2], matches[3]) + replacement.Contracted + safeSliceString(src, matches[4], matches[5])
	}
}

func contractionContractCallback(src string) func(string) string {
	return func(match string) string {
		replacement := contractionContractMap[strings.TrimSpace(strings.ToLower(match))]
		matches := contractionContractRegexp.FindStringSubmatchIndex(match)

		if len(matches) < 6 {
			return match
		}

		return safeSliceString(src, matches[2], matches[3]) + replacement.Contracted + safeSliceString(src, matches[4], matches[5])
	}
}

func fixPunctSpaceCallback(match string) string {
	if match == " : " {
		return ": "
	}

	for _, b := range match {
		if b != ' ' {
			return string(b)
		}
	}

	return match
}

func iCapsCallback(src string) func(string) string {
	return func(match string) string {
		matches := iCapsRegexp.FindStringSubmatchIndex(match)
		return src[matches[2]:matches[3]] + "I" + src[matches[4]:matches[5]]
	}
}

func firstWordCapCallback(match string) string {
	runes := []rune(match)

	if len(runes) == 1 {
		return string(unicode.ToUpper(runes[0]))
	}

	lastPunc := runes[0]
	lastI := len(runes) - 1
	capitalized := unicode.ToUpper(runes[lastI])

	if unicode.IsSpace(runes[1]) {
		runes[lastI] = capitalized
		return string(runes)
	}

	return string([]rune{lastPunc, ' ', capitalized})
}

func multiPunctCallback(match string) string {
	first := match[0]

	switch first {
	case '.':
		return "..."
	case '!':
		return "!!"
	case '$':
		return "$$$"
	case '=':
		return "=="
	default:
		return string(first)
	}
}

// FormatCacheJanitor cleans up the message format cache in the background.
func FormatCacheJanitor() {
	for range time.Tick(time.Minute * 2) {
		now := time.Now()

		formatCache.ForEach(func(key string, results *FormatResults) bool {
			if now.Sub(results.creationTime).Hours() > 2.9 {
				formatCache.Delete(key)
			}

			return true
		})

		if formatCache.Len() > 8192 {
			formatCache = sfNewSize(64)
		}
	}
}

// RuneReplace replaces any occurrence of rune `from` to `to` in `runes`.
func RuneReplace(runes []rune, from, to rune) []rune {
	for idx, r := range runes {
		if r == from {
			runes[idx] = to
		}
	}

	return runes
}

// ByteReplace replaces any occurrence of byte `from` to `to` in `bytes`.
func ByteReplace(bytes []byte, from, to byte) []byte {
	for idx, b := range bytes {
		if b == from {
			bytes[idx] = to
		}
	}

	return bytes
}

// IsPunct returns whether a rune is an ending punctuation symbol or not.
func IsPunct(r rune) bool {
	return r == '.' || r == '!' || r == '?' || r == '>' || r == '"' || r == '•' ||
		r == '·' || r == '〞' || r == '〟' || r == '。' || r == ')'
}

// FirstWordIsQuestion returns whether the first word of a []rune indicates a question.
func FirstWordIsQuestion(runes []rune) bool { // what, have, has, or are
	rLen := len(runes)

	return rLen >= 4 &&
		(((runes[0] == 'w' || runes[0] == 'W') && (runes[1] == 'h' || runes[1] == 'H') &&
			(runes[2] == 'a' || runes[2] == 'A') && (runes[3] == 't' || runes[3] == 'T')) ||
			((runes[0] == 'h' || runes[0] == 'H') && (runes[1] == 'a' || runes[1] == 'A') &&
				(runes[2] == 'v' || runes[2] == 'V') && (runes[3] == 'e' || runes[3] == 'E')) ||
			((runes[0] == 'h' || runes[0] == 'H') && (runes[1] == 'a' || runes[1] == 'A') &&
				(runes[2] == 's' || runes[2] == 'S')) ||
			((runes[0] == 'a' || runes[0] == 'A') && (runes[1] == 'r' || runes[1] == 'R') &&
				(runes[2] == 'e' || runes[2] == 'E')) ||
			((runes[0] == 'h' || runes[0] == 'H') && (runes[1] == 'o' || runes[1] == 'O') &&
				(runes[2] == 'w' || runes[2] == 'W'))) &&
		(unicode.IsSpace(runes[3]) || (rLen > 4 && unicode.IsSpace(runes[4])))
}

// Format formats a message for chat storage.
func Format(input string, isCJK bool) string {
	if input == "" {
		return ""
	}

	if cached, ok := formatCache.GetOK(input); ok {
		return cached.result
	}

	result := FormatMsg(input, isCJK)
	formatCache.Set(input, &FormatResults{
		result:       result,
		creationTime: time.Now(),
	})
	return result
}

// FormatMsg formats a message, bypassing the format cache.
func FormatMsg(input string, isCJK bool) string {
	runes := []rune(input)
	lastCharIsPunc := IsPunct(runes[len(runes)-1])

	if !lastCharIsPunc {
		chosenChar := '.'

		if isCJK {
			chosenChar = '。'
		} else if FirstWordIsQuestion(runes) {
			chosenChar = '?'
		}

		runes = append(runes, chosenChar)
	}

	result := string(runes)
	result = newLineRegexp.ReplaceAllLiteralString(result, " ")
	result = strings.Replace(result, "\u200b", "", -1)
	result = urlRegexp.ReplaceAllLiteralString(result, "")
	result = singleSpaceRegexp.ReplaceAllLiteralString(result, " ")
	result = leadingTrailingSpaceRegexp.ReplaceAllLiteralString(result, "")
	result = fixPunctSpaceRegexp.ReplaceAllStringFunc(result, fixPunctSpaceCallback)
	result = iCapsRegexp.ReplaceAllStringFunc(result, iCapsCallback(result))
	result = contractionAposRegexp.ReplaceAllStringFunc(result, contractionAposCallback(result))
	result = contractionContractRegexp.ReplaceAllStringFunc(result, contractionContractCallback(result))
	result = firstWordCapRegexp.ReplaceAllStringFunc(result, firstWordCapCallback)
	result = multiPunctRegexp.ReplaceAllStringFunc(result, multiPunctCallback)

	return result
}
