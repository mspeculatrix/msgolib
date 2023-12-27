package smartparallel

import "github.com/tarm/serial"

const (
	// Terminator : byte used to terminate sent messages
	Terminator = 0
	// SerialCommandChar : ASCII command code - to preceed command bytes
	SerialCommandChar = 1
	// CmdPing : used to check if SmartParallel is alive and connected
	CmdPing = 1
	// CmdAckDisable : disable use of ACK in printing
	CmdAckDisable = 2
	// CmdAckEnable : enable use of ACK in printing
	CmdAckEnable = 3
	// CmdAutofeedDisable : disable use of printer's AUTOFEED function
	CmdAutofeedDisable = 4
	// CmdAutoFeedEnable : enable use of printer's AUTOFEED function
	CmdAutoFeedEnable = 5
	// CmdPrtModeNormal : print in standard 80-column mde
	CmdPrtModeNormal = 8
	// CmdPrtModeCond : print in 132-column condensed mode
	CmdPrtModeCond = 9
	// CmdPrtModeDbl : print in 40-column double-width mode
	CmdPrtModeDbl = 10
	// CmdLineEndNormal : don't add anything to ends of lines (default)
	CmdLineEndNormal = 16
	// CmdLinefeedLF : add linefeeds (ASCII 10) to ends of lines
	CmdLinefeedLF = 17
	// CmdLinefeedCR : add carriage returns (ASCII 13) to ends of lines
	CmdLinefeedCR = 18
	// CmdLinefeedCRLF : add both LF and CR to ends of lines
	CmdLinefeedCRLF = 19
	// CmdReportState : request status report from SmartParallel
	CmdReportState = 32
	// CmdReportAck : check if use of ACK is enabled
	CmdReportAck = 33
	// CmdReportAutofeed : check if AUTOFEED is enabled
	CmdReportAutofeed = 34
	// ReadBufSize : Default size for read buffer
	ReadBufSize = 1024 // bytes
	// DefaultColumns : Default number of columns for printer
	DefaultColumns = 80 // because it's an Epson MX-80
)

var (
	// Init : byte sequence to initialise printer (for Epson)
	Init = []byte{27, 64} // ESC @
	// LineEnd : to send at end of each line of text
	LineEnd = []byte{13, 10} // CR and LF
	// TransmitEnd : to terminate each line
	TransmitEnd = []byte{Terminator} // as byte array for easy appending
	// SetTabs : code to set tab positions
	SetTabs = []byte{1, 64}
)

/*
CheckSerialInput : pull next nul-terminated string from serial port
*/
func CheckSerialInput(sPort *serial.Port, readBuf []byte) (int, string) {
	readBuf = readBuf[:0]  // empty out read buffer but retain it in memory
	buf := make([]byte, 1) // temp buffer for each read
	charIdx := 0
	done := false
	for !done {
		n, err := sPort.Read(buf) // read one char
		if n > 0 {
			// we've received something, even though it might not be all.
			// Note that, although we've received chars, err might be
			// non-nil as well, which is why we're testing for them
			// separately.
			// SmartParallel terminates outgoing serial messages with a
			// newline/linefeed - ASCII 10, hex 0x0A.
			if buf[0] == 10 {
				done = true
			} else {
				readBuf = append(readBuf, buf[0])
				charIdx++
				if charIdx == ReadBufSize {
					done = true
				}
			}
		}
		if err != nil {
			done = true
			// there can be an error, such as EOF, even when n is not 0
			if err.Error() == "EOF" {
				//verbosePrintln("!EOF")
			} else {
				//verbosePrintln("Error read input", err.Error())
			}
		}
	}
	return charIdx, string(readBuf[:charIdx])
}
