%{

package parser

import (
    "text/scanner"
    "strings"
    "io/ioutil"
    "fmt"
    "errors"
    "strconv"

    "../nodes"
    "../types"
)

// TODO: fix nodes positions
// TODO: somehow add comments

%}

%union {
    script *nodes.Script

    pos scanner.Position

    nodes []nodes.Node
    node nodes.Node
    text string

    _type types.Type
    value types.Value

    integerValue struct{
        Value int32
        IsHex bool
    }
    floatValue float64
    stringValue string

    integer int

    comment *nodes.Comment
    identifier *nodes.Identifier
    constant *nodes.Constant

    variable *nodes.Variable
    variables []*nodes.Variable

    function *nodes.Function
    functions []*nodes.Function

    state *nodes.State
    states []*nodes.State

    blockStatement *nodes.BlockStatement
    statement nodes.Statement
    statements []nodes.Statement

    expression nodes.Expression
    expressions []nodes.Expression

    listItem *nodes.ListItemExpression

    _struct *nodes.Struct
}

%token <text> COMMENT C_STYLE_COMMENT
%token INTEGER FLOAT STRING KEY VECTOR ROTATION LIST

%token <integerValue> INTEGER_CONSTANT
%token <floatValue> FLOAT_CONSTANT
%token <stringValue> STRING_CONSTANT

%token <text> IDENTIFIER

%token <pos> DEFAULT STATE
%token <pos> JUMP RETURN IF ELSE FOR DO WHILE
%token <pos> INC_OP DEC_OP
%token ADD_ASSIGN SUB_ASSIGN MUL_ASSIGN DIV_ASSIGN MOD_ASSIGN
%token EQ NEQ LEQ GEQ
%token BOOLEAN_AND BOOLEAN_OR
%token SHIFT_LEFT SHIFT_RIGHT
%token <pos> ARRAY_BRACES

%token BOOLEAN
%token <pos> TRUE FALSE

%token <pos> INCLUDE PRAGMA

%token <pos> STRUCT
%token <pos> SWITCH CASE BREAK CONTINUE
%token <pos> CONST
%token <pos> DELETE

%token <pos> '{' '@' '-' '!' '~' '(' '<' '[' '#'

%nonassoc LOWER_THAN_ELSE
%nonassoc ELSE

%nonassoc INTEGER_CONSTANT FLOAT_CONSTANT
%right '=' MUL_ASSIGN DIV_ASSIGN MOD_ASSIGN ADD_ASSIGN SUB_ASSIGN
%left 	BOOLEAN_AND BOOLEAN_OR
%left	'|'
%left	'^'
%left	'&'
%left	EQ NEQ
%left	'<' LEQ '>' GEQ
%left	SHIFT_LEFT SHIFT_RIGHT
%left 	'+' '-'
%left	'*' '/' '%'
%right	'!' '~' INC_OP DEC_OP
%right LIST_ITEM_PREC '['
%nonassoc INITIALIZER


%type <script> lscript_program
%type <nodes> globals

%type <nodes> preproc

%type <nodes> global
%type <comment> comment

%type <identifier> identifier
%type <_type> typename
%type <variables> variable_declaration
%type <variables> variable_declarations
%type <function> function

%type <variables> function_arguments
%type <variable> function_argument

%type <constant> constant

%type <function> event
%type <functions> events

%type <states> states
%type <states> other_states
%type <state> default_state
%type <state> state

%type <blockStatement> block_statement
%type <statement> empty_statement
%type <statement> statement

%type <statements> statements_with_comments

%type <expressions> for_expression_list
%type <expressions> next_for_expression_list
%type <expression> lvalue_identifiers
%type <listItem> list_item_expression
%type <expression> lvalue
%type <expression> expression
%type <expression> unary_expression
%type <expression> typecast_expression
%type <expression> unary_postfix_expression
%type <expressions> list_values

%type <_struct> struct

%type <expression> struct_expression
%type <variable> struct_expression_variable
%type <variables> struct_expression_variables

%type <variable> pre_variable
%type <variables> pre_variables

%%

