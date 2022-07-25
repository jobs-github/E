<p id="id_top"></p>

- [Q - a tiny go-like script](#q---a-tiny-go-like-script)
  - [Features](#features)
  - [Quick Start](#quick-start)
    - [run interactive interpreter](#run-interactive-interpreter)
    - [run scripts](#run-scripts)
    - [conditional expression](#conditional-expression)
    - [recursion](#recursion)
    - [closure](#closure)
    - [for loop](#for-loop)
    - [iterator](#iterator)
    - [defer](#defer)
    - [script parameters injection](#script-parameters-injection)
    - [eval AST](#eval-ast)
    - [dump & load AST as json](#dump--load-ast-as-json)
    - [simple module management](#simple-module-management)
    - [simple file I/O](#simple-file-io)
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
    - [open](#open)
    - [remove](#remove)
  - [types](#types)
    - [null](#null)
    - [boolean](#boolean)
    - [integer](#integer)
    - [string](#string)
    - [array](#array)
    - [array_iter](#array_iter)
    - [hash](#hash)
    - [hash_iter](#hash_iter)
    - [builtin](#builtin)
    - [function](#function)
    - [method](#method)
    - [file](#file)
  - [License](#license)
  - [Author](#author)
  - [More](#more)

# Q - a tiny go-like script #

Q is a tiny, dynamic interpreted language for Go.

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
* assign statement
* conditional expression
* go-like for loop
* object member call
* iterator
* defer
* script parameters injection
* eval AST
* dump & load AST as json
* simple module management
* simple file I/O

[back to top](#id_top)

## Quick Start ##

### run interactive interpreter ###

    chmod +x /make.sh
    ./make.sh
    ./Q

[back to top](#id_top)

### run scripts ###

    chmod +x /make.sh
    ./make.sh
    ./Q scripts/hello.qs

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

    func main() {
        var a = [1,2,3,4,5];
        var result = map(a, double);
        println(result);

        var rc = sum([1,2,3,4,5]);
        println(rc);
    };

    main();

[back to top](#id_top)

### [closure](scripts/closure.qs) ###

    func sub(x) {
        return func(y) {
            return x - y;
        };
    };

    var map = func(arr, fn) {
        var iter = func(arr, accumulated) {
            if (len(arr) == 0) {
                return accumulated;
            } else {
                return iter(arr.tail(), accumulated.push(fn(arr.first())));
            }
        };
        return iter(arr, []);
    };
    var double = func(x) {
        return x * 2;
    };

    var reduce = func(arr, initial, fn) {
        var iter = func(arr, result) {
            if (len(arr) == 0) {
                return result;
            } else {
                return iter(arr.tail(), fn(result, arr.first()));
            }
        };
        return iter(arr, initial);
    };
    var add = func(x, y) {
        return x + y;
    };
    var sum = func(arr) {
        return reduce(arr, 0, add);
    };

    var main = func() {
        var subFn = sub(10);
        var r = subFn(5);
        println(r);

        var a = [1,2,3,4,5];
        var result = map(a, double);
        println(result);

        var rc = sum([1,2,3,4,5]);
        println(rc);
    };

    main();

[back to top](#id_top)

### [for loop](scripts/for.qs) ###

    func test_for1() {
        var i = 0;
        var result = 0;
        for ;; {
            if (i > 4) {
                break;
            }
            result = result + i;
            i = i + 1;
        }
        println("test_for1: ", result);
    };

    func test_for2() {
        var result = 0;
        for var i = 0; i < 5; i = i + 1 {
            result = result + i;
        }
        println("test_for2: ", result);
    };

    func test_for3() {
        var result = 0;
        for var i = 0; i < 5;  {
            result = result + i;
            i = i + 1;
        }
        println("test_for3: ", result);
    };

    func test_for4() {
        var i = 0;
        var result = 0;
        for ; i < 5;  {
            result = result + i;
            i = i + 1;
        }
        println("test_for4: ", result);
    };

    func test_for5() {
        var i = 0;
        var result = 0;
        for i = 0; i < 5; i = i + 1 {
            result = result + i;
        }
        println("test_for5: ", result);
    };

    func test_for6() {
        var i = 0;
        var result = 0;
        for {
            if (i > 4) {
                break;
            }
            result = result + i;
            i = i + 1;
        }
        println("test_for6: ", result);
    };

    func main() {
        var fns = [test_for1, test_for2, test_for3, test_for4, test_for5, test_for6];
        var sz = fns.len();
        for var i = 0; i < sz; i = i + 1 {
            fns[i]();
        }
    };

    main();

[back to top](#id_top)

### [iterator](scripts/iterator.qs) ###

    func test_array_iter() {
        var arr = [10,20,30];
        var result = [];
        for var iter = arr.iter(); iter != null; iter = iter.next() {
            var idx = iter.key();
            var val = iter.value();
            result.push(idx);
            result.push(val);
        };
        println(result);
    };

    func test_hash_iter() {
        var h = {"k1": "v1", "k2": "v2", "k3": "v3"};
        var result = [];
        for var iter = h.iter(); iter != null; iter = iter.next() {
            var key = iter.key();
            var val = iter.value();
            result.push(key);
            result.push(val);
        };
        println(result);
    };

    func main() {
        test_array_iter();
        test_hash_iter();
    };

    main();

[back to top](#id_top)

### [defer](scripts/defer.qs) ###

    func test_defer(x) {
        defer func() {
            println(x);
        }();
        defer println(x * x * x);
        println(x * 2);
    };

    func main() {
        test_defer(2);
    };

    main();

[back to top](#id_top)

### [script parameters injection](scripts/hello.qs) ###

    chmod +x /make.sh
    ./make.sh
    ./Q scripts/hello.qs hello hello world

[back to top](#id_top)

### eval AST ###

    package main

    import (
        "github.com/jobs-github/escript/eval"
    )

    func main() {
        e := eval.New(false)
        code := `var r = (1 > 0) ? (1 + 1) : (10 % 3); println(r);`
        node, _ := e.LoadAst(code)
        e.EvalAst(node, e.NewEnv(nil, nil))
    }

[back to top](#id_top)

### dump & load AST as json ###

dump AST as json:  

    chmod +x /make.sh
    ./make.sh
    ./Q --dump scripts/conditional.qs > scripts/conditional.json

load & eval AST from json:  

    chmod +x /make.sh
    ./make.sh
    ./Q --load scripts/conditional.json

[back to top](#id_top)

### [simple module management](scripts/import.qs) ###

    /* scripts/import.qs */
    import "utils"
    import libmath "math"

    func hello() {
        println("hello");
    };

    func main() {
        var fn = utils.hello;
        hello();
        fn();
        utils.hello();
        println(utils.EOF);
        var r = libmath.add(1,2);
        println(r)
    };

    main();

    /* scripts/math.qs */
    func add(x, y) {
        return x + y;
    };

    /* scripts/utils.qs */
    func hello() {
        println("utils.hello");
    };

    var EOF = "EOF";

[back to top](#id_top)

### [simple file I/O](scripts/file.qs) ###

    func test_write(url) {
        println("test_write");
        var f = open(url, "w+");
        if (!f) {
            return;
        }
        defer f.close()
        var r = f.write("hello world");
        if (!r) {
            println("write fail")
        }
    };

    func test_read(url) {
        println("test_read");
        var f = open(url, "r");
        if (!f) {
            return;
        }
        defer f.close()
        var sz = f.size();
        println("size: ", sz);
        var data = f.read();
        println(data);
    };

    func test_seek(url) {
        println("test_seek");
        var f = open(url, "r");
        if (!f) {
            return;
        }
        defer f.close()
        var off = f.size() / 2;
        println("offset: ", off);
        f.seek(off);
        var data = f.read();
        println(data);
    };

    func main() {
        var url = "tmp.txt";
        defer remove(url);
        
        test_write(url);
        test_read(url);
        test_seek(url);
    };

    main();

[back to top](#id_top)

## [builtin function](builtin/builtin.go) ##

### type ###

    var s = "EOF";
    println(type(s));

[back to top](#id_top)

### str ###

    var arr = [1,2,3];
    println(str(arr));

[back to top](#id_top)

### int ###

    var is = "123";
    println(int(is));

[back to top](#id_top)

### print ###

    print("123 ", 456, " ", [7,8,9]);

[back to top](#id_top)

### println ###

    println("aaa ", 111, " ", {"k1": 1, "k2": "2"});

[back to top](#id_top)

### printf ###

    var testi = 999;
    var tests = "2022";
    printf("%v\n%v\n", testi, tests);

[back to top](#id_top)

### sprintf ###

    var testi2 = 9999;
    var tests2 = "2222";
    var ms = sprintf("%v-%v", testi, tests);
    println(ms);

[back to top](#id_top)

### [getenv](scripts/env.qs) ###

    export test_env=HAHA; ./Q scripts/env.qs

[back to top](#id_top)

### [loads](scripts/json.qs) ###

    var s = "{\"k3\":{\"k31\":true,\"k32\":[1,2,3]},\"k2\":\"2\",\"k1\":1}"; 
    var obj = loads(s);
    println(obj);

[back to top](#id_top)

### [dumps](scripts/json.qs) ###

    var obj1 = { "k1": null, "k2": 123 };
    var objstr = dumps(obj1);
    println(objstr);

[back to top](#id_top)

### [open](scripts/file.qs) ###

    var f = open(url, "w+");
    if (!f) {
        return;
    }

[back to top](#id_top)

### [remove](scripts/file.qs) ###

    defer remove(url);

[back to top](#id_top)

## [types](object/def.go) ##

### [null](object/null.go) ###

method  |comment
--------|-------
not     |!

    >> var n = null;
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

    >> var b = true;
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

    >> var i = 123
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

    >> var s = "123"
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
set     |set value of index
len     |length of array
index   |get value by index
not     |!
iter    |get iterator
first   |first value
last    |last value
tail    |remove first value and return rest
push    |append value

    >> var arr = [1,2,3,4,5]
    [1, 2, 3, 4, 5]
    >> arr.push("hello")
    [1, 2, 3, 4, 5, hello]
    >> arr.set(0, 100)

    >> arr
    [100, 2, 3, 4, 5, hello]
    >> arr.len()
    6
    >> arr[0]
    100
    >> arr.not()
    false
    >> !arr
    false
    >> arr.iter()
    array_iter
    >> arr.first()
    100
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
### [array_iter](object/arrayiter.go) ###

method  |comment
--------|-------
not     |!
next    |move to next item
key     |index of current item
value   |value of current item

    >> var arr = [1,2,3]
    [1, 2, 3]
    >> var iter = arr.iter()
    array_iter
    >> iter.key()
    0
    >> iter.value()
    1
    >> iter = iter.next()
    array_iter
    >> iter.value()
    2

[back to top](#id_top)
### [hash](object/hash.go) ###

method  |comment
--------|-------
set     |set value by key
len     |length of hash pairs
index   |get value by key
not     |!
iter    |get iterator

    >> var h = {"k1": 1, "k2": "bbb", "k3": [1,2,3]}
    {k1: 1, k2: bbb, k3: [1, 2, 3]}
    >> h.set("k1", 100)

    >> h["k1"]
    100
    >> h.len()
    3
    >> h.not()
    false
    >> h.iter()
    hash_iter
    >> h["k1"] = 1000

    >> h["k1"]
    1000
    >> h.index("k1")
    1000
    >> !h
    false
    >> len(h)
    3

[back to top](#id_top)
### [hash_iter](object/hashiter.go) ###

method  |comment
--------|-------
not     |!
next    |move to next pair
key     |key of current pair
value   |value of current pair

    >> var h = {"k1":1, "k2":2,"k3":"3"}
    {k1: 1, k2: 2, k3: 3}
    >> var iter = h.iter()
    hash_iter
    >> iter.key()
    k1
    >> iter.value()
    1
    >> iter = iter.next()
    hash_iter
    >> iter.key()
    k2
    >> iter.value()
    2

[back to top](#id_top)
### [builtin](object/builtin.go) ###

method  |comment
--------|-------
not     |!

    >> var l = len
    <built-in function len>
    >> var s = "123"
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

    >> var fn = func(x, y) { return x + y; }
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

    >> var arr = [1,2,3]
    [1, 2, 3]
    >> var fn = arr.index
    <built-in method index of array object>
    >> fn(1)
    2
    >> !fn
    false
    >> fn.not()
    false
    >> var l = arr.len
    <built-in method len of array object>
    >> l()
    3

[back to top](#id_top)

### [file](object/file.go) ###

method  |comment
--------|-------
not     |!
close   |close file
seek    |sets the offset for the next Read or Write on file to offset
size    |get file size
read    |read file data as string
write   |write string to file

    >> var f = open("/usr/local/go/src/os/file.go")
    <open file '/usr/local/go/src/os/file.go', mode 'r' at &{0xc000126120}>
    >> f.size()
    22064
    >> f.seek(22000)
    true
    >> f.read()
    se(); err1 != nil && err == nil {
            err = err1
        }
        return err
    }

    >> f.close()
    >> var f = open("/tmp/tmp.txt", "w+")
    <open file '/tmp/tmp.txt', mode 'w+' at &{0xc000126180}>
    >> f.write("hello world")
    true
    >> f.close()
    >> var f = open("/tmp/tmp.txt", "r")
    <open file '/tmp/tmp.txt', mode 'r' at &{0xc0001261e0}>
    >> f.read()
    hello world
    >> f.close()
    >> remove("/tmp/tmp.txt")

## License ##

Q is licensed under [New BSD License](https://opensource.org/licenses/BSD-3-Clause), a very flexible license to use.

[back to top](#id_top)

## Author ##

* chengzhuo (jobs, yao050421103@163.com)  

[back to top](#id_top)

## More ##

- [Writing An Interpreter In Go](https://interpreterbook.com/)  
- Why named Q ? - Just name it.  

[back to top](#id_top)