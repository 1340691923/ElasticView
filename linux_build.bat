cmd /k "cd vue && npm run build:stage && cd .. &&set goos=linux&& go build  -o ElasticView"
echo "build success"
