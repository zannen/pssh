package expand

import (
	"strconv"
	"strings"
)

func Expand(e string) ([]string, error) {

	pos := 0
	var sb strings.Builder
	for strings.Contains(e[pos:], "[") && strings.Contains(e[pos:], "]") {
		openPos := strings.Index(e[pos:], "[") + 1
		closePos := strings.Index(e[pos:], "]") + 1
		sb.WriteString(e[pos : pos+openPos])
		sb.WriteString(strings.Replace(e[pos+openPos:pos+closePos], ",", "ᚁ", 999))
		pos += closePos
	}
	if pos > 0 {
		sb.WriteString(e[pos:])
		e = sb.String()
	}

	// "a;b;c" -> ["a", "b", "c"]
	if strings.Contains(e, ";") {
		results := make([]string, 0)
		for _, item := range strings.Split(e, ";") {
			r, err := Expand(item)
			if err != nil {
				return nil, err
			}
			results = append(results, r...)
		}
		return results, nil
	}

	// "host{x,y,z}.com" -> ["hostx.com", "hosty.com", "hostz.com"]
	if strings.Contains(e, "{") && strings.Contains(e, "}") {
		openCurlyPos := strings.Index(e, "{")
		closeCurlyPos := strings.Index(e, "}")
		prefix := e[:openCurlyPos]
		mid := e[openCurlyPos+1 : closeCurlyPos]
		postfix := e[closeCurlyPos+1:]
		results := make([]string, 0)
		for _, item := range strings.Split(mid, ",") {
			sb.Reset()
			sb.WriteString(prefix)
			sb.WriteString(item)
			sb.WriteString(postfix)
			r, err := Expand(sb.String())
			if err != nil {
				return nil, err
			}
			results = append(results, r...)
		}
		return results, nil
	}

	// "host[1-3,5,7-9].com" -> ["host1.com", "host2.com", "host3.com", "host5.com", "host7.com", "host8.com", "host9.com"]
	if strings.Contains(e, "[") && strings.Contains(e, "]") {
		openSqPos := strings.Index(e, "[")
		closeSqPos := strings.Index(e, "]")
		prefix := e[:openSqPos]
		mid := e[openSqPos+1 : closeSqPos]
		postfix := e[closeSqPos+1:]
		results := make([]string, 0)
		for _, item := range strings.Split(mid, "ᚁ") {
			if strings.Contains(item, "-") {
				dashPos := strings.Index(item, "-")
				from, err := strconv.Atoi(item[:dashPos])
				if err != nil {
					return nil, err
				}
				to, err := strconv.Atoi(item[dashPos+1:])
				if err != nil {
					return nil, err
				}
				for i := from; i <= to; i++ {
					sb.Reset()
					sb.WriteString(prefix)
					sb.WriteString(strconv.Itoa(i))
					sb.WriteString(postfix)
					r, err := Expand(sb.String())
					if err != nil {
						return nil, err
					}
					results = append(results, r...)
				}
			} else {
				sb.Reset()
				sb.WriteString(prefix)
				sb.WriteString(item)
				sb.WriteString(postfix)
				r, err := Expand(sb.String())
				if err != nil {
					return nil, err
				}
				results = append(results, r...)
			}
		}
		return results, nil
	}

	return []string{e}, nil
}
