var r = func (x) {
    x * 2;
}(2);

println(r);

func add(x) {
    func(y) {
        x + y;
    };
};

var fn = add(1);
println(fn(2));