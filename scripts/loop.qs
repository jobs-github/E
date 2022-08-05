const n = 10;
const cond = func(i){ i < 10 };
const next = func(i){ i + 1 };

const r0 = for(0, cond, next, func(i, st){
    state(false, println(i + n))
}, state(false, 0));
println(r0);

const r1 = for(0, cond, next, func(i, st){
    state(false, i + st.value())
}, state(false, 0));

println(r1);

const r2 = for(0, cond, next, func(i, st){
    i < 5 ? state(false, i + st.value()) : state(true, i + st.value())
}, state(false, 0));

println(r2);