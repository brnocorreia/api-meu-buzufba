import { FastifyInstance } from "fastify";
import { ZodTypeProvider } from "fastify-type-provider-zod";
import z from "zod";
import { prisma } from "../lib/prisma";

export async function getAllRoutes(app: FastifyInstance) {
  app.withTypeProvider<ZodTypeProvider>().get(
    "/routes",
    {
      schema: {
        summary: "Get all routes",
        tags: ["routes"],
        response: {
          200: z.array(
            z.object({
              routeId: z.number(),
              routeName: z.string(),
              tripLength: z.number(),
              departureLocation: z.string(),
              arrivalLocation: z.string(),
              statusCd: z.string(),
              createdAt: z.date(),
              updatedAt: z.date(),
              departures: z.array(
                z.object({
                  departureTime: z.string(),
                  obs: z.string().nullable(),
                  statusCd: z.string(),
                })
              ),
            })
          ),
        },
      },
    },
    async (request, reply) => {
      const routes = await prisma.route.findMany({
        include: {
          departures: {
            select: {
              departureTime: true,
              obs: true,
              statusCd: true,
            },
          },
        },
      });
      return reply.status(200).send(routes);
    }
  );
}
