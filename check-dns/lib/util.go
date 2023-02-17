package checkdns

import (
	"strings"
)

// ref: https://github.com/miekg/dns/blob/d8fbd0a7551ae5476df0eb003b247322761c6e82/types.go#L524
func sprintTxt(txt []string) string {
	var out strings.Builder
	for i, s := range txt {
		out.Grow(3 + len(s))
		if i > 0 {
			out.WriteString(` "`)
		} else {
			out.WriteByte('"')
		}
		for j := 0; j < len(s); {
			b, n := nextByte(s, j)
			if n == 0 {
				break
			}
			writeTXTStringByte(&out, b)
			j += n
		}
		out.WriteByte('"')
	}
	return out.String()
}

func nextByte(s string, offset int) (byte, int) {
	if offset >= len(s) {
		return 0, 0
	}
	if s[offset] != '\\' {
		// not an escape sequence
		return s[offset], 1
	}
	switch len(s) - offset {
	case 1: // dangling escape
		return 0, 0
	case 2, 3: // too short to be \ddd
	default: // maybe \ddd
		if isDigit(s[offset+1]) && isDigit(s[offset+2]) && isDigit(s[offset+3]) {
			return dddStringToByte(s[offset+1:]), 4
		}
	}
	// not \ddd, just an RFC 1035 "quoted" character
	return s[offset+1], 2
}

func isDigit(b byte) bool { return b >= '0' && b <= '9' }

func dddStringToByte(s string) byte {
	_ = s[2] // bounds check hint to compiler; see golang.org/issue/14808
	return byte((s[0]-'0')*100 + (s[1]-'0')*10 + (s[2] - '0'))
}

func writeTXTStringByte(s *strings.Builder, b byte) {
	switch {
	case b == '"' || b == '\\':
		s.WriteByte('\\')
		s.WriteByte(b)
	case b < ' ' || b > '~':
		s.WriteString(escapeByte(b))
	default:
		s.WriteByte(b)
	}
}

func escapeByte(b byte) string {
	if b < ' ' {
		return escapedByteSmall[b*4 : b*4+4]
	}

	b -= '~' + 1
	// The cast here is needed as b*4 may overflow byte.
	return escapedByteLarge[int(b)*4 : int(b)*4+4]
}

const (
	escapedByteSmall = "" +
		`\000\001\002\003\004\005\006\007\008\009` +
		`\010\011\012\013\014\015\016\017\018\019` +
		`\020\021\022\023\024\025\026\027\028\029` +
		`\030\031`
	escapedByteLarge = `\127\128\129` +
		`\130\131\132\133\134\135\136\137\138\139` +
		`\140\141\142\143\144\145\146\147\148\149` +
		`\150\151\152\153\154\155\156\157\158\159` +
		`\160\161\162\163\164\165\166\167\168\169` +
		`\170\171\172\173\174\175\176\177\178\179` +
		`\180\181\182\183\184\185\186\187\188\189` +
		`\190\191\192\193\194\195\196\197\198\199` +
		`\200\201\202\203\204\205\206\207\208\209` +
		`\210\211\212\213\214\215\216\217\218\219` +
		`\220\221\222\223\224\225\226\227\228\229` +
		`\230\231\232\233\234\235\236\237\238\239` +
		`\240\241\242\243\244\245\246\247\248\249` +
		`\250\251\252\253\254\255`
)

// ref: https://github.com/miekg/dns/blob/d8fbd0a7551ae5476df0eb003b247322761c6e82/types.go#L455
func sprintName(s string) string {
	var dst strings.Builder

	for i := 0; i < len(s); {
		if s[i] == '.' {
			if dst.Len() != 0 {
				dst.WriteByte('.')
			}
			i++
			continue
		}

		b, n := nextByte(s, i)
		if n == 0 {
			// Drop "dangling" incomplete escapes.
			if dst.Len() == 0 {
				return s[:i]
			}
			break
		}
		if isDomainNameLabelSpecial(b) {
			if dst.Len() == 0 {
				dst.Grow(len(s) * 2)
				dst.WriteString(s[:i])
			}
			dst.WriteByte('\\')
			dst.WriteByte(b)
		} else if b < ' ' || b > '~' { // unprintable, use \DDD
			if dst.Len() == 0 {
				dst.Grow(len(s) * 2)
				dst.WriteString(s[:i])
			}
			dst.WriteString(escapeByte(b))
		} else {
			if dst.Len() != 0 {
				dst.WriteByte(b)
			}
		}
		i += n
	}
	if dst.Len() == 0 {
		return s
	}
	return dst.String()
}

func isDomainNameLabelSpecial(b byte) bool {
	switch b {
	case '.', ' ', '\'', '@', ';', '(', ')', '"', '\\':
		return true
	}
	return false
}
