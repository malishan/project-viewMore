/*
The library provides different types of logging such as infoLog, traceLog, warningLog, errorLog and fatalLog.
The logging can be done either on standard ouput or a specified logFile.
There is a provision of logging different log types differently based on its environment meant setup.
To specify a certain log type as file logging enable its value as 1 in env var.
*/

package loglib

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"project/project-viewMore/apicontext"
	"strings"
)

const (
	//INFO level 0
	INFO = iota
	//TRACE level 1
	TRACE
	//ERROR level 2
	ERROR
	//FATAL level 3
	FATAL
)

// FieldsMap map of key value pair to log
type FieldsMap map[string]interface{}

var (
	maxLogLevel = FATAL

	info     *log.Logger
	trace    *log.Logger
	errorlog *log.Logger
	fatal    *log.Logger
)

func init() {

	infoLog, traceLog, errorLog, fatalLog := os.Stdout, os.Stdout, os.Stderr, os.Stderr

	loginit(infoLog, traceLog, errorLog, fatalLog)
}

func loginit(infoHandle, traceHandle, errorHandle, fatalHandle io.Writer) {
	info = log.New(infoHandle, "INFO|", log.LUTC|log.LstdFlags|log.Lshortfile)

	trace = log.New(traceHandle, "TRACE|", log.LUTC|log.LstdFlags|log.Lshortfile)

	errorlog = log.New(errorHandle, "ERROR|", log.LUTC|log.LstdFlags|log.Lshortfile)

	fatal = log.New(fatalHandle, "FATAL|", log.LUTC|log.LstdFlags|log.Lshortfile)
}

func generatePrefix(ctx apicontext.CustomContext) string {
	return strings.Join([]string{ctx.UserID, ctx.UserName, ctx.Email}, ":")
}

func generateTrackingIDs(ctx apicontext.CustomContext) string {
	var retString string
	requestID := ctx.RequestID

	if requestID != "" {
		retString += "requestId=" + requestID
	}
	return retString
}

func doLog(cLog *log.Logger, level, callDepth int, v ...interface{}) {
	if level > maxLogLevel {
		cLog.SetOutput(os.Stderr)
	}

	cLog.Output(callDepth, fmt.Sprintln(v...))
}

//Info dedicated for logging valuable information
func infoLog(v ...interface{}) {
	doLog(info, INFO, 1, v...)
}

//Trace system gives facility to helps you isolate your system problems by monitoring selected events Ex. entry and exit
func traceLog(v ...interface{}) {
	doLog(trace, TRACE, 1, v...)
}

//Error logging error
func errorLog(v ...interface{}) {
	doLog(errorlog, ERROR, 1, v...)
}

//Fatal logging error
func fatalLog(v ...interface{}) {
	doLog(fatal, FATAL, 1, v...)
	os.Exit(1)
}

//GenericInfo generates info log
func GenericInfo(ctx apicontext.CustomContext, infoMessage string, fields FieldsMap) {
	prefix := generatePrefix(ctx)
	trackingIDs := generateTrackingIDs(ctx)
	fieldsBytes, _ := json.Marshal(fields)
	fieldsString := string(fieldsBytes)
	msg := fmt.Sprintf("|%s|%s|",
		prefix,
		trackingIDs)
	if fields != nil && len(fields) > 0 {
		infoLog(msg, infoMessage, "|", fieldsString)
	} else {
		infoLog(msg, infoMessage)
	}
}

//GenericTrace generates trace log
func GenericTrace(ctx apicontext.CustomContext, traceMessage string, fields FieldsMap) {
	prefix := generatePrefix(ctx)
	trackingIDs := generateTrackingIDs(ctx)
	msg := fmt.Sprintf("|%s|%s|",
		prefix,
		trackingIDs)
	if fields != nil && len(fields) > 0 {
		fieldsBytes, _ := json.Marshal(fields)
		fieldsString := string(fieldsBytes)
		traceLog(msg, traceMessage, "|", fieldsString)
	} else {
		traceLog(msg, traceMessage)
	}
}

//GenericError generates error log
func GenericError(ctx apicontext.CustomContext, e error, fields FieldsMap) {
	prefix := generatePrefix(ctx)
	trackingIDs := generateTrackingIDs(ctx)
	msg := ""
	if e != nil {
		msg = fmt.Sprintf("|%s|%s|%s", prefix, trackingIDs, e.Error())
	} else {
		msg = fmt.Sprintf("|%s|%s", prefix, trackingIDs)
	}

	if fields != nil && len(fields) > 0 {
		fieldsBytes, _ := json.Marshal(fields)
		fieldsString := string(fieldsBytes)
		errorLog(msg, "|", fieldsString)
	} else {
		errorLog(msg)
	}
}

//GenericFatalLog generates fatal log and then exits with os.Exit(1)
func GenericFatalLog(ctx apicontext.CustomContext, e error, fields FieldsMap) {
	prefix := generatePrefix(ctx)
	trackingIDs := generateTrackingIDs(ctx)
	msg := ""
	if e != nil {
		msg = fmt.Sprintf("|%s|%s|%s", prefix, trackingIDs, e.Error())
	} else {
		msg = fmt.Sprintf("|%s|%s", prefix, trackingIDs)
	}

	if fields != nil && len(fields) > 0 {
		fieldsBytes, _ := json.Marshal(fields)
		fieldsString := string(fieldsBytes)
		fatalLog(msg, "|", fieldsString)
	} else {
		fatalLog(msg)
	}
}
