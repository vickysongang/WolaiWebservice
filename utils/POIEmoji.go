// POIEmoji
package utils

import (
	"fmt"
	"regexp"
	"unicode/utf8"
)

// 过滤 emoji 表情
func FilterEmoji(content string) string {
	new_content := ""
	for _, value := range content {
		r, size := utf8.DecodeRuneInString(string(value))
		if size <= 3 {
			ok, _ := regexp.MatchString(`U\+E[0-9a-fA-F]{3}`, fmt.Sprintf("%U", r))
			if ok {
				new_content += emojiToUnicode(r)
			} else {
				new_content += string(value)
			}
		} else {
			new_content += emojiToUnicode(r)
		}
	}
	return new_content
}

func emojiToUnicode(runeValue rune) (unicodeValue string) {
	rInt := int(runeValue)
	if rInt < 128 {
		unicodeValue = string(runeValue)
	} else {
		//暂时用［表情］代替
		unicodeValue = "[表情]"
		//		unicodeValue = fmt.Sprintf("%U", runeValue)
	}
	return
}
