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

	log.WithFields(logrus.Fields{
		"omg":    true,
		"number": 122,
	}).Warn("The group's number increased tremendously!")

	assert.Equal(t, buf.String()[206:], `400,"labels":{"number":"122","omg":"true"},"textPayload":"The group's number increased tremendously!"}`)
}
