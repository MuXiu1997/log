package log

import "github.com/sirupsen/logrus"

var _ logrus.Formatter = &NullFormatter{}

type NullFormatter struct {
}

func (n *NullFormatter) Format(_ *logrus.Entry) ([]byte, error) {
	return nil, nil
}
