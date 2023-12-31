package test

import (
	"fmt"
	"github.com/fogleman/gg"
	"image/color"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
	"unicode"
)

func Test_img(t *testing.T) {
	str := "const axios = require(\"axios\")\n\n//{\"post_type\":\"request\",\"request_type\":\"friend\",\"time\":1679386108,\"self_id\":2673893724,\"user_id\":1978150028,\"comment\":\"信息\",\"flag\":\"1679386108000000\"}\nexports.addfriends = async ({flag, comment}) => {\n    console.log(\"commit:\", comment)\n    const {data: res} = await axios({\n        method: \"post\",\n        url: \"http://127.0.0.1:5000/set_friend_add_request\",\n        data: {\n            flag: flag,\n            approve: true,\n        }\n    })\n    console.log(res)\n\n}\n"
	var width float64
	var s string
	list := make([]string, 0)
	// 判断文件是否存在
	if _, err := os.Stat("../config/tr.png"); err != nil {
		//不存在
		fmt.Println("不存在")
		width := 1960
		height := 1080
		times := time.Now().Format("2006-01-02 15:04:05")
		fmt.Println(times)
		dc := gg.NewContext(width, height)
		dc.SetHexColor("#FF8B13")
		dc.DrawRectangle(0, 0, float64(width), float64(height))
		dc.Fill()
		//字体
		face, err := gg.LoadFontFace("../config/t.ttf", 60)
		if err != nil {
			log.Panicln(err)
			return
		}
		f, err := gg.LoadFontFace("../config/t.ttf", 40)
		if err != nil {
			log.Panicln(err)
			return
		}
		dc.SetFontFace(face)
		dc.SetHexColor("#27374D")
		//加载图片
		image, err := gg.LoadImage("../config/3.png")
		if err != nil {
			return
		}

		//时间
		//dc.DrawStringAnchored(times, 210, float64(height-70), 0.5, 0.5)

		dc.DrawImageAnchored(image, width-170, height-80, 0.5, 0.5)

		dc.SetColor(color.RGBA{249, 251, 231, 150})
		dc.SetFontFace(f)

		rand.Seed(time.Now().UnixMicro())
		for i := 0; i < 10; i++ {
			fmt.Println()
			dc.Push()
			dc.RotateAbout(gg.Radians(40), float64(width/2), float64(height/2))
			dc.DrawStringAnchored("@GoBat", float64(rand.Int63n(1920)), float64(rand.Int63n(1080)), 0.5, 0.5)
			dc.Pop()
		}
		err = dc.SavePNG("../config/t.png")
		if err != nil {
			log.Panicln(err)
		}

	} else {
		//图片存在
		fmt.Println("存在")
		img, err := gg.LoadImage("../config/t.png")
		if err != nil {
			log.Panicln(err)
		}
		face, err := gg.LoadFontFace("../config/t.ttf", 70)
		if err != nil {
			log.Panicln(err)
			return
		}
		dc := gg.NewContextForImage(img)
		dc.SetFontFace(face)
		dc.SetHexColor("#333")
		wd, _ := dc.MeasureString(str)
		fmt.Println("wd", wd)
		if wd < 1960 {
			//不满1行
			dc.DrawString(str, 40, 350)
		} else {
			var h float64
			var w float64
			for i, r := range str {

				w, h = dc.MeasureString(string(r))
				wd, _ := dc.MeasureString(str[i:])
				width += w
				s += string(r)
				if width >= 1960 {
					list = append(list, s)
					width = 0
					s = ""
				}

				if wd < 1800 && (i == len(str)-3 && !unicode.IsLetter(rune(str[i : i+1][0])) && !unicode.IsSymbol(rune(str[i : i+1][0])) && !unicode.IsNumber(rune(str[i : i+1][0])) && !unicode.IsSpace(rune(str[i : i+1][0]))) || i == len(str)-1 {
					fmt.Println(i, len(str), s, unicode.IsLetter(rune(str[i : i+1][0])), str[i : i+1][0])
					list = append(list, s)
				}

			}
			for i, s := range list {
				fmt.Println(i, s)
				if i == 0 {
					dc.DrawString(s, 40, 80+float64(i)*h)
					continue
				}
				dc.DrawString(s, 0, 80+float64(i)*h)
			}
		}
		err = dc.SavePNG("../config/f.png")
		if err != nil {
			log.Panicln(err)
		}
	}
}
