import { IBase } from "@allape/gocrud";
import { ILV } from "@allape/gocrud-react/src/helper/antd.tsx";

export type DatasourceType = "dufs" | "local";

export default interface IDatasource extends IBase {
  name: string;
  type: DatasourceType;
  cwd: string;
}

export interface IFileInfo {
  name: string;
  isDir: boolean;
  size: number;
  mtime: number;
}

export const DatasourceTypes: ILV<DatasourceType>[] = [
  {
    value: "dufs",
    label: "Dufs",
  },
  {
    value: "local",
    label: "Local",
  },
];
