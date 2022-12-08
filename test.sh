for file in ./scripts/*.es
do
    echo RUN "$file"
    ./escript $file
    echo ""
done