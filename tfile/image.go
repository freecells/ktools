/*
 * @Author: Keven
 * @version: v1.0.1
 * @Date: 2021-03-08 15:02:40
 * @LastEditors: Keven
 * @LastEditTime: 2021-11-02 11:06:59
 */
package tfile

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrBucket       = errors.New("Invalid Image bucket!")
	ErrSize         = errors.New("Invalid Image size!")
	ErrInvalidImage = errors.New("Invalid Image Type!")
)

type ImageConf struct {
	MaxWidth  int
	MaxHeight int
}

func B64ImageSave(savePath, data string, config ImageConf) (suffix string, err error) {
	idx := strings.Index(data, ";base64,")
	if idx < 0 {
		return "", ErrInvalidImage
	}
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data[idx+8:]))
	buff := bytes.Buffer{}
	_, err = buff.ReadFrom(reader)
	if err != nil {
		return "", err
	}
	imgCfg, suffix, err := image.DecodeConfig(bytes.NewReader(buff.Bytes()))
	if err != nil {
		return "", err
	}

	if imgCfg.Width > config.MaxWidth || imgCfg.Height > config.MaxHeight {
		return "", ErrSize
	}

	saveFullPath := savePath + "." + suffix

	ioutil.WriteFile(saveFullPath, buff.Bytes(), 0666)

	return
}

//AcceptImgType 系统接受的图片类型
func AcceptImgType(imgType string) bool {

	acceptTypes := []string{"webp", "jpeg", "jpg", "png", "gif"}

	for _, val := range acceptTypes {

		if val == imgType {
			return true
		}
	}

	return false
}

//随机文件名
func RandomFileName(fileName string) (fname, ftype string, err error) {

	sfs := strings.Split(fileName, ".")

	if len(sfs) >= 2 {
		ftype = sfs[len(sfs)-1]
	} else {

		err = errors.New("文件格式不正确")
		return
	}

	fname = uuid.New().String() + "." + ftype

	return
}
