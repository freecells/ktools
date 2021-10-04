/*
 * @Author: Keven
 * @version: v1.0.1
 * @Date: 2021-04-09 20:32:12
 * @LastEditors: Keven
 * @LastEditTime: 2021-04-09 20:53:10
 */
package tfile

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
)

func HttpGet(imgUrl, savePath string) (fileName string, err error) {

	res, err := http.DefaultClient.Get(imgUrl)

	if err != nil || res.StatusCode != 200 {
		return "", err
	}

	imgType := res.Header.Get("Content-Type")
	splitArr := strings.Split(imgType, "/")
	imgType = splitArr[len(splitArr)-1]

	if AcceptImgType(imgType) {

		fileName = uuid.New().String() + "." + imgType

		// savePath like storage/public/afiles/
		aimg, err := os.Create(savePath + fileName)
		defer aimg.Close()

		if err != nil {
			return "", err
		}

		_, err = io.Copy(aimg, res.Body)

		if err != nil {
			return "", err
		}
	} else {
		err = errors.New("图片格式不正确")
		return "", err
	}

	return
}
