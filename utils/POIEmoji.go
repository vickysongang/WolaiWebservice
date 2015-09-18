// POIEmoji
package utils

import "unicode/utf8"

// 过滤 emoji 表情
func FilterEmoji(content string) string {
	newContent := ""
	for _, value := range content {
		r, size := utf8.DecodeRuneInString(string(value))
		if size <= 3 {
			newContent += string(value)
		} else {
			newContent += emojiToUnicode(r)
		}
	}
	return newContent
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
