package nodes

type ExpressionType int

const (
	ExpressionUnknown ExpressionType = iota

	ExpressionBraces       // (...)
	ExpressionUnary        // -a
	ExpressionBinary       // a + b
	ExpressionTypecast     // (integer)a
	ExpressionFunctionCall /// a(1, 2, 3)
	ExpressionVector       // <a, b, c>
	ExpressionRotation     // <a, b, c, d>
	ExpressionList         // [a, b, c, d]
	ExpressionLValue       // identifier or vector.x or rotation.x
	ExpressionConstant     // 123

    ExpressionListItem     // q[1], q[1, 3]
    ExpressionStruct       // q{ a: 5, b: 10 }

    ExpressionDelete       // delete q[5]
    ExpressionLength       // #q
)

type Expression Node
