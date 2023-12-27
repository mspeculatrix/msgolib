/*
Package fileutils
Library: msgolib
Offered up under GPL 3.0 but absolutely not guaranteed fit for use.
This is code created by an amateur dilettante, so use at your own risk.
Github: https://github.com/mspeculatrix
Blog: https://mansfield-devine.com/speculatrix/
*/

package fileutils

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

/******************************************************************************
 *****   CONFIG FILES                                                     *****
 ******************************************************************************/

// ReadConfigFile - reads a configuration file where settings are stored as
// key/value pairs separated by '='. Returns a map.
// Param filepath should be a relative or absolute path + filename to the
// config file.
// Lines starting with #, ; or / are ignored as comments.
func ReadConfigFile(filepath string) (data map[string]string, err error) {
	data = make(map[string]string)
	fh, err := os.Open(filepath)

	if err != nil {
		return data, fmt.Errorf("readcfgfile : %v", err)
	}
	defer fh.Close()                // make sure the file is closed whatever
	scanner := bufio.NewScanner(fh) // to read line-by-line
	for scanner.Scan() {            // iterate over lines in file
		// Not sure if the TrimSpace() is necessary, but let's be cautious.
		line := strings.TrimSpace(scanner.Text()) // get next line
		// Ignore blank lines and comments, otherwise...
		if len(line) > 0 && !isComment(line) {
			items := strings.Split(line, "=")
			key := strings.TrimSpace(items[0])
			switch len(items) {
			case 1:
				// for some reason there was only a key, no value
				data[key] = ""
			case 2:
				// this is what we're expecting
				data[key] = strings.TrimSpace(items[1])
			default:
				// the value may itself have contained one or more '='. Use
				// the first item as the key and stitch back together
				// the rest as the value with '=' reinstated.
				data[key] = strings.TrimSpace(strings.Join(items[1:], "="))
			}
		}
	}
	return data, err
}

// WriteConfigFile - writes a map to a file in k=v format.
// A timestamp entry is added automatically.
func WriteConfigFile(filepath string, data map[string]string) (lineCount int, err error) {
	fh, err := os.Create(filepath)
	if err != nil {
		return lineCount, err
	}
	defer fh.Sync()
	defer fh.Close()
	_, err = fh.WriteString("timestamp=" + FileTimestamp() + "\n")
	if err != nil {
		return lineCount, fmt.Errorf("error writing string : %v", err)
	}
	lineCount = 1
	for k, v := range data {
		_, err = fh.WriteString(k + "=" + v + "\n")
		if err != nil {
			return lineCount, fmt.Errorf("writing data string : %v", err)
		}
		lineCount++
	}
	return lineCount, err
}

/******************************************************************************
 *****   PID FILES                                                        *****
 ******************************************************************************/

// ReadPIDFile reads PID value from a file as a string
func ReadPIDFile(filepath string) (pidStr string, err error) {
	pidStr = ""
	if FileExists(filepath) {
		dat, err := ioutil.ReadFile(filepath)
		if err != nil {
			return pidStr, err
		}
		pidStr = string(dat)
		pidStr = strings.TrimSpace(pidStr)
	}
	return pidStr, err
}

// WritePIDToFile writes PID of current program to file.
// Returns string version of that number.
func WritePIDToFile(filepath string) (string, error) {
	fh, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer fh.Close()
	defer fh.Sync()
	pidStr := strconv.Itoa(os.Getpid())
	_, err = fh.WriteString(pidStr)
	return pidStr, err
}

/******************************************************************************
 *****   LOG FILES                                                        *****
 ******************************************************************************/

// WriteToLogFile writes to a a simple log file. Adds a given line of text to
// the file, creating the file if necessary.
func WriteToLogFile(filepath string, logdata string, addTimestamp bool) error {
	fh, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		return err
	}
	defer fh.Close()
	defer fh.Sync()
	if addTimestamp {
		_, err = fh.Write([]byte(FileTimestamp() + " "))
		if err != nil {
			return err
		}
	}
	_, err = fh.Write([]byte(logdata + "\n"))
	return err
}

/******************************************************************************
 *****   MISC                                                             *****
 ******************************************************************************/

// FileExists - checks if a file exists and isn't a dir.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// FileTimestamp - returns a string suitable for timestamping files.
func FileTimestamp() string {
	now := time.Now()
	tsFmt := "%d-%02d-%02d %02d:%02d:%02d"
	return fmt.Sprintf(tsFmt, now.Year(), now.Month(), now.Day(), now.Hour(),
		now.Minute(), now.Second())
}

// isComment() checks to see if a supplied string (assumed to be a line from a
// file) starts with a character that would qualify it as a comment line.
func isComment(ln string) bool {
	comment := false
	commentChars := []string{"#", ";", "/"}
	for _, testchar := range commentChars {
		if ln[0:1] == testchar {
			comment = true
		}
	}
	return comment
}
