%{
package minifp

import (
        "text/scanner"
        )

func lexpos(yylex yyLexer) scanner.Position {
  return yylex.(*parser).sc.Pos()
}

%}

%union {
  astlist []ASTNode
  ast ASTNode
  assign *ASTAssign
  assignlist []*ASTAssign
  arglist []string
  ident string
}

%start main

%token <ident> tokIdent
%token <ident> tokLetrec tokIn
%token <ast> tokLiteral
%token <ident> tokArrow

%type<astlist> main toplevelExprList
%type<ast> expr toplevelExpr
%type<assign> binding
%type<assignlist> bindingList
%type<arglist> arglist

%%

main: toplevelExprList { yylex.(*parser).result = $1 }

toplevelExprList: toplevelExpr { $$ = []ASTNode{$1} }
  | toplevelExprList ';' toplevelExpr { $$ = append($1, $3)}

toplevelExpr: tokIdent arglist '=' expr {
    rhs := $4
    if len($2) > 0 {
      rhs = newLambda(lexpos(yylex), $2, $4)
    }
    $$ = &ASTAssign{pos: lexpos(yylex), Sym: InternSymbol($1), Expr: rhs}
  }
  | expr { $$ = $1 }

arglist: { $$ = nil }
  | arglist tokIdent { $$ = append($1, $2) }

expr: tokLiteral
  | expr expr { $$ = &ASTApply{pos: lexpos(yylex), Head:$1, Tail:$2} }
  | tokIdent { $$ = &ASTVar{pos: lexpos(yylex), Sym: InternSymbol($1)} }
  | expr '+' expr { $$ = &ASTApplyBuiltin{pos: lexpos(yylex), Op: BuiltinOpAdd, Args: []ASTNode{$1, $3} } }
  | expr '-' expr { $$ = &ASTApplyBuiltin{pos: lexpos(yylex), Op: BuiltinOpSub, Args: []ASTNode{$1, $3} } }
  | expr '*' expr { $$ = &ASTApplyBuiltin{pos: lexpos(yylex), Op: BuiltinOpMul, Args: []ASTNode{$1, $3} } }
  | '(' expr ')' { $$ = $2 }
  | '\\' arglist tokArrow expr { $$ = newLambda(lexpos(yylex), $2, $4) }
  | tokLetrec bindingList tokIn expr { $$ = &ASTLetrec{pos: lexpos(yylex), Bindings: $2, Body: $4} }

bindingList:
  binding { $$ = []*ASTAssign{$1} }
  | bindingList ';' binding { $$ = append($1, $3) }

binding: tokIdent '=' expr {$$ = &ASTAssign{pos: lexpos(yylex), Sym: InternSymbol($1), Expr: $3}}
