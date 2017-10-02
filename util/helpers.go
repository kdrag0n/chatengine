package util

import (
	container_list "container/list"
	"encoding/binary"
	"net"
)

// GetElemAt retrieves the element at targetIdx of list.
func GetElemAt(targetIdx int, list *container_list.List) *container_list.Element {
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

// IP2Int converts a net.IP to an uint32 IP.
func IP2Int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

// Int2IP converts an uint32 IP to a net.IP.
func Int2IP(nn uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip
}