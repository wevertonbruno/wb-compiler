package parser

import (
	"fmt"
	"github.com/wevertonbruno/wb-compiler/analyzers/lexer"
	"github.com/wevertonbruno/wb-compiler/analyzers/reader"
	"github.com/wevertonbruno/wb-compiler/ast"
	"strconv"
	"testing"
)

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5\n", 5, "+", 5},
		{"5 - 5\n", 5, "-", 5},
		{"5 * 5\n", 5, "*", 5},
		{"5 / 5\n", 5, "/", 5},

		{"5.5 + 7.8\n", 5.5, "+", 7.8},
		{"0.4 - 5.4\n", 0.4, "-", 5.4},
		{"5.5 * 5.1\n", 5.5, "*", 5.1},
		{"5.1 / 5.1\n", 5.1, "/", 5.1},

		{"5 > 5\n", 5, ">", 5},
		{"5 < 5\n", 5, "<", 5},
		{"5 == 5\n", 5, "==", 5},
		{"5 != 5\n", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		r := reader.NewInput(tt.input)
		l := lexer.NewLexer(r)
		p := NewParser(l)
		program := p.Parse()
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExprStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}
		exp, ok := stmt.Expr.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expr)
		}
		if !testLiteral(t, exp.Left, tt.leftValue) {
			return
		}
		if exp.Token.Spelling != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Token.Spelling)
		}
		if !testLiteral(t, exp.Right, tt.rightValue) {
			return
		}
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5\n", "!", 5},
		{"-15\n", "-", 15},
	}
	for _, tt := range prefixTests {
		r := reader.NewInput(tt.input)
		l := lexer.NewLexer(r)
		p := NewParser(l)
		program := p.Parse()
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExprStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}
		exp, ok := stmt.Expr.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expr)
		}
		if exp.Token.Spelling != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Token.Spelling)
		}

		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"a + b * c", "(a + (b * c))"},
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
	}
	for _, tt := range tests {
		r := reader.NewInput(tt.input)
		l := lexer.NewLexer(r)
		p := NewParser(l)
		program := p.Parse()
		if actual := program.String(); actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := []string{
		`if (x < y){
			x 
		} else { 
			y 
		}`,
		`if (x < y){
			x 
		} else {
			if (x + y) { 
				y 
			} else {
				z
			}
		}`,
	}
	for _, in := range input {
		r := reader.NewInput(in)
		l := lexer.NewLexer(r)
		p := NewParser(l)
		program := p.Parse()
		if len(program.Statements) != 1 {
			t.Fatalf("program.Body does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExprStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}
		exp, ok := stmt.Expr.(*ast.IfExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T",
				stmt.Expr)
		}
		if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
			return
		}
		if len(exp.TrueBlockCondition.Statements) != 1 {
			t.Errorf("consequence is not 1 statements. got=%d\n",
				len(exp.TrueBlockCondition.Statements))
		}
		_, ok = exp.TrueBlockCondition.Statements[0].(*ast.ExprStatement)
		if !ok {
			t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
				exp.TrueBlockCondition.Statements[0])
		}
		if exp.FalseBlockCondition != nil {
			_, ok := exp.FalseBlockCondition.Statements[0].(*ast.ExprStatement)
			if !ok {
				t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
					exp.FalseBlockCondition.Statements[0])
			}
		}
	}
}

