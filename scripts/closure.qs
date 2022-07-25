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