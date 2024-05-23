import express from 'express';
import React from 'react';
import { renderToString } from 'react-dom/server';
import AppServer from '../src/AppServer';

const app = express();
const PORT = process.env.PORT || 3001;

function getHTML(content) {
    return `
    <!DOCTYPE html>
    <html lang="en">
      <head>
        <meta charset="UTF-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <title>React SSR</title>
      </head>
      <body>
        <div id="root">${content}</div>
      </body>
    </html>
  `;
}

app.get('/', (req, res) => {
  const content = renderToString(<AppServer />);
  res.send(getHTML(content));
});

app.listen(PORT, () => {
  console.log(`Server is listening on port ${PORT}`);
});