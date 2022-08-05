<p id="id_top"></p>

- [escript - a embedded script for Go](#escript---a-embedded-script-for-go)
  - [Features](#features)
  - [Quick Start](#quick-start)
    - [run interactive interpreter](#run-interactive-interpreter)
    - [run scripts](#run-scripts)
    - [conditional expression](#conditional-expression)
    - [recursion](#recursion)
    - [closure](#closure)
    - [script parameters injection](#script-parameters-injection)
    - [eval AST](#eval-ast)
    - [dump & load AST as json](#dump--load-ast-as-json)
  - [builtin function](#builtin-function)
    - [type](#type)
    - [str](#str)
    - [int](#int)
    - [print](#print)
    - [println](#println)
    - [printf](#printf)
    - [sprintf](#sprintf)
    - [getenv](#getenv)
    - [loads](#loads)
    - [dumps](#dumps)
  - [types](#types)
    - [null](#null)
    - [boolean](#boolean)
    - [integer](#integer)
    - [string](#string)
    - [array](#array)
    - [hash](#hash)
    - [builtin](#builtin)
    - [function](#function)
    - [method](#method)
  - [License](#license)
  - [Author](#author)
  - [More](#more)

# escript - a embedded script for Go #

escript is a tiny, expression-based embedded language for Go.

    package main

    import (
        "github.com/jobs-github/escript/eval"
    )

    func main() {
        e := eval.New(false)
        code := `println("hello world");`
        node, _ := e.LoadAst(code)
        e.EvalAst(node, e.NewEnv(nil, nil))
    }

[back to top](#id_top)

## Features ##

Not only ([Writing An Interpreter In Go](https://interpreterbook.com/))  
* go-like syntax (also borrow some from python/c)
* variable bindings
* arithmetic expressions
* built-in functions
* first-class and higher-order functions (closure)
* data types
    * integer
    * boolean
    * string
    * array
    * hash

But also:  
* conditional expression
* object member call
* script parameters injection
* eval AST
* dump & load AST as json

[back to top](#id_top)

## Quick Start ##

### run interactive interpreter ###

    chmod +x /make.sh
    ./make.sh
    ./escript

[back to top](#id_top)

### run scripts ###

    chmod +x /make.sh
    ./make.sh
    ./escript scripts/hello.qs

[back to top](#id_top)

### [conditional expression](scripts/conditional.qs) ###

    (1 > 0) ? (1 + 1) : (10 % 3)

[back to top](#id_top)

### [recursion](scripts/mapreduce.qs) ###

    func map_iter(arr, accumulated, fn) {
        return (len(arr) == 0) ? accumulated : map_iter(arr.tail(), accumulated.push(fn(arr.first())), fn);
    };
    func map(arr, fn) {
        return map_iter(arr, [], fn);
    };
    func double(x) {
        return x * 2;
    };

    func reduce(arr, result, fn) {
        return (len(arr) == 0) ? result : reduce(arr.tail(), fn(result, arr.first()), fn);
    };
    func add(x, y) {
        return x + y;
    };
    func sum(arr) {
        return reduce(arr, 0, add);
    };

    const a = [1,2,3,4,5];
    const result = map(a, double);
    println(result);

    const rc = sum([1,2,3,4,5]);
    println(rc);

[back to top](#id_top)

### [closure](scripts/closure.qs) ###

    func sub(x) {
        return func(y) {
            return x - y;
        };
    };

    const map = func(arr, fn) {
        const iter = func(arr, accumulated) {
            if (len(arr) == 0) {
                return accumulated;
            } else {
                return iter(arr.tail(), accumulated.push(fn(arr.first())));
            }
        };
        return iter(arr, []);
    };
    const double = func(x) {
        return x * 2;
    };

    const reduce = func(arr, initial, fn) {
        const iter = func(arr, result) {
            if (len(arr) == 0) {
                return result;
            } else {
                return iter(arr.tail(), fn(result, arr.first()));
            }
        };
        return iter(arr, initial);
    };
    const add = func(x, y) {
        return x + y;
    };
    const sum = func(arr) {
        return reduce(arr, 0, add);
    };

    const subFn = sub(10);
    const r = subFn(5);
    println(r);

    const a = [1,2,3,4,5];
    const result = map(a, double);
    println(result);

    const rc = sum([1,2,3,4,5]);
    println(rc);

[back to top](#id_top)

### [script parameters injection](scripts/hello.qs) ###

    chmod +x /make.sh
    ./make.sh
    ./escript scripts/hello.qs hello hello world

[back to top](#id_top)

### eval AST ###

    package main

    import (
        "github.com/jobs-github/escript/eval"
    )

    func main() {
        e := eval.New(false)
        code := `const r = (1 > 0) ? (1 + 1) : (10 % 3); println(r);`
        node, _ := e.LoadAst(code)
        e.EvalAst(node, e.NewEnv(nil, nil))
    }

[back to top](#id_top)

### dump & load AST as json ###

dump AST as json:  

    chmod +x /make.sh
    ./make.sh
    ./escript --dump scripts/conditional.qs > scripts/conditional.json

load & eval AST from json:  

    chmod +x /make.sh
    ./make.sh
    ./escript --load scripts/conditional.json

[back to top](#id_top)

## [builtin function](builtin/builtin.go) ##

### type ###

    const s = "EOF";
    println(type(s));

[back to top](#id_top)

### str ###

    const arr = [1,2,3];
    println(str(arr));

[back to top](#id_top)

### int ###

    const is = "123";
    println(int(is));

[back to top](#id_top)

### print ###

    print("123 ", 456, " ", [7,8,9]);

[back to top](#id_top)

### println ###

    println("aaa ", 111, " ", {"k1": 1, "k2": "2"});

[back to top](#id_top)

### printf ###

    const testi = 999;
    const tests = "2022";
    printf("%v\n%v\n", testi, tests);

[back to top](#id_top)

### sprintf ###

    const testi2 = 9999;
    const tests2 = "2222";
    const ms = sprintf("%v-%v", testi, tests);
    println(ms);

[back to top](#id_top)

### [getenv](scripts/env.qs) ###

    export test_env=HAHA; ./escript scripts/env.qs

[back to top](#id_top)

### [loads](scripts/json.qs) ###

    const s = "{\"k3\":{\"k31\":true,\"k32\":[1,2,3]},\"k2\":\"2\",\"k1\":1}"; 
    const obj = loads(s);
    println(obj);

[back to top](#id_top)

### [dumps](scripts/json.qs) ###

    const obj1 = { "k1": null, "k2": 123 };
    const objstr = dumps(obj1);
    println(objstr);

[back to top](#id_top)

## [types](object/def.go) ##

### [null](object/null.go) ###

method  |comment
--------|-------
not     |!

    >> const n = null;
    >> n.not();
    true
    >> !n;
    true

[back to top](#id_top)

### [boolean](object/boolean.go) ###

method  |comment
--------|-------
not     |!
opposite|-
int     |convert to int

    >> const b = true;
    true
    >> b.not()
    false
    >> b.opposite()
    -1
    >> b.int()
    1
    >> !b
    false
    >> -b
    -1
    >> int(b)
    1

[back to top](#id_top)

### [integer](object/integer.go) ###

method  |comment
--------|-------
not     |!
opposite|-
int     |convert to int

    >> const i = 123
    123
    >> i.not()
    false
    >> i.opposite()
    -123
    >> i.int()
    123
    >> !i
    false
    >> -i
    -123
    >> int(i)
    123

[back to top](#id_top)
### [string](object/string.go) ###

method  |comment
--------|-------
len     |length of string
index   |get value by index
not     |!
int     |convert to int

    >> const s = "123"
    123
    >> s.len()
    3
    >> s.index(1)
    2
    >> s.not()
    false
    >> s.int()
    123
    >> len(s)
    3
    >> s[1]
    2
    >> !s
    false
    >> int(s)
    123

[back to top](#id_top)
### [array](object/array.go) ###

method  |comment
--------|-------
len     |length of array
index   |get value by index
not     |!
first   |first value
last    |last value
tail    |remove first value and return rest
push    |append value

    >> const arr = [1,2,3,4,5]
    [1, 2, 3, 4, 5]
    >> arr.push("hello")
    [1, 2, 3, 4, 5, hello]

    >> arr.len()
    6
    >> arr[0]
    1
    >> arr.not()
    false
    >> !arr
    false
    >> arr.first()
    1
    >> arr.last()
    hello
    >> arr.tail()
    [2, 3, 4, 5, hello]
    >> len(arr)
    6
    >> arr[1] = 200

    >> arr
    [100, 200, 3, 4, 5, hello]

[back to top](#id_top)

### [hash](object/hash.go) ###

method  |comment
--------|-------
len     |length of hash pairs
index   |get value by key
not     |!

    >> const h = {"k1": 1, "k2": "bbb", "k3": [1,2,3]}
    {k1: 1, k2: bbb, k3: [1, 2, 3]}

    >> h["k1"]
    1
    >> h.len()
    3
    >> h.not()
    false
    >> h.index("k1")
    1000
    >> !h
    false
    >> len(h)
    3

[back to top](#id_top)

### [builtin](object/builtin.go) ###

method  |comment
--------|-------
not     |!

    >> const l = len
    <built-in function len>
    >> const s = "123"
    123
    >> l(s)
    3
    >> l.not()
    false
    >> !l
    false

[back to top](#id_top)
### [function](object/function.go) ###

method  |comment
--------|-------
not     |!

    >> const fn = func(x, y) { return x + y; }
    func (x, y) {
    return (x + y);
    }
    >> fn(1,2)
    3
    >> !fn
    false
    >> fn.not()
    false

[back to top](#id_top)
### [method](object/objectfunc.go) ###

method  |comment
--------|-------
not     |!

    >> const arr = [1,2,3]
    [1, 2, 3]
    >> const fn = arr.index
    <built-in method index of array object>
    >> fn(1)
    2
    >> !fn
    false
    >> fn.not()
    false
    >> const l = arr.len
    <built-in method len of array object>
    >> l()
    3

[back to top](#id_top)

## License ##

escript is licensed under [New BSD License](https://opensource.org/licenses/BSD-3-Clause), a very flexible license to use.

[back to top](#id_top)

## Author ##

* chengzhuo (jobs, yao050421103@163.com)  

[back to top](#id_top)

## More ##

- [Writing An Interpreter In Go](https://interpreterbook.com/)  

[back to top](#id_top)