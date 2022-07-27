var r = func (x) {
    return x * 2;
}(2);

println(r);

func add(x) {
    return func(y) {
        return x + y;
    };
};

var fn = add(1);
println(fn(2));