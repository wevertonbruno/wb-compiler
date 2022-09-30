package parser

import (
	"fmt"
	"github.com/wevertonbruno/wb-compiler/analyzers/lexer"
	"github.com/wevertonbruno/wb-compiler/analyzers/token"
	"github.com/wevertonbruno/wb-compiler/ast"
	"strconv"
)

const (
	_ byte = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL

	parserError   = "parser error. %v"
	expectedError = "expected %v, got %v"
)

var (
	//precedence table
	precedences = map[token.Kind]byte{
		token.OP_EQ:     EQUALS,
		token.OP_NOTEQ:  EQUALS,
		token.OP_LT:     LESSGREATER,
		token.OP_GT:     LESSGREATER,
		token.OP_LTE:    LESSGREATER,
		token.OP_GTE:    LESSGREATER,
		token.OP_PLUS:   SUM,
		token.OP_MINUS:  SUM,
		token.OP_DIVIDE: PRODUCT,
		token.OP_MULTI:  PRODUCT,
	}
)

type (
	Parser struct {
		lexer        *lexer.Lexer
		currentToken token.Token
		peekToken    token.Token

		prefixParseFn map[token.Kind]prefixParseFn
		infixParseFn  map[token.Kind]infixParseFn
	}

	prefixParseFn func() ast.Expr
	infixParseFn  func(ast.Expr) ast.Expr
)

func NewParser(lex *lexer.Lexer) *Parser {
	parser := &Parser{
		lexer:        lex,
		currentToken: lex.GetToken(),
		peekToken:    lex.GetToken(),
	}
	parser.prefixParseFn = make(map[token.Kind]prefixParseFn)
	parser.registerPrefix(token.IDENTIFIER, parser.parseIdentifier)
	parser.registerPrefix(token.INTLIT, parser.parseIntegerLiteral)
	parser.registerPrefix(token.DECIMALLIT, parser.parseDecimalLiteral)
	parser.registerPrefix(token.TRUE, parser.parseBoolean)
	parser.registerPrefix(token.FALSE, parser.parseBoolean)
	parser.registerPrefix(token.NOT, parser.parsePrefixExpr)
	parser.registerPrefix(token.OP_MINUS, parser.parsePrefixExpr)
	parser.registerPrefix(token.L_BRACKET, parser.parseGroupedExpr)
	parser.registerPrefix(token.IF, parser.parseIfExpr)

	parser.infixParseFn = make(map[token.Kind]infixParseFn)
	parser.registerInfix(token.OP_PLUS, parser.parseInfixExpr)
	parser.registerInfix(token.OP_MINUS, parser.parseInfixExpr)
	parser.registerInfix(token.OP_MULTI, parser.parseInfixExpr)
	parser.registerInfix(token.OP_DIVIDE, parser.parseInfixExpr)
	parser.registerInfix(token.OP_EQ, parser.parseInfixExpr)
	parser.registerInfix(token.OP_NOTEQ, parser.parseInfixExpr)
	parser.registerInfix(token.OP_LT, parser.parseInfixExpr)
	parser.registerInfix(token.OP_LTE, parser.parseInfixExpr)
	parser.registerInfix(token.OP_GT, parser.parseInfixExpr)
	parser.registerInfix(token.OP_GTE, parser.parseInfixExpr)

	return parser
}

func (p *Parser) nextToken(ignoreNewLine bool) {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.GetToken()
	if ignoreNewLine && p.currentToken.Match(token.NEWLINE) {
		for p.currentToken.Match(token.NEWLINE) {
			p.currentToken = p.peekToken
			p.peekToken = p.lexer.GetToken()
		}
	}
}

func (p *Parser) match(kinds ...token.Kind) {
	var errMessage string
	for i, k := range kinds {
		if p.currentToken.Match(k) {
			p.nextToken(false)
			return
		}
		errMessage += k.Name()
		if i < len(kinds) {
			errMessage += " or "
		}
	}

	p.abort(fmt.Sprintf(expectedError, errMessage, p.currentToken.Kind.Name()))
}

func (p *Parser) expectedPeek(kind token.Kind) {
	if p.peekToken.Match(kind) {
		p.nextToken(false)
	} else {
		p.abort(fmt.Sprintf(expectedError, kind.Name(), p.currentToken.Kind.Name()))
	}
}

func (p *Parser) check(k token.Kind) bool {
	return p.currentToken.Match(k)
}

func (p *Parser) checkPeek(k token.Kind) bool {
	return p.peekToken.Match(k)
}

func (p *Parser) abort(message string) error {
	panic(fmt.Errorf(parserError, message))
}

func (p *Parser) registerPrefix(k token.Kind, fn prefixParseFn) {
	p.prefixParseFn[k] = fn
}

func (p *Parser) registerInfix(k token.Kind, fn infixParseFn) {
	p.infixParseFn[k] = fn
}

func (p *Parser) peekPrecedence() byte {
	if p, ok := precedences[p.peekToken.Kind]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) currentPrecedence() byte {
	if p, ok := precedences[p.currentToken.Kind]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) Parse() *ast.Prog {
	prog := newProgNode()
	p.ignoreNewLines()
	for !p.currentToken.Match(token.EOF) {
		stmt := p.parseBlockItem()
		if stmt != nil {
			prog.BlockItems = append(prog.BlockItems, stmt)
		}
		p.nextToken(true)
	}
	fmt.Println(prog)
	return prog
}

