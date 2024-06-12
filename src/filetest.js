import fs from "fs";
import fsFileTree from "fs-file-tree";

function buildFileTree(path) {
    const tree = {};
    const files = fs.readdirSync(path);
    files.forEach((file) => {
        const filePath = `${path}/${file}`;
        if (fs.statSync(filePath).isDirectory()) {
            tree[file] = buildFileTree(filePath);
        } else {
            tree[file] = file;
        }
    });
    return tree;
}

console.log(buildFileTree("./test"));