lscript_program:
    globals states {
        $$ = yylex.(*Lexer).script
        $$.Globals = append($$.Globals, $1...)
        for _, child := range $2 {
            $$.Globals = append($$.Globals, child)
        }
    }
|   globals {
        $$ = yylex.(*Lexer).script
        $$.Globals = append($$.Globals, $1...)
    }
|   states {
        $$ = yylex.(*Lexer).script
        for _, child := range $1 {
            $$.Globals = append($$.Globals, child)
        }
    }

globals:
    global {
        $$ = $1
    }
|   preproc {
        $$ = $1
    }
|   globals preproc {
        $$ = append($1, $2...)
    }
|   globals global {
        $$ = append($1, $2...)
    }

preproc:
    '#' INCLUDE STRING_CONSTANT {
        filename := strings.Trim($3, "\"")
        script, err := ParseFile(filename)
        if err != nil {
            yylex.Error(err.Error())
            Nerrs++
            goto ret1
        } else {
            comment := &nodes.Comment{
                Text: fmt.Sprintf("#include \"%s\"", filename),
            }
            comment.At = $1

            $$ = append([]nodes.Node{ comment }, script.Globals...)
        }
    }
|   '#' PRAGMA IDENTIFIER {
        switch $3 {
        case "skip_unused":
            yylex.(*Lexer).script.SkipUnused = true
        case "no_skip_unused":
            yylex.(*Lexer).script.SkipUnused = false
        default:
            yylex.Error("Invalid pragma " + $3)
        }
    }

global:
    comment {
        $$ = []nodes.Node{ $1 }
    }
|   variable_declaration {
        for _, child := range $1 {
            $$ = append($$, child)
        }
    }
|   function {
        $$ = []nodes.Node{ $1 }
    }
|   struct {
        $$ = []nodes.Node{ $1 }
    }
|   empty_statement {
        $$ = []nodes.Node{ $1 }
    }


comment:
    COMMENT {
        $$ = &nodes.Comment{
            Text: $1,
            IsCStyle: false,
        }
        $$.SetPosition(yylex.(*Lexer).LastPos)
    }

|   C_STYLE_COMMENT {
        $$ = &nodes.Comment{
            Text: $1,
            IsCStyle: true,
        }
        $$.SetPosition(yylex.(*Lexer).LastPos)
    }


identifier:
    IDENTIFIER {
        $$ = &nodes.Identifier{
            Name: $1,
        }
        $$.SetPosition(yylex.(*Lexer).LastPos)
    }

typename:
    INTEGER { $$ = types.Integer }
|   FLOAT { $$ = types.Float }
|   STRING { $$ = types.String }
|   KEY { $$ = types.Key }
|   VECTOR { $$ = types.Vector }
|   ROTATION { $$ = types.Rotation }
|   LIST { $$ = types.List }
|   BOOLEAN { $$ = types.Boolean }

constant:
    INTEGER_CONSTANT {
        $$ = &nodes.Constant{
            Value: &types.IntegerValue{
                At: yylex.(*Lexer).LastPos,
                Value: $1.Value,
                IsHex: $1.IsHex,
            },
        }
        $$.SetPosition(yylex.(*Lexer).LastPos)
    }
|   FLOAT_CONSTANT {
        $$ = &nodes.Constant{
            Value: &types.FloatValue{
                At: yylex.(*Lexer).LastPos,
                Value: $1,
            },
        }
        $$.SetPosition(yylex.(*Lexer).LastPos)
    }
|   STRING_CONSTANT {
        $$ = &nodes.Constant{
            Value: &types.StringValue{
                At: yylex.(*Lexer).LastPos,
                Value: $1,
            },
        }
        $$.SetPosition(yylex.(*Lexer).LastPos)
    }
|   TRUE {
        $$ = &nodes.Constant{
            Value: &types.BooleanValue{
                At: yylex.(*Lexer).LastPos,
                Value: true,
            },
        }
        $$.SetPosition(yylex.(*Lexer).LastPos)
    }
|   FALSE {
        $$ = &nodes.Constant{
            Value: &types.BooleanValue{
                At: yylex.(*Lexer).LastPos,
                Value: false,
            },
        }
        $$.SetPosition(yylex.(*Lexer).LastPos)
    }


