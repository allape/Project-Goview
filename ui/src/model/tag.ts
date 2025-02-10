import { IBase } from "@allape/gocrud";

export default interface ITag extends IBase {
  name: string;
  key: string;
  color: string;
}
