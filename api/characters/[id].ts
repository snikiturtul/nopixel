import type { VercelRequest, VercelResponse } from "@vercel/node";
import client from "../client";

export default async (request: VercelRequest, response: VercelResponse) => {
  const character_name = (request.query.id as string).toLowerCase();
  const character = await client.characters.findFirst({
    select: {
      id: true,
      accounts: true,
      name: true,
      search_phrases: true,
    },
    where: {
      search_phrases: {
        hasEvery: character_name.split("_"),
      },
    },
  });
  if (character) {
    response.status(200).json({
      ...character,
      account: character.accounts,
      accounts: undefined,
      link:
        character.accounts.platform === 1
          ? `twitch.tv/${character.accounts.username}`
          : undefined,
    });
  } else response.status(404).send("character not found");
};
