/*
 * @Author: Feng
 * @version: v1.0.0
 * @Date: 2020-07-03 13:50:49
 * @LastEditors: Feng
 * @LastEditTime: 2020-09-03 08:30:19
 */
package ktcp

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/rs/zerolog/log"
)

type Ktcp struct {
}

var conns = map[string]net.Conn{}

//Send tcp message
func Send(ip, mes string) (res string, err error) {

	// log.Debug().Msg("Cmd is: " + mes)

	conn := conns[ip]

	if conn == nil {

		conn, err = reconnect(ip)

		if err != nil {
			return
		}

	}

	//send mes
	readData := make([]byte, 1024)

	hexData, _ := hex.DecodeString(mes)

	_, err = conn.Write(hexData)

	if err != nil {

		conn, err = reconnect(ip)

		if err != nil {
			return
		}
	}

	//listen for reply
	conn.SetReadDeadline(time.Now().Add(time.Second * 2))
	rint, err := conn.Read(readData) //error reconnect

	if err != nil {

		_, err1 := reconnect(ip)

		err = errors.New("读取命令超时 tcp timeout, reconnect success, " + err.Error() + "," + err1.Error())
		return

	}

	if rint > 0 {

		readData = readData[:rint]
	} else {

		err = errors.New("读取数据为空！")
		return
	}

	res = fmt.Sprintf("%x", readData)

	return
}

func reconnect(ip string) (conn net.Conn, err error) {

	conn, err = net.Dial("tcp", ip)

	if err != nil {

		return
	}

	conns[ip] = conn

	//send hello msg because the device will send a login mes please enter pwd ...
	tdata := make([]byte, 30)
	conn.Write([]byte("hi"))
	time.Sleep(time.Millisecond * 500)
	conn.Read(tdata)

	// fmt.Println("reconnect ok", ip)
	log.Info().Msgf("Reconnect to IP： %s success", ip)
	return
}
