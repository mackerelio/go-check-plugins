// +build windows

package eventlog

//go:generate go run $GOROOT/src/syscall/mksyscall_windows.go -output zsyscall_windows.go syscall_windows.go

//sys   ClearEventLog(eventLog syscall.Handle, backupFileName *uint16) (err error) = advapi32.ClearEventLogW
//sys   CloseEventLog(eventLog syscall.Handle) (err error) = advapi32.CloseEventLog
//sys   FormatMessage(flags uint32, source syscall.Handle, messageID uint32, languageID uint32, buffer *byte, bufferSize uint32, arguments uintptr) (numChars uint32, err error) = kernel32.FormatMessageW
//sys   GetNumberOfEventLogRecords(eventLog syscall.Handle, numberOfRecords *uint32) (err error) = advapi32.GetNumberOfEventLogRecords
//sys   GetOldestEventLogRecord(eventLog syscall.Handle, oldestRecord *uint32) (err error) = advapi32.GetOldestEventLogRecord
//sys   LoadLibraryEx(filename *uint16, file syscall.Handle, flags uint32) (handle syscall.Handle, err error) = kernel32.LoadLibraryExW
//sys   OpenEventLog(uncServerName *uint16, sourceName *uint16) (handle syscall.Handle, err error) = advapi32.OpenEventLogW
//sys   ReadEventLog(eventLog syscall.Handle, readFlags ReadFlag, recordOffset uint32, buffer *byte, numberOfBytesToRead uint32, bytesRead *uint32, minNumberOfBytesNeeded *uint32) (err error) = advapi32.ReadEventLogW

type EVENTLOGRECORD struct {
	Length              uint32
	Reserved            uint32
	RecordNumber        uint32
	TimeGenerated       uint32
	TimeWritten         uint32
	EventID             uint32
	EventType           uint16
	NumStrings          uint16
	EventCategory       uint16
	ReservedFlags       uint16
	ClosingRecordNumber uint32
	StringOffset        uint32
	UserSidLength       uint32
	UserSidOffset       uint32
	DataLength          uint32
	DataOffset          uint32
}

type ReadFlag uint32

const (
	EVENTLOG_SEQUENTIAL_READ ReadFlag = 1 << iota
	EVENTLOG_SEEK_READ
	EVENTLOG_FORWARDS_READ
	EVENTLOG_BACKWARDS_READ
)

type EventType uint16

const (
	EVENTLOG_SUCCESS    EventType = 0
	EVENTLOG_ERROR_TYPE           = 1 << (iota - 1)
	EVENTLOG_WARNING_TYPE
	EVENTLOG_INFORMATION_TYPE
	EVENTLOG_AUDIT_SUCCESS
	EVENTLOG_AUDIT_FAILURE
)

func (et EventType) String() string {
	switch et {
	case EVENTLOG_SUCCESS:
		return "Success"
	case EVENTLOG_ERROR_TYPE:
		return "Error"
	case EVENTLOG_AUDIT_FAILURE:
		return "Audit Failure"
	case EVENTLOG_AUDIT_SUCCESS:
		return "Audit Success"
	case EVENTLOG_INFORMATION_TYPE:
		return "Information"
	case EVENTLOG_WARNING_TYPE:
		return "Warning"
	default:
		return "Unknown"
	}
}

type SIDType uint32

const (
	SidTypeUser SIDType = 1 + iota
	SidTypeGroup
	SidTypeDomain
	SidTypeAlias
	SidTypeWellKnownGroup
	SidTypeDeletedAccount
	SidTypeInvalid
	SidTypeUnknown
	SidTypeComputer
	SidTypeLabel
)

func (st SIDType) String() string {
	switch st {
	case SidTypeUser:
		return "User"
	case SidTypeGroup:
		return "Group"
	case SidTypeDomain:
		return "Domain"
	case SidTypeAlias:
		return "Alias"
	case SidTypeWellKnownGroup:
		return "Well Known Group"
	case SidTypeDeletedAccount:
		return "Deleted Account"
	case SidTypeInvalid:
		return "Invalid"
	case SidTypeUnknown:
		return "Unknown"
	case SidTypeComputer:
		return "Unknown"
	case SidTypeLabel:
		return "Label"
	default:
		return "Unknown"
	}
}

const (
	DONT_RESOLVE_DLL_REFERENCES         uint32 = 0x0001
	LOAD_LIBRARY_AS_DATAFILE            uint32 = 0x0002
	LOAD_WITH_ALTERED_SEARCH_PATH       uint32 = 0x0008
	LOAD_IGNORE_CODE_AUTHZ_LEVEL        uint32 = 0x0010
	LOAD_LIBRARY_AS_IMAGE_RESOURCE      uint32 = 0x0020
	LOAD_LIBRARY_AS_DATAFILE_EXCLUSIVE  uint32 = 0x0040
	LOAD_LIBRARY_SEARCH_DLL_LOAD_DIR    uint32 = 0x0100
	LOAD_LIBRARY_SEARCH_APPLICATION_DIR uint32 = 0x0200
	LOAD_LIBRARY_SEARCH_USER_DIRS       uint32 = 0x0400
	LOAD_LIBRARY_SEARCH_SYSTEM32        uint32 = 0x0800
	LOAD_LIBRARY_SEARCH_DEFAULT_DIRS    uint32 = 0x1000
)
