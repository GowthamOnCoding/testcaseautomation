package cmd

import (
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"

	"os"
)

func DownloadFileFromLinuxServer(server string, user string, password string, remoteFilepath string, localFilepath string) error {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", server, config)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := sftp.NewClient(conn)
	if err != nil {
		return err
	}
	defer client.Close()

	srcFile, err := client.Open(remoteFilepath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(localFilepath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = srcFile.WriteTo(dstFile)
	if err != nil {
		return err
	}
	return nil
}
