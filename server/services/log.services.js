import fs from "fs/promises"; // gebruik fs/promises voor asynchrone file operaties

// Helper function to find the newest log files
const getNewestLogFiles = async (dir) => {
  const files = await fs.readdir(dir);

  const stdoutFiles = files.filter(
    (file) =>
      file.startsWith("STUDENTENTUIN01-") &&
      file.includes("-stdout-") &&
      file.endsWith(".txt")
  );

  const stderrFiles = files.filter(
    (file) =>
      file.startsWith("STUDENTENTUIN01-") &&
      file.includes("-stderr-") &&
      file.endsWith(".txt")
  );

  stdoutFiles.sort((a, b) => {
    const timestampA = parseInt(a.split("-")[3].replace(".txt", ""));
    const timestampB = parseInt(b.split("-")[3].replace(".txt", ""));
    return timestampB - timestampA;
  });

  stderrFiles.sort((a, b) => {
    const timestampA = parseInt(a.split("-")[3].replace(".txt", ""));
    const timestampB = parseInt(b.split("-")[3].replace(".txt", ""));
    return timestampB - timestampA;
  });

  const newestStdout =
    stdoutFiles.length > 0 ? path.join(dir, stdoutFiles[0]) : null;
  const newestStderr =
    stderrFiles.length > 0 ? path.join(dir, stderrFiles[0]) : null;

  return { newestStdout, newestStderr };
}
