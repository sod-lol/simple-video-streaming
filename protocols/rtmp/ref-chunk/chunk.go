package chunk

import (
	"bufio"

	"github.com/sirupsen/logrus"
)

//Chunk is basic packet in rtmp protocol that transfer over network
type Chunk struct {
	BasicHeader       *BasicHeader
	MessageHeader     *MessageHeader
	ExtendedTimeStamp *ExtendedTimeStamp
}

//NewChunk return chunk
func NewChunk() *Chunk {
	return &Chunk{
		BasicHeader:       NewBasicHeader(),
		MessageHeader:     NewMessageHeader(),
		ExtendedTimeStamp: NewExtendedTimeStamp(),
	}
}

func (c *Chunk) Read(reader *bufio.Reader) {
	logrus.Debug("[Debug] start reading chunk headers")
	c.BasicHeader.Read(reader)
	c.MessageHeader.Read(reader, c.BasicHeader.Fmt)
	if c.MessageHeader.TtimeStamp == 0xffffff || c.MessageHeader.TimeStampDelta == 0xffffff {
		c.ExtendedTimeStamp.hasExtendedTimeStamp = true
	}
	c.ExtendedTimeStamp.Read(reader)

	logrus.Debug("[Debug] finished reading chunk headers")
}

//ReadPayload is function that read payload that in chunk
func (c *Chunk) ReadPayload(reader *bufio.Reader, chunkData []byte) {
	if _, err := reader.Read(chunkData); err != nil {
		logrus.Errorf("[Error] %v occured during reading chunk with chunk stream ID %v", err, c.BasicHeader.CSID)
	}
}

func (c *Chunk) Write(writer *bufio.Writer, chunkData []byte) {
	logrus.Debug("[Debug] start writing chunk")
	c.BasicHeader.Write(writer)
	c.MessageHeader.Write(writer, c.BasicHeader.Fmt)
	if c.MessageHeader.TtimeStamp == 0xffffff || c.MessageHeader.TimeStampDelta == 0xffffff {
		c.ExtendedTimeStamp.hasExtendedTimeStamp = true
	}
	c.ExtendedTimeStamp.Write(writer)

	if _, err := writer.Write(chunkData); err != nil {
		logrus.Errorf("[Error] %v occured during writing chunk with chunk stream ID %v", err, c.BasicHeader.CSID)
	}
	logrus.Debug("[Debug] finished writing chunk")
}