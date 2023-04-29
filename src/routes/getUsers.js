import got from 'got';
import { JSDOM } from 'jsdom';
import capitalise from '../utils/capitalise.js';

const schema = {
	description: 'Get a user',
	tags: ['users'],
	summary: 'Get a user',
	querystring: {
		type: 'object',
		required: ['name'],
		properties: {
			name: {
				type: 'string',
				description: 'Name of the user',
			},
		},
	},
	response: {
		200: {
			description: 'Successful response',
			type: 'object',
			properties: {
				name: {
					type: 'string',
					description: 'Name of the user',
				},
				school: {
					type: 'string',
					description: 'School of the user',
				},
				role: {
					type: 'string',
					description: 'Role of the user',
				},
				country: {
					type: 'string',
					description: 'Country of the user',
				},
				solvedProblems: {
					type: 'array',
					description: 'Solved problems of the user',
					items: {
						type: 'string',
					},
				},
			},
		},
	},
};

export default async function(fastify, opts, done) {
	fastify.get('/getUser', {
		schema,
	}, async (request, reply) => {
		const queryname = request.query.name;
		const res = await got(`https://codebreaker.xyz/profile/${queryname}`);
		const { document } = (new JSDOM(res.body)).window;

		const name = document.querySelector('#profile-card > div > table > tbody > tr:nth-child(1) > td').textContent.trim();
		const school = document.querySelector('#profile-card > div > table > tbody > tr:nth-child(3) > td').textContent.trim();
		const role = document.querySelector('#profile-card > div > table > tbody > tr:nth-child(5) > td').textContent.trim();
		const country = document.querySelector('#profile-card > div > table > tbody > tr:nth-child(7) > td').textContent.trim();

		const solvedProblems = [];
		const problemsTable = document.querySelector('body > div.container-fluid.p-4 > div > div.col-sm-7 > table > tbody');
		for (const row of problemsTable.rows) {
			solvedProblems.push(row.cells[0].textContent.trim());
		}

		reply.send({
			name: capitalise(name, true),
			school: capitalise(school, true),
			role: capitalise(role),
			country: capitalise(country, true),
			solvedProblems,
		});
	});
	done();
}