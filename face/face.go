package face

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/freecells/ktools/tmath"

	"github.com/disintegration/imaging"
	pigo "github.com/esimov/pigo/core"
	"github.com/fogleman/gg"
	"golang.org/x/crypto/ssh/terminal"
)

const banner = `
┌─┐┬┌─┐┌─┐
├─┘││ ┬│ │
┴  ┴└─┘└─┘

Go (Golang) Face detection library.
    Version: %s

`

// pipeName is the file name that indicates stdin/stdout is being used.
const pipeName = "-"

// Version indicates the current build version.
var Version string

var (
	dc        *gg.Context
	fd        *faceDetector
	plc       *pigo.PuplocCascade
	flpcs     map[string][]*pigo.FlpCascade
	imgParams *pigo.ImageParams
)

var (
	eyeCascades  = []string{"lp46", "lp44", "lp42", "lp38", "lp312"}
	mouthCascade = []string{"lp93", "lp84", "lp82", "lp81"}
)

// faceDetector struct contains Pigo face detector general settings.
type faceDetector struct {
	angle         float64
	cascadeFile   string
	destination   string
	minSize       int
	maxSize       int
	shiftFactor   float64
	scaleFactor   float64
	iouThreshold  float64
	puploc        bool
	puplocCascade string
	flploc        bool
	flplocDir     string
	markDetEyes   bool
}

// coord holds the detection coordinates
type coord struct {
	Row   int `json:"x,omitempty"` //	YYYYYYYYYYYYYYYY
	Col   int `json:"y,omitempty"` //   XXXXXXXXXXXXXXXX
	Scale int `json:"size,omitempty"`
}

// detection holds the detection points of the various detection types
type detection struct {
	FacePoints     coord   `json:"face,omitempty"`
	EyePoints      []coord `json:"eyes,omitempty"`
	LandmarkPoints []coord `json:"landmark_points,omitempty"`
}

var FaceDets []detection

//FaceDet 面部点位获取
func FaceDet(imgPath string) (allFacePoints []float64) {

	imgPaths := strings.Split(imgPath, "/")

	imgName := imgPaths[len(imgPaths)-1]

	var (
		// Flags
		source        = imgPath                      //flag.String("in", pipeName, "Source image")
		destination   = "storage/app/out/" + imgName //flag.String("out", pipeName, "Destination image")
		cascadeFile   = "libs/cascade/facefinder"    //flag.String("cf", "", "Cascade binary file")
		puploc        = true                         //flag.Bool("pl", false, "Pupils/eyes localization")
		puplocCascade = "libs/cascade/puploc"        // flag.String("plc", "", "Pupil localization cascade file")
		flploc        = true                         // flag.Bool("flp", false, "Use facial landmark points localization")
		flplocDir     = "libs/cascade/lps"           //flag.String("flpdir", "", "The facial landmark points base directory")
		minSize       = 100                          // flag.Int("min", 20, "Minimum size of face")
		maxSize       = 1000                         // flag.Int("max", 1000, "Maximum size of face")
		shiftFactor   = 0.1                          //flag.Float64("shift", 0.1, "Shift detection window by percentage")
		angle         = 0.0                          //flag.Float64("angle", 0.0, "0.0 is 0 radians and 1.0 is 2*pi radians")
		iouThreshold  = 0.2                          //flag.Float64("iou", 0.2, "Intersection over union (IoU) threshold")
		scaleFactor   = 1.1                          // flag.Float64("scale", 1.1, "Scale detection window by percentage")
		isCircle      = false                        //flag.Bool("circle", false, "Use circle as detection marker")
		markEyes      = false                        // flag.Bool("mark", true, "Mark detected eyes")
		// jsonf         = ""                              //ssflag.String("json", "", "Output the detection points into a json file")
	)

	// Progress indicator
	s := new(spinner)
	s.start("Processing...")
	start := time.Now()

	fd = &faceDetector{
		angle:         angle,
		destination:   destination,
		cascadeFile:   cascadeFile,
		minSize:       minSize,
		maxSize:       maxSize,
		shiftFactor:   shiftFactor,
		scaleFactor:   scaleFactor,
		iouThreshold:  iouThreshold,
		puploc:        puploc,
		puplocCascade: puplocCascade,
		flploc:        flploc,
		flplocDir:     flplocDir,
		markDetEyes:   markEyes,
	}

	var dst io.Writer
	if fd.destination != "empty" {
		if fd.destination == pipeName {
			if terminal.IsTerminal(int(os.Stdout.Fd())) {
				log.Fatalln("`-` should be used with a pipe for stdout")
			}
			dst = os.Stdout
		} else {
			fileTypes := []string{".jpg", ".jpeg", ".png"}
			ext := filepath.Ext(fd.destination)

			if !inSlice(ext, fileTypes) {
				log.Fatalf("Output file type not supported: %v", ext)
			}

			fn, err := os.OpenFile(fd.destination, os.O_CREATE|os.O_WRONLY, 0755)
			if err != nil {
				log.Fatalf("Unable to open output file: %v", err)
			}
			defer fn.Close()
			dst = fn
		}
	}

	faces, err := fd.detectFaces(source)
	if err != nil {
		log.Fatalf("Detection error: %v", err)
	}

	dets, err := fd.drawFaces(faces, isCircle) //可能会检测到多个人脸
	if err != nil {
		log.Fatalf("Error creating the image output: %s", err)
	}

	if fd.destination != "empty" {
		if err := fd.encodeImage(dst); err != nil {
			log.Fatalf("Error encoding the output image: %v", err)
		}
	}

	// fmt.Println("Dets", dets)

	allFacePoints = dealFaceCode(dets)

	s.stop()

	log.Printf("\nDone in: \x1b[92m%.2fs\n", time.Since(start).Seconds())

	return
}

