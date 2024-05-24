import express from "express";

import router from "./server/router.js";

const app = express();
const port = 3001 || process.env.PORT;

app.set("view engine", "ejs");
app.use(express.static("public"));
app.use(router);

app.listen(port, () => {
  console.log(`Server listening on port ${port}`);
});
