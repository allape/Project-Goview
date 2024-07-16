import type { LV } from '../config/lv';
import { BASE_URL } from '../config/server';
import ajax, { type IBaseModel } from './http';

export interface IFile {
  isDir: boolean;
  name: string;
  size: number;

  _displayName: string;
  _path: string;
  _preview: string;
}

export type Type = 'dufs' | 'local';

export interface IDatasource extends IBaseModel {
  name: string;
  type: Type;
  cwd: string;
}

export const Types: LV[] = [
  { label: 'Dufs', value: 'dufs' },
  { label: 'Local', value: 'local' },
];

export function getAll(): Promise<IDatasource[]> {
  return ajax('/datasource/all');
}

export function save(ds: IDatasource): Promise<IDatasource> {
  return ajax('/datasource/save', {
    method: 'POST',
    body: JSON.stringify(ds),
  });
}

export function stat(dsid: IDatasource['id'], file: string): Promise<IFile> {
  return ajax(`/datasource/stat/${dsid}/${file}`);
}

export async function ls(dsid: IDatasource['id'], cwd: string): Promise<IFile[]> {
  const files = await ajax<IFile[]>(`/datasource/ls/${dsid}/${cwd}`);
  return files.map((file: IFile) => {
    file._displayName = `${file.name}${file.isDir ? '/' : ''}`;
    file._path = `${cwd || '/'}${encodeURIComponent(file.name)}`;
    file._preview = `${BASE_URL}/preview/get/${dsid}/${file._path}?t=${Date.now()}`;
    return file;
  });
}

export function cat(dsid: IDatasource['id'], file: string): Promise<Blob> {
  return fetch(`/datasource/cat/${dsid}/${file}`).then(res => res.blob());
}
