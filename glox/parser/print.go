package parser

import (
	"fmt"
	"strings"

	"github.com/Martin-Martinez4/crafting-interpreters/glox/errorhandling"
)

type AstPrinter struct {
}

func (astp *AstPrinter) Print(expr Expr) string {

	s, ok := expr.Accept(astp).(string)
	if !ok {
		errorhandling.ReportAndExit(0, "", "could not create string from expr")
	}
	return s
}

func (astp *AstPrinter) VisitBinary(expr *Binary) any {
	return astp.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}
func (astp *AstPrinter) VisitGrouping(expr *Grouping) any {
	return astp.parenthesize("group", expr.Expression)
}
func (astp *AstPrinter) VisitLiteral(expr *Literal) any {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%+v", expr.Value)
}
func (astp *AstPrinter) VisitUnary(expr *Unary) any {
	return astp.parenthesize(expr.Operator.Lexeme, expr.Right)
}
func (astp *AstPrinter) VisitVariable(expr *Variable) any {
	return nil
}
func (astp *AstPrinter) VisitAssign(expr *Assign) any {
	return nil
}
func (astp *AstPrinter) VisitLogical(expr *Logical) any {
	return nil
}
func (astp *AstPrinter) VisitCall(expr *CallExpr) any {
	return nil
}
func (astp *AstPrinter) VisitGet(expr *Get) any {
	return nil
}
func (astp *AstPrinter) VisitSet(expr *Set) any {
	return nil
}

func (astp *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	var ss strings.Builder

	ss.WriteString("(")
	ss.WriteString(name)
	for i := 0; i < len(exprs); i++ {
		ss.WriteString(" ")
		s, ok := exprs[i].Accept(astp).(string)
		if !ok {
			errorhandling.ReportAndExit(0, "parenthesize", "could not create string from result of Accept method")
		}
		ss.WriteString(s)
	}
	ss.WriteString(")")
	return ss.String()
}
