package validator

import (
	"regexp"
	"unicode/utf8"
)

func CheckValidPhoneNumber(phoneNumber string) error {
	reg := regexp.MustCompile(`^(?:010|020)-?\d{3,4}-?\d{4}$`)
	isValid := reg.MatchString(phoneNumber)

	if !isValid {
		return ErrInvalid
	}

	return nil
}

func CheckValidNickname(nickname string) error {
	length := utf8.RuneCountInString(nickname)

	if length < 2 || length > 10 {
		return ErrInvalid
	}

	// 한글, 영문, 숫자로만 구성되었는지 확인
	reg := regexp.MustCompile(`^[가-힣A-Za-z0-9]+$`)
	isValid := reg.MatchString(nickname)

	if !isValid {
		return ErrInvalid
	}

	return nil
}

func CheckValidPassword(password string) error {
	// 길이 검사
	if len(password) < 8 || len(password) > 20 {
		return ErrInvalid
	}

	// 각 요구사항을 정규식으로 검사
	reEnglish := regexp.MustCompile(`[A-Za-z]`)                  // 영문자 존재 여부
	reNumber := regexp.MustCompile(`[0-9]`)                      // 숫자 존재 여부
	reSpecial := regexp.MustCompile(`[~!@#$%^&*]`)               // 특수문자(~!@#$%^&*) 존재 여부
	reOnlyValid := regexp.MustCompile(`^[A-Za-z0-9~!@#$%^&*]+$`) // 허용된 문자 이외의 문자가 있는지

	containsEnglish := reEnglish.MatchString(password)
	containsNumber := reNumber.MatchString(password)
	containsSpecial := reSpecial.MatchString(password)
	containsOnlyValid := reOnlyValid.MatchString(password)

	isValid := containsEnglish && containsNumber && containsSpecial && containsOnlyValid

	if !isValid {
		return ErrInvalid
	}

	return nil
}
