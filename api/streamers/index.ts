import type { VercelRequest, VercelResponse } from "@vercel/node";
import client from "../client";

export default async (request: VercelRequest, response: VercelResponse) => {
  const streamers = await client.accounts.findMany({
    select: { id: true, platform_id: true, username: true },
    skip:
      Number((request.query.page as string) ?? 0) *
      Number((request.query.limit as string) ?? 25),
    take: Number((request.query.limit as string) ?? 25),
    where: { NOT: { platform_id: null } },
  });
  response.status(200).json(
    streamers.map(streamer => ({
      ...streamer,
      platform_id: Number(streamer.platform_id) ?? streamer.platform_id,
    }))
  );
};
