package simpleHTTP

import (
	"adapter/modules/adapter"
	"adapter/modules/prompb"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"io/ioutil"
	"net/http"
)

func RemoteWrtie(w http.ResponseWriter, r *http.Request) {
	compressed, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// resolve snappy
	reqBuf, err := snappy.Decode(nil, compressed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// resolve json
	var wreq prompb.WriteRequest

	if err := proto.Unmarshal(reqBuf, &wreq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// resolve data
	adapter.RemoteWriter(wreq)
	if _, err := w.Write([]byte("ok")); err != nil {
		return
	}
}

func RemoteRead(w http.ResponseWriter, r *http.Request) {
	compressed, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// snappy
	reqBuf, err := snappy.Decode(nil, compressed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// resolve json
	var rreq prompb.ReadRequest
	if err := proto.Unmarshal(reqBuf, &rreq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	naiveData := adapter.RemoteReader(rreq)
	data, _ := proto.Marshal(naiveData)
	// sender
	w.Header().Set("Content-Type", "application/x-protobuf")
	w.Header().Set("Content-Encoding", "snappy")
	compressed = snappy.Encode(nil, data)
	if _, err := w.Write(compressed); err != nil {
		return
	}
}
