import { IBase } from "@allape/gocrud";
import { IBaseSearchParams } from "@allape/gocrud/src/model.ts";
import IDatasource from "./datasource.ts";

export default interface IPreview extends IBase {
  datasourceId: string;
  key: string;
  digest: string;
  cover: string;
  ffprobeInfo: string;
}

export interface IPreviewSearchParams extends IBaseSearchParams {
  datasourceId?: IDatasource["id"];
  key?: IPreview["key"];
  digest?: IPreview["digest"];
  ffprobeInfo?: IPreview["ffprobeInfo"];
}
