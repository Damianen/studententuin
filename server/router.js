import { Router } from 'express';
const router = Router();

router.get('/', (req, res, next) => {
    res.render("index.ejs", {
        content: "react aint workin"
    });
});

export default router;