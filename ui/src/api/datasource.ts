import Crudy, { aapi } from "@allape/gocrud-react";
import { SERVER_URL } from "@allape/gocrud-react/src/config";
import IDatasource, { IFileInfo } from "../model/datasource.ts";

export const DatasourceCrudy = new Crudy<IDatasource>(
  `${SERVER_URL}/datasource`,
);

export function readDir(
  id: IDatasource["id"],
  wd: string,
): Promise<IFileInfo[]> {
  return aapi.get(
    `${SERVER_URL}/datasource/readdir/${id}${wd}`,
  );
}

export function getFileURLFromDatasource(id: IDatasource['id'], filename: string): string {
  return `${SERVER_URL}/datasource/fetch/${id}${filename}`;
}
