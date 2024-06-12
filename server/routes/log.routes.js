import { Router } from "express";
import { getNewestLogFiles } from "../services/log.services.js";
import path from 'path';
import fs from "fs/promises"; // gebruik fs/promises voor asynchrone file operaties
const router = Router();
const __dirname = path.resolve();

router.get("/api/logs", async (req, res) => {
    console.log("api");
    res.sendStatus(200);
});

const getLatestLogs = async (io) => {
    try {
        const logsDir = path.resolve(__dirname, "./logs");
        console.log(`Looking for logs in: ${logsDir}`);

        try {
            await fs.access(logsDir);
        } catch (err) {
            console.error(`Directory does not exist: ${logsDir}`, err);
            return { error: "Logs directory does not exist" };
        }

        const { newestStdout, newestStderr } = await getNewestLogFiles(logsDir);

        if (!newestStdout && !newestStderr) {
            return { error: "No log files found" };
        }

        const stdoutData = newestStdout
            ? await fs.readFile(newestStdout, "utf8")
            : null;
        const stderrData = newestStderr
            ? await fs.readFile(newestStderr, "utf8")
            : null;

        return { stdout: stdoutData, stderr: stderrData };
    } catch (err) {
        console.error("Error reading log files:", err);
        return { error: err }
    }
}

export default getLatestLogs;
