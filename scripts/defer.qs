func test_defer(x) {
    defer func() {
        println(x);
    }();
    defer println(x * x * x);
    println(x * 2);
};

func test_file() {
    var f = open("/usr/local/go/src/os/file.go");
    defer func() {
        f.close();
        println("file closed");
    }();
    println(f);
};

func main() {
    test_defer(2);
    test_file();
};

main();