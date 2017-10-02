package main

import (
	container_list "container/list"

	"github.com/abiosoft/ishell"
)

func getElemAt(targetIdx int, list *container_list.List) *container_list.Element {
	half := list.Len() / 2

	if targetIdx <= half {
		currentIdx := 0

		for e := list.Front(); e != nil; e = e.Next() {
			if currentIdx == targetIdx {
				return e
			}

			currentIdx++
		}
	} else {
		currentIdx := list.Len() - 1

		for e := list.Back(); e != nil; e = e.Prev() {
			if currentIdx == targetIdx {
				return e
			}

			currentIdx--
		}
	}

	return nil
}

func ctxPrintfAntivet(c *ishell.Context, fmt string, val ...interface{}) {
	c.Printf(fmt, val)
}
