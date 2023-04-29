import got from 'got';
import { JSDOM } from 'jsdom';

const schema = {
	description: 'Gets all problems',
	tags: ['problems'],
	summary: 'Gets all problems',
	response: {
		200: {
			description: 'Successful response',
			type: 'array',
			items: {
				type: 'object',
				properties: {
					problemid: {
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
					tags: {
						type: 'array',
						description: 'Tags of the problem',
						items: {
							type: 'string',
						},
					},
					type: {
						type: 'string',
						description: 'Type of the problem',
					},
					acs: {
						type: 'number',
						description: 'Number of accepted submissions',
					},
				},
			},
		},
	},
};

export default async function(fastify, opts, done) {
	fastify.get('/getProblems', {
		schema,
	}, async (request, reply) => {
		const res = await got('https://codebreaker.xyz/problems');
		const { document } = (new JSDOM(res.body)).window;

		const problems = [];
		const problemTable = document.querySelector('#myTable > tbody');
		for (const row of problemTable.rows) {
			const problemid = row.cells[0].textContent.trim();
			const tagAndTitle = row.cells[2].textContent.trim();
			const [tag, title] = tagAndTitle.split('\n\t\t\t\t\t\t\n\t\t\t\t\t\t\t');
			// remove first 7 chars of tag and split by ', '
			const tags = tag.slice(7).split(', ');
			const source = row.cells[3].textContent.trim();
			const type = row.cells[4].textContent.trim();
			const acs = parseInt(row.cells[5].textContent.trim());

			problems.push({
				problemid,
				title,
				source,
				type,
				acs,
				tags,
			});
		}

		return problems;
	});
	done();
}