const express = require('express');
const bodyParser = require('body-parser');
const { exec } = require('child_process');
const crypto = require('crypto');

const app = express();
const port = process.env.PORT || 3001;
const secret = process.env.GITHUB_WEBHOOK_SECRET;

app.use(bodyParser.urlencoded({ extended: false }));
app.use(bodyParser.json());

app.post('/webhook', (req, res) => {
  const event = req.headers['x-github-event'];
  const branch = req.body.ref.split('/').pop();
  const signature = req.headers['x-hub-signature'];

  const hmac = crypto.createHmac('sha1', secret);
  const digest = Buffer.from('sha1=' + hmac.update(JSON.stringify(req.body)).digest('hex'), 'utf8');
  const checksum = Buffer.from(signature, 'utf8');

  if (checksum.length !== digest.length || !crypto.timingSafeEqual(digest, checksum)) {
    console.log('Invalid signature');
    res.status(400).send('Invalid signature');
    return;
  }

  if (event === 'push' && branch === 'main') {
    console.log('Received push event to main branch');
    // Stop the server
    exec('cd ../. && pm2 stop all', (error, stdout, stderr) => {
      if (error) {
        console.error(`Error stopping server: ${error}`);
        res.status(500).send('Error stopping server');
      } else {
        console.log(`Stopped server: ${stdout}`);
        // Pull the latest code from GitHub
        exec('cd ../. git pull', (error, stdout, stderr) => {
          if (error) {
            console.error(`Error pulling code from GitHub: ${error}`);
            res.status(500).send('Error pulling code from GitHub');
          } else {
            console.log(`Pulled code from GitHub: ${stdout}`);
            // Install dependencies
            exec('cd ../. yarn install', (error, stdout, stderr) => {
              if (error) {
                console.error(`Error installing dependencies: ${error}`);
                res.status(500).send('Error installing dependencies');
              } else {
                console.log(`Installed dependencies: ${stdout}`);
                // Start the server
                exec('cd ../. yarn start', (error, stdout, stderr) => {
                  if (error) {
                    console.error(`Error starting server: ${error}`);
                    res.status(500).send('Error starting server');
                  } else {
                    console.log(`Started server: ${stdout}`);
                    res.status(200).send('Server restarted');
                  }
                });
              }
            });
          }
        });
      }
    });
  } else {
    console.log(`Received unsupported event or branch: ${event} ${branch}`);
    res.status(200).send('OK');
  }
});

app.listen(port, '0.0.0.0', () => {
  console.log(`Webhook receiver listening at http://0.0.0.0:${port}`);
});

