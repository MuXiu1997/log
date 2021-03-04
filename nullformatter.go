package log

import "github.com/sirupsen/logrus"

type NullFormatter struct {
}

func (n *NullFormatter) Format(_ *logrus.Entry) ([]byte, error) {
	return nil, nil
}
