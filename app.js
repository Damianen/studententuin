const express = require('express')
const app = express()
const port = 3000

app.set("view options", {layout: false});
app.use(express.static(__dirname + '/public'));

app.get('/', function(req, res) {
    res.render('index.html');
});

app.listen(port, () => {
  console.log(`Example app listening on port ${port}`)
})