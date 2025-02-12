import Crudy, { aapi } from "@allape/gocrud-react";
import { SERVER_URL } from "@allape/gocrud-react/src/config";
import IDatasource from "../model/datasource.ts";
import IPreview from "../model/preview.ts";

export const PreviewCrudy = new Crudy<IPreview>(`${SERVER_URL}/preview`);

export function generatePreview(
  datasourceId: IDatasource["id"],
  filename: string,
): Promise<IPreview> {
  return aapi.get(`${SERVER_URL}/preview/from-ds/${datasourceId}${filename}`, {
    method: "PUT",
  });
}

export function getPreviewURLByDatasource(
  id: IDatasource["id"],
  filename: string,
): string {
  return `${SERVER_URL}/preview/by-ds/${id}${filename}`;
}

export function getPreviewURLByKey(key: IPreview["key"]): string {
  return `${SERVER_URL}/preview/by-key/${encodeURIComponent(key)}`;
}
