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