package gquiz

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Knetic/govaluate"
)

type QuizExecutor struct {
	ui      UI
	qResult *QResult
}

func NewQuizExecutor(ui UI, qr *QResult) *QuizExecutor {
	if qr == nil {
		qr = &QResult{}
	}
	return &QuizExecutor{
		ui:      ui,
		qResult: qr,
	}
}

func (qe *QuizExecutor) Execute(qGraph *QGraph) (*QResult, error) {
	root := qGraph.FindRootNode()
	if root == nil {
		return nil, errors.New("no root node found")
	}

	var currentNode *QNode
	currentNode = root
	for currentNode != nil {
		nextNodeName, err := qe.HandleNode(currentNode)

		if err != nil {
			return nil, err
		}
		if nextNodeName != "" {
			currentNode = qGraph.FindNode(nextNodeName)
		} else {
			currentNode = nil
		}
	}

	return qe.qResult, nil
}

// HandleNode will return the next node name.
func (qe *QuizExecutor) HandleNode(qNode *QNode) (string, error) {
	currentVars := map[string]interface{}{}
	for i := range qNode.Questions {
		// TODO error handling.
		q := &qNode.Questions[i]
		answer, _ := qe.HandleQuestion(&qNode.Questions[i])
		if q.Type != "" {
			switch q.Type {
			case "int":
				intV, err := strconv.Atoi(answer)
				if err != nil {
					return "", fmt.Errorf("value %s not a valid int value", answer)
				}
				currentVars[q.VarName] = intV
				break
			default:
				return "", fmt.Errorf("type: %s not supported", q.Type)
			}
		} else {
			currentVars[q.VarName] = answer
		}
		if q.Persistent {
			(*qe.qResult)[q.VarName] = answer
		}
	}

	if len(qNode.Transitions) == 0 {
		return "", nil
	}
	if len(qNode.Transitions) == 1 {
		return qNode.Transitions[0].Name, nil
	}

	// evaluate the result to find the path to go.
	for i := range qNode.Transitions {
		t := qNode.Transitions[i]

		expression, _ := govaluate.NewEvaluableExpression(t.Condition)
		result, _ := expression.Evaluate(currentVars)
		if result.(bool) {
			return t.Name, nil
		}
	}
	return "", errors.New("no valid trasition found")
}

func (qe *QuizExecutor) HandleQuestion(q *Question) (string, error) {
	var answer string
	// use the values in the old vars as the default value.
	var defaultValue string
	if (*qe.qResult)[q.VarName] != "" {
		defaultValue = (*qe.qResult)[q.VarName]
	} else {
		defaultValue = q.Default
	}
	if defaultValue != "" {
		qe.ui.Println(fmt.Sprintf("%s(%s)", q.Description, defaultValue))
	} else {
		qe.ui.Println(q.Description)
	}
	if q.Candidates != nil && len(q.Candidates) > 0 {
		for true {
			for i := range q.Candidates {
				qe.ui.Println(fmt.Sprintf("%d.%s -- %s", (i + 1), q.Candidates[i].Value, q.Candidates[i].Description))
			}
			input := qe.ui.GetInput()
			choice, err := strconv.Atoi(input)
			if err != nil {
				qe.ui.Println("Please make the choice, 1,2...")
				continue
			}
			if choice < 1 || choice > len(q.Candidates) {
				qe.ui.Println("choice out of range.")
				continue
			}
			answer = q.Candidates[choice-1].Value
			break
		}
	}
	if answer == "" {
		answer = qe.ui.GetInput()
		if answer == "" && defaultValue != "" {
			answer = defaultValue
		}
	}
	return answer, nil
}
