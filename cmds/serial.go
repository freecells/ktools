/*
 * @Author: Feng
 * @version: v1.0.0
 * @Date: 2020-07-10 15:56:51
 * @LastEditors: Keven
 * @LastEditTime: 2021-09-28 13:36:37
 */
package cmds

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/snksoft/crc"
	serial "github.com/tarm/goserial"
)

//CRC16MODBUS  crc 校验参数
var CRC16MODBUS = &crc.Parameters{Width: 16, Polynomial: 0x8005, Init: 0xFFFF, ReflectIn: true, ReflectOut: true, FinalXor: 0x0}

/*SerialSend 发送 串口命令，返回数据
 * @params: cmd 命令字串，checkCrc 是否检查crc, crcType检查的crc 类型，crcLH 是否大端开头
 * @return:res 返回数据，err 错误
 */
func SerialSend(cmd, serName, serBaud string, checkCrc, crcLH bool, crcType *crc.Parameters) (res string, err error) {

	baud, _ := strconv.Atoi(serBaud)

	cfg := &serial.Config{
		Name:        serName,
		Baud:        baud,
		ReadTimeout: 3 /*毫秒*/}

	iorwc, err := serial.OpenPort(cfg)

	if err != nil {

		log.Error().Caller().Msg(err.Error())
		return
	}

	defer iorwc.Close()

	buffer := make([]byte, 500)

	//发命令之前清空缓冲区
	num, err := iorwc.Read(buffer)

	//发命令数据类型为[]byte
	byteHex, err := hex.DecodeString(cmd)

	if err != nil {

		log.Error().Caller().Msg(err.Error())
		return
	}

	num, err = iorwc.Write(byteHex)

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	tryNum := 2
	for i := 0; i < 100; i++ {

		num, err = iorwc.Read(buffer)
		if num > 0 {
			res += fmt.Sprintf("%x", string(buffer[:num]))
		} else {
			tryNum--
		}

		//查找读到信息的结尾标志
		if tryNum == 0 {
			break
		}
	}

	if len(res) <= 0 {
		res = ""
		return
	}

	//校验crc
	if checkCrc {
		if !checkCrc2(res, crcType, crcLH) {
			res = ""
			err = errors.New("校验crc失败")
		}
	}

	return
}

//checkCrc2 check crc with crctype
func checkCrc2(rdata string, crcType *crc.Parameters, crcLH bool) (res bool) {

	checkCode := rdata[len(rdata)-4:]

	val := rdata[:len(rdata)-4]

	hexByte, _ := hex.DecodeString(val)

	crcCode := crc.CalculateCRC(crcType, hexByte)
	// crcCode := crc.CalculateCRC(crc.CRC16MODBUS, []byte(val))

	crcCodeStr := fmt.Sprintf("%04X", crcCode)

	if crcLH {

		crcCodeStr = crcCodeStr[2:] + crcCodeStr[:2]
	}

	// fmt.Println(crcCodeStr)

	res = false

	checkCode = strings.ToUpper(checkCode)

	if strings.EqualFold(crcCodeStr, checkCode) {
		res = true
	}

	return
}
