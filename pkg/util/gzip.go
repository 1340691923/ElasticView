package util

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

func GzipCompress(data string) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	gzW := gzip.NewWriter(buf)
	_, err := gzW.Write([]byte(data))
	if err != nil {
		return nil, err
	}
	gzW.Close()
	return buf.Bytes(), err
}

func GzipUnCompress(data []byte) ([]byte, error) {
	gzR, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(gzR)
	return b, err
}