pre_variable:
    identifier {
        $$ = &nodes.Variable{
            Name: $1,
        }
        $$.SetPosition($1.Position())
    }
|   identifier '=' expression {
        $$ = &nodes.Variable{
            Name: $1,
            RValue: $3,
        }
        $$.SetPosition($1.Position())
    }
|   identifier '{' '}' {
        $$ = &nodes.Variable{
            Name: $1,
            RValue: &nodes.StructExpression{
            },
        }
        $$.RValue.SetPosition($2)
        $$.SetPosition($1.Position())
    }
|   identifier '{' struct_expression_variables '}' {
        $$ = &nodes.Variable{
            Name: $1,
            RValue: &nodes.StructExpression{
                Fields: $3,
            },
        }
        $$.RValue.SetPosition($2)
        $$.SetPosition($1.Position())
    }

pre_variables:
    pre_variable {
        $$ = append($$, $1)
    }
|   pre_variables ',' pre_variable {
        $$ = append($1, $3)
    }

variable_declaration:
    typename pre_variables ';' {
        for _, child := range $2 {
            child.Type = $1
            $$ = append($$, child)
        }
    }
|   identifier pre_variables ';' {
        for _, child := range $2 {
            child.Type = types.Type($1.Name)
            if child.RValue == nil {
                child.RValue = &nodes.StructExpression{
                    Name: $1,
                }
                child.RValue.SetPosition(child.Position())
            }
            $$ = append($$, child)
        }
    }
|   identifier ARRAY_BRACES pre_variables ';' {
        for _, child := range $3 {
            child.Type = types.Type($1.Name + "[]")
            $$ = append($$, child)
        }
    }
|   CONST pre_variables ';' {
        for _, child := range $2 {
            child.IsConstant = true
            $$ = append($$, child)
        }
    }

variable_declarations:
    variable_declaration {
        $$ = $1
    }
|   variable_declarations variable_declaration {
        $$ = append($1, $2...)
    }

struct:
    STRUCT identifier '{' variable_declarations '}' {
        $$ = &nodes.Struct{
            Name: $2,
            Fields: $4,
        }
        $$.SetPosition($1)
    }

function:
    identifier '(' ')' block_statement {
        $$ = &nodes.Function{
            Name: $1,
            Type: types.Void,
            Body: $4,
        }
        $$.SetPosition($1.Position())
    }
|   typename identifier '(' ')' block_statement {
        $$ = &nodes.Function{
            Name: $2,
            Type: $1,
            Body: $5,
        }
        $$.SetPosition($2.Position())
    }
|   identifier '(' function_arguments ')' block_statement {
        $$ = &nodes.Function{
            Name: $1,
            Type: types.Void,
            Arguments: $3,
            Body: $5,
        }
        $$.SetPosition($1.Position())
    }
|   typename identifier '(' function_arguments ')' block_statement {
        $$ = &nodes.Function{
            Name: $2,
            Type: $1,
            Arguments: $4,
            Body: $6,
        }
        $$.SetPosition($2.Position())
    }

function_arguments:
    function_argument {
        $$ = append($$, $1)
    }
|   function_argument ',' function_arguments {
        $$ = append($$, $1)
        $$ = append($$, $3...)
    }

function_argument:
    typename identifier {
        $$ = &nodes.Variable{
            Name: $2,
            Type: $1,
            IsArgument: true,
        }
        $$.SetPosition($2.Position())
    }
|   identifier identifier {
        $$ = &nodes.Variable{
            Name: $2,
            Type: types.Type($1.Name),
            IsArgument: true,
        }
        $$.SetPosition($2.Position())
    }


event:
    identifier '(' ')' block_statement {
        $$ = &nodes.Function{
            Name: $1,
            Type: types.Void,
            Body: $4,
        }
        $$.SetPosition($1.Position())
    }
|   identifier '(' function_arguments ')' block_statement {
        $$ = &nodes.Function{
            Name: $1,
            Type: types.Void,
            Arguments: $3,
            Body: $5,
        }
        $$.SetPosition($1.Position())
    }

