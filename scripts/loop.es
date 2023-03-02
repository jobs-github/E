const n = 10;

const r0 = loop(10, func(i){
    state(println(i + n))
});
println(r0);

const r1 = loop(10, func(i){
    state(i + i)
});
println(r1);

const r2 = loop(10, func(i){
    i < 5 ? state(i + i) : state(i + i, true)
});
println(r2);