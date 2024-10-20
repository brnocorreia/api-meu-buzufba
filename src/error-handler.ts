import { FastifyInstance } from "fastify";

import { ZodError } from "zod";
import { BadRequest, NotFound } from "./routes/_errors/errors";

type FastifyErrorHandler = FastifyInstance["errorHandler"];

export const errorHandler: FastifyErrorHandler = (error, request, reply) => {
  if (error instanceof ZodError) {
    return reply.status(400).send({
      message: "Error during validation",
      errors: error.flatten().fieldErrors,
    });
  }

  if (error instanceof BadRequest) {
    return reply.status(400).send({ message: error.message });
  }

  if (error instanceof NotFound) {
    return reply.status(404).send({ message: error.message });
  }

  return reply.status(500).send({ message: "Internal Server Error!" });
};
