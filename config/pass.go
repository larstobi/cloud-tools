package config

import (
	"bufio"
	"os"
	"os/exec"
	"fmt"
	"io"
	"strconv"
	"regexp"
	"bytes"
)

// GetPasswordFor will lookup key in pass - the linux password store
// and ask for GPG password
//
func GetPasswordFor(key string) string {

	cmd := exec.Command("pass", key)

	// Ask for gpg password if necessary
	cmd.Stdin = os.Stdin

	// capture the output and error pipes
	stdout, _ := cmd.StdoutPipe()
	//stderr, _ := cmd.StderrPipe()

	cmd.Start()
	defer cmd.Wait()

	buff := bufio.NewScanner(stdout)
	var password string

	for buff.Scan() {
		password += buff.Text()
	}

	return password

}

// GeneratePasswordFor will generate password of given length in given password storage dir
func GeneratePasswordFor(passwordStorageDirectory string, passName string, passLength int) (string) {

	environment := getIsolatedPassEnvironment(passwordStorageDirectory)

	cmd := exec.Command("pass", "generate", "-f", "-n", passName, strconv.Itoa(passLength))

	cmd.Env = environment
	// Enable when debugging
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr

	output, _ := cmd.Output()

	outputString := string(output)

	// After removing removing the ctrl bytes, this is all that's left of the colour codes in the string
	prefix := regexp.QuoteMeta("[1m[37mThe generated password for [4m" + passName + "[24m is:[0m[1m[93m")
	suffix := regexp.QuoteMeta("[0m")
	r := regexp.MustCompile(prefix + "(.*)" + suffix)

	return r.FindStringSubmatch(stripCtlAndExtFromBytes(outputString))[1]

}

func stripCtlAndExtFromBytes(str string) string {
	b := make([]byte, len(str))
	var bl int
	for i := 0; i < len(str); i++ {
		c := str[i]
		if c >= 32 && c < 127 {
			b[bl] = c
			bl++
		}
	}
	return string(b[:bl])
}

// InsertPasswordFor will insert password pass store in given storage dir
func InsertPasswordFor(passwordStorageDirectory string, passName string, password string) {

	environment := getIsolatedPassEnvironment(passwordStorageDirectory)

	cmd := exec.Command("pass", "insert", "-e", passName)
	cmd.Env = environment
	stdin , _ := cmd.StdinPipe()
	cmd.Start()
	defer cmd.Wait()
	io.Copy(stdin, bytes.NewBufferString(password))
	defer stdin.Close()
}

func getIsolatedPassEnvironment(passwordStorageDirectory string) ([]string) {
	var environment []string
	environment = append(environment,
		fmt.Sprintf("PATH=%s", os.Getenv("PATH")),
		fmt.Sprintf("HOME=%s", os.Getenv("HOME")),
		fmt.Sprintf("PASSWORD_STORE_DIR=%s", passwordStorageDirectory),
		fmt.Sprintf("PASSWORD_STORE_GIT=%s", passwordStorageDirectory))
	return environment
}