Hello!

I am beginner developer and SL scripter. Because LSL is too dumb language, and development is complicated, i decided to write some translator, to extends its syntax.

Additional feature - it formats source code to make it well-readable and good-looking.

Because i have not too much experience, source code is really bad and weird... Please forgive me for that. But it works, and that is good.

Also please forgive me for my bad english.

Great thanks to [Makopo/lslint](https://github.com/Makopo/lslint "Makopo/lslint") and [Sei_Lisa/kwdb](https://bitbucket.org/Sei_Lisa/kwdb "Sei_Lisa/kwdb") projects!

**Usage**

`lsl_ext <filename.lslx>`

**Limitations and issues**

This is not lint tool. It does not make full validation of source code.

It does not optimizes your code. I plan to add it later, but now - no.

Also it can show error messages with not so accurate column position... but error is somewhere nearby.

If there is not `builtins.txt` file in current folder, application will download [its latest version](https://bitbucket.org/api/1.0/repositories/Sei_Lisa/kwdb/raw/default/outputs/builtins.txt "latest version").

**Compiling**

This application is written in `go` and requires [`goyacc`](https://godoc.org/golang.org/x/tools/cmd/goyacc "`goyacc`") tool (which should be installed automatically).

To build application simply run
```sh
go generate
```
And again, please forgive me for source code quality!

# Features

### Lazy lists
```c
list test = [1, 2, 3, 4, 5, llGetPos()];
test[1] = "qqq";
test[2] += 10;
test[3] = test[4] + 20;
llOwnerSay(test[5]);
llSetPos(test[6]);
delete test[1, 2];
```
&darr;&darr;&darr;
```lsl
list test = [1, 2, 3, 4, 5, llGetPos()];

test = llListReplaceList(test, ["qqq"], 1, 1);
test = llListReplaceList(test, [llList2Integer(test, 2) + 10], 2, 2);
test = llListReplaceList(test, [llList2Float(test, 4) + 20.0], 3, 3);

llOwnerSay(llList2String(test, 5));
llSetPos((vector)llList2String(test, 6));

test = llDeleteSubList(test, 1, 2);
```
Almost same as in Firestorm, but better, as i think. It supports all assignment operators `=`, `+=`, `-=`, `*=`, `/=`, `%=`, and can auto-detect required type (depending on nearest constant value or variable; but `test[1] + test[2]` will be detected as strings).

Typecasting:

`(string)test[1]` &rarr; `llList2String(test, 1)`

`(key)test[1]` &rarr; `llList2Key(test, 1)`

`(integer)test[1]` &rarr; `llList2Integer(test, 1)`

`(float)test[1]` &rarr; `llList2Float(test, 1)`

`(vector)test[1]` &rarr; `(vector)llList2String(test, 1)`

`(rotation)test[1]` &rarr; `(rotation)llList2String(test, 1)`

`(list)test[1, 5]` &rarr; `llList2List(test, 1, 5)`

Unfortunately, you cannot use it with vector or rotation fields assignment, like `test[2].x = 10;`. Sorry!

##### Usage with strings
```c
string s = "qwertyuiop";
s[5, 10] = "123";
llOwnerSay(s[1, 3]);
delete s[3];
```
&darr;&darr;&darr;
```lsl
string s = "qwertyuiop";

s = llInsertString(llDeleteSubString(s, 5, 10), 5, "123");

llOwnerSay(llGetSubString(s, 1, 3));
llSetPos((vector)llList2String(test, 3));

s = llDeleteSubString(s, 3, 3);
```

Square braces can be used for strings too.

### Structures
```c
struct TestStruct {
    string name = "<unknown>";
    integer access;
}

TestStruct one{name: "One"};
TestStruct two{access: -10};

TestStruct[] arr;

Dump(TestStruct s) {
    llOwnerSay(s.name + " " + s.access);
}

TestFunc() {
    one.access = 5;;

    arr += one + two;
    arr += TestStruct{ name: "Three", access: 20 };
    arr += TestStruct{ name: "Four" };;

    integer i = 2;
    arr[i + 1].access = 1;

    Dump(arr[0]);
    Dump(arr[1]);

    delete arr[3];
    delete arr[0, 1];
}
```
&darr;&darr;&darr;
```lsl
list one = ["One", 0];
list two = ["<unknown>", -10];
list arr;

Dump(list s) {
    llOwnerSay(llList2String(s, 0) + " " + llList2String(s, 1));
}

TestFunc() {
    one = llListReplaceList(one, [5], 1, 1);

    arr += one + two;
    arr += ["Three", 20];
    arr += ["Four", 0];

    integer i = 2;

    arr = llListReplaceList(arr, [1], 2 * (i + 1) + 1, 2 * (i + 1) + 1);

    Dump(llList2List(arr, 2 * 0, 2 * 0 + 1));
    Dump(llList2List(arr, 2 * 1, 2 * 1 + 1));

    arr = llDeleteSubList(arr, 2 * 3, 2 * 3 + 1);
    arr = llDeleteSubList(arr, 2 * 0, 2 * 1 + 1);
}
```
I missed this feature so much in LSL!

You can define default values for structure fields, and use arrays of structures together with lazy lists.
Unfortunately, you cannot use fields of list and struct types. I just don't know how to handle it, 

### Switch/case

```c
int i = 5;
switch (i + 1) {
case 1, 3, 5:
    llOwnerSay("1, 3, 5");
case 2:
    llOwnerSay("2");
case 4:
    llOwnerSay("4");
default:
    llOwnerSay("default");
}
```
&darr;&darr;&darr;
```lsl
// switch(i + 1)
if (((i + 1) == 1) || ((i + 1) == 3) || ((i + 1) == 5)) {
    llOwnerSay("1, 3, 5");
} else if ((i + 1) == 2) {
    llOwnerSay("2");
} else if ((i + 1) == 4) {
    llOwnerSay("4");
} else {
    llOwnerSay("default");
}
```

Simple-in-use switch statement.

#### Legacy
```c
#pragma legacy_switch

int i = 5;
switch (i + 1) {
case 1, 3, 5:
    llOwnerSay("1, 3, 5");
    break;
case 2:
    llOwnerSay("2");
case 4:
    llOwnerSay("4");
    break;
default:
    llOwnerSay("default");
}
```
&darr;&darr;&darr;
```lsl
integer i = 5;

// switch(i + 1)
if (((i + 1) == 1) || ((i + 1) == 3) || ((i + 1) == 5)) jump case_1894;
else if ((i + 1) == 2) jump case_1949;
else if ((i + 1) == 4) jump case_1980;
else jump case_2023;

@case_1894; // case 1, 3, 5:
    llOwnerSay("1, 3, 5");

    jump switch_end_1876; // break

@case_1949; // case 2:
    llOwnerSay("2");

@case_1980; // case 4:
    llOwnerSay("4");

    jump switch_end_1876; // break

@case_2023; // default:
    llOwnerSay("default");

@switch_end_1876;
```
C-style `switch` operator. Almost like in Firestorm, but colon `:` is required, and no need in braces `{...}`.

Use `#pragma legacy_switch` to enable it.

### Break/continue in loops
```c
integer i;
for (;; i++) {
    if (i > 5) break;
    if (i < 3) continue;
    llOwnerSay((string)i);
}
```
&darr;&darr;&darr;
```lsl
for (; TRUE; i++) {
    if (i > 5) {
        jump for_end_1068; // break
    }

    if (i < 3) {
        jump for_body_end_1068; // continue
    }

    llOwnerSay((string)i);

    @for_body_end_1068;
}
@for_end_1068;
```
Can be used in `for`, `while`, `do` loops. If there aren't any `break`/`continue` statements inside, labels will not be added.

### Miscellaneous
#### Include
```c
#include "filename.lslx"
```
Pastes specified script as part of current script.

#### Pragma options
```c
#pragma legacy_switch
#pragma no_legacy_switch
#pragma skip_unused
#pragma no_skip_unused
```
Specify options for translator.

`legacy_switch` enabled C-style (Firestorm-like) `switch` statement support.

`skip_unused` - don't print unused variables and functions. Doesn't works well now, will be reworked in future.

*Warning!!!* Those options can be renamed or changed in future!

#### Aliases
`boolean` &rarr; `integer`

`true` &rarr; `TRUE`

`false` &rarr; `FALSE`

`int` &rarr; `integer`

Those aliases has no any special purpose.

#### Constants
```c
const CONSTANT_NAME = "value";
const ANOTHER_CONSTANT = 123;
```
&darr;&darr;&darr;
```lsl
string CONSTANT_NAME = "value";
integer ANOTHER_CONSTANT = 123;
```

Constants cannot be assigned. Has no any special purpose too.

#### Multiple variables declaration
```c
integer a = 5, b = 10, c;
```
&darr;&darr;&darr;
```lsl
integer a = 5;
integer b = 10;
integer c;
```

#### List/string length
```c
TestStruct[] arr;
llOwnerSay((string)#some_list);
llOwnerSay((string)#"123");
llOwnerSay((string)#arr);
```
&darr;&darr;&darr;
```lsl
list arr;

llOwnerSay((string)llGetListLength(some_list));
llOwnerSay((string)llStringLength("123"));
llOwnerSay((string)(llGetListLength(arr) / 2));
```
Maybe it is clumsy... i don't really like this... but okay, it it much easier than those long function names.

Also works with lists of structures, returning count of items. If you want to get real list length, use `llGetListLength(arr)` or `#((list)arr)`.

#### Code formatting

Scripts are being parsed into syntax tree and displayed back, with custom formatting. It is pretty accurate, but can skip empty lines. To add empty line, put `;` as statement:
```lsl
llOwnerSay("123");;
```
Also some lists will be pretty-printed using newlines:
```lsl
llSetPrimitiveParams([PRIM_COLOR, ALL_SIDES, ZERO_VECTOR, 1, PRIM_COLOR, 3, <1.0, 1.0, 1.0>, 1.0]);
llHTTPRequest("https://google.com", [HTTP_METHOD, "GET", HTTP_VERBOSE_THROTTLE, FALSE, HTTP_BODY_MAXLENGTH, 16384], "");
```
&darr;&darr;&darr;
```lsl
llSetPrimitiveParams([
    PRIM_COLOR, ALL_SIDES, ZERO_VECTOR, 1,
    PRIM_COLOR, 3, <1.0, 1.0, 1.0>, 1.0
]);
llHTTPRequest("https://google.com", [
    HTTP_METHOD, "GET",
    HTTP_VERBOSE_THROTTLE, FALSE,
    HTTP_BODY_MAXLENGTH, 16384
], "");
```
