import React from 'react';
import ReactDOM from 'react-dom/client';
import '../public/style.css';

fetch('http://localhost:3001/api/user/r@struijlaart.n')
    .then(response => {
      if (!response.ok) {
        throw new Error('Network response was not ok ' + response.statusText);
      }
      return response.json();
    })
    .then(data => {
      console.dir("succes");
    })
    .catch(error => {
      console.dir(error);
    });

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
    <div className="bg-primary-green">
    <h1>Hello World</h1>
  </div>
);