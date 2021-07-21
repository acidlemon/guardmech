package db

import (
	"context"
)

type Context = context.Context

type relationRow interface {
	TargetSeqID() int64
}

func compareSeqID(ids []int64, rows []relationRow) (added, deleted []int64) {
	// added
	if len(ids) != 0 {
		for _, v := range ids {
			found := false
			for _, w := range rows {
				if v == w.TargetSeqID() {
					found = true
					break
				}
			}
			if !found {
				added = append(added, v)
			}
		}
	}

	// deleted
	if len(rows) != 0 {
		for _, v := range rows {
			found := false
			for _, w := range ids {
				if v.TargetSeqID() == w {
					found = true
					break
				}
			}
			if !found {
				deleted = append(deleted, v.TargetSeqID())
			}
		}
	}

	return
}
