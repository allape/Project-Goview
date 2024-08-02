import type { LV } from '../config/lv';
import { BASE_URL } from '../config/server';
import ajax, { type IBaseModel } from './http';
import type { IPreview } from './preview';

export type FileKey = string;

export interface IFile {
  isDir: boolean;
  name: string;
  size: number;

  _displayName: string;
  _path: string;
  _cover: string;
  _fileURL: string;
}

export interface IPreviewFile {
  key: FileKey;
  stat: IFile;
  preview?: IPreview;
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

export function stat(dsid: IDatasource['id'], file: string): Promise<IPreviewFile> {
  return ajax(`/datasource/stat/${dsid}/${file}`);
}

export async function ls(dsid: IDatasource['id'], cwd: string): Promise<IPreviewFile[]> {
  const files = await ajax<IPreviewFile[]>(`/datasource/ls/${dsid}/${cwd}`);
  return files.map((pf: IPreviewFile) => {
    const file = pf.stat;
    file._displayName = `${file.name}${file.isDir ? '/' : ''}`;
    file._path = `${cwd}${encodeURIComponent(file._displayName)}`;
    file._cover = pf.preview?.cover ? `${BASE_URL}/preview/static/${pf.preview.cover}` : `${BASE_URL}/preview/image/no-preview.jpg`;
    file._fileURL = `${BASE_URL}/datasource/cat/${dsid}/${file._path}`;
    return pf;
  });
}

export function cat(dsid: IDatasource['id'], file: string): Promise<Blob> {
  return fetch(`/datasource/cat/${dsid}/${file}`).then(res => res.blob());
}
