package gquiz_test

import (
	"io/ioutil"

	"github.com/andyliuliming/gquiz"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type MockUI struct {
	mockInputs   []string
	currentIndex int
}

func NewMockUI(mockInputs []string) *MockUI {
	return &MockUI{
		mockInputs:   mockInputs,
		currentIndex: 0,
	}
}
func (l *MockUI) Println(message string) {
	println(message)
}

func (l *MockUI) GetInput() string {
	result := l.mockInputs[l.currentIndex]
	l.currentIndex += 1
	return result
}

var _ = Describe("QuizExecutor", func() {
	Describe("Execute", func() {
		var (
			qGraph gquiz.QGraph
		)
		BeforeEach(func() {
			content, err := ioutil.ReadFile("./sample.yml")
			Expect(err).To(BeNil())
			quizBuilder := gquiz.QuizBuilder{}
			qGraph, err = quizBuilder.BuildQGraph(content)
			Expect(err).To(BeNil())
		})
		Context("Everything OK", func() {
			It("should execute sucessfully", func() {
				mockUI := NewMockUI(
					[]string{
						"", "", "",
						"1",
						"1",
						"1",
						"",
					},
				)
				quizExecutor := gquiz.NewQuizExecutor(mockUI, nil)
				qResult, err := quizExecutor.Execute(&qGraph)
				Expect(err).To(BeNil())
				Expect((*qResult)["admin_name"]).To(Equal("kluser"))
				Expect(qResult).NotTo(BeNil())
			})
		})

		Context("When have old values.", func() {
			It("old values should been used.", func() {
				adminName := "kluser2"
				qr := &gquiz.QResult{"admin_name": adminName}

				mockUI := NewMockUI(
					[]string{
						"", "", "",
						"1",
						"1",
						"1",
						"",
					},
				)
				quizExecutor := gquiz.NewQuizExecutor(mockUI, qr)
				qResult, err := quizExecutor.Execute(&qGraph)
				Expect(err).To(BeNil())
				Expect((*qResult)["admin_name"]).To(Equal("kluser2"))
				Expect(qResult).NotTo(BeNil())
			})
		})
	})
})
