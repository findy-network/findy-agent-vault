package resolver

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/glog"

	"github.com/findy-network/findy-agent-api/resolver"
	"github.com/findy-network/findy-agent-api/tools/data"
	"github.com/lainio/err2"
)

const maxPatchSize = 100

type PaginationParams struct {
	first  *int
	last   *int
	before *string
	after  *string
}

func logPaginationRequest(prefix string, params *PaginationParams) {
	var first, last, before, after string
	if params.first != nil {
		first = fmt.Sprintf(", first: %d", *params.first)
	}
	if params.last != nil {
		last = fmt.Sprintf(", last: %d", *params.last)
	}
	if params.before != nil {
		before = fmt.Sprintf(", before: %s", *params.before)
	}
	if params.after != nil {
		after = fmt.Sprintf(", after: %s", *params.after)
	}
	glog.V(2).Infof("%s%s%s%s%s", prefix, after, before, first, last)

}

func parseCursor(cursor string) (int64, error) {
	plain, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return 0, errors.New(resolver.ErrorCursorInvalid)
	}

	parts := strings.Split(string(plain), ":")
	if len(parts) != 2 {
		return 0, errors.New(resolver.ErrorCursorInvalid)
	}

	value, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, errors.New(resolver.ErrorCursorInvalid)
	}

	return value, nil
}

func validateFirstAndLast(first, last *int) error {
	if first == nil && last == nil {
		return errors.New(resolver.ErrorFirstLastMissing)
	}
	if (first != nil && (*first < 1 || *first > maxPatchSize)) ||
		(last != nil && (*last < 1 || *last > maxPatchSize)) {
		return errors.New(resolver.ErrorFirstLastInvalid)
	}
	return nil
}

func pick(items *data.Items, pagination *PaginationParams) (afterIndex int, beforeIndex int, err error) {
	defer err2.Return(&err)

	afterIndex = 0
	beforeIndex = items.Count() - 1

	err2.Check(validateFirstAndLast(pagination.first, pagination.last))

	if pagination.after != nil || pagination.before != nil {
		var afterVal, beforeVal int64
		if pagination.after != nil {
			afterVal, err = parseCursor(*pagination.after)
			err2.Check(err)
		}
		if pagination.before != nil {
			beforeVal, err = parseCursor(*pagination.before)
			err2.Check(err)
		}
		for index := 0; index < items.Count(); index++ {
			created := items.CreatedForIndex(index)
			if afterVal > 0 && created <= afterVal {
				afterIndex = index + 1
			}
			if beforeVal > 0 && created < beforeVal {
				beforeIndex = index
			}
			if (beforeVal > 0 && created > beforeVal) ||
				(beforeVal == 0 && created > afterVal) {
				break
			}
		}
	}

	if pagination.first != nil {
		afterPlusFirst := afterIndex + (*pagination.first - 1)
		if beforeIndex > afterPlusFirst {
			beforeIndex = afterPlusFirst
		}
	} else if pagination.last != nil {
		beforeMinusLast := beforeIndex - (*pagination.last - 1)
		if afterIndex < beforeMinusLast {
			afterIndex = beforeMinusLast
		}
	}
	return afterIndex, beforeIndex + 1, nil
}
