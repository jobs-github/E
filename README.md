<p id="id_top"></p>

- [escript - an embedded script for Go](#escript---a-embedded-script-for-go)
  - [Features](#features)
  - [Quick Start](#quick-start)
    - [run interactive interpreter](#run-interactive-interpreter)
    - [run scripts](#run-scripts)
    - [conditional expression](#conditional-expression)
    - [recursion](#recursion)
    - [closure](#closure)
    - [eval](#eval)
    - [embedded eval](#embedded-eval)
    - [dump \& load AST as json](#dump--load-ast-as-json)
  - [builtin function](#builtin-function)
    - [type](#type)
    - [str](#str)
    - [print](#print)
    - [println](#println)
    - [printf](#printf)
    - [sprintf](#sprintf)
    - [loads](#loads)
    - [dumps](#dumps)
    - [loop](#loop)
    - [map](#map)
    - [reduce](#reduce)
    - [filter](#filter)
    - [range](#range)
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

# escript - an embedded script for Go #

escript is a tiny, expression-based embedded language for Go.

    package main

    import (
        "fmt"
        "github.com/jobs-github/escript"
    )

    func main() {
        code := `const r = (1 > 0) ? (1 + 1) : (10 % 3); println(r);`
        r, _ := escript.NewState(code)
        res, _ := r.Run(nil)
        fmt.Println(res)
    }

[back to top](#id_top)

## Features ##

* go-like syntax (also borrow some from python/c)
* variable bindings
* arithmetic expressions
* built-in functions
* first-class and higher-order functions (closure)
* conditional expression
* object member call
* eval AST
* dump & load AST as json
* data types
    * integer
    * boolean
    * string
    * array
    * hash

[back to top](#id_top)

## Quick Start ##

### run interactive interpreter ###

    chmod +x /make.sh
    ./make.sh
    ./escript

    
**Note** : you can change `useVM` in `repl/main.go` to choose default interpreter.  

[back to top](#id_top)

### run scripts ###

    chmod +x /make.sh
    ./make.sh
    ./escript scripts/hello.es

[back to top](#id_top)

### [conditional expression](scripts/conditional.es) ###

    (1 > 0) ? (1 + 1) : (10 % 3)

[back to top](#id_top)

### [recursion](scripts/mapreduce.es) ###

    func map_iter(arr, accumulated, fn) {
        (arr.len() == 0) ? accumulated : map_iter(arr.tail(), accumulated.push(fn(arr.first())), fn);
    };
    func my_map(arr, fn) {
        map_iter(arr, [], fn);
    };
    func double(x) {
        x * 2;
    };

    func my_reduce(arr, result, fn) {
        (arr.len() == 0) ? result : my_reduce(arr.tail(), fn(result, arr.first()), fn);
    };
    func add(x, y) {
        x + y;
    };
    func sum(arr) {
        my_reduce(arr, 0, add);
    };

    const a = [1,2,3,4,5];
    const result = my_map(a, double);
    println(result);

    const rc = sum([1,2,3,4,5]);
    println(rc);

[back to top](#id_top)

### [closure](scripts/closure.es) ###

    const r = func (x) {
        x * 2;
    }(2);

    println(r);

    func add(x) {
        func(y) {
            x + y;
        };
    };

    const fn = add(1);
    println(fn(2));

[back to top](#id_top)

### eval ###

    package main

    import (
        "fmt"
        "github.com/jobs-github/escript"
    )

    func main() {
        code := `const r = (1 > 0) ? (1 + 1) : (10 % 3); println(r);`
        r, _ := escript.NewState(code)
        res, _ := r.Run(nil)
        fmt.Println(res)
    }

[back to top](#id_top)

### embedded eval ###

    package main

    import (
        "fmt"
        "github.com/jobs-github/escript"
        "github.com/jobs-github/escript/object"
    )

    func main() {
        T := func() (object.Object, error) { return object.True, nil }
        F := func() (object.Object, error) { return object.False, nil }
        r, _ := escript.NewState(`($a || $b) && ($c || $d);`)
        res, _ := r.Run(object.Symbols{"a": T, "b": T, "c": T, "d": F})
        fmt.Println(res)
    }

### dump & load AST as json ###

dump AST as json:  

    chmod +x /make.sh
    ./make.sh
    ./escript --dump scripts/conditional.es > scripts/conditional.json

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

### [loads](scripts/json.es) ###

    const s = "{\"k3\":{\"k31\":true,\"k32\":[1,2,3]},\"k2\":\"2\",\"k1\":1}"; 
    const obj = loads(s);
    println(obj);

[back to top](#id_top)

### [dumps](scripts/json.es) ###

    const obj1 = { "k1": null, "k2": 123 };
    const objstr = dumps(obj1);
    println(objstr);

### [loop](scripts/loop.es) ###

    loop(10, func(i){
        i < 5 ? println(i) : println(i * i)
    });

### [map](scripts/array.es) ###

    const arr = [1,2,3];
    const r = map(arr, func(i, item) { item * 2 });
    println(r);

### [reduce](scripts/array.es) ###

    const arr = [1,2,3];
    const acc = reduce(arr, func(x, y) { x + y }, 0);
    println(acc);

### [filter](scripts/array.es) ###

    const arr = [1,2,3];
    const fr = filter(arr, func(i, item) { item > 1 });
    println(fr);

### [range](scripts/array.es) ###

    println(range(10, func(i) { 
        sprintf("key_%v", (i % 2 == 0) ? i * 2 : i)
    }));

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
neg|-
int     |convert to int

    >> const b = true;
    true
    >> b.not()
    false
    >> b.neg()
    -1
    >> b.int()
    1
    >> !b
    false
    >> -b
    -1

[back to top](#id_top)

### [integer](object/integer.go) ###

method  |comment
--------|-------
not     |!
neg|-
int     |convert to int

    >> const i = 123
    123
    >> i.not()
    false
    >> i.neg()
    -123
    >> i.int()
    123
    >> !i
    false
    >> -i
    -123

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
    >> s[1]
    2
    >> !s
    false

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

[back to top](#id_top)

### [hash](object/hash.go) ###

method  |comment
--------|-------
len     |length of hash pairs
index   |get value by key
keys    |get keys
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
    >> h.keys()
    [k1, k2, k3]

[back to top](#id_top)

### [builtin](object/builtin.go) ###

method  |comment
--------|-------
not     |!

    >> const t = type
    <built-in function type>
    >> const s = "123"
    123
    >> t(s)
    string
    >> t.not()
    false
    >> !t
    false

[back to top](#id_top)
### [function](object/function.go) ###

method  |comment
--------|-------
not     |!

    >> const fn = func(x, y) { x + y; }
    closure[0xc000425260]
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

- [Writing A Compiler In Go](https://compilerbook.com/)  

[back to top](#id_top)