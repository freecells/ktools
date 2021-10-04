/*
 * @Author: Feng
 * @version: v1.0.0
 * @Date: 2020-07-03 13:50:49
 * @LastEditors: Keven
 * @LastEditTime: 2021-10-04 11:03:26
 */

package cmds

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/freecells/ktools/datadeal"
	"github.com/freecells/ktools/datetime"
	"github.com/freecells/ktools/ktcp"

	"github.com/rs/zerolog/log"
)

type cmdFormat struct {
	Status   string
	Unset    string
	Pdata    string
	OpenDoor string
}

//cmd 命令模板
var cmdFormats = cmdFormat{
	Status:   "%s0300000034",
	Unset:    "%s0300010052",
	Pdata:    "%s03000600ff",
	OpenDoor: "%s0500%s00ff",
}

type ErrAll struct {
	Ingrid     int
	Charging   string
	ChargeFull string
	Haserr     int
	Battery    string
	No2        string
	No         string
	So2        string
	H2         string
	H2s        string
	Co2        string
	Co         string
	O2         string
	Ch4        string
}

/*Cstatus *
* todo:need test
* @msg:
* @param Cstatusl
* @return:
 */
func Cstatus(ip, cmd string) (cdoors []string, cerrs []ErrAll, err error) {

	//data for 42 door length = 97
	//02035C0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000EF78

	//data for 48 door： len = 218 020368f0ffefffff7f0000020000000100020002000200020002000200020002000200020002000100020002000200020002000000020002000200020002000200020002000200020002000200010002000200020002000200020002000200000000000000000000000100a453

	data, err := ktcp.Send(ip, cmd)
	log.Info().Caller().Msgf("柜子仪器状态cmd is: %s , data：%s", cmd, data)

	if err != nil {
		log.Error().Caller().Msgf("柜子 仪器状态读取错误：%s，返回数据是 : %s", err.Error(), data)
		return
	}

	if len(data) < 218 {
		err = errors.New("柜子状态数据 格式错误 data is: " + data)
		return
	}

	doorStatus := data[6:18]

	endNum := len(data) - 4
	errStatus := data[22:endNum] //获得最终的数值

	doorOC := []string{}

	for i := 0; i < 6; i++ {
		door8 := doorStatus[i*2 : i*2+2]
		// door8 = door8[2:] + door8[:2]
		doorint, _ := strconv.ParseInt(door8, 16, 0)
		door8bit := fmt.Sprintf("%08b", doorint)
		door8bit = datadeal.ReverseString(door8bit)

		for i := 0; i < 8; i++ {
			if door8bit[i:i+1] == "0" {
				doorOC = append(doorOC, "open")
			} else {
				doorOC = append(doorOC, "close")
			}
		}
	}
	cdoors = doorOC

	errAlls := []ErrAll{}
	for i := 0; i < 48; i++ {
		err1 := errStatus[i*4 : i*4+4]
		err1 = err1[2:] + err1[:2]

		errint, _ := strconv.ParseInt(err1, 16, 0)
		errbit := fmt.Sprintf("%016b", errint)
		errbit = datadeal.ReverseString(errbit)

		ingrid := 1
		if errbit[0:1] == "0" && errbit[1:2] == "0" {
			ingrid = 0
		}

		haserr := 0
		if errbit[2:3] == "1" || errbit[3:4] == "1" || errbit[4:5] == "1" || errbit[5:6] == "1" {
			haserr = 1
		}
		if errbit[6:7] == "1" || errbit[7:8] == "1" || errbit[8:9] == "1" {
			haserr = 1
		}
		if errbit[9:10] == "1" || errbit[10:11] == "1" || errbit[11:12] == "1" {
			haserr = 1
		}

		errall := ErrAll{
			Ingrid:     ingrid,
			Charging:   errbit[0:1],
			ChargeFull: errbit[1:2],
			Haserr:     haserr,
			Battery:    errbit[2:3],
			No2:        errbit[3:4],
			No:         errbit[4:5],
			So2:        errbit[5:6],
			H2:         errbit[6:7],
			H2s:        errbit[7:8],
			Co2:        errbit[8:9],
			Co:         errbit[9:10],
			O2:         errbit[10:11],
			Ch4:        errbit[11:12],
		}

		errAlls = append(errAlls, errall)
	}

	cerrs = errAlls

	return

}

