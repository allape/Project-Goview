import type { IDatasource } from './datasource';
import ajax, { type IBaseModel } from './http';

export interface IPreview extends IBaseModel {
  id: number;
  datasourceId: IDatasource['id'];
  key: string;
  digest: string;
  cover: string;
}


export function gen(dsid: IDatasource['id'], file: string): Promise<IPreview> {
  return ajax(`/preview/gen/${dsid}/${file}`);
}
