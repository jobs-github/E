var n = 10;
var cond = func(i){ i < 10 };
var next = func(i){ i + 1 };

var r0 = for(0, cond, next, func(i, st){
    state(false, println(i + n))
}, state(false, 0));
println(r0);

var r1 = for(0, cond, next, func(i, st){
    state(false, i + st.value())
}, state(false, 0));

println(r1);

var r2 = for(0, cond, next, func(i, st){
    i < 5 ? state(false, i + st.value()) : state(true, i + st.value())
}, state(false, 0));

println(r2);