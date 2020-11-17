# Connector

Пакет для открытия SSH и SFTP соединений.

Пример:

```
package main

import (
	"fmt"
	"log"
	connector "github.com/NuclearLouse/utilities-connectSSH_SFTP"
)

func main() {
	c := &connector.Credentials{
		Host:           "111.222.333.444",
		AuthMethod:     "key",
		User:           "barmaley",
		PrivateKeyFile: "C:\\Users\\barmaley\\.ssh\\my_server\\private.key",
		TimeOut:        5,
	}

	connect, err := connector.NewSSH(c)
	if err != nil {
		log.Fatalln("ssh client:", err)
	}
	defer connect.client.Close()

	client, err := connect.ClientSFTP()
	if err != nil {
		log.Fatalln("sftp client:", err)
	}
	defer client.Close()

	dirs, err := client.ReadDir("/home/barmaley")
	if err != nil {
		log.Fatalln("read home dir:", err)
	}
	for _, dir := range dirs {
		fmt.Printf("%-25s\tdir:%t\n", dir.Name(), dir.IsDir())
	}
}

```
