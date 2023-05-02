import got from 'got';
import { JSDOM } from 'jsdom';

const schema = {
	description: 'Get a problem',
	tags: ['problems'],
	summary: 'Get a problem',
	querystring: {
		type: 'object',
		required: ['problemId'],
		properties: {
			problemId: {
				type: 'string',
				description: 'ID of the problem',
			},
		},
	},
	response: {
		200: {
			description: 'Successful response',
			type: 'object',
			properties: {
				problemId: {
					type: 'string',
					description: 'ID of the problem',
				},
				title: {
					type: 'string',
					description: 'Title of the problem',
				},
				statement: {
					type: 'string',
					description: 'Problem statement, or link to PDF if it is in PDF format',
				},
				timeLimit: {
					type: 'string',
					description: 'Time limit of the problem',
				},
				memoryLimit: {
					type: 'string',
					description: 'Memory limit of the problem',
				},
				acs: {
					type: 'number',
					description: 'Number of accepted submissions',
				},
				source: {
					type: 'string',
					description: 'Source of the problem',
				},
				subtasks: {
					type: 'array',
					description: 'Subtasks of the problem',
					items: {
						type: 'object',
						properties: {
							subtaskId: {
								type: 'string',
								description: 'ID of the subtask',
							},
							score: {
								type: 'number',
								description: 'Score of the subtask',
							},
						},
					},
				},
			},
		},
	},
};

export default async function(fastify, opts, done) {
	fastify.get('/getProblem', {
		schema,
	}, async (request, reply) => {
		const queryproblemid = request.query.problemId;
		const res = await got(`https://codebreaker.xyz/problem/${queryproblemid}`);
		const { document } = (new JSDOM(res.body)).window;

		// see if it is in PDF format
		// if there is an iframe present in the page, it is in PDF format
		const iframe = document.querySelector('iframe');
		let statement;
		if (iframe) {
			statement = iframe.src;
		}
		else {
			statement = document.getElementById('statement').innerHTML;
		}

		const infoBox = document.querySelector('#problem-right-column > div > div').textContent.split('\n').map((line) => line.trim()).filter((line) => line !== '');
		// time limit
		const timeLimit = infoBox[0].split(':')[1].trim();
		// memory limit
		const memoryLimit = infoBox[1].split(':')[1].trim();
		// acs
		const acs = parseInt(infoBox[2].split(':')[1].trim());
		// source
		const source = infoBox[6];

		// subtasks
		const subTaskTable = document.querySelector('#problem-right-column > table > tbody');
		const subtasks = [];
		if (subTaskTable) {
			const subTaskRows = subTaskTable.querySelectorAll('tr');
			for (const subTaskRow of subTaskRows) {
				const subTaskCells = subTaskRow.querySelectorAll('td');
				const subtaskId = subTaskCells[0].textContent.trim();
				const score = parseInt(subTaskCells[1].textContent.trim());
				subtasks.push({
					subtaskId,
					score,
				});
			}
		}


		const problem = {
			problemId: queryproblemid,
			title: document.querySelector('#statement > h3 > b').textContent.trim(),
			statement,
			timeLimit,
			source,
			memoryLimit,
			acs,
			subtasks,
		};

		reply.send(problem);
	});
	done();
}