events:
    event {
        $$ = append($$, $1)
    }
|   event events {
        $$ = append($$, $1)
        $$ = append($$, $2...)
    }

states:
    default_state {
        $$ = append($$, $1)
    }
|   default_state other_states {
        $$ = append($$, $1)
        $$ = append($$, $2...)
    }

other_states:
    state {
        $$ = append($$, $1)
    }
|   state other_states {
        $$ = append($$, $1)
        $$ = append($$, $2...)
    }

default_state:
    DEFAULT '{' events '}' {
        $$ = &nodes.State{
            Events: $3,
        }
        $$.SetPosition($1)
    }
|   DEFAULT '{' '}' {
        $$ = &nodes.State{
        }
        $$.SetPosition($1)
    }

state:
    STATE identifier '{' events '}' {
        $$ = &nodes.State{
            Name: $2,
            Events: $4,
        }
        $$.SetPosition($1)
    }
|   STATE identifier '{' '}' {
        $$ = &nodes.State{
            Name: $2,
        }
        $$.SetPosition($1)
    }

block_statement:
    '{' '}' {
        $$ = &nodes.BlockStatement{
        }
        $$.SetPosition($1)
    }
|   '{' statements_with_comments '}' {
        $$ = &nodes.BlockStatement{
            Children: $2,
        }
        $$.SetPosition($1)
    }

statements_with_comments:
    statement {
        $$ = append($$, $1)
    }
|   variable_declaration {
        for _, child := range $1 {
            $$ = append($$, child)
        }
    }
|   comment {
        $$ = append($$, $1)
    }
|   statement statements_with_comments {
        $$ = append($$, $1)
        $$ = append($$, $2...)
    }
|   variable_declaration statements_with_comments {
        for _, child := range $1 {
            $$ = append($$, child)
        }
        $$ = append($$, $2...)
    }
|   comment statements_with_comments {
        $$ = append($$, $1)
        $$ = append($$, $2...)
    }

empty_statement:
    ';' {
        $$ = &nodes.EmptyStatement{
        }
        $$.SetPosition(yylex.(*Lexer).LastPos)
    }

statement:
    empty_statement {
        $$ = $1
    }
|   STATE identifier ';' {
        $$ = &nodes.StateStatement{
            Name: $2,
        }
        $$.SetPosition($1)
    }
|   STATE DEFAULT ';' {
        $$ = &nodes.StateStatement{
        }
        $$.SetPosition($1)
    }
|   '@' identifier ';' {
        $$ = &nodes.LabelStatement{
            Name: $2,
        }
        $$.SetPosition($1)
    }
|   JUMP identifier ';' {
        $$ = &nodes.JumpStatement{
            Name: $2,
        }
        $$.SetPosition($1)
    }
|   RETURN expression ';' {
        $$ = &nodes.ReturnStatement{
            Value: $2,
        }
        $$.SetPosition($1)
    }
|   RETURN ';' {
        $$ = &nodes.ReturnStatement{
        }
        $$.SetPosition($1)
    }
|   expression ';' {
        $$ = &nodes.ExpressionStatement{
            Expression: $1,
        }
        $$.SetPosition($1.Position())
    }
|   block_statement {
        $$ = $1
    }
|   IF '(' expression ')' statement %prec LOWER_THAN_ELSE {
        $$ = &nodes.IfStatement{
            If: $3,
            Then: $5,
        }
        $$.SetPosition($1)
    }
|   IF '(' expression ')' statement ELSE statement {
        $$ = &nodes.IfStatement{
            If: $3,
            Then: $5,
            Else: $7,
        }
        $$.SetPosition($1)
    }
|   FOR '(' for_expression_list ';' expression ';' for_expression_list ')' statement {
        $$ = &nodes.ForStatement{
            Init: $3,
            Condition: $5,
            Loop: $7,
            Body: $9,
        }
        $$.SetPosition($1)
    }
|   FOR '(' for_expression_list ';' ';' for_expression_list ')' statement {
        $$ = &nodes.ForStatement{
            Init: $3,
            Loop: $6,
            Body: $8,
        }
        $$.SetPosition($1)
    }
