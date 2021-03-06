package gquiz_test

import (
	"io/ioutil"
	"os"

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
						"", "", "", "",
						"1",
						"1",
						"1",
						"",
					},
				)
				quizExecutor := gquiz.NewQuizExecutor(mockUI, nil)
				qResult, err := quizExecutor.Execute(&qGraph)
				Expect(err).To(BeNil())
				Expect(qResult["admin_name"]).To(Equal("kluser"))
				Expect(qResult).NotTo(BeNil())
			})
		})

		Context("When have old values.", func() {
			It("old values should been used.", func() {
				adminName := "kluser2"
				qr := gquiz.QResult{"admin_name": adminName}

				mockUI := NewMockUI(
					[]string{
						"", "", "", "",
						"1",
						"1",
						"1",
						"",
					},
				)
				quizExecutor := gquiz.NewQuizExecutor(mockUI, qr)
				qResult, err := quizExecutor.Execute(&qGraph)
				Expect(err).To(BeNil())
				Expect(qResult["admin_name"]).To(Equal("kluser2"))
				Expect(qResult).NotTo(BeNil())
			})
		})

		Context("When have default value for choice.", func() {
			It("default value for choice should been used.", func() {
				adminName := "kluser2"
				qr := gquiz.QResult{"admin_name": adminName}
				repoAddress := "github.com/Microsoft/kunlun"
				mockUI := NewMockUI(
					[]string{
						"", "", "", "", "",
						repoAddress,
						"",
						"1",
						"20",
						"",
						"",
					},
				)
				quizExecutor := gquiz.NewQuizExecutor(mockUI, qr)
				qResult, err := quizExecutor.Execute(&qGraph)
				Expect(err).To(BeNil())
				Expect(qResult["project_source_code_path"]).To(Equal(repoAddress))
				Expect(qResult).NotTo(BeNil())
			})
		})

		Context("When constant value set", func() {
			It("constant value should be used", func() {
				mockUI := NewMockUI(
					[]string{
						"", "", "", "",
						"1",
						"1",
						"1",
					},
				)
				quizExecutor := gquiz.NewQuizExecutor(mockUI, nil)
				qResult, err := quizExecutor.Execute(&qGraph)
				Expect(err).To(BeNil())
				Expect(qResult["final_artifact"]).To(Equal("small_kubernetes.yml"))
				Expect(qResult).NotTo(BeNil())
			})
		})

		Context("When the environment is set.", func() {
			It("value in env should been used.", func() {
				iaas := "other_iaas"
				os.Setenv("KL_IAAS", iaas)
				mockUI := NewMockUI(
					[]string{
						"", "", "", "",
						"1",
						"1",
						"1",
						"",
					},
				)
				quizExecutor := gquiz.NewQuizExecutor(mockUI, nil)
				qResult, err := quizExecutor.Execute(&qGraph)
				Expect(err).To(BeNil())
				Expect(qResult["iaas"]).To(Equal(iaas))
				Expect(qResult).NotTo(BeNil())
			})
		})
	})
})
