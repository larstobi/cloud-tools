package config

import (
	"bufio"
	"os"
	"os/exec"
)

// Will lookup key in pass - the linux password store
// and ask for GPG password
//
func GetPasswordFor(key string) string {

	cmd := exec.Command("pass", key)
	//cmd.Env = ()
	// Ask for gpg password if necessary
	cmd.Stdin = os.Stdin

	// capture the output and error pipes
	stdout, _ := cmd.StdoutPipe()
	//stderr, _ := cmd.StderrPipe()

	cmd.Start()

	// Don't let main() exit before our command has finished running
	// doesn't block
	defer cmd.Wait()

	buff := bufio.NewScanner(stdout)
	var password string

	for buff.Scan() {
		password += buff.Text()
	}

	return password

}