//处理点位
func dealFaceCode(dets []detection) (points []float64) {

	//选择最佳点位 集合
	bestFaceIndex := 0
	tempScale := 0
	for i := 0; i < len(dets); i++ {
		scale := dets[i].FacePoints.Scale
		if scale > tempScale {
			tempScale = scale
			bestFaceIndex = i
		}
	}
	bestDet := dets[bestFaceIndex]

	// fmt.Println(bestDet)
	//========end==========
	if len(bestDet.LandmarkPoints) < 15 {
		return
	}

	// fmt.Printf("the point length is %d,point is %s", len(bestDet.LandmarkPoints), bestDet.LandmarkPoints)
	//纠正点位 //点位放大到目标大小
	if len(bestDet.LandmarkPoints) == 15 {

		fangDa := 1.0

		tgSize := 380.0

		faceSize := float64(bestDet.FacePoints.Scale)

		fangDa = tgSize / faceSize

		leftEye := bestDet.EyePoints[0]

		rightEye := bestDet.EyePoints[1]

		//计算 脸部倾斜角度
		angle := tmath.CalcAngle(float64(leftEye.Col), float64(leftEye.Row), float64(rightEye.Col), float64(rightEye.Row))
		fmt.Printf("the angle is %f", angle)
		// if angle > 1 {}
		if leftEye.Row < rightEye.Row {
			angle = -angle //负角度 为顺时针旋转
		}

		setPoints := []float64{}

		// myFacePoint := []coord{}

		mx := bestDet.FacePoints.Col - (bestDet.FacePoints.Scale / 2)

		my := bestDet.FacePoints.Row - (bestDet.FacePoints.Scale / 2)
		// mx, my = 0, 0

		// setPoints = append(setPoints, float64((bestDet.FacePoints.Col-mx)*fangDa), float64((bestDet.FacePoints.Row-my)*fangDa))

		//执行特征点 左上角平移与旋转 放大
		for _, val2 := range bestDet.LandmarkPoints {

			val2.Row -= my

			val2.Col -= mx

			//平移
			x1, y1 := float64(val2.Col), float64(val2.Row)

			if angle > 1 {

				//旋转
				x1, y1 = tmath.XytoX1Y1(x1, y1, angle)
			}

			//放大
			setPoints = append(setPoints, x1*fangDa, y1*fangDa)

		}

		//瞳孔位置判定有问题  检测范围跳动很大。。。

		if len(setPoints) < 20 {
			setPoints = []float64{}
		}

		DrawImg(setPoints)
		points = setPoints

	}
	//=====end========================================

	return
}

