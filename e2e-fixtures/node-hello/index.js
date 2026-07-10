const http = require('http');

const port = Number(process.env.PORT || 3000);
http
	.createServer((req, res) => {
		console.log(`request ${req.method} ${req.url}`);
		res.end('hello from stt-sample-node\n');
	})
	.listen(port, () => console.log(`stt-sample-node listening on ${port}`));
