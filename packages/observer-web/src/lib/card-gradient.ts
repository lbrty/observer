const CARD_GRADIENTS = [
  "card-grad-indigo",
  "card-grad-foam",
  "card-grad-gold",
  "card-grad-rose",
  "card-grad-sky",
  "card-grad-violet",
  "card-grad-amber",
  "card-grad-teal",
] as const;

export type CardGradient = (typeof CARD_GRADIENTS)[number];

export function cardGradient(index: number): CardGradient {
  return CARD_GRADIENTS[index % CARD_GRADIENTS.length];
}
