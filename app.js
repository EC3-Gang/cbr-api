'use strict';

import { join } from 'path';
import AutoLoad from '@fastify/autoload';
// define __dirname for es6 modules
const __dirname = new URL('.', import.meta.url).pathname;

export const options = {};

export default async function(fastify, opts) {
	fastify.register(AutoLoad, {
		dir: join(__dirname, 'plugins'),
		options: Object.assign({}, opts),
	});

	// This loads all plugins defined in routes
	// define your routes in one of these
	fastify.register(AutoLoad, {
		dir: join(__dirname, 'routes'),
		options: Object.assign({}, opts),
	});
}
