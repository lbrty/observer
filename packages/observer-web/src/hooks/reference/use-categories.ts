import type { Category, CreateCategoryInput, UpdateCategoryInput } from "@/types/reference";

import { makeReferenceHooks } from "./make-reference-hooks";

const {
  useList: useCategories,
  useCreate: useCreateCategory,
  useUpdate: useUpdateCategory,
  useDelete: useDeleteCategory,
} = makeReferenceHooks<Category, CreateCategoryInput, UpdateCategoryInput>("categories");

export { useCategories, useCreateCategory, useUpdateCategory, useDeleteCategory };
