import express from 'express';
import bodyParser from 'body-parser';
import { exec } from 'child_process';

const app = express();
const port = 3001;

app.use(bodyParser.json());

app.post('/github-webhook', (req, res) => {
	console.log('Received a webhook request!');
	// Pull the latest changes from GitHub and restart the app
	exec('cd ../. && git pull', (error, stdout, stderr) => {
		if (error) {
			console.error(`Error executing git pull: ${error}`);
			res.status(500).send('Error executing git pull');
		}
		else {
			console.log(`Git pull completed: ${stdout}`);
			// eslint-disable-next-line no-shadow
			exec('cd ../. && yarn install && nohup pm2 restart api', (error, stdout, stderr) => {
				if (error) {
					console.error(`Error executing yarn install: ${error}`);
					res.status(500).send('Error executing yarn install');
				}
				else {
					res.status(200).send('OK');
				}
			});
		}
	});
});

// Start the Express.js app
app.listen(port, '0.0.0.0', () => {
	console.log(`Webhook server started on port ${port}`);
});
