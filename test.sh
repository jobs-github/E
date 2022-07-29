for file in ./scripts/*.qs
do
    echo RUN "$file"
    ./escript $file
    echo ""
done