// detectFaces run the detection algorithm over the provided source image.
func (fd *faceDetector) detectFaces(source string) ([]pigo.Detection, error) {
	var srcFile io.Reader
	if source == pipeName {
		if terminal.IsTerminal(int(os.Stdin.Fd())) {
			log.Fatalln("`-` should be used with a pipe for stdin")
		}
		srcFile = os.Stdin
	} else {
		file, err := os.Open(source)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		srcFile = file
	}

	src, err := pigo.DecodeImage(srcFile)
	if err != nil {
		return nil, err
	}

	pixels := pigo.RgbToGrayscale(src)
	cols, rows := src.Bounds().Max.X, src.Bounds().Max.Y

	dc = gg.NewContext(cols, rows)
	dc.DrawImage(src, 0, 0)

	imgParams = &pigo.ImageParams{
		Pixels: pixels,
		Rows:   rows,
		Cols:   cols,
		Dim:    cols,
	}

	cParams := pigo.CascadeParams{
		MinSize:     fd.minSize,
		MaxSize:     fd.maxSize,
		ShiftFactor: fd.shiftFactor,
		ScaleFactor: fd.scaleFactor,
		ImageParams: *imgParams,
	}

	cascadeFile, err := ioutil.ReadFile(fd.cascadeFile)
	if err != nil {
		return nil, err
	}

	p := pigo.NewPigo()
	// Unpack the binary file. This will return the number of cascade trees,
	// the tree depth, the threshold and the prediction from tree's leaf nodes.
	classifier, err := p.Unpack(cascadeFile)
	if err != nil {
		return nil, err
	}

	if fd.puploc {
		pl := pigo.NewPuplocCascade()

		cascade, err := ioutil.ReadFile(fd.puplocCascade)
		if err != nil {
			return nil, err
		}
		plc, err = pl.UnpackCascade(cascade)
		if err != nil {
			return nil, err
		}

		if fd.flploc {
			flpcs, err = pl.ReadCascadeDir(fd.flplocDir)
			if err != nil {
				return nil, err
			}
		}
	}

	// Run the classifier over the obtained leaf nodes and return the detection results.
	// The result contains quadruplets representing the row, column, scale and detection score.
	faces := classifier.RunCascade(cParams, fd.angle)

	// Calculate the intersection over union (IoU) of two clusters.
	faces = classifier.ClusterDetections(faces, fd.iouThreshold)

	return faces, nil
}

