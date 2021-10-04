package ktools

type DataSets struct {
	Label           string    `json:"label"`
	BackgroundColor string    `json:"backgroundColor"`
	BorderColor     string    `json:"borderColor"`
	Fill            bool      `json:"fill"`
	Data            []float32 `json:"data"`
}
type ChartData struct {
	Labels   []string   `json:"labels"`
	DataSets []DataSets `json:"datasets"`
}

type ChartOption struct {
	Responsive bool      `json:"responsive"`
	Title      Title     `json:"title"`
	Tooltips   Tooltips  `json:"tooltips"`
	Hover      Hover     `json:"hover"`
	Scales     Scales    `json:"scales"`
	Animation  Animation `json:"animation"`
}

type Title struct {
	Display bool   `json:"display"`
	Text    string `json:"text"`
}

type Tooltips struct {
	Mode      string `json:"mode"`
	Intersect bool   `json:"intersect"`
}

type Hover struct {
	Mode      string `json:"mode"`
	Intersect bool   `json:"intersect"`
}

type Scales struct {
	X XY `json:"x"`
	Y XY `json:"y"`
}

type XY struct {
	Stacked    bool       `json:"stacked"`
	Display    bool       `json:"display"`
	ScaleLabel ScaleLabel `json:"scaleLabel"`
}

type ScaleLabel struct {
	Display     bool   `json:"display"`
	LabelString string `json:"labelString"`
}

type Animation struct {
	Duration int    `json:"duration"`
	Easing   string `json:"easing"`
}

type ChatSet struct {
	Type    string      `json:"type"`
	Data    ChartData   `json:"data"`
	Options ChartOption `json:"options"`
}

var BarOption = ChartOption{
	Responsive: true,
	Title: Title{
		Text:    "bar chart",
		Display: true,
	},
	Tooltips: Tooltips{
		Mode:      "index",
		Intersect: false,
	},
	Hover: Hover{
		Mode:      "nearest",
		Intersect: true,
	},
	Animation: Animation{
		Duration: 1000,
		Easing:   "easeInQuart",
	},

	Scales: Scales{
		X: XY{
			Stacked: false,
			Display: true,
			ScaleLabel: ScaleLabel{
				Display:     true,
				LabelString: "时间",
			},
		},
		Y: XY{
			Stacked: false,
			Display: true,
			ScaleLabel: ScaleLabel{
				Display:     true,
				LabelString: "数值",
			},
		},
	},
}
