const express = require('express');

const app = express();

app.get('/', (req, res) => {
  res.send('Hello World!');
});

module.exports = app;

// Start the server only if this file is run directly
if (require.main === module) {
  const port = process.env.PORT || 3000;
  app.listen(port, () => {
    console.log(`listening on port ${port}.\n`);
  });
}