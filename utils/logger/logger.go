package logger

import (
	"bullshape/utils"

	"github.com/sirupsen/logrus"
)

type requestLog struct {
	requestID string
}

func NewLogger(reqID string) logrus.FieldLogger {
	if len(reqID) == 0 {
		reqID = utils.NewUUIDV4()
	}
	log := logrus.StandardLogger()
	log.SetFormatter(&logrus.TextFormatter{})
	log.SetReportCaller(true)

	var retLogger logrus.FieldLogger = log
	return retLogger.WithField("req ID:", reqID)

}