// drawFaces marks the detected faces with a circle in case isCircle is true, otherwise marks with a rectangle.
func (fd *faceDetector) drawFaces(faces []pigo.Detection, isCircle bool) ([]detection, error) {
	var (
		qThresh float32 = 5.0
		perturb         = 63
	)

	var (
		detections     []detection
		eyesCoords     []coord
		landmarkCoords []coord
		puploc         *pigo.Puploc
	)

	for _, face := range faces {
		if face.Q > qThresh {
			if isCircle {
				dc.DrawArc(
					float64(face.Col),
					float64(face.Row),
					float64(face.Scale/2),
					0,
					2*math.Pi,
				)
			} else {
				dc.DrawRectangle(
					float64(face.Col-face.Scale/2),
					float64(face.Row-face.Scale/2),
					float64(face.Scale),
					float64(face.Scale),
				)
			}
			faceCoord := &coord{
				Row:   face.Row,
				Col:   face.Col,
				Scale: face.Scale,
			}

			dc.SetLineWidth(2.0)
			dc.SetStrokeStyle(gg.NewSolidPattern(color.RGBA{R: 255, G: 0, B: 0, A: 255}))
			dc.Stroke()

			if fd.puploc && face.Scale > 50 {
				rect := image.Rect(
					face.Col-face.Scale/2,
					face.Row-face.Scale/2,
					face.Col+face.Scale/2,
					face.Row+face.Scale/2,
				)
				rows, cols := rect.Max.X-rect.Min.X, rect.Max.Y-rect.Min.Y
				ctx := gg.NewContext(rows, cols)
				faceZone := ctx.Image()

				// left eye
				puploc = &pigo.Puploc{
					Row:      face.Row - int(0.075*float32(face.Scale)),
					Col:      face.Col - int(0.175*float32(face.Scale)),
					Scale:    float32(face.Scale) * 0.25,
					Perturbs: perturb,
				}
				leftEye := plc.RunDetector(*puploc, *imgParams, fd.angle, false)
				if leftEye.Row > 0 && leftEye.Col > 0 {
					if fd.angle > 0 {
						drawDetections(ctx,
							float64(cols/2-(face.Col-leftEye.Col)),
							float64(rows/2-(face.Row-leftEye.Row)),
							float64(leftEye.Scale),
							color.RGBA{R: 255, G: 0, B: 0, A: 255}, //瞳孔
							fd.markDetEyes,
						)
						angle := (fd.angle * 180) / math.Pi
						rotated := imaging.Rotate(faceZone, 2*angle, color.Transparent)
						final := imaging.FlipH(rotated)

						dc.DrawImage(final, face.Col-face.Scale/2, face.Row-face.Scale/2)
					} else {
						drawDetections(dc,
							float64(leftEye.Col),
							float64(leftEye.Row),
							float64(leftEye.Scale),
							color.RGBA{R: 255, G: 0, B: 0, A: 255}, //瞳孔
							fd.markDetEyes,
						)
					}
					eyesCoords = append(eyesCoords, coord{
						Row:   leftEye.Row,
						Col:   leftEye.Col,
						Scale: int(leftEye.Scale),
					})
				}

				// right eye
				puploc = &pigo.Puploc{
					Row:      face.Row - int(0.075*float32(face.Scale)),
					Col:      face.Col + int(0.185*float32(face.Scale)),
					Scale:    float32(face.Scale) * 0.25,
					Perturbs: perturb,
				}

				rightEye := plc.RunDetector(*puploc, *imgParams, fd.angle, false)
				if rightEye.Row > 0 && rightEye.Col > 0 {
					if fd.angle > 0 {
						drawDetections(ctx,
							float64(cols/2-(face.Col-rightEye.Col)),
							float64(rows/2-(face.Row-rightEye.Row)),
							float64(rightEye.Scale),
							color.RGBA{R: 255, G: 0, B: 0, A: 255}, //瞳孔位置
							fd.markDetEyes,
						)
						// convert radians to angle
						angle := (fd.angle * 180) / math.Pi
						rotated := imaging.Rotate(faceZone, 2*angle, color.Transparent)
						final := imaging.FlipH(rotated)

						dc.DrawImage(final, face.Col-face.Scale/2, face.Row-face.Scale/2)
					} else {
						drawDetections(dc,
							float64(rightEye.Col),
							float64(rightEye.Row),
							float64(rightEye.Scale),
							color.RGBA{R: 255, G: 0, B: 0, A: 255}, //瞳孔位置
							fd.markDetEyes,
						)
					}
					eyesCoords = append(eyesCoords, coord{
						Row:   rightEye.Row,
						Col:   rightEye.Col,
						Scale: int(rightEye.Scale),
					})
				}

				if fd.flploc {
					for _, eye := range eyeCascades {
						for _, flpc := range flpcs[eye] {
							flp := flpc.FindLandmarkPoints(leftEye, rightEye, *imgParams, perturb, false)
							if flp.Row > 0 && flp.Col > 0 {
								drawDetections(dc,
									float64(flp.Col),
									float64(flp.Row),
									float64(flp.Scale*0.5),
									color.RGBA{R: 0, G: 100, B: 222, A: 255}, // 眼睛
									false,
								)
							}
							//左眼 点位
							landmarkCoords = append(landmarkCoords, coord{
								Row:   flp.Row,
								Col:   flp.Col,
								Scale: int(flp.Scale),
							})

							flp = flpc.FindLandmarkPoints(leftEye, rightEye, *imgParams, perturb, true)
							if flp.Row > 0 && flp.Col > 0 {
								drawDetections(dc,
									float64(flp.Col),
									float64(flp.Row),
									float64(flp.Scale*0.5),
									color.RGBA{R: 0, G: 100, B: 222, A: 255}, // 眼睛
									false,
								)
							}

							//右眼点位
							landmarkCoords = append(landmarkCoords, coord{
								Row:   flp.Row,
								Col:   flp.Col,
								Scale: int(flp.Scale),
							})
						}
					}

					for _, mouth := range mouthCascade {
						for _, flpc := range flpcs[mouth] {
							flp := flpc.FindLandmarkPoints(leftEye, rightEye, *imgParams, perturb, false)
							if flp.Row > 0 && flp.Col > 0 {
								drawDetections(dc,
									float64(flp.Col),
									float64(flp.Row),
									float64(flp.Scale*0.5),
									color.RGBA{R: 230, G: 230, B: 50, A: 255}, // 鼻子 嘴
									false,
								)
							}
							landmarkCoords = append(landmarkCoords, coord{
								Row:   flp.Row,
								Col:   flp.Col,
								Scale: int(flp.Scale),
							})
						}
					}
					flp := flpcs["lp84"][0].FindLandmarkPoints(leftEye, rightEye, *imgParams, perturb, true)
					if flp.Row > 0 && flp.Col > 0 {
						drawDetections(dc,
							float64(flp.Col),
							float64(flp.Row),
							float64(flp.Scale*0.5),
							color.RGBA{R: 160, G: 0, B: 255, A: 255}, //右嘴角
							false,
						)
						landmarkCoords = append(landmarkCoords, coord{
							Row:   flp.Row,
							Col:   flp.Col,
							Scale: int(flp.Scale),
						})
					}
				}
			}
			detections = append(detections, detection{
				FacePoints:     *faceCoord,
				EyePoints:      eyesCoords,
				LandmarkPoints: landmarkCoords,
			})
		}
	}
	return detections, nil
}

