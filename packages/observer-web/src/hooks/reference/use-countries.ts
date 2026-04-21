import type { Country, CreateCountryInput, UpdateCountryInput } from "@/types/reference";

import { makeReferenceHooks } from "./make-reference-hooks";

const {
  useList: useCountries,
  useCreate: useCreateCountry,
  useUpdate: useUpdateCountry,
  useDelete: useDeleteCountry,
} = makeReferenceHooks<Country, CreateCountryInput, UpdateCountryInput>("countries");

export { useCountries, useCreateCountry, useUpdateCountry, useDeleteCountry };
