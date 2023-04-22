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