|   DO statement WHILE '(' expression ')' ';' {
        $$ = &nodes.DoStatement{
            Body: $2,
            Condition: $5,
        }
        $$.SetPosition($1)
    }
|   WHILE '(' expression ')' statement {
        $$ = &nodes.WhileStatement{
            Condition: $3,
            Body: $5,
        }
        $$.SetPosition($1)
    }
|   SWITCH '(' expression ')' block_statement {
        $$ = &nodes.SwitchStatement{
            Expression: $3,
            Block: $5,
        }
        $$.SetPosition($1)
    }
|   CASE next_for_expression_list ':' {
        $$ = &nodes.CaseStatement{
            Expressions: $2,
        }
        $$.SetPosition($1)
    }
|   DEFAULT ':' {
        $$ = &nodes.CaseStatement{}
        $$.SetPosition($1)
    }
|   BREAK ';' {
        $$ = &nodes.BreakStatement{}
        $$.SetPosition($1)
    }
|   CONTINUE ';' {
        $$ = &nodes.ContinueStatement{}
        $$.SetPosition($1)
    }



for_expression_list:
    {
    }
|   next_for_expression_list {
        $$ = $1
    }

next_for_expression_list:
    expression {
        $$ = append($$, $1)
    }
|   expression ',' next_for_expression_list {
        $$ = append($$, $1)
        $$ = append($$, $3...)
    }


lvalue_identifiers:
    identifier {
        $$ = &nodes.LValueExpression{
            Name: $1,
        }
        $$.SetPosition($1.Position())
    }
|   identifier '.' identifier {
        $$ = &nodes.LValueExpression{
            Name: $1,
            Item: $3,
        }
        $$.SetPosition($1.Position())
    }

list_item_expression:
    lvalue_identifiers '[' expression ']' {
        $$ = &nodes.ListItemExpression{
            LValue: $1,
            IsRange: false,
            StartIndex: $3,
        }
        $$.SetPosition($1.Position())
    }
|   lvalue_identifiers '[' expression ',' expression ']' {
        $$ = &nodes.ListItemExpression{
            LValue: $1,
            IsRange: true,
            StartIndex: $3,
            EndIndex: $5,
        }
        $$.SetPosition($1.Position())
    }

lvalue:
    lvalue_identifiers {
        $$ = $1
    }
|   list_item_expression {
        $$ = $1
    }
|   lvalue_identifiers '[' expression ']' '.' identifier {
        $$ = &nodes.ListItemExpression{
            LValue: $1,
            IsRange: false,
            StartIndex: $3,
            Item: $6,
        }
        $$.SetPosition($1.Position())
    }

expression:
    unary_expression {
        $$ = $1
    }
|   typecast_expression {
        $$ = $1
    }
|   lvalue '=' expression {
        $$ = &nodes.BinaryExpression{
            Operator: "=",
            LValue: $1,
            RValue: $3,
        }
        $$.SetPosition($1.Position())
    }
|   lvalue ADD_ASSIGN expression {
        $$ = &nodes.BinaryExpression{
            Operator: "+=",
            LValue: $1,
            RValue: $3,
        }
        $$.SetPosition($1.Position())
    }
|   lvalue SUB_ASSIGN expression {
        $$ = &nodes.BinaryExpression{
            Operator: "-=",
            LValue: $1,
            RValue: $3,
        }
        $$.SetPosition($1.Position())
    }
|   lvalue MUL_ASSIGN expression {
        $$ = &nodes.BinaryExpression{
            Operator: "*=",
            LValue: $1,
            RValue: $3,
        }
        $$.SetPosition($1.Position())
    }
|   lvalue DIV_ASSIGN expression {
        $$ = &nodes.BinaryExpression{
            Operator: "/=",
            LValue: $1,
            RValue: $3,
        }
        $$.SetPosition($1.Position())
    }
|   lvalue MOD_ASSIGN expression {
        $$ = &nodes.BinaryExpression{
            Operator: "%=",
            LValue: $1,
            RValue: $3,
        }
        $$.SetPosition($1.Position())
    }
