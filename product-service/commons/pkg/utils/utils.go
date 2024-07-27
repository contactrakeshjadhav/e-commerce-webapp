package utils

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/commons/pkg/core/model"
	"github.com/pkg/errors"
)

const (
	HOSTNAME_LABEL_MAX_LEN = 63
	HOSTNAME_MAX_LEN       = 255
	EMPTY_SPACE            = ""

	patternAnythingNotLetterNumberOrHyphen = `[^a-zA-Z0-9-]+`
	patternOnlyDigitsOrLetters             = `[a-zA-Z0-9]`

	hyphenChar string = "-"
	dotChar    string = "."
)

var (
	ErrHostnameLabelCantBeEmpty                 error = errors.New("hostname label can't be empty")
	ErrHostnameCantBeEmpty                      error = errors.New("hostname can't be empty")
	ErrHostnameLabelMaxLenExceeded              error = errors.New("labels must be 63 characters or less")
	ErrHostnameMaxLenExceeded                   error = errors.New("must be 255 characters or less")
	ErrHostnameLabelMustStartWithLetterOrNumber error = errors.New("each label must start with a letter or number")
	ErrHostnameLabelMustEndWithLetterOrNumber   error = errors.New("each label must end with a letter or number")
	ErrHostnameLabelInvalidCharacter            error = errors.New("characters on each label must be a letter, number, or hyphen")
)

const RANDOM_BYTES_COUNT int = 4

type ObjectType string

const (
	ManifestItem ObjectType = "ManifestItem"
	Project      ObjectType = "Project"
	CDRFile      ObjectType = "CDRFile"
	Template     ObjectType = "Template"
	Workspace    ObjectType = "Workspace"
)

func ValidateObjectType(objectType string) bool {
	ot := ObjectType(objectType)
	switch ot {
	case ManifestItem:
		return true
	case Project:
		return true
	case CDRFile:
		return true
	case Template:
		return true
	case Workspace:
		return true
	default:
		return false
	}
}

// GetRandomHex creates a new hexadecimal string
// return error an any given operation fail
func GetRandomHex() (string, error) {
	b := make([]byte, RANDOM_BYTES_COUNT)
	n, err := rand.Read(b)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate random hex string")
	}
	if n != RANDOM_BYTES_COUNT {
		return "", errors.Errorf("failed to generate random hex string: got only %d bytes, %d expected", n, RANDOM_BYTES_COUNT)
	}
	return hex.EncodeToString(b), nil
}

// Contains Returns true if string exists in a slice of strings
func Contains(s []string, target string) bool {
	for _, item := range s {
		if item == target {
			return true
		}
	}

	return false
}

// RemoveDuplicates Removes exact duplicate values from slice
func RemoveDuplicates(s []string) []string {
	checkUniq := make(map[string]struct{})
	newValue := make([]string, 0)

	for _, item := range s {
		checkUniq[item] = struct{}{}
	}
	for item := range checkUniq {
		newValue = append(newValue, item)
	}
	return newValue
}

// GenerateValidHostnameLabel Generates valid hostname label from given string
// The values for the host name label must conform to the following rules:
// - maximum of 63 characters.
// - label must start and end with a letter or number.
// - remaining characters must be a letter, number, or hyphen.
func GenerateValidHostnameLabel(s string) (string, error) {
	if len(s) == 0 {
		return "", ErrHostnameLabelCantBeEmpty
	}
	// maximum of HOSTNAME_LABEL_MAX_LEN characters
	retval := TruncateString(s, HOSTNAME_LABEL_MAX_LEN)

	// replace spaces with hyphen
	retval = strings.ReplaceAll(retval, " ", hyphenChar)

	// avoid repeated hyphens
	retval = removeRepeatedHyphens(retval)

	// remove special characters except hyphen
	reg, err := regexp.Compile(patternAnythingNotLetterNumberOrHyphen)
	if err != nil {
		return "", err
	}
	retval = reg.ReplaceAllString(retval, "")

	if len(retval) == 0 {
		return "", ErrHostnameLabelCantBeEmpty
	}

	// must start with a letter or number
	retval, err = startWithLetterOrNumber(retval)
	if err != nil {
		return "", err
	}

	if len(retval) == 0 {
		return "", ErrHostnameLabelCantBeEmpty
	}

	// must end with a letter or number
	last := len(retval) - 1
	if string(retval[last]) == hyphenChar {
		retval = retval[:last]
	}

	return retval, nil
}

