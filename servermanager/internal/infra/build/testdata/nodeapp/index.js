const http = require('http');

const port = Number(process.env.PORT || 3000);
http
	.createServer((req, res) => {
		console.log(`request ${req.method} ${req.url}`);
		res.end('hello from stt-build-fixture\n');
	})
	.listen(port, () => console.log(`fixture listening on ${port}`));