|   expression EQ expression {
        $$ = &nodes.BinaryExpression{
            Operator: "==",
            LValue: $1,
            RValue: $3,
        }
        $$.SetPosition($1.Position())
    }
|   expression NEQ expression {
        $$ = &nodes.BinaryExpression{
            Operator: "!=",
            LValue: $1,
            RValue: $3,
        }
        $$.SetPosition($1.Position())
    }
|   expression LEQ expression {
        $$ = &nodes.BinaryExpression{
            Operator: "<=",
            LValue: $1,
            RValue: $3,
        }
        $$.SetPosition($1.Position())
    }
|   expression GEQ expression {
        $$ = &nodes.BinaryExpression{
            Operator: ">=",
            LValue: $1,
            RValue: $3,
        }
        $$.SetPosition($1.Position())
    }
|   expression '<' expression {
        $$ = &nodes.BinaryExpression{
            Operator: "<",
            LValue: $1,
            RValue: $3,
        }
        $$.SetPosition($1.Position())
    }
|   expression '>' expression {
        $$ = &nodes.BinaryExpression{
            Operator: ">",
            LValue: $1,
            RValue: $3,
        }
        $$.SetPosition($1.Position())
    }
|   expression '+' expression {
        $$ = &nodes.BinaryExpression{
            Operator: "+",
            LValue: $1,
            RValue: $3,
        }
        $$.SetPosition($1.Position())
    }
|   expression '-' expression {
        $$ = &nodes.BinaryExpression{
            Operator: "-",
            LValue: $1,
            RValue: $3,
        }
        $$.SetPosition($1.Position())
    }
|   expression '*' expression {
        $$ = &nodes.BinaryExpression{
            Operator: "*",
            LValue: $1,
            RValue: $3,
        }
        $$.SetPosition($1.Position())
    }
|   expression '/' expression {
        $$ = &nodes.BinaryExpression{
            Operator: "/",
            LValue: $1,
            RValue: $3,
        }
        $$.SetPosition($1.Position())
    }
|   expression '%' expression {
        $$ = &nodes.BinaryExpression{
            Operator: "%",
            LValue: $1,
            RValue: $3,
        }
        $$.SetPosition($1.Position())
    }
|   expression '&' expression {
        $$ = &nodes.BinaryExpression{
            Operator: "&",
            LValue: $1,
            RValue: $3,
        }
        $$.SetPosition($1.Position())
    }
|   expression '|' expression {
        $$ = &nodes.BinaryExpression{
            Operator: "|",
            LValue: $1,
            RValue: $3,
        }
        $$.SetPosition($1.Position())
    }
|   expression '^' expression {
        $$ = &nodes.BinaryExpression{
            Operator: "^",
            LValue: $1,
            RValue: $3,
        }
        $$.SetPosition($1.Position())
    }
|   expression BOOLEAN_AND expression {
        $$ = &nodes.BinaryExpression{
            Operator: "&&",
            LValue: $1,
            RValue: $3,
        }
        $$.SetPosition($1.Position())
    }
|   expression BOOLEAN_OR expression {
        $$ = &nodes.BinaryExpression{
            Operator: "||",
            LValue: $1,
            RValue: $3,
        }
        $$.SetPosition($1.Position())
    }
|   expression SHIFT_LEFT expression {
        $$ = &nodes.BinaryExpression{
            Operator: "<<",
            LValue: $1,
            RValue: $3,
        }
        $$.SetPosition($1.Position())
    }
|   expression SHIFT_RIGHT expression {
        $$ = &nodes.BinaryExpression{
            Operator: ">>",
            LValue: $1,
            RValue: $3,
        }
        $$.SetPosition($1.Position())
    }


unary_expression:
    '-' expression {
        $$ = &nodes.UnaryExpression{
            Operator: "-",
            RValue: $2,
        }
        $$.SetPosition($1)
    }
|   '!' expression {
        $$ = &nodes.UnaryExpression{
            Operator: "!",
            RValue: $2,
        }
        $$.SetPosition($1)
    }
|   '~' expression {
        $$ = &nodes.UnaryExpression{
            Operator: "~",
            RValue: $2,
        }
        $$.SetPosition($1)
    }
