import { Router } from "express";
import httpServer from "../../server.js";
import { Server } from "socket.io";
import { getNewestLogFiles } from "../services/log.services.js";
const router = Router();

router.get("/api/logs", async (req, res) => {

    const io = new Server(httpServer, {});
    try {
        const logsDir = path.resolve(__dirname, "../logs");
        console.log(`Looking for logs in: ${logsDir}`);

        try {
            await fs.access(logsDir);
        } catch (err) {
            console.error(`Directory does not exist: ${logsDir}`, err);
            return res.status(404).json({ error: "Logs directory does not exist" });
        }

        const { newestStdout, newestStderr } = await getNewestLogFiles(logsDir);

        if (!newestStdout && !newestStderr) {
            return res.status(404).json({ error: "No log files found" });
        }

        const stdoutData = newestStdout
            ? await fs.readFile(newestStdout, "utf8")
            : null;
        const stderrData = newestStderr
            ? await fs.readFile(newestStderr, "utf8")
            : null;

        res.json({ stdout: stdoutData, stderr: stderrData });
    } catch (err) {
        console.error("Error reading log files:", err);
        res
            .status(500)
            .json({ error: "Error reading log files", details: err.message });
    }
});

export default router;
