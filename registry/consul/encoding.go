package consul

import (
	"bytes"
	"compress/zlib"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"

	"github.com/Allenxuxu/stark/registry"
)

func encode(buf []byte) string {
	var b bytes.Buffer
	defer b.Reset()

	w := zlib.NewWriter(&b)
	if _, err := w.Write(buf); err != nil {
		return ""
	}
	w.Close()

	return hex.EncodeToString(b.Bytes())
}

func decode(d string) []byte {
	hr, err := hex.DecodeString(d)
	if err != nil {
		return nil
	}

	br := bytes.NewReader(hr)
	zr, err := zlib.NewReader(br)
	if err != nil {
		return nil
	}

	rbuf, err := ioutil.ReadAll(zr)
	if err != nil {
		return nil
	}
	zr.Close()

	return rbuf
}

func encodeEndpoints(en []*registry.Endpoint) []string {
	var tags []string
	for _, e := range en {
		if b, err := json.Marshal(e); err == nil {
			tags = append(tags, "e-"+encode(b))
		}
	}
	return tags
}

func decodeEndpoints(tags []string) []*registry.Endpoint {
	var en []*registry.Endpoint

	// use the first format you find
	var ver byte

	for _, tag := range tags {
		if len(tag) == 0 || tag[0] != 'e' {
			continue
		}

		// check version
		if ver > 0 && tag[1] != ver {
			continue
		}

		var e *registry.Endpoint
		var buf []byte

		// Old encoding was plain
		if tag[1] == '=' {
			buf = []byte(tag[2:])
		}

		// New encoding is hex
		if tag[1] == '-' {
			buf = decode(tag[2:])
		}

		if err := json.Unmarshal(buf, &e); err == nil {
			en = append(en, e)
		}

		// set version
		ver = tag[1]
	}
	return en
}

func encodeVersion(v string) []string {
	return []string{"v=" + v}
}

func decodeVersion(tags []string) (string, bool) {
	for _, tag := range tags {
		if len(tag) < 2 || tag[0] != 'v' {
			continue
		}

		if tag[1] == '=' {
			return tag[2:], true
		}
	}
	return "", false
}
