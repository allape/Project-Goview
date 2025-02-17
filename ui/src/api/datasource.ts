import Crudy, { aapi } from "@allape/gocrud-react";
import { SERVER_URL } from "@allape/gocrud-react/src/config";
import IDatasource, { IFileInfo } from "../model/datasource.ts";
import IPreview from "../model/preview.ts";
import { URLString } from "./common.ts";

export const DatasourceCrudy = new Crudy<IDatasource>(
  `${SERVER_URL}/datasource`,
);

export function readDir(
  id: IDatasource["id"],
  wd: string,
): Promise<IFileInfo[]> {
  return aapi.get(`${SERVER_URL}/datasource/readdir/${id}${wd}`);
}

export function getFileURLFromDatasource(
  id: IDatasource["id"],
  filename: string,
): URLString {
  return `${SERVER_URL}/datasource/by-ds/${id}${filename}`;
}

export function getFileURLByKey(key: IPreview["key"]): URLString {
  return `${SERVER_URL}/datasource/by-key/${encodeURIComponent(key)}`;
}
