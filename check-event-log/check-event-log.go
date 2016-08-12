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
	Log           string `long:"log" description:"Event names (comma separated)"`
	Code          string `long:"code" description:"Event codes (comma separated)"`
	Type          string `long:"type" description:"Event types (comma separated)"`
	Source        string `long:"source" description:"Event source (comma separated)"`
	ReturnContent bool   `short:"r" long:"return" description:"Return matched line"`
	StateDir      string `short:"s" long:"state-dir" default:"/var/mackerel-cache/check-log" value-name:"DIR" description:"Dir to keep state files under"`
	NoState       bool   `long:"no-state" description:"Don't use state file and read whole logs"`
	Verbose       bool   `long:"verbose" description:"Verbose output"`

	logList    []string
	codeList   []string
	typeList   []string
	sourceList []string
}

func (opts *logOpts) prepare() error {
	opts.logList = strings.Split(opts.Log, ",")
	opts.codeList = strings.Split(opts.Code, ",")
	opts.typeList = strings.Split(opts.Type, ",")
	opts.sourceList = strings.Split(opts.Source, ",")
	return nil
}

func main() {
	ckr := run(os.Args[1:])
	ckr.Name = "Event Log"
	ckr.Exit()
}

func regCompileWithCase(ptn string, caseInsensitive bool) (*regexp.Regexp, error) {
	if caseInsensitive {
		ptn = "(?i)" + ptn
	}
	return regexp.Compile(ptn)
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

	for _, f := range opts.sourceList {
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
	msg := fmt.Sprintf("%d warnings, %d criticals for pattern /%s/.", 0, 0, "")
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
		log.Print(err)
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

func (opts *logOpts) searchLog(event string) (int64, int64, string, error) {
	stateFile := getStateFile(opts.StateDir, event)
	recordNumber := int64(0)
	if !opts.NoState {
		s, err := getBytesToSkip(stateFile)
		if err != nil {
			return 0, 0, "", err
		}
		recordNumber = s
	}

	providerName := event

	ptr := syscall.StringToUTF16Ptr(providerName)
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

	flags := eventlog.EVENTLOG_FORWARDS_READ | eventlog.EVENTLOG_SEQUENTIAL_READ

	if recordNumber > 0 && oldnum <= uint32(recordNumber) {
		if uint32(recordNumber)+1 == oldnum+num {
			return 0, 0, "", nil
		}
		flags = eventlog.EVENTLOG_FORWARDS_READ | eventlog.EVENTLOG_SEEK_READ
	} else {
		recordNumber = 0
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
				log.Println(err)
				break
			}
		}

		r := *(*eventlog.EVENTLOGRECORD)(unsafe.Pointer(&buf[0]))
		if opts.Verbose {
			println(r.RecordNumber)
			println(time.Unix(int64(r.TimeGenerated), 0).String())
			println(time.Unix(int64(r.TimeWritten), 0).String())
			println(r.EventID)
		}

		sourceName, sourceNameOff := bytesToString(buf[unsafe.Sizeof(eventlog.EVENTLOGRECORD{}):])
		computerName, _ := bytesToString(buf[unsafe.Sizeof(eventlog.EVENTLOGRECORD{})+uintptr(sourceNameOff+2):])
		if opts.Verbose {
			println(sourceName)
			println(computerName)
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
		message, err := getResourceMessage(providerName, sourceName, r.EventID, argsptr)
		if err == nil {
			if opts.Verbose {
				println(message)
			}
		}

		lastNumber = r.RecordNumber
	}

	if !opts.NoState {
		err = writeBytesToSkip(stateFile, int64(lastNumber))
		if err != nil {
			log.Printf("writeByteToSkip failed: %s\n", err.Error())
		}
	}

	return 0, 0, "", nil
}

var stateRe = regexp.MustCompile(`^([A-Z]):[/\\]`)

func getStateFile(stateDir, f string) string {
	return filepath.Join(stateDir, stateRe.ReplaceAllString(f, `$1`+string(filepath.Separator)))
}

func getBytesToSkip(f string) (int64, error) {
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

func writeBytesToSkip(f string, num int64) error {
	err := os.MkdirAll(filepath.Dir(f), 0755)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(f, []byte(fmt.Sprintf("%d", num)), 0644)
}
