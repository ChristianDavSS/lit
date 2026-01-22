package domain

// Point represents a position in the source code.
type Point struct {
	Row    uint
	Column uint
}

// Feedback represents a feedback message based on a metric value.
type Feedback struct {
	Value   int
	Message string
}

// FunctionData represents the analysis data of a function.
type FunctionData struct {
	Name                     string
	TotalParams, Complexity  int
	StartPosition            Point
	StartByte, EndByte, Size uint
	Feedback                 string
}

// Messages contains the rules for generating feedback.
var Messages = map[string][]Feedback{
	"parameters": {
		{Value: 5, Message: "   WARNING: The function takes more parameters than recommended.\n"},
		{Value: 8, Message: "   WARNING: The function is dangerous and has many responsibilities.\n    You may want to make it cleaner/readable.\n"},
		{Value: 10, Message: "   POTENTIAL ERROR: The function takes too much parameters itself. You should simplify it.\n"},
	},
	"complexity": {
		{Value: 10, Message: "   WARNING: The function it's a little bit complex.\n"},
		{Value: 15, Message: "   POTENTIAL ERROR: The function it's very complex itself. You may simplify or break it down.\n"},
		{Value: 20, Message: "   POTENTIAL ERROR: The function it's too complex and unreadable, taking too much decisions.\n    You must modify it and make it cleaner!.\n"},
	},
	"size": {
		{Value: 85, Message: "   WARNING: The functions might be long. You might want to make it smaller.\n"},
		{Value: 120, Message: "   WARNING: The function is longer than recommended. You should consider a refactor.\n"},
		{Value: 150, Message: "   POTENTIAL ERROR: The function is too long. You must refactor now!.\n"},
	},
}

// AddInitialData populates the initial data of a FunctionData.
func (f *FunctionData) AddInitialData(name string, totalParams int, startByte, endByte, size uint, startPos Point) {
	f.Name = name
	f.TotalParams = totalParams
	f.StartByte = startByte
	f.EndByte = endByte
	f.Size = size
	f.StartPosition = startPos
}

// IsTargetInRange checks if the function is within the given byte range.
func (f *FunctionData) IsTargetInRange(startByte, endByte uint) bool {
	return f.StartByte <= startByte && f.EndByte >= endByte
}

// SetFunctionFeedback generates feedback based on the function metrics.
func (f *FunctionData) SetFunctionFeedback() {
	for key, value := range Messages {
		var msg string
		for _, item := range value {
			if f.getValue(key) >= item.Value {
				msg = item.Message
			}
		}
		f.Feedback += msg
	}
}

func (f *FunctionData) getValue(key string) int {
	dict := map[string]int{
		"parameters": f.TotalParams,
		"complexity": f.Complexity,
		"size":       int(f.Size),
	}

	return dict[key]
}
