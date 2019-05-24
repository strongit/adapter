package adapter

import (
	"adapter/lib"
	"adapter/modules/conf"
	"adapter/modules/logger"
	"adapter/modules/prompb"
	"adapter/modules/tikv"
	"bytes"
	"go.uber.org/zap/buffer"
	"log"
	"strconv"
	"time"
)


var (
	buffers = buffer.NewPool()
)

func RemoteWriter(data prompb.WriteRequest) {
	strtime := time.Now().UnixNano()
	var timsSeriesCnt int64

	for _, oneDoc := range data.Timeseries {
		labels := oneDoc.Labels
		samples := oneDoc.Samples
		log.Println("Naive write data:", labels, samples, len(samples))

		//build index and return labelID
		labelID := buildIndex(labels, samples)
		//log.Println("LabelID:", labelID)

		//write timeseries data
		writeTimeseriesData(labelID, samples)

		labelsByte := lib.GetBytes(labels)
		SaveOriDoc(labelID, labelsByte)

	}

	timsSeriesCnt = int64(len(data.Timeseries))
	endtime := time.Now().UnixNano()
	difftime := endtime - strtime

	if timsSeriesCnt > 0 {
		logger.Println("writedata.log", "info", "write statistics data", string(difftime), timsSeriesCnt, difftime/timsSeriesCnt)
	}
}

//build md5 data and store to kv if not exist
func buildIndex(labels []*prompb.Label, samples []*prompb.Sample) string {
	//make md
	//key type key#value
	buf := buffers.Get()
	defer buf.Free()
	for _, v := range labels {
		buf.AppendString(v.Name)
		buf.AppendString("#")
		buf.AppendString(v.Value)
	}

	labelBytes := buf.Bytes()
	labelID := lib.MakeMDByByte(labelBytes)
	buf.Reset()

	//labels index
	for _, v := range labels {
		//key type index:label:__name__#latency
		buf.AppendString("index:label:")
		buf.AppendString(v.Name)
		buf.AppendString("#")
		buf.AppendString(v.Value)
		key := buf.String()
		buf.Reset()
		//log.Println("Write label md:", key, labelID)

		//key type index:status:__name__#latency+labelID
		buf.AppendString("index:status:")
		buf.AppendString(v.Name)
		buf.AppendString("#")
		buf.AppendString(v.Value)
		buf.AppendString("+")
		buf.AppendString(labelID)
		indexStatus := buf.Bytes()

		indexStatusKey, err := tikv.Get([]byte(indexStatus))
		if err != nil {
			logger.Println("buildIndex.log","error", "buildIndex get indexStatus err", err)
		}
		//log.Println("indexStatus:", indexStatusKey)

		//not in index
		if "" == indexStatusKey.Value {
			err := tikv.Puts([]byte(indexStatus), []byte("1"))
			if err != nil {
				logger.Println("buildIndex.log","error","buildIndex puts indexStatus err", err)
			}

			//wtire tikv
			oldKey, err := tikv.Get([]byte(key))
			if err != nil {
				logger.Println("buildIndex.log","error","buildIndex get index:label err", err)
			}
			if oldKey.Value == "" {
				err := tikv.Puts([]byte(key), []byte(labelID))
				if err != nil {
					logger.Println("buildIndex.log","error","buildIndex puts index:label err", err)
				}
			} else {
				b := bytes.NewBufferString(oldKey.Value)
				b.WriteString(labelID)
				v := b.Bytes()

				err := tikv.Puts([]byte(key), v)
				if err != nil {
					logger.Println("buildIndex.log","error","buildIndex puts index:label2 err", err)
				}
			}
		}

		buf.Reset()
	}

	buf.Reset()
	buf.AppendString("index:timeseries:")

	now := time.Now().UnixNano() / int64(time.Millisecond)

	interval := int64(conf.RunTimeInfo.TimeInterval * 1000 * 60)
	now = (now / interval) * interval

	buf.AppendString(labelID)
	buf.AppendString(":")
	buf.AppendString(strconv.FormatInt(now, 10))

	timeIndexBytes := buf.Bytes()

	//timeseries index
	for _, v := range samples {
		oldKey, err := tikv.Get(timeIndexBytes)
		if err != nil {
			logger.Println("buildIndex.log","error","buildIndex get index:timeseries err", err)
		}
		//log.Println("Timeseries indexStatus:", oldKey)
		if oldKey.Value == "" {
			err := tikv.Puts(timeIndexBytes, lib.Int64ToBytes(v.Timestamp))
			if err != nil {
				logger.Println("buildIndex.log","error","buildIndex puts index:timeseries err", err)
			}
		} else {
			bs := buffers.Get()
			bs.AppendString(oldKey.Value)
			bs.AppendString(strconv.FormatInt(v.Timestamp, 10))
			v := bs.Bytes()

			err := tikv.Puts(timeIndexBytes, v)

			if err != nil {
				logger.Println("buildIndex.log","error","buildIndex puts index:timeseries2 err", err)
			}
			bs.Free()
		}
	}

	return labelID
}

func writeTimeseriesData(labelID string, samples []*prompb.Sample) {
	buf := buffers.Get()
	defer buf.Free()
	for _, v := range samples {
		//key type timeseries:doc:labelMD#timestamp
		buf.AppendString("timeseries:doc:")
		buf.AppendString(labelID)
		buf.AppendString(":")
		buf.AppendString(strconv.FormatInt(v.Timestamp, 10))
		key := buf.Bytes()
		
		//write to tikv
		err := tikv.Puts(key, []byte(strconv.FormatFloat(v.Value, 'E', -1, 64)))
		if err != nil {
			logger.Println("writeTimeseriesData.log","error","writeTimeseriesData puts timeseries:doc err", err)
		}
		//log.Println("Write timeseries:", string(key), strconv.FormatFloat(v.Value, 'E', -1, 64))
		buf.Reset()
	}
}

func SaveOriDoc(labelID string, originalMsg []byte) {
	buf := buffers.Get()
	defer buf.Free()
	buf.AppendString("doc:")
	buf.AppendString(labelID)
	key := buf.Bytes()

	err := tikv.Puts(key, originalMsg)
	if err != nil {
		logger.Println("SaveOriDoc.log","error","SaveOriDoc puts timeseries:doc err", err)
	}
	//log.Println("Write meta:", string(key), string(originalMsg))
}
