package validator

import (
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/purplior/podoroot/domain/assistant"
	"github.com/purplior/podoroot/domain/assisterform"
)

var (
	ErrInvalidAssistantTitle       = errors.New("잘못된 형식의 제목이에요")
	ErrInvalidAssistantDescription = errors.New("잘못된 형식의 설명이에요")
	ErrInvalidTags                 = errors.New("잘못된 형식의 태그에요")
	ErrInvalidFields               = errors.New("잘못된 형식의 입력이에요")
	ErrInvalidQueryMessages        = errors.New("잘못된 형식의 질의문이에요")
)

func checkValidAssistantTitle(
	title string,
) bool {
	// 1. 길이 검사 (3 ~ 20 글자)
	length := utf8.RuneCountInString(title)
	if length < 3 || length > 20 {
		return false
	}

	// 2. 허용된 문자만 사용되었는지 확인
	//    허용: 한국어(가-힣,ㄱ-ㅎ,ㅏ-ㅣ), 영어 대소문자(a-zA-Z), 숫자(0-9), 공백(\s)
	allowedChars := regexp.MustCompile(`^[가-힣ㄱ-ㅎㅏ-ㅣa-zA-Z0-9\s]+$`)
	if !allowedChars.MatchString(title) {
		return false
	}

	// 3. 연속된 두 개 이상의 공백이 없어야 함
	if regexp.MustCompile(`\s{2}`).MatchString(title) {
		return false
	}

	// 4. 한국어나 영어 문자가 최소 한 글자 이상 포함
	hasKoreanOrEnglish := regexp.MustCompile(`[가-힣ㄱ-ㅎㅏ-ㅣa-zA-Z]`)

	return hasKoreanOrEnglish.MatchString(title)
}

func checkValidAssistantDescription(
	description string,
) bool {
	// 1. 길이 검사 (10~150자)
	length := utf8.RuneCountInString(description)
	if length < 10 || length > 150 {
		return false
	}

	// 2. 허용된 문자만 사용되었는지 확인
	//    허용 문자: 한국어(가-힣,ㄱ-ㅎ,ㅏ-ㅣ), 영어 대소문자(a-zA-Z), 숫자(0-9),
	//              마침표(.), 콤마(,) 및 공백(\s)
	allowedChars := regexp.MustCompile(`^[가-힣ㄱ-ㅎㅏ-ㅣa-zA-Z0-9.,!?_~\s]+$`)
	if !allowedChars.MatchString(description) {
		return false
	}

	// 3. 연속된 두 개 이상의 공백이 없어야 함
	if regexp.MustCompile(`\s{2}`).MatchString(description) {
		return false
	}

	// 4. 한국어나 영어 문자가 최소 한 글자 이상 포함
	hasKoreanOrEnglish := regexp.MustCompile(`[가-힣ㄱ-ㅎㅏ-ㅣa-zA-Z]`)

	return hasKoreanOrEnglish.MatchString(description)
}

func checkValidAssistantTags(
	tags []string,
) bool {
	isTag := len(tags) > 0 && len(tags[0]) > 0

	return isTag
}

func checkValidAssisterFieldName(input string) bool {
	re := regexp.MustCompile(`^[가-힣ㄱ-ㅎㅏ-ㅣa-zA-Z0-9.,\s]{1,15}$`)

	return re.MatchString(input)
}

func checkValidAssistantFields(
	fields []assisterform.AssisterField,
) bool {
	checkingMap := map[string]bool{}

	for _, field := range fields {
		name := strings.TrimSpace(field.Name)

		if len(name) == 0 {
			return false
		}
		if !checkValidAssisterFieldName(name) {
			return false
		}

		v, ok := checkingMap[name]
		if v && ok {
			return false
		}

		if field.Type == assisterform.AssisterFieldType_ParagraphGroup {
			mapChildren := field.Option["children"].([]interface{})
			children := make([]assisterform.AssisterField, len(mapChildren))
			for i, mapChild := range mapChildren {
				child, err := assisterform.MakeAssisterFieldFromMap(mapChild.(map[string]interface{}))
				if err != nil {
					return false
				}

				children[i] = child
			}

			isChildrenValid := checkValidAssistantFields(children)
			if !isChildrenValid {
				return false
			}
		}

		checkingMap[name] = true
	}

	return true
}

func checkValidAssistantQueryMessages(
	fields []assisterform.AssisterField,
	queryMessages []assisterform.AssisterQueryMessage,
) bool {
	names := make([]string, len(fields))
	for i, field := range fields {
		names[i] = field.Name
	}

	for _, queryMessage := range queryMessages {
		if len(queryMessage.Content) < 1 {
			return false
		}

		if len(strings.Split(queryMessage.Content, "\n")) > 200 {
			return false
		}
	}

	for _, name := range names {
		isInclude := false
		for _, queryMessage := range queryMessages {
			if strings.Contains(queryMessage.Content, name) {
				isInclude = true
			}
		}
		if !isInclude {
			return false
		}
	}

	return true
}

func CheckValidAssistantRegisterRequest(
	request assistant.RegisterOneRequest,
) error {
	if !checkValidAssistantTitle(request.Title) {
		return ErrInvalidAssistantTitle
	}
	if !checkValidAssistantDescription(request.Description) {
		return ErrInvalidAssistantDescription
	}
	if !checkValidAssistantTags(request.Tags) {
		return ErrInvalidTags
	}
	if !checkValidAssistantFields(request.Fields) {
		return ErrInvalidFields
	}
	if !checkValidAssistantQueryMessages(
		request.Fields,
		request.QueryMessages,
	) {
		return ErrInvalidQueryMessages
	}

	return nil
}

func CheckValidAssistantUpdateRequest(
	request assistant.UpdateOneRequest,
) error {
	if !checkValidAssistantTitle(request.Title) {
		return ErrInvalidAssistantTitle
	}
	if !checkValidAssistantDescription(request.Description) {
		return ErrInvalidAssistantDescription
	}
	if !checkValidAssistantTags(request.Tags) {
		return ErrInvalidTags
	}
	if !checkValidAssistantFields(request.Fields) {
		return ErrInvalidFields
	}
	if !checkValidAssistantQueryMessages(
		request.Fields,
		request.QueryMessages,
	) {
		return ErrInvalidQueryMessages
	}

	return nil
}
