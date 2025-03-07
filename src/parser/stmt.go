package parser

import (
	"github.com/pvarma-05/MowaLang/src/ast"
	"github.com/pvarma-05/MowaLang/src/lexer"
)

func parse_stmt(p *parser) ast.Stmt {
	stmt_fn, exists := stmt_lu[p.currentTokenKind()]

	if exists {
		return stmt_fn(p)
	}
	expression := parse_expr(p, default_bp)
	p.expect(lexer.SEMI_COLON)

	return ast.ExprStmt{
		Expression: expression,
	}
}
