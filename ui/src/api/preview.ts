import { aapi } from "@allape/gocrud-react";
import { SERVER_URL } from "@allape/gocrud-react/src/config";
import IDatasource from "../model/datasource.ts";
import IPreview from "../model/preview.ts";

export function generatePreview(
  datasourceId: IDatasource["id"],
  filename: string,
): Promise<IPreview> {
  return aapi.get(`${SERVER_URL}/preview/${datasourceId}${filename}`, {
    method: "PUT",
  });
}
