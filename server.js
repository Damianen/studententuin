import express from "express";

import router from './server/router.js';
import userRoutes from './server/routes/user.routes.js';

const app = express();
const port = process.env.PORT || 3001;

app.set("view engine", "ejs");
app.use(express.static("public"));
app.use(router);
app.use(userRoutes);

app.listen(port, () => {
  console.log(`Server listening on port ${port}`);
});
