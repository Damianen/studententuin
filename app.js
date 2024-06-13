import path from "path";
import fs from "fs";
import express from "express";

import render from "./views/App.jsx";

const PORT = process.env.PORT || 3000;
const app = express();
const __dirname = path.dirname(fileURLToPath(import.meta.url));

app.get("/", (req, res) => {
    fs.readFile(path.resolve("./public/index.html"), "utf8", (err, data) => {
        if (err) {
            console.error(err);
            return res.status(500).send("An error occurred");
        }

        render();

        return res.send(data);
    });
});

app.use(
    express.static(path.resolve(__dirname, ".", "dist"), { maxAge: "30d" })
);

app.listen(PORT, () => {
    console.log(`Server is listening on port ${PORT}`);
});

