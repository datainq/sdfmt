package sdfmt

import (
	"bytes"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/sirupsen/logrus"
)

func TestStackdriverFormatter_Format(t *testing.T) {
	log := logrus.New()
	log.SetFormatter(&StackdriverFormatter{})

	var buf bytes.Buffer
	log.SetOutput(&buf)

	log.WithFields(logrus.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")
	assert.Equal(t, buf.String()[:123], `{"labels":{"animal":"walrus","size":"10"},"message":"A group of walrus emerges from the ocean","severity":300,"timestamp":"`)

	buf.Reset()
	log.WithFields(logrus.Fields{
		"omg":    true,
		"number": 122,
	}).Warn("The group's number increased tremendously!")
	assert.Equal(t, buf.String()[:123], `{"labels":{"number":"122","omg":"true"},"message":"The group's number increased tremendously!","severity":400,"timestamp":"`)
}
