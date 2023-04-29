import Fastify from 'fastify';
import fastifySwagger from '@fastify/swagger';
import fastifySwaggerUi from '@fastify/swagger-ui';
import fastifyAllow from 'fastify-allow';
import fs from 'fs/promises';
const fileUrl = new URL('../package.json', import.meta.url);
const { version } = JSON.parse(await fs.readFile(fileUrl));


const fastify = Fastify({
	logger: true,
});

await fastify.register(fastifyAllow);

await fastify.register(fastifySwagger, {
	exposeRoute: true,
	swagger: {
		swagger: '2.0',
		info: {
			title: 'Codebreaker.xyz API',
			description: 'Unofficial Codebreaker.xyz API',
			version,
		},
	},
});

await fastify.register(fastifySwaggerUi, {
	routePrefix: '/',
});

fastify.get('/helloworld', async (request, reply) => {
	return { hello: 'world' };
});

// register all routes under ./routes recursively
const files = await fs.readdir('./src/routes');
for (const file of files) {
	if (file.endsWith('.js')) {
		const route = await import(`./routes/${file}`);
		await fastify.register(route.default, { prefix: '/api' });
	}
}


const start = async () => {
	try {
		await fastify.listen({
			port: 3000,
		});
	}
	catch (err) {
		fastify.log.error(err);
		process.exit(1);
	}
};

start();

