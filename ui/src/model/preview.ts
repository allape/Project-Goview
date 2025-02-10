import { IBase } from '@allape/gocrud';

export default interface IPreview extends IBase {
  datasourceId: string;
  key: string;
  digest: string;
  cover: string;
  ffprobeInfo: string;
}
