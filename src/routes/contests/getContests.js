import got from 'got';
import { JSDOM } from 'jsdom';

const schema = {
	description: 'Gets all contests',
	tags: ['contests'],
	summary: 'Gets all contests',
	response: {
		200: {
			description: 'Successful response',
			type: 'array',
			items: {
				type: 'object',
				properties: {
					name: {
						type: 'string',
						description: 'Name of the contest',
					},
					id: {
						type: 'string',
						description: 'ID of the contest',
					},
					start: {
						type: 'string',
						description: 'Start time of the contest',
					},
					end: {
						type: 'string',
						description: 'End time of the contest',
					},
					duration: {
						type: 'string',
						description: 'Duration of the contest',
					},
					type: {
						type: 'string',
						description: 'Type of the contest',
					},
				},
			},
		},
	},
};

export default function(fastify, opts, done) {
	fastify.get('/getContests', { schema }, async (request, reply) => {
		const res = await got('https://codebreaker.xyz/contests');
		const { document } = (new JSDOM(res.body)).window;
		const contestTypes = [
			'ongoing',
			'future',
			'past',
			'practice',
			'collections',
		];

		const contests = [];

		const contestsList = [...document.querySelectorAll('#myTable > tbody')];
		contestsList.pop();

		for (let i = 0; i < contestsList.length; i++) {
			const contest = contestsList[i].children;
			for (let j = 0; j < contest.length; j++) {
				const contestInfo = contest[j].children;
				const contestName = contestInfo[0].textContent.trim();
				const contestId = contestInfo[0].querySelector('a').href.split('/')[4];
				const contestStart = contestInfo[1].textContent.trim();
				const contestEnd = contestInfo[2].textContent.trim();
				const contestDuration = contestInfo[3].textContent.trim();

				contests.push({
					name: contestName,
					id: contestId,
					start: contestStart,
					end: contestEnd,
					duration: (
						// if contestDuration is None, return NA
						// if not, return it appended with ' minutes'
						contestDuration === 'None' ? 'NA' : contestDuration + ' minutes'
					),
					type: contestTypes[i],
				});
			}
		}

		reply.send(contests);
	});
	done();
}