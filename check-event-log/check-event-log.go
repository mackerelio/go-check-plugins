package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows/registry"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
	"github.com/mackerelio/go-check-plugins/check-event-log/internal/eventlog"
)

type logOpts struct {
	Log            string `long:"log" description:"Event Names (comma separated)"`
	ID             string `long:"id" description:"Event IDs (comma separated)"`
	Type           string `long:"type" description:"Event Types (comma separated)"`
	SourcePattern  string `long:"source-pattern" description:"Event Source (regexp pattern)"`
	MessagePattern string `long:"message-pattern" description:"Message Pattern (regexp pattern)"`
	ReturnContent  bool   `short:"r" long:"return" description:"Return matched line"`
	StateDir       string `short:"s" long:"state-dir" default:"/var/mackerel-cache/check-event-log" value-name:"DIR" description:"Dir to keep state files under"`
	NoState        bool   `long:"no-state" description:"Don't use state file and read whole logs"`
	FailFirst      bool   `long:"fail-first" description:"Count errors on first seek"`
	Verbose        bool   `long:"verbose" description:"Verbose output"`

	logList        []string
	idList         []int64
	typeList       []string
	sourcePattern  *regexp.Regexp
	messagePattern *regexp.Regexp
}

func stringList(s string) []string {
	l := strings.Split(s, ",")
	if len(l) == 0 || l[0] == "" {
		return []string{}
	}
	return l
}

func (opts *logOpts) prepare() error {
	opts.logList = stringList(opts.Log)
	if len(opts.logList) == 0 || opts.logList[0] == "" {
		opts.logList = []string{"Application"}
	}
	for _, id := range stringList(opts.ID) {
		negate := int64(1)
		if id != "" && id[0] == '!' {
			negate = -1
			id = id[1:]
		}
		i, err := strconv.Atoi(id)
		if err != nil {
			return err
		}
		opts.idList = append(opts.idList, int64(i)*negate)
	}
	opts.typeList = stringList(opts.Type)

	var err error
	opts.sourcePattern, err = regexp.Compile(opts.SourcePattern)
	if err != nil {
		return err
	}
	opts.messagePattern, err = regexp.Compile(opts.MessagePattern)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	ckr := run(os.Args[1:])
	ckr.Name = "Event Log"
	ckr.Exit()
}

func parseArgs(args []string) (*logOpts, error) {
	opts := &logOpts{}
	_, err := flags.ParseArgs(opts, args)
	return opts, err
}

func run(args []string) *checkers.Checker {
	opts, err := parseArgs(args)
	if err != nil {
		os.Exit(1)
	}

	err = opts.prepare()
	if err != nil {
		return checkers.Unknown(err.Error())
	}

	checkSt := checkers.OK
	warnNum := int64(0)
	critNum := int64(0)
	errorOverall := ""

	for _, f := range opts.logList {
		w, c, errLines, err := opts.searchLog(f)
		if err != nil {
			return checkers.Unknown(err.Error())
		}
		warnNum += w
		critNum += c
		if opts.ReturnContent {
			errorOverall += errLines
		}
	}
	msg := fmt.Sprintf("%d warnings, %d criticals.", warnNum, critNum)
	return checkers.NewChecker(checkSt, msg)
}

func bytesToString(b []byte) (string, uint32) {
	var i int
	s := make([]uint16, len(b)/2)
	for i = range s {
		s[i] = uint16(b[i*2]) + uint16(b[(i*2)+1])<<8
		if s[i] == 0 {
			s = s[0:i]
			break
		}
	}
	return string(utf16.Decode(s)), uint32(i * 2)
}

func getResourceMessage(providerName, sourceName string, eventID uint32, argsptr uintptr) (string, error) {
	regkey := fmt.Sprintf(
		"SYSTEM\\CurrentControlSet\\Services\\EventLog\\%s\\%s",
		providerName, sourceName)
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, regkey, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer key.Close()

	val, _, err := key.GetStringValue("EventMessageFile")
	if err != nil {
		return "", err
	}
	val, err = registry.ExpandString(val)
	if err != nil {
		return "", err
	}

	handle, err := eventlog.LoadLibraryEx(syscall.StringToUTF16Ptr(val), 0,
		eventlog.DONT_RESOLVE_DLL_REFERENCES|eventlog.LOAD_LIBRARY_AS_DATAFILE)
	if err != nil {
		return "", err
	}
	defer syscall.CloseHandle(handle)

	msgbuf := make([]byte, 1<<16)
	numChars, err := eventlog.FormatMessage(
		syscall.FORMAT_MESSAGE_FROM_SYSTEM|
			syscall.FORMAT_MESSAGE_FROM_HMODULE|
			syscall.FORMAT_MESSAGE_ARGUMENT_ARRAY,
		handle,
		eventID,
		0,
		&msgbuf[0],
		uint32(len(msgbuf)),
		argsptr)
	if err != nil {
		return "", err
	}
	message, _ := bytesToString(msgbuf[:numChars*2])
	message = strings.Replace(message, "\r", "", -1)
	message = strings.TrimSuffix(message, "\n")
	return message, nil
}