|   INC_OP lvalue {
        $$ = &nodes.UnaryExpression{
            Operator: "++",
            RValue: $2,
        }
        $$.SetPosition($1)
    }
|   DEC_OP lvalue {
        $$ = &nodes.UnaryExpression{
            Operator: "--",
            RValue: $2,
        }
        $$.SetPosition($1)
    }
|   DELETE list_item_expression {
        $$ = &nodes.DeleteExpression{
            RValue: $2,
        }
        $$.SetPosition($1)
    }
|   unary_postfix_expression {
        $$ = $1
    }
|   '#' unary_expression {
        $$ = &nodes.LengthExpression{
            RValue: $2,
        }
        $$.SetPosition($1)
    }
|   '(' expression ')' {
        $$ = &nodes.BracesExpression{
            Child: $2,
        }
        $$.SetPosition($1)
    }

typecast_expression:
    '(' typename ')' unary_expression {
        $$ = &nodes.TypecastExpression{
            Type: $2,
            Child: $4,
        }
        $$.SetPosition($1)
    }

unary_postfix_expression:
    '<' expression ',' expression ',' expression '>' %prec INITIALIZER {
        $$ = &nodes.VectorExpression{
            X: $2,
            Y: $4,
            Z: $6,
        }
        $$.SetPosition($1)
    }
|   '<' expression ',' expression ',' expression ',' expression '>' %prec INITIALIZER {
        $$ = &nodes.RotationExpression{
            X: $2,
            Y: $4,
            Z: $6,
            S: $8,
        }
        $$.SetPosition($1)
    }
|  ARRAY_BRACES %prec INITIALIZER {
        $$ = &nodes.ListExpression{
        }
        $$.SetPosition($1)
    }
|  '[' list_values ']' %prec INITIALIZER {
        $$ = &nodes.ListExpression{
            Values: $2,
        }
        $$.SetPosition($1)
    }
|   struct_expression %prec INITIALIZER {
        $$ = $1
    }
|   lvalue {
        $$ = $1
    }
|   lvalue INC_OP {
        $$ = &nodes.UnaryExpression{
            Operator: "++",
            RValue: $1,
            IsPostfix: true,
        }
        $$.SetPosition($1.Position())
    }
|   lvalue DEC_OP {
        $$ = &nodes.UnaryExpression{
            Operator: "--",
            RValue: $1,
            IsPostfix: true,
        }
        $$.SetPosition($1.Position())
    }
|   identifier '(' ')' {
        $$ = &nodes.FunctionCallExpression{
            Name: $1,
        }
        $$.SetPosition($1.Position())
    }
|   identifier '(' list_values ')' {
        $$ = &nodes.FunctionCallExpression{
            Name: $1,
            Arguments: $3,
        }
        $$.SetPosition($1.Position())
    }
|   constant {
        $$ = $1
    }

list_values:
    expression {
        $$ = append($$, $1)
    }
|   expression ',' list_values {
        $$ = append($$, $1)
        $$ = append($$, $3...)
    }

struct_expression:
    identifier '{' '}' {
        $$ = &nodes.StructExpression{
            Name: $1,
        }
        $$.SetPosition($1.Position())
    }
|   identifier '{' struct_expression_variables '}' {
        $$ = &nodes.StructExpression{
            Name: $1,
            Fields: $3,
        }
        $$.SetPosition($1.Position())
    }

struct_expression_variables:
    struct_expression_variable {
        $$ = append($$, $1)
    }
|   struct_expression_variable ',' struct_expression_variables {
        $$ = append($$, $1)
        $$ = append($$, $3...)
    }

struct_expression_variable:
    identifier ':' expression {
        $$ = &nodes.Variable{
            Name: $1,
            RValue: $3,
        }
        $$.SetPosition($1.Position())
    }

%%

