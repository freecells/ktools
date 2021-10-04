/*
 * @Author: Keven
 * @version: v1.0.1
 * @Date: 2021-09-28 13:31:03
 * @LastEditors: Keven
 * @LastEditTime: 2021-09-28 13:31:03
 */
package cmds

import (
	"encoding/hex"
	"fmt"

	"github.com/snksoft/crc"
)

//CmdWithCrc return cmd string with crc
func CmdWithCrc(cmd string) (cmdcrc string) {

	hexByte, _ := hex.DecodeString(cmd)

	crcCode := crc.CalculateCRC(CRC16MODBUS, hexByte)

	crcCodeStr := fmt.Sprintf("%04X", crcCode)

	crcCodeStr = crcCodeStr[2:] + crcCodeStr[:2]

	cmdcrc = cmd + crcCodeStr

	return
}
