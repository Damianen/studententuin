import express, { static as _static } from 'express';
import { fileURLToPath } from 'url'
import path from 'path'
const app = express();
const port = 3000;
const __dirname = path.dirname(fileURLToPath(import.meta.url));

app.set("view options", {layout: false});
app.use(_static(path.join(__dirname, 'public')));

app.get('/', function(req, res) {
    res.render('index.html');
});

app.listen(port, () => {
  console.log(`Example app listening on port ${port}`);
});