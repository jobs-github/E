import "utils"
import libmath "math"

func hello() {
    println("hello");
};

func main() {
    var fn = utils.hello;
    hello();
    fn();
    utils.hello();
    println(utils.EOF);
    var r = libmath.add(1,2);
    println(r)
};

main();