import dotenv from "dotenv";

export function config() {
  try {
    const result = dotenv.config();

    if (result.error) {
      throw result.error;
    }

    return result.parsed;
  } catch (err) {
    throw err;
  }
}
