package helpers

import "fmt"

func SyncAccount(accountName string) error {
	commands := [][]string{
		[]string{"sudo", "adduser", "--disabled-password", "--gecos", "\"\"", accountName},
		[]string{"sudo", "chmod", "750", fmt.Sprintf("/home/%s", accountName)},
		[]string{"sudo", "adduser", "nginx", accountName},
		[]string{"sudo", "adduser", "apache", accountName},

		// rforce nologin as shell
		[]string{"sudo", "chsh", "-s", "/usr/sbin/nologin", accountName},
	}

	dirs := []string{"apps", "domains", ".ssh"}
	for _, dir := range dirs {
		commands = append(commands, []string{"sudo", "mkdir", "-p", fmt.Sprintf("/home/%s/%s", accountName, dir)})
		commands = append(commands, []string{"sudo", "chmod", "750", fmt.Sprintf("/home/%s/%s", accountName, dir)})
		commands = append(commands, []string{"sudo", "chown", fmt.Sprintf("%s:%s", accountName, accountName), fmt.Sprintf("/home/%s/%s", accountName, dir)})
	}

	for _, command := range commands {
		_, errOut, err := ExecuteShellCommand(command[0], command[1:])
		if err != nil {
			return fmt.Errorf("error %s (%s) while executing %s", errOut, err, command)
		}
	}

	// remove authorized_keys and ignore output
	ExecuteShellCommand("sudo", []string{"rm", fmt.Sprintf("/home/%s/.ssh/authorized_keys", accountName)})

	return nil
}