func newProgNode() *ast.Prog {
	return &ast.Prog{
		BlockItems: []ast.BlockItem{},
	}
}

func (p *Parser) parseBlockItem() ast.BlockItem {
	switch p.currentToken.Kind {
	case token.VAR:
		return p.parseVarStatement()
	default:
		return p.parseStatement()
	}
}

func (p *Parser) parseStatement() ast.Stmt {
	switch p.currentToken.Kind {
	case token.RETURN:
		return p.parseReturnStatement()
	case token.L_BRACE:
		return p.parseBlockStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseNewLine() {
	p.match(token.NEWLINE)
	for p.check(token.NEWLINE) {
		p.nextToken(false)
	}
}

func (p *Parser) advanceIgnoringNewLines() {
	if !p.check(token.NEWLINE) {
		for p.checkPeek(token.NEWLINE) {
			p.nextToken(false)
		}
	}
}

func (p *Parser) ignoreNewLines() {
	for p.check(token.NEWLINE) {
		p.nextToken(false)
	}
}

func (p *Parser) checkSeparator() bool {
	return p.currentToken.Match(token.SEMICOLON) ||
		p.currentToken.Match(token.NEWLINE)
}

func (p *Parser) peekSeparator() bool {
	return p.peekToken.Match(token.SEMICOLON) ||
		p.peekToken.Match(token.NEWLINE)
}

func (p *Parser) parseVarStatement() *ast.Declaration {
	stmt := &ast.Declaration{Token: p.currentToken}
	p.expectedPeek(token.IDENTIFIER)
	stmt.ID = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Spelling}
	p.expectedPeek(token.ASSIGN)
	p.nextToken(false)
	stmt.Value = p.parseExpression(LOWEST)
	if p.checkSeparator() {
		p.nextToken(false)
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	ret := &ast.ReturnStatement{
		Token: p.currentToken,
	}
	for !p.checkSeparator() {
		p.nextToken(false)
	}

	return ret
}

func (p *Parser) parseExpressionStatement() *ast.ExprStatement {
	expr := &ast.ExprStatement{Token: p.currentToken}
	expr.Expr = p.parseExpression(LOWEST)
	if p.peekSeparator() {
		p.nextToken(false)
	}
	return expr
}

func (p *Parser) parseExpression(precedence byte) ast.Expr {
	prefix := p.prefixParseFn[p.currentToken.Kind]
	if prefix == nil {
		p.abort(fmt.Sprintf("no prefix parse function for %s found", p.currentToken.Spelling))
		return nil
	}
	leftExpr := prefix()

	for !p.peekSeparator() && precedence < p.peekPrecedence() {
		infix := p.infixParseFn[p.peekToken.Kind]
		if infix == nil {
			return leftExpr
		}
		p.nextToken(false)
		leftExpr = infix(leftExpr)
	}
	return leftExpr
}

func (p *Parser) parseIdentifier() ast.Expr {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Spelling}
}

func (p *Parser) parseBoolean() ast.Expr {
	return &ast.Boolean{Token: p.currentToken, Value: p.check(token.TRUE)}
}

func (p *Parser) parseIntegerLiteral() ast.Expr {
	lit := &ast.IntegerLiteral{Token: p.currentToken}

	value, err := strconv.ParseInt(p.currentToken.Spelling, 0, 64)
	if err != nil {
		p.abort(fmt.Sprintf("could not parse %s as Integer", p.currentToken.Spelling))
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseDecimalLiteral() ast.Expr {
	lit := &ast.DecimalLiteral{Token: p.currentToken}

	value, err := strconv.ParseFloat(p.currentToken.Spelling, 0)
	if err != nil {
		p.abort(fmt.Sprintf("could not parse %s as Float", p.currentToken.Spelling))
	}
	lit.Value = value
	return lit
}

func (p *Parser) parsePrefixExpr() ast.Expr {
	expr := &ast.PrefixExpression{Token: p.currentToken}
	p.nextToken(false)
	expr.Right = p.parseExpression(PREFIX)
	return expr
}

func (p *Parser) parseInfixExpr(left ast.Expr) ast.Expr {
	expr := &ast.InfixExpression{
		Token: p.currentToken,
		Left:  left,
	}
	precedence := p.currentPrecedence()
	p.nextToken(false)
	expr.Right = p.parseExpression(precedence)

	return expr
}

func (p *Parser) parseGroupedExpr() ast.Expr {
	p.nextToken(false)
	expr := p.parseExpression(LOWEST)
	p.expectedPeek(token.R_BRACKET)
	p.nextToken(false)
	return expr
}

func (p *Parser) parseIfExpr() ast.Expr {
	ifExpr := &ast.IfExpression{Token: p.currentToken}
	p.nextToken(false)

	ifExpr.Condition = p.parseExpression(LOWEST)
	ifExpr.TrueBlockCondition = p.parseStatement()
	if p.checkPeek(token.ELSE) {
		p.nextToken(false)
		ifExpr.FalseBlockCondition = p.parseBlockStatement()
	}

	return ifExpr
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.currentToken}
	block.Statements = []ast.BlockItem{}
	p.nextToken(true)
	for !p.check(token.R_BRACE) && !p.check(token.EOF) {
		stmt := p.parseBlockItem()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken(true)
	}

	return block
}
