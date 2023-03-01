const n = 10;
const cond = func(i){ i < 10 };
const next = func(i){ i + 1 };

const r0 = for(state(0), cond, next, func(i, st){
    state(println(i + n))
});
println(r0);

const r1 = for(state(0), cond, next, func(i, st){
    state(i + st.value())
});

println(r1);

const r2 = for(state(0), cond, next, func(i, st){
    i < 5 ? state(i + st.value()) : state(i + st.value(), true)
});

println(r2);