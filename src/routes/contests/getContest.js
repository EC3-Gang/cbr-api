import got from 'got';
import { JSDOM } from 'jsdom';

const schema = {
	description: 'Get a contest',
	tags: ['contests'],
	summary: 'Get a contest',
	querystring: {
		type: 'object',
		required: ['contestid'],
		properties: {
			contestid: {
				type: 'string',
				description: 'ID of the contest',
			},
		},
	},
	response: {
		200: {
			description: 'Successful response',
			type: 'object',
			properties: {
				contestId: {
					type: 'string',
					description: 'ID of the contest',
				},
				name: {
					type: 'string',
					description: 'Name of the contest',
				},
				description: {
					type: 'string',
					description: 'Description of the contest',
				},
				problems: {
					type: 'array',
					description: 'Problems in the contest',
					items: {
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
							source: {
								type: 'string',
								description: 'Source of the problem',
							},
							type: {
								type: 'string',
								description: 'Type of the problem',
							},
						},
					},
				},
			},
		},
	},
};

export default function(fastify, opts, done) {
	fastify.get('/getContest', { schema }, async (request, reply) => {
		const contestId = request.query.contestid;
		const res = await got(`https://codebreaker.xyz/contest/${contestId}`);
		const { document } = (new JSDOM(res.body)).window;

		const name = document.querySelector('body > div:nth-child(4) > h1').textContent.trim().split(' ').slice(0, -1).join(' ');
		const description = document.querySelector('#contestdescription > div').textContent.trim();

		const problemTable = document.querySelector('#myTable > tbody');

		const problems = [];
		for (const problem of problemTable.children) {
			const problemId = problem.children[0].textContent.trim();
			const title = problem.children[2].textContent.trim();
			const source = problem.children[3].textContent.trim();
			const type = problem.children[4].textContent.trim();

			problems.push({
				problemId,
				title,
				source,
				type,
			});
		}

		reply.send({
			contestId,
			name,
			description,
			problems,
		});
	});
	done();
}