package strgen

import "strings"

// 한글 초성 리스트
var cho = []rune{
	'ㄱ', 'ㄲ', 'ㄴ', 'ㄷ', 'ㄸ', 'ㄹ', 'ㅁ', 'ㅂ', 'ㅃ',
	'ㅅ', 'ㅆ', 'ㅇ', 'ㅈ', 'ㅉ', 'ㅊ', 'ㅋ', 'ㅌ', 'ㅍ', 'ㅎ',
}

func ExtractInitialChar(str string) string {
	if str == "" {
		return ""
	}
	r := []rune(str)[0]

	// 한글 범위 U+AC00(가) ~ U+D7A3(힣)
	if r >= 0xAC00 && r <= 0xD7A3 {
		// 한글 음절 코드 분석
		syllableIndex := r - 0xAC00
		choseongIndex := syllableIndex / (21 * 28)
		return string(cho[choseongIndex])
	} else {
		// 한글이 아닌 경우(여기서는 영어라고 가정),
		// 첫 글자를 대문자로 변환하여 반환
		return strings.ToUpper(string(r))
	}
}