func (fd *faceDetector) encodeImage(dst io.Writer) error {
	var err error
	img := dc.Image()

	switch dst.(type) {
	case *os.File:
		ext := filepath.Ext(dst.(*os.File).Name())
		switch ext {
		case "", ".jpg", ".jpeg":
			err = jpeg.Encode(dst, img, &jpeg.Options{Quality: 100})
		case ".png":
			err = png.Encode(dst, img)
		default:
			err = errors.New("unsupported image format")
		}
	default:
		err = jpeg.Encode(dst, img, &jpeg.Options{Quality: 100})
	}
	return err
}

type spinner struct {
	stopChan chan struct{}
}

// Start process
func (s *spinner) start(message string) {
	s.stopChan = make(chan struct{}, 1)

	go func() {
		for {
			for _, r := range `⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏` {
				select {
				case <-s.stopChan:
					return
				default:
					fmt.Fprintf(os.Stderr, "\r%s%s %c%s", message, "\x1b[35m", r, "\x1b[39m")
					time.Sleep(time.Millisecond * 100)
				}
			}
		}
	}()
}

// End process
func (s *spinner) stop() {
	s.stopChan <- struct{}{}
}

// inSlice checks if the item exists in the slice.
func inSlice(item string, slice []string) bool {
	for _, it := range slice {
		if it == item {
			return true
		}
	}
	return false
}

// drawDetections is a helper function to draw the detection marks
func drawDetections(ctx *gg.Context, x, y, r float64, c color.RGBA, markDet bool) {
	ctx.DrawArc(x, y, r*0.15, 0, 2*math.Pi)
	ctx.SetFillStyle(gg.NewSolidPattern(c))
	ctx.Fill()

	if markDet {
		ctx.DrawRectangle(x-(r*1.5), y-(r*1.5), r*3, r*3)
		ctx.SetLineWidth(2.0)
		ctx.SetStrokeStyle(gg.NewSolidPattern(color.RGBA{R: 255, G: 255, B: 0, A: 255}))
		ctx.Stroke()
	}
}

func DrawImg(points []float64) {

	dc := gg.NewContext(1000, 720)

	dc.SetRGB(255, 255, 255)
	dc.Fill()
	dc.FillPreserve()
	dc.Stroke()

	for i := 0; i < len(points); i += 2 {

		dc.SetRGB(0, 0, 255)

		dc.DrawPoint(points[i], points[i+1], 5)
		dc.SetRGB(0, 0, 255)
		dc.Fill()
	}

	imgName := strconv.FormatInt(time.Now().UTC().UnixNano(), 10) + ".png"

	dc.SavePNG("out/0_" + imgName)
}
