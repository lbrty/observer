import type {
  CreateOfficeInput,
  Office,
  UpdateOfficeInput,
} from "@/types/reference";

import { makeReferenceHooks } from "./make-reference-hooks";

const {
  useList: useOffices,
  useCreate: useCreateOffice,
  useUpdate: useUpdateOffice,
  useDelete: useDeleteOffice,
} = makeReferenceHooks<Office, CreateOfficeInput, UpdateOfficeInput>("offices");

export { useOffices, useCreateOffice, useUpdateOffice, useDeleteOffice };