type unZeroValid struct {
	ZeroDay  int
	ValidDay int
}
type UnSetDay struct {
	DoorNo int
	YqType string
	CH4    unZeroValid
	O2     unZeroValid
	CO     unZeroValid
	CO2    unZeroValid
	H2S    unZeroValid
	H2     unZeroValid
	SO2    unZeroValid
	NO     unZeroValid
	NO2    unZeroValid
}

func CunSetDay(ip, mes string) (unset UnSetDay, err error) {

	data, err := ktcp.Send(ip, mes)
	if err != nil {
		log.Debug().Caller().Msgf("获取未标零标校数据 错误：%s，返回数据是 : %s", err.Error(), data)

		return
	}

	// fmt.Println("unset days:", data)

	doorData := data[6:8]
	yqTypeData := data[8:12]
	unsetData := data[12:86]

	tdoorno, _ := strconv.ParseInt(doorData, 16, 0)
	unset.DoorNo = int(tdoorno) + 2
	unset.YqType = yqTypeData

	unset.CH4.ZeroDay, unset.CH4.ValidDay = dealUnset(unsetData[:8])
	unset.O2.ZeroDay, unset.O2.ValidDay = dealUnset(unsetData[8:16])
	unset.CO.ZeroDay, unset.CO.ValidDay = dealUnset(unsetData[16:24])
	unset.CO2.ZeroDay, unset.CO2.ValidDay = dealUnset(unsetData[24:32])
	unset.H2S.ZeroDay, unset.H2S.ValidDay = dealUnset(unsetData[32:40])
	unset.H2.ZeroDay, unset.H2.ValidDay = dealUnset(unsetData[40:48])
	unset.SO2.ZeroDay, unset.SO2.ValidDay = dealUnset(unsetData[48:56])
	unset.NO.ZeroDay, unset.NO.ValidDay = dealUnset(unsetData[56:64])
	unset.NO2.ZeroDay, unset.NO2.ValidDay = dealUnset(unsetData[64:72])

	return
}

//处理 未标 时间
func dealUnset(unsetData string) (zero, valid int) {

	Unzero := unsetData[:4]
	Unzero = Unzero[:2] + Unzero[2:]
	tdata, _ := strconv.ParseInt(Unzero, 16, 0)
	zero = int(tdata)

	Unvalid := unsetData[4:]
	Unvalid = Unvalid[:2] + Unvalid[2:]
	tdata, _ = strconv.ParseInt(Unvalid, 16, 0)
	valid = int(tdata)

	return
}

type PositionData struct {
	PositionID int
	Date       time.Time
	CH4        int
	O2         int
	CO         int
}
type YqData struct {
	DataID int
	DoorNo int
	Datas  []PositionData
}

/**
 * Cyqdatas 处理 仪器定位点数据
 */
func Cyqdatas(ip, mes string) (yqData YqData, err error) {

	data, err := ktcp.Send(ip, mes)
	if err != nil {
		return
	}

	// fmt.Println("yq datas:", data)
	// log.Debug().Caller().Msgf("仪器定位数据：%s", data)

	if len(data) < 200 {
		err = errors.New("cmd is: " + mes + "：返回的数据格式不正确，数据是：" + data)
		return
	}

	didData := data[4:8]
	didData = didData[2:] + didData[:2]
	dnoData := data[8:10]
	//0203 7400 27 0f86
	//0203 7400 0b 2001
	pData := data[14:238]

	t, _ := strconv.ParseInt(didData, 16, 0)

	yqData.DataID = int(t)

	t, _ = strconv.ParseInt(dnoData, 16, 0)

	yqData.DoorNo = int(t) + 2

	for i := 0; i < 8; i++ {
		/**
		* 020374000e2301
		0b00140517102c34000000d20000
		0800140517102c34000000d20000
		0100140517102c36000000d20000
		1100140517102c38000000d20000
		0b00140517102c3a000000d20000
		0800140517102c3a000000d20000
		1100140517102d02000000d20000
		0b00140517102d02000000d20000
		4a71
		*/
		var pdata PositionData
		oneData := pData[i*28 : i*28+28]

		if oneData[16:28] == "000000000000" {
			return
		}

		id := oneData[:4]
		idRev := ""
		if id[:1] == "8" {
			idRev = id[2:] + id[1:2]

		} else {

			idRev = id[2:] + id[:2]
		}

		tid, _ := strconv.ParseInt(idRev, 16, 0)
		pdata.PositionID = int(tid)

		date := oneData[4:16]

		pdata.Date = dealDate(date)

		ch4 := oneData[16:20]
		ch4 = ch4[:2] + ch4[2:]
		t, _ := strconv.ParseInt(ch4, 16, 0)
		pdata.CH4 = int(t)

		o2 := oneData[20:24]
		o2 = o2[:2] + o2[2:]
		t, _ = strconv.ParseInt(o2, 16, 0)

		pdata.O2 = int(t)
		co := oneData[24:28]
		co = co[:2] + co[2:]
		t, _ = strconv.ParseInt(co, 16, 0)
		pdata.CO = int(t)

		// log.Debug().Msgf("解析的仪器数据是: %v", pdata)
		yqData.Datas = append(yqData.Datas, pdata)

		//check read end
		idint, _ := strconv.ParseInt(id, 16, 0)
		idbit := fmt.Sprintf("%016b", idint)

		if idbit[:1] == "1" {
			break
		}
	}

	return
}

