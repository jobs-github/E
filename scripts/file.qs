func test_write(url) {
    println("test_write");
    var f = open(url, "w+");
    if (!f) {
        return;
    }
    defer f.close()
    var r = f.write("hello world");
    if (!r) {
        println("write fail")
    }
};

func test_read(url) {
    println("test_read");
    var f = open(url, "r");
    if (!f) {
        return;
    }
    defer f.close()
    var sz = f.size();
    println("size: ", sz);
    var data = f.read();
    println(data);
};

func test_seek(url) {
    println("test_seek");
    var f = open(url, "r");
    if (!f) {
        return;
    }
    defer f.close()
    var off = f.size() / 2;
    println("offset: ", off);
    f.seek(off);
    var data = f.read();
    println(data);
};

func main() {
    var url = "tmp.txt";
    defer remove(url);
    
    test_write(url);
    test_read(url);
    test_seek(url);
};

main();