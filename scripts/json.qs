const s = "{\"k3\":{\"k31\":true,\"k32\":[1,2,3]},\"k2\":\"2\",\"k1\":1}"; 
const obj = loads(s);
println(obj);
const ss = dumps(obj);
println(ss);

const obj1 = { "k1": null, "k2": 123 };
const objstr = dumps(obj1);
println(objstr);
const obj2 = loads(objstr);
println(obj2);
println(obj1 == obj2);