func removeRepeatedHyphens(s string) string {
	var retval strings.Builder
	// iterate characteres on given string
	for i := 0; i < len(s); i++ {
		if string(s[i]) != hyphenChar {
			// WriteByte appends the byte to buffer.
			// The returned error is always nil.
			retval.WriteByte(s[i])
			continue
		}

		// append if consecutive value is not a hyphen
		if i == len(s)-1 || s[i] != s[i+1] {
			retval.WriteByte(s[i])
		}
	}
	return retval.String()
}

func startWithLetterOrNumber(s string) (string, error) {
	reg, err := regexp.Compile(patternOnlyDigitsOrLetters)
	if err != nil {
		return "", err
	}

	return removeFirstCharIf(s, reg), nil
}

func removeFirstCharIf(s string, r *regexp.Regexp) string {
	if len(s) == 0 {
		return s
	}

	retval := s
	if !r.MatchString(string(retval[0])) {
		retval = retval[1:]
		return removeFirstCharIf(retval, r)
	}

	return retval
}

// TruncateString Truncate string based on given length
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}

	return s[:maxLen]
}

// RemoveWhitespaces Remove all whitespaces from string
func RemoveWhitespaces(s string) string {
	retval := make([]rune, 0, len(s))
	for _, ch := range s {
		if !unicode.IsSpace(ch) {
			retval = append(retval, ch)
		}
	}
	return string(retval)
}

// ValidateHostname Validates given hostname.
// It returns an error if hostname does not meet the following rules:
// - labels must be 63 characters or less
// - each label must start and end with a letter or number
// - remaining characters on each label must be a letter, number, or hyphen
// - whole hostname must not exceed length of 255 characters
func ValidateHostname(h string) error {
	if len(h) == 0 {
		return ErrHostnameCantBeEmpty
	}

	if len(h) > HOSTNAME_MAX_LEN {
		return ErrHostnameMaxLenExceeded
	}

	// remove dot suffix (if present)
	if strings.HasSuffix(h, dotChar) {
		h = h[:len(h)-1]
	}

	labels := strings.Split(h, ".")
	// validate each label
	for _, l := range labels {
		if len(l) == 0 {
			return ErrHostnameLabelCantBeEmpty
		}

		// validate label len
		if len(l) > HOSTNAME_LABEL_MAX_LEN {
			return ErrHostnameLabelMaxLenExceeded
		}

		// must start and end with a letter or number
		reg, err := regexp.Compile(patternOnlyDigitsOrLetters)
		if err != nil {
			return err
		}

		if !reg.MatchString(string(l[0])) {
			return ErrHostnameLabelMustStartWithLetterOrNumber
		}

		if !reg.MatchString(string(l[len(l)-1])) {
			return ErrHostnameLabelMustEndWithLetterOrNumber
		}

		// remaining characters must be a letter, number, or hyphen
		reg, err = regexp.Compile(patternAnythingNotLetterNumberOrHyphen)
		if err != nil {
			return err
		}

		if reg.MatchString(l) {
			return ErrHostnameLabelInvalidCharacter
		}
	}

	return nil
}

func ExecuteHTTPRequest(ctx context.Context, client *http.Client, method string, url string, payload interface{}) (*http.Response, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.Wrapf(model.ErrInternalServerFail, "failed to marshall request body: %v", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, errors.Wrapf(model.ErrInternalServerFail, "failed to create http request: %v", err)
	}

	reqIDValue := ctx.Value(model.ReqIDKey)
	reqID, ok := reqIDValue.(string)
	if !ok {
		return nil, errors.Wrapf(model.ErrInternalServerFail, "request id is of unexpected type: %v", reflect.TypeOf(reqIDValue))
	}

	tokenValue := ctx.Value(model.TokenKey)
	token, ok := tokenValue.(string)
	if !ok {
		return nil, errors.Wrapf(model.ErrInternalServerFail, "user token is of unexpected type: %v", reflect.TypeOf(reqIDValue))
	}

	req.Header.Add(model.ReqIDKey.String(), reqID)
	req.Header.Add("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(model.ErrInternalServerFail, "failed to execute http request: %v", err)
	}
	return res, nil
}

func GetUniqueElement[T any](oldElement, newElement []T) ([]T, error) {

	elementMap := make(map[interface{}]struct{})
	var uniqueElement []T
	for _, val := range oldElement {
		elementMap[val] = struct{}{}
	}

	for _, val := range newElement {
		if _, ok := elementMap[val]; !ok {
			uniqueElement = append(uniqueElement, val)
		}
	}

	if len(uniqueElement) == 0 {
		return nil, errors.New("products are already exits for the object")
	}

	return uniqueElement, nil
}
