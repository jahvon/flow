package expr_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jahvon/flow/internal/services/expr"
)

func TestExpr(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Expr Suite")
}

var _ = Describe("Expr", func() {
	Describe("IsTruthy", func() {
		It("should evaluate truthy expressions correctly", func() {
			tests := []struct {
				expr     string
				env      *expr.ExpressionData
				expected bool
			}{
				{"true", nil, true},
				{"false", nil, false},
				{"1", nil, true},
				{"0", nil, false},
				{`"true"`, nil, true},
				{`"false"`, nil, false},
			}

			for _, test := range tests {
				result, err := expr.IsTruthy(test.expr, test.env)
				Expect(err).NotTo(HaveOccurred())
				By("testing expression: " + test.expr)
				Expect(result).To(Equal(test.expected))
			}
		})
	})

	Describe("Evaluate", func() {
		It("should evaluate expressions correctly", func() {
			tests := []struct {
				expr     string
				env      *expr.ExpressionData
				expected interface{}
			}{
				{"1 + 1", nil, 2},
				{"true && false", nil, false},
				{`"hello" + " " + "world"`, nil, "hello world"},
			}

			for _, test := range tests {
				result, err := expr.Evaluate(test.expr, test.env)
				Expect(err).NotTo(HaveOccurred())
				By("testing expression: " + test.expr)
				Expect(result).To(Equal(test.expected))
			}
		})
	})

	Describe("EvaluateString", func() {
		It("should evaluate string expressions correctly", func() {
			tests := []struct {
				expr     string
				env      *expr.ExpressionData
				expected string
			}{
				{`"hello"`, nil, "hello"},
				{`"foo" + "bar"`, nil, "foobar"},
			}

			for _, test := range tests {
				result, err := expr.EvaluateString(test.expr, test.env)
				Expect(err).NotTo(HaveOccurred())
				By("testing expression: " + test.expr)
				Expect(result).To(Equal(test.expected))
			}
		})
	})
})
