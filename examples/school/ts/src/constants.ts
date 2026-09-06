import { checkIntegerRange, readGlobalValues } from "./settings.js";

// SchoolConstants Holds shared Meaning, independent of Deployment.
export type SchoolConstants = Readonly<{
  minimumAgeYears: number;
}>;

export const SchoolValues = readGlobalValues<SchoolConstants>(
  new URL("../../constants/school.json", import.meta.url), {
  minimumAgeYears: checkIntegerRange(1, 150),
  },
);