func dealDate(hexdate string) (dateTime time.Time) {
	//140517102b3a
	dateStr := ""

	for i := 0; i < len(hexdate)/2; i++ {
		ty := hexdate[i*2 : i*2+2]
		tyy, _ := strconv.ParseInt(ty, 16, 0)
		tyyStr := strconv.Itoa(int(tyy))
		if len(tyyStr) == 1 {
			tyyStr = "0" + tyyStr
		}
		dateStr += tyyStr

	}

	//不处理秒 一分钟以内数据过滤 秒置0
	dateStr = dateStr[:10] + "00"

	dateTime = datetime.ParesLoc("060102150405", dateStr)
	return
}

func CopenDoor(ip, mes string) (err error) {

	_, err = ktcp.Send(ip, mes)
	if err != nil {
		return
	}
	// fmt.Println(data)
	return
}

type Cabinet2Cmd struct {
}

var ipFormat = "192.168.0.%d:23"

func (c2 *Cabinet2Cmd) OpenDoorCmd(cno, doorno int) (err error) {

	cnoStr := fmt.Sprintf("%02x", cno)
	doornoStr := fmt.Sprintf("%02x", doorno)

	cmd := fmt.Sprintf(cmdFormats.OpenDoor, cnoStr, doornoStr)

	cmdcrc := CmdWithCrc(cmd)

	ipAddr := fmt.Sprintf(ipFormat, cno)

	err = CopenDoor(ipAddr, cmdcrc)

	return
}

func (c2 *Cabinet2Cmd) StatusCmd(cno int) (doors []string, errs []ErrAll, err error) {

	cnoStr := fmt.Sprintf("%02x", cno)

	cmd := fmt.Sprintf(cmdFormats.Status, cnoStr)

	cmdCrc := CmdWithCrc(cmd)

	ipAddr := fmt.Sprintf(ipFormat, cno)

	doors, errs, err = Cstatus(ipAddr, cmdCrc)

	return
}
func (c2 *Cabinet2Cmd) UnsetDayCmd(cno int) (unsetDay UnSetDay, err error) {
	cnoStr := fmt.Sprintf("%02x", cno)

	cmd := fmt.Sprintf(cmdFormats.Unset, cnoStr)

	cmdCrc := CmdWithCrc(cmd)

	ipAddr := fmt.Sprintf(ipFormat, cno)

	unsetDay, err = CunSetDay(ipAddr, cmdCrc)

	return
}
func (c2 *Cabinet2Cmd) YqdatasCmd(cno int) (yqdata YqData, err error) {
	cnoStr := fmt.Sprintf("%02x", cno)

	cmd := fmt.Sprintf(cmdFormats.Pdata, cnoStr)

	cmdCrc := CmdWithCrc(cmd)

	ipAddr := fmt.Sprintf(ipFormat, cno)

	yqdata, err = Cyqdatas(ipAddr, cmdCrc)

	return
}
