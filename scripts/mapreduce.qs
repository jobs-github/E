func map_iter(arr, accumulated, fn) {
    (len(arr) == 0) ? accumulated : map_iter(arr.tail(), accumulated.push(fn(arr.first())), fn);
};
func map(arr, fn) {
    map_iter(arr, [], fn);
};
func double(x) {
    x * 2;
};

func reduce(arr, result, fn) {
    (len(arr) == 0) ? result : reduce(arr.tail(), fn(result, arr.first()), fn);
};
func add(x, y) {
    x + y;
};
func sum(arr) {
    reduce(arr, 0, add);
};

var a = [1,2,3,4,5];
var result = map(a, double);
println(result);

var rc = sum([1,2,3,4,5]);
println(rc);
