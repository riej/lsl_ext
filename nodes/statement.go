package nodes

type StatementType int

const (
	StatementUnknown StatementType = iota

	StatementEmpty  // ;
	StatementBlock  // { ... }
	StatementReturn // return ...;
	StatementLabel  // @label;
	StatementJump   // jump label;
	StatementIf     // if (...) ... else ...
	StatementFor    // for (...; ...; ...) ...
	StatementDo     // do ... while (...);
	StatementWhile  // while (...) ...
	StatementState  // state ...;

	StatementVariable
	StatementExpression
	StatementComment

    StatementSwitch
    StatementCase
    StatementBreak
    StatementContinue
)

type Statement Node
