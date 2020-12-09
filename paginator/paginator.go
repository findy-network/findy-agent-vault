package paginator

import (
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/lainio/err2"
)

const (
	ErrorFirstLastMissing = "you must provide a `first` or `last` value to properly paginate the objects"
	ErrorFirstLastInvalid = "you must provide a valid `first` or `last` value in range 1-100"
	ErrorCursorInvalid    = "cursor value is invalid"
)

const (
	cursorPartsCount = 2
	maxPatchSize     = 100
)

type Params struct {
	First  *int
	Last   *int
	Before *string
	After  *string
	Object interface{}
}

type BatchInfo struct {
	Count  int
	Tail   bool
	Before uint64
	After  uint64
}

func LogRequest(prefix string, params *Params) {
	var first, last, before, after string
	if params.First != nil {
		first = fmt.Sprintf(", first: %d", *params.First)
	}
	if params.Last != nil {
		last = fmt.Sprintf(", last: %d", *params.Last)
	}
	if params.Before != nil {
		before = fmt.Sprintf(", before: %s", *params.Before)
	}
	if params.After != nil {
		after = fmt.Sprintf(", after: %s", *params.After)
	}
	utils.LogMed().Infof("%s%s%s%s%s", prefix, after, before, first, last)
}

func CreateCursor(created uint64, object interface{}) string {
	typeName := reflect.TypeOf(object).Name()
	return base64.StdEncoding.EncodeToString(
		[]byte(typeName + ":" + strconv.FormatUint(created, 10)),
	)
}

func ParseCursor(cursor string, object interface{}) (uint64, error) {
	plain, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return 0, errors.New(ErrorCursorInvalid)
	}

	parts := strings.Split(string(plain), ":")
	if len(parts) != cursorPartsCount {
		return 0, errors.New(ErrorCursorInvalid)
	}

	value, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return 0, errors.New(ErrorCursorInvalid)
	}

	if parts[0] != reflect.TypeOf(object).Name() {
		return 0, errors.New(ErrorCursorInvalid)
	}

	return value, nil
}

func ValidateFirstAndLast(first, last *int) (count int, valid bool, err error) {
	if first == nil && last == nil {
		return 0, false, errors.New(ErrorFirstLastMissing)
	}
	if first != nil {
		if *first < 1 || *first > maxPatchSize {
			return 0, false, errors.New(ErrorFirstLastInvalid)
		}
		return *first, false, nil
	}
	if last != nil && (*last < 1 || *last > maxPatchSize) {
		return 0, false, errors.New(ErrorFirstLastInvalid)
	}
	return *last, true, nil
}

func Validate(prefix string, params *Params) (info *BatchInfo, err error) {
	defer err2.Return(&err)
	LogRequest(prefix, params)

	var count int
	tail := false
	var before, after uint64

	count, tail, err = ValidateFirstAndLast(params.First, params.Last)
	err2.Check(err)

	if params.After != nil {
		after, err = ParseCursor(*params.After, params.Object)
		err2.Check(err)
	}
	if params.Before != nil {
		after, err = ParseCursor(*params.Before, params.Object)
		err2.Check(err)
	}

	info = &BatchInfo{
		Count:  count,
		Tail:   tail,
		After:  after,
		Before: before,
	}

	return
}
