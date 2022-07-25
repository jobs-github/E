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