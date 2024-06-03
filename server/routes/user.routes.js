import { Router } from "express";
import userService from "../services/user.service.js";

const router = Router();

router.get("/api/user/:userEmail", userService.getByEmail);

export default router;
