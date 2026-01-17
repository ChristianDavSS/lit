package languages

type Feedback struct {
	Value   int
	Message string
}

var Messages = map[string][]Feedback{
	"parameters": {
		Feedback{Value: 5, Message: "   WARNING: The function takes more parameters than recommended.\n"},
		Feedback{Value: 8, Message: "   WARNING: The function is dangerous and has many responsibilities.\n    You may want to make it cleaner/readable.\n"},
		Feedback{Value: 10, Message: "   POTENTIAL ERROR: The function takes too much parameters itself. You should simplify it.\n"},
	},
	"complexity": {
		Feedback{Value: 10, Message: "   WARNING: The function it's a little bit complex.\n"},
		Feedback{Value: 15, Message: "   POTENTIAL ERROR: The function it's very complex itself. You may simplify or break it down.\n"},
		Feedback{Value: 20, Message: "   POTENTIAL ERROR: The function it's too complex and unreadable, taking too much decisions.\n    You must modify it and make it cleaner!.\n"},
	},
	"size": {
		Feedback{Value: 85, Message: "   WARNING: The functions might be long. You might want to make it smaller.\n"},
		Feedback{Value: 120, Message: "   WARNING: The function is longer than recommended. You should consider a refactor.\n"},
		Feedback{Value: 150, Message: "   POTENTIAL ERROR: The function is too long. You must refactor now!.\n"},
	},
}
