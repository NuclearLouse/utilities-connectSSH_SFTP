package connector

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// Credentials ...AuthMethod : "key", "password", "keyboard". Port default ":22"
type Credentials struct {
	Host           string
	Port           string
	AuthMethod     string
	User           string
	Password       string
	PrivateKeyFile string
	TimeOut        int64
}

// SSH ...
type SSH struct {
	Client *ssh.Client
}

// NewSSH ...
func NewSSH(c *Credentials) (*SSH, error) {

	cfg := &ssh.ClientConfig{

		User:            c.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Duration(c.TimeOut) * time.Second,
	}

	switch c.AuthMethod {
	case "key":
		privateKey, err := ioutil.ReadFile(c.PrivateKeyFile)
		if err != nil {
			return nil, err
		}
		signer, err := ssh.ParsePrivateKey(privateKey)
		if err != nil {
			return nil, err
		}
		cfg.Auth = []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		}
	case "password":
		cfg.Auth = []ssh.AuthMethod{
			ssh.Password(c.Password),
		}
	case "keyboard":
		cfg.Auth = []ssh.AuthMethod{
			ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) ([]string, error) {
				// Just sends the password back for all questions
				answers := make([]string, len(questions))
				for i := range answers {
					answers[i] = c.Password
				}
				return answers, nil
			}),
		}
	default:
		err := errors.New("unsupported authentication method")
		return nil, err
	}

	if c.Port == "" {
		c.Port = "22"
	}

	client, err := ssh.Dial("tcp", c.Host+":"+c.Port, cfg)
	if err != nil {
		return nil, err
	}

	return &SSH{Client: client}, nil
}

// ClientSFTP ...
func (s *SSH) ClientSFTP() (*sftp.Client, error) {
	client, err := sftp.NewClient(s.Client)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// NewSession ...
func (s *SSH) NewSession() (*ssh.Session, error) {
	session, err := s.Client.NewSession()
	if err != nil {
		return nil, err
	}

	return session, nil
}

// RunCommand ...
// TODO: дописать выполнение команд с возвратом результатов
func RunCommand(session *ssh.Session, cmd string) error {
	sessStdOut, err := session.StdoutPipe()
	if err != nil {
		return err
	}
	go io.Copy(os.Stdout, sessStdOut)
	sessStderr, err := session.StderrPipe()
	if err != nil {
		return err
	}
	go io.Copy(os.Stderr, sessStderr)
	err = session.Run(cmd) // eg., /usr/bin/whoami
	if err != nil {
		return err
	}
	// var stdoutBuf bytes.Buffer
	// session.Stdout = &stdoutBuf
	// session.Run(cmd)

	// return hostname + ": " + stdoutBuf.String()

	return nil
}
