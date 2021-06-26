package server

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/Zioyi/zedis/sdb"
	"net"
	"strings"
)


func exce(db sdb.Db, cmdString string) (interface{}, error) {
	subs := strings.Split(cmdString, "\r\n")
	validSubs := make([]string, 0)
	for i := 2; i < len(subs); i += 2 {
		// validation
		validSubs = append(validSubs, subs[i])
	}

	cmd := validSubs[0]
	switch cmd {
	case "set":
		fmt.Printf("handle set, %v\n", validSubs)
		err := db.Set(validSubs[1], validSubs[2])
		if err != nil {
			fmt.Printf("set data err, %v\n", err)
		}
		return nil, nil
	case "get":
		fmt.Printf("handle get, %v\n", validSubs)
		val, err := db.Get(validSubs[1])
		if err != nil {
			fmt.Printf("get data err, %v\n", err)
			return nil, err
		}
		return val, nil
	}

	return nil, errors.New("invalid cmd")

}

func process(db sdb.Db,conn net.Conn) {
	defer conn.Close()

	for {
		reader := bufio.NewReader(conn)
		var buf [128]byte
		n, err := reader.Read(buf[:])
		if err != nil {
			fmt.Printf("read from conn failed, err:%v\n", err)
			break
		}

		recv := string(buf[:n])
		fmt.Printf("收到的数据：\n-----start-----\n%v\n-----end-----\n", recv)
		res, err := exce(db, recv)
		fmt.Println(res, err)
		if err != nil {
			_, err = conn.Write([]byte(fmt.Sprintf("-%s\r\n", err)))
			break
		}
		if res == nil {
			_, err = conn.Write([]byte("+OK\r\n"))
		} else {
			s := res.(string)
			fmt.Printf("got %s\n", s)
			_, err = conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(s), s)))
		}


		if err != nil {
			fmt.Printf("write from conn falied, err:%v\n", err)
			break
		}
	}
}

func Run() {
	mDb := sdb.NewMemoryDb()
	var db sdb.Db
	db = mDb

	listen, err := net.Listen("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Printf("listen failed, err: %v\n", err)
		return
	}
	fmt.Printf("Fake Redis Sever Running\n")
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("accept failed, err:%v\n", err)
			return
		}
		process(db, conn)
	}
}