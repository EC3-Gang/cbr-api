import Fastify from 'fastify';
import fastifySwagger from '@fastify/swagger';
import fastifySwaggerUi from '@fastify/swagger-ui';
import fastifyAllow from 'fastify-allow';
import * as dotenv from 'dotenv';
import fs from 'fs/promises';
import { join } from 'path';
import console from 'consola';
const fileUrl = new URL('../package.json', import.meta.url);
const { version } = JSON.parse(await fs.readFile(fileUrl));

dotenv.config();


const fastify = Fastify({
	logger: {
		transport: {
			target: 'pino-pretty',
			options: {
				translateTime: 'SYS:HH:MM:ss Z',
			},
		},
	},
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


const registerAllRoutes = async (path) => {
	// register all routes under ./routes recursively
	const files = await fs.readdir(path);
	for (const file of files) {
		if (file.endsWith('.js')) {
			console.info(`Registering ${path}/${file}`);
			const route = await import(join(process.cwd(), `./${path}/${file}`));
			await fastify.register(route.default, { prefix: '/api' });
		}
		else {
			await registerAllRoutes(`${path}/${file}`);
		}
	}
};

await registerAllRoutes('./src/routes');


const start = async () => {
	try {
		fastify.listen({
			port: process.env.PORT || 3000,
			host: process.env.HOST || '0.0.0.0',
		});
	}
	catch (err) {
		fastify.log.error(err);
		process.exit(1);
	}
};

start();

