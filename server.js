import express from "express";

import router from "./server/router.js";

const app = express();
const port = process.env.PORT || 3001;

app.set("view engine", "ejs");
app.use(express.static("public"));
app.use(router);

app.listen(port, () => {
  console.log(`Server listening on port ${port}`);
});
