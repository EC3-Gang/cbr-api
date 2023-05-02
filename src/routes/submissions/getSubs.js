import got from 'got';

const schema = {
	description: 'Gets all submissions for a problem (takes a while on questions with more subs)',
	tags: ['submissions'],
	summary: 'Gets all submissions for a problem',
	// optional querystring (problem id)
	querystring: {
		type: 'object',
		required: ['problemId'],
		properties: {
			problemId: {
				type: 'string',
				description: 'ID of the problem',
			},
			ac: {
				type: 'boolean',
				description: 'Whether to only get AC submissions',
			},
		},
	},
	response: {
		200: {
			description: 'Successful response',
			type: 'array',
			items: {
				type: 'object',
				properties: {
					id: {
						type: 'string',
						description: 'ID of the submission',
					},
					date: {
						type: 'string',
						description: 'Date of the submission in ISO 8601 format',
					},
					problemId: {
						type: 'string',
						description: 'Problem ID of the submission',
					},
					user: {
						type: 'string',
						description: 'Username of the user who submitted',
					},
					language: {
						type: 'string',
						description: 'Language of the submission',
					},
					maxTime: {
						type: 'string',
						description: 'Maximum time of the submission',
					},
					maxMemory: {
						type: 'string',
						description: 'Maximum memory of the submission',
					},
					score: {
						type: 'string',
						description: 'Score of the submission',
					},
				},
			},
		},
	},
};

export default function(fastify, opts, done) {
	fastify.get('/getSubmissions', { schema }, async (request, reply) => {
		const { problemId } = request.query;
		// get json from localhost:3002/attempts
		const { body } = await got(`http://localhost:3002/attempts?problem=${problemId}`, {
			responseType: 'json',
		});

		/* JSON returned:
		[
			{
			  id: 128736,
				submission: '2022-06-22T20:04:22Z',
				username: 'errorgorn',
				problem: 'manhattancompass',
				score: 100,
				language: 'C++17',
				max_time: 0.458,
				max_memory: 33
			}
		]
		*/
		console.log(body);
		// rename keys in the body and -1s in max_time and max_memory to N/A, and sort by sub id
		const submissions = body.map(submission => {
			return {
				id: submission.id,
				date: submission.submission,
				problemId: submission.problem,
				user: submission.username,
				language: submission.language,
				maxTime: submission.max_time === -1 ? 'N/A' : submission.max_time,
				maxMemory: submission.max_memory === -1 ? 'N/A' : submission.max_memory,
				score: submission.score,
			};
		}).sort((a, b) => b.id - a.id);

		// if ac is true, filter out non-AC submissions
		if (request.query.ac) {
			reply.send(submissions.filter(submission => submission.score === 100));
		}
		else {
			reply.send(submissions);
		}
	});
	done();
}