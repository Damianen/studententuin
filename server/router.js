import { Router } from 'express';
const router = Router();

router.get('/', (req, res, next) => {
    res.render("index.ejs");
});

router.get('/requestForm', (req, res, next) => {
    res.render("index.ejs");
});

router.get('/login', (req, res, next) => {
    res.render("index.ejs");
});

router.get('/handleiding', (req, res, next) => {
    res.render("index.ejs");
});

router.get('/manage', (req, res, next) => {
    res.render("index.ejs");
});

router.get('/about', (req, res, next) => {
    res.render("index.ejs");
});

router.get('/contact', (req, res, next) => {
    res.render("index.ejs");
});

export default router;