var keywords = map[string]int{
    "integer": INTEGER,
    "int": INTEGER,
    "float": FLOAT,
    "string": STRING,
    "key": KEY,
    "vector": VECTOR,
    "rotation": ROTATION,
    "quaternion": ROTATION,
    "list": LIST,

    "default": DEFAULT,
    "state": STATE,

    "jump": JUMP,
    "return": RETURN,
    "if": IF,
    "else": ELSE,
    "for": FOR,
    "do": DO,
    "while": WHILE,

    "++": INC_OP,
    "--": DEC_OP,

    "+=": ADD_ASSIGN,
    "-=": SUB_ASSIGN,
    "*=": MUL_ASSIGN,
    "/=": DIV_ASSIGN,
    "%=": MOD_ASSIGN,

    "==": EQ,
    "!=": NEQ,
    "<=": LEQ,
    ">=": GEQ,

    "&&": BOOLEAN_AND,
    "||": BOOLEAN_OR,

    "<<": SHIFT_LEFT,
    ">>": SHIFT_RIGHT,

    "boolean": BOOLEAN,
    "true": TRUE,
    "false": FALSE,
    "TRUE": TRUE,
    "FALSE": FALSE,

    "include": INCLUDE,
    "pragma": PRAGMA,

    "struct": STRUCT,

    "[]": ARRAY_BRACES,

    "switch": SWITCH,
    "case": CASE,
    "break": BREAK,
    "continue": CONTINUE,

    "const": CONST,

    "delete": DELETE,
}

type Lexer struct {
    scanner.Scanner

    LastPos scanner.Position

    error error
    script *nodes.Script
}

func (self *Lexer) Lex(lval *yySymType) int {
    var err error

    token := self.Scan()
    text := self.TokenText()

    lval.pos = self.Pos()
    lval.pos.Offset -= len(text)
    lval.pos.Column -= len(text)

    self.LastPos = lval.pos

    if token == scanner.EOF {
        return 0
    }
/*
    if text == "#" {
        token = self.Scan()
        text += self.TokenText()
    }
*/

    if strings.HasPrefix(text, "//") {
        lval.text = strings.TrimPrefix(text, "//")
        return COMMENT
    } else if strings.HasPrefix(text, "/*") {
        lval.text = strings.TrimPrefix(strings.TrimSuffix(text, "*/"), "/*")
        return C_STYLE_COMMENT
    }

    if value, ok := keywords[text]; ok {
        lval.text = text
        return value
    }

    if token > 0 && self.Peek() > 0 {
        text2 := string(token) + string(self.Peek())

        if value, ok := keywords[text2]; ok {
            self.Scan()
            lval.text = text2
            return value
        }
    }

    switch token {
    case scanner.Ident:
        lval.text = text
        return IDENTIFIER
    case scanner.Int:
        var value int64

        if strings.HasPrefix(text, "0x") {
            value, err = strconv.ParseInt(text[2:], 16, 32)
        } else {
            value, err = strconv.ParseInt(text, 10, 32)
        }

        if err != nil {
            self.Error(fmt.Sprintf("Invalid integer value: %s", text))
            return 0
        }

        lval.integerValue.Value = int32(value)
        lval.integerValue.IsHex = strings.HasPrefix(text, "0x")
        return INTEGER_CONSTANT
    case scanner.Float:
        lval.floatValue, err = strconv.ParseFloat(text, 64)
        if err != nil {
            self.Error(fmt.Sprintf("Invalid float value: %s", text))
            return 0
        }
        return FLOAT_CONSTANT
    case scanner.String:
        lval.stringValue = text
        return STRING_CONSTANT
    }

    if int(token) > 0 {
        return int(token)
    }

    return 0
}

func (self *Lexer) Error(error string) {
    self.error = errors.New(fmt.Sprintf("%s %s", self.Pos().String(), error))
    self.script = nil
}

func ParseFile(filename string) (*nodes.Script, error) {
    source, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lexer := &Lexer{}
	lexer.Init(strings.NewReader(string(source)))
	lexer.Filename = filename
    lexer.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanFloats | scanner.ScanStrings | scanner.ScanComments

    yyErrorVerbose = true

    lexer.script = &nodes.Script{
        Filename: filename,
    }
    lexer.script.At = scanner.Position{
        Filename: filename,
        Offset: 0,
        Line: 1,
        Column: 1,
    }

    yyParse(lexer)

    return lexer.script, lexer.error
}
