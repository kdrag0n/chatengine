package main

import (
	"fmt"
	"bytes"
	"time"
)

var (
	byteUnits = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
)

func renderType(buf *bytes.Buffer, v interface{}) {
	switch val := v.(type) {
	case codeSect:
		buf.WriteString(`<pre style="display:inline;"><code>`)
		buf.WriteString(string(val))
		buf.WriteString(`</code></pre>`)
	case []codeSect:
		buf.WriteString(`<ul>`)
		for _, i := range val {
			buf.WriteString(`<li><pre style="display:inline;"><code>`)
			buf.WriteString(string(i))
			buf.WriteString(`</code></pre></li>`)
		}
		buf.WriteString(`</ul>`)
	case string:
		buf.WriteString(val)
	case []string:
		buf.WriteString(`<ul>`)
		for _, i := range val {
			buf.WriteString(`<li>`)
			buf.WriteString(i)
			buf.WriteString(`</li>`)
		}
		buf.WriteString(`</ul>`)
	case time.Time:
		buf.WriteString(val.Format(preciseTimeFmt))
	case []interface{}:
		buf.WriteString(`<ul>`)
		for _, i := range val {
			buf.WriteString(`<li>`)
			renderType(buf, i)
			buf.WriteString(`</li>`)
		}
		buf.WriteString(`</ul>`)
	case []genMap:
		buf.WriteString(`<ul>`)
		for _, m := range val {
			buf.WriteString(`<li><ul>`)
			for kk, vv := range m {
				buf.WriteString(`<li>`)
				buf.WriteString(kk)
				buf.WriteString(`: `)
				renderType(buf, vv)
				buf.WriteString(`</li>`)
			}
			buf.WriteString(`</ul></li>`)
		}
		buf.WriteString(`</ul>`)
	case map[string]string:
		buf.WriteString(`<ul>`)
		for kk, vv := range val {
			buf.WriteString(`<li>`)
			buf.WriteString(fmt.Sprintf("%s: %s", kk, vv))
			buf.WriteString(`</li>`)
		}
		buf.WriteString(`</ul>`)
	case map[string]interface{}:
		buf.WriteString(`<ul>`)
		for kk, vv := range val {
			buf.WriteString(`<li>`)
			buf.WriteString(kk)
			buf.WriteString(`: `)
			renderType(buf, vv)
			buf.WriteString(`</li>`)
		}
		buf.WriteString(`</ul>`)
	case map[interface{}]interface{}:
		buf.WriteString(`<ul>`)
		for kk, vv := range val {
			buf.WriteString(`<li>`)
			renderType(buf, kk)
			buf.WriteString(`: `)
			renderType(buf, vv)
			buf.WriteString(`</li>`)
		}
		buf.WriteString(`</ul>`)
	case []byte:
		buf.WriteString(`bytes <code>`)
		buf.Write(val)
		buf.WriteString(`</code>`)
	case float32, float64:
		buf.WriteString(fmt.Sprintf("%.4f", val))
	default:
		buf.WriteString(fmt.Sprintf("%+v", val))
	}
}

func byteNum(bytes uint64) string {
	smartNum := float32(bytes)
	unitIdx := 0

	for smartNum >= 1000 {
		smartNum /= 1000
		unitIdx++
	}

	return fmt.Sprintf("%.4f %s (%d bytes)", smartNum, byteUnits[unitIdx], bytes)
}