func TestDeclarationParsing(t *testing.T) {
	input := []struct {
		code  string
		id    string
		value interface{}
	}{
		{`var x = 1`, "x", 1},
		{`var x = true`, "x", true},
		{`var x = 1.5`, "x", 1.5},
	}
	for _, in := range input {
		r := reader.NewInput(in.code)
		l := lexer.NewLexer(r)
		p := NewParser(l)
		program := p.Parse()
		if len(program.Statements) != 1 {
			t.Fatalf("program.Body does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}
		decl, ok := program.Statements[0].(*ast.DeclStatement)
		if !ok {
			t.Fatalf("Statement isn't Declaration. got=%T\n",
				program.Statements[0])
		}
		if !testLiteral(t, decl.Value, in.value) {
			return
		}

		if !testLiteral(t, decl.ID, in.id) {
			return
		}
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `func (
					x, 
					y
				) { 
				x + y; 
				}`
	r := reader.NewInput(input)
	l := lexer.NewLexer(r)
	p := NewParser(l)
	program := p.Parse()
	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExprStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	function, ok := stmt.Expr.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T",
			stmt.Expr)
	}
	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n",
			len(function.Parameters))
	}
	testLiteral(t, function.Parameters[0], "x")
	testLiteral(t, function.Parameters[1], "y")
	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n",
			len(function.Body.Statements))
	}
	bodyStmt, ok := function.Body.Statements[0].(*ast.ExprStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T",
			function.Body.Statements[0])
	}
	testInfixExpression(t, bodyStmt.Expr, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "func() {};", expectedParams: []string{}},
		{input: "func(x) {};", expectedParams: []string{"x"}},
		{input: "func(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}
	for _, tt := range tests {
		r := reader.NewInput(tt.input)
		l := lexer.NewLexer(r)
		p := NewParser(l)
		program := p.Parse()
		stmt := program.Statements[0].(*ast.ExprStatement)
		function := stmt.Expr.(*ast.FunctionLiteral)
		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d\n",
				len(tt.expectedParams), len(function.Parameters))
		}
		for i, ident := range tt.expectedParams {
			testLiteral(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"
	r := reader.NewInput(input)
	l := lexer.NewLexer(r)
	p := NewParser(l)
	program := p.Parse()
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExprStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	exp, ok := stmt.Expr.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
			stmt.Expr)
	}
	if !testIdentifier(t, exp.Function, "add") {
		return
	}
	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}
	testLiteral(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func testLiteral(t *testing.T, il ast.Expr, value interface{}) bool {
	switch v := value.(type) {
	case bool:
		return testBooleanLiteral(t, il, v)
	case float64:
		return testDecimalLiteral(t, il, v)
	case int64:
		return testIntegerLiteral(t, il, v)
	case int:
		return testIntegerLiteral(t, il, int64(v))
	case string:
		return testIdentifier(t, il, v)
	}
	t.Errorf("type of exp not handled. got=%T", il)
	return false
}

func testIntegerLiteral(t *testing.T, il ast.Expr, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}
	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}
	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value,
			integ.TokenLiteral())
		return false
	}
	return true
}

func testDecimalLiteral(t *testing.T, il ast.Expr, value float64) bool {
	integ, ok := il.(*ast.DecimalLiteral)
	if !ok {
		t.Errorf("il not *ast.DecimalLiteral. got=%T", il)
		return false
	}
	if integ.Value != value {
		t.Errorf("expr.Value not %f. got=%f", value, integ.Value)
		return false
	}
	if integ.TokenLiteral() != fmt.Sprintf("%s", strconv.FormatFloat(value, 'f', -1, 64)) {
		t.Errorf("expr.TokenLiteral not %f. got=%s", value,
			integ.TokenLiteral())
		return false
	}
	return true
}

func testBooleanLiteral(t *testing.T, il ast.Expr, value bool) bool {
	integ, ok := il.(*ast.Boolean)
	if !ok {
		t.Errorf("il not *ast.Boolean. got=%T", il)
		return false
	}
	if integ.Value != value {
		t.Errorf("expr.Value not %v. got=%v", value, integ.Value)
		return false
	}
	if integ.TokenLiteral() != fmt.Sprintf("%v", value) {
		t.Errorf("expr.TokenLiteral not %v. got=%s", value,
			integ.TokenLiteral())
		return false
	}
	return true
}

func testInfixExpression(t *testing.T, exp ast.Expr, left interface{},
	operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.OperatorExpression. got=%T(%s)", exp, exp)
		return false
	}
	if !testLiteral(t, opExp.Left, left) {
		return false
	}
	if opExp.Token.Kind.Name() != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Token.Kind.Name())
		return false
	}
	if !testLiteral(t, opExp.Right, right) {
		return false
	}
	return true
}

func testIdentifier(t *testing.T, exp ast.Expr, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}
	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value,
			ident.TokenLiteral())
		return false
	}
	return true
}
