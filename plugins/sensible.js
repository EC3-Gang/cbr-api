'use strict';

import fp from 'fastify-plugin';
import fastifySensible from '@fastify/sensible';

export default fp(async function(fastify, opts) {
	fastify.register(fastifySensible, {
		errorHandler: false,
	});
});