func (opts *logOpts) searchLog(eventName string) (warnNum, critNum int64, errLines string, err error) {
	stateFile := getStateFile(opts.StateDir, eventName)
	recordNumber := int64(0)
	if !opts.NoState {
		s, err := getLastOffset(stateFile)
		if err != nil {
			return 0, 0, "", err
		}
		recordNumber = s
	}

	ptr := syscall.StringToUTF16Ptr(eventName)
	h, err := eventlog.OpenEventLog(nil, ptr)
	if err != nil {
		log.Fatal(err)
	}
	defer eventlog.CloseEventLog(h)

	var num, oldnum uint32

	eventlog.GetNumberOfEventLogRecords(h, &num)
	if err != nil {
		log.Fatal(err)
	}
	eventlog.GetOldestEventLogRecord(h, &oldnum)
	if err != nil {
		log.Fatal(err)
	}

	flags := eventlog.EVENTLOG_FORWARDS_READ | eventlog.EVENTLOG_SEEK_READ

	if recordNumber > 0 && oldnum <= uint32(recordNumber) {
		if uint32(recordNumber)+1 == oldnum+num {
			return 0, 0, "", nil
		}
		recordNumber++
	} else {
		recordNumber = 1
	}

	size := uint32(1)
	buf := []byte{0}

	var readBytes uint32
	var nextSize uint32
	var lastNumber uint32
	for i := uint32(recordNumber); i < oldnum+num; i++ {
		err = eventlog.ReadEventLog(
			h,
			flags,
			i,
			&buf[0],
			size,
			&readBytes,
			&nextSize)
		if err != nil {
			if err != syscall.ERROR_INSUFFICIENT_BUFFER {
				break
			}
			buf = make([]byte, nextSize)
			size = nextSize
			err = eventlog.ReadEventLog(
				h,
				eventlog.EVENTLOG_FORWARDS_READ|eventlog.EVENTLOG_SEQUENTIAL_READ,
				i,
				&buf[0],
				size,
				&readBytes,
				&nextSize)
			if err != nil {
				log.Printf("eventlog.ReadEventLog: %v", err)
				break
			}
		}

		r := *(*eventlog.EVENTLOGRECORD)(unsafe.Pointer(&buf[0]))
		if opts.Verbose {
			log.Printf("RecordNumber=%v", r.RecordNumber)
			log.Printf("TimeGenerated=%v", time.Unix(int64(r.TimeGenerated), 0).String())
			log.Printf("TimeWritten=%v", time.Unix(int64(r.TimeWritten), 0).String())
			log.Printf("EventID=%v", r.EventID)
		}
		lastNumber = i

		if len(opts.idList) > 0 {
			found := false
			for _, id := range opts.idList {
				if id > 0 && uint32(id) == r.EventID {
					found = true
					break
				} else if id <= 0 && uint32(-id) != r.EventID {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		tn := eventlog.EventType(r.EventType).String()
		if opts.Verbose {
			log.Printf("EventType=%v", tn)
		}
		tn = strings.ToLower(tn)
		if len(opts.typeList) > 0 {
			found := false
			for _, typ := range opts.typeList {
				if typ == tn {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		sourceName, sourceNameOff := bytesToString(buf[unsafe.Sizeof(eventlog.EVENTLOGRECORD{}):])
		computerName, _ := bytesToString(buf[unsafe.Sizeof(eventlog.EVENTLOGRECORD{})+uintptr(sourceNameOff+2):])
		if opts.Verbose {
			log.Printf("SourceName=%v", sourceName)
			log.Printf("ComputerName=%v", computerName)
		}

		if opts.sourcePattern != nil {
			if !opts.sourcePattern.MatchString(sourceName) {
				continue
			}
		}

		off := uint32(0)
		args := make([]*byte, uintptr(r.NumStrings)*unsafe.Sizeof((*uint16)(nil)))
		for n := 0; n < int(r.NumStrings); n++ {
			args[n] = &buf[r.StringOffset+off]
			_, boff := bytesToString(buf[r.StringOffset+off:])
			off += boff + 2
		}

		var argsptr uintptr
		if r.NumStrings > 0 {
			argsptr = uintptr(unsafe.Pointer(&args[0]))
		}
		message, err := getResourceMessage(eventName, sourceName, r.EventID, argsptr)
		if err == nil {
			if opts.Verbose {
				log.Printf("Message=%v", message)
			}
			if opts.messagePattern != nil {
				if !opts.messagePattern.MatchString(message) {
					continue
				}
			}
		}
		if opts.ReturnContent {
			errLines += sourceName + ":" + strings.Replace(message, "\n", "", -1) + "\n"
		}
		switch tn {
		case "error":
			critNum++
		case "audit failure":
			critNum++
		case "warning":
			warnNum++
		}
	}

	if !opts.NoState {
		err = writeLastOffset(stateFile, int64(lastNumber))
		if err != nil {
			log.Printf("writeLastOffset failed: %s\n", err.Error())
		}
	}

	if recordNumber == 1 && !opts.FailFirst {
		if strings.Contains(errLines, "テストエラー") {
			println("================================================")
		}
		return 0, 0, "", nil
	}
	return warnNum, critNum, errLines, nil
}

var stateRe = regexp.MustCompile(`^([A-Z]):[/\\]`)

func getStateFile(stateDir, f string) string {
	return filepath.ToSlash(filepath.Join(stateDir, stateRe.ReplaceAllString(f, `$1`+string(filepath.Separator))))
}

func getLastOffset(f string) (int64, error) {
	_, err := os.Stat(f)
	if err != nil {
		return 0, nil
	}
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return 0, err
	}
	i, err := strconv.ParseInt(strings.Trim(string(b), " \r\n"), 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func writeLastOffset(f string, num int64) error {
	err := os.MkdirAll(filepath.Dir(f), 0755)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(f, []byte(fmt.Sprintf("%d", num)), 0